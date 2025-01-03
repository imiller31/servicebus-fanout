package streams

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/go-shuttle/v2"
	"github.com/imiller31/servicebus-fanout/protos"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var _ protos.StreamServiceServer = StreamServer{}

type StreamServer struct {
	protos.UnimplementedStreamServiceServer

	serverCtx         context.Context
	RegisteredStreams *StreamManager
	Responses         *ResponseManager
	sbClient          *azservicebus.Client
}

func NewStreamServer(ctx context.Context, sbClient *azservicebus.Client) *StreamServer {
	return &StreamServer{
		serverCtx:         ctx,
		RegisteredStreams: NewStreamManager(),
		Responses:         NewResponseManager(),
		sbClient:          sbClient,
	}
}

func (s StreamServer) GetRegisteredStream(streamType streamType) grpc.BidiStreamingServer[protos.ClientRequest, protos.NotificationForwarderWrapper] {
	// downstream client could have multiple streams open, so we need to select one randomly to send the response to
	if stream, ok := s.RegisteredStreams.GetRandomStream(streamType); ok {
		return stream
	}
	return nil
}

func (s StreamServer) GetMessages(stream grpc.BidiStreamingServer[protos.ClientRequest, protos.NotificationForwarderWrapper]) error {

	req, err := stream.Recv()
	if err != nil {
		logrus.Errorf("failed to receive message: %v", err)
		return err
	}

	clName := clientName(req.ClientName)
	clType := streamType(req.ClientType)

	if !req.IsRegistration {
		logrus.Errorf("received message that is not a registration request")
		return errors.New("received message that is not a registration request")
	}

	s.RegisteredStreams.SetStream(streamType(req.ClientType), clientName(req.ClientName), stream)
	logrus.Infof("registered stream for client %s of type %s, keeping stream open", req.ClientName, req.ClientType)
	for {
		req, err := stream.Recv()
		if err != nil {
			logrus.Errorf("stream for client %s closed with error: %s", clName, err)
			s.RegisteredStreams.DeleteStream(clType, clName)
			logrus.Infof("deleted stream for client %s", clName)
			return err
		}

		msgResponseChannel, ok := s.Responses.Get(req.MessageId)
		if !ok {
			logrus.Errorf("no response channel found for message id %s", req.MessageId)
			continue
		}
		msgResponseChannel <- req
		logrus.Tracef("sent message to response channel for message id %s", req.MessageId)
	}

}

func (s StreamServer) HandleMessage(ctx context.Context, settler shuttle.MessageSettler, msg *azservicebus.ReceivedMessage) {
	logrus.Debugf("received message: %s", msg)

	var msgType string
	var ok bool
	if msgType, ok = msg.ApplicationProperties["type"].(string); !ok {
		logrus.Errorf("message does not have a type property, completing")
		if err := settler.CompleteMessage(ctx, msg, nil); err != nil {
			logrus.Errorf("failed to settle message: %v", err)
			return
		}
	}

	stream := s.GetRegisteredStream(streamType(msgType))
	if stream == nil {
		logrus.Errorf("no stream registered for message type %s", msgType)
		if err := settler.CompleteMessage(ctx, msg, nil); err != nil {
			logrus.Errorf("failed to settle message: %v", err)
		}
		return
	}
	// we need to register a response channel for this message id
	responseChannel := make(chan *protos.ClientRequest)

	// defer cleaning up the response channel and message entry
	defer close(responseChannel)
	defer s.Responses.Delete(msg.MessageID)

	s.Responses.Set(msg.MessageID, responseChannel)
	logrus.Tracef("registered response channel for message id %s", msg.MessageID)

	forwardMsg := &protos.NotificationForwarderWrapper{
		Data:      msg.Body,
		MessageId: msg.MessageID,
	}

	// now we can send forward the message to the client for some processing
	if err := stream.Send(forwardMsg); err != nil {
		logrus.Errorf("failed to send message to client: %v", err)
	}

	logrus.Infof("sent message %s to client for processing", msg.MessageID)

	// wait for a response from the client
	resp := <-responseChannel
	logrus.Infof("received response from client: %s", resp)

	if resp.Error != nil && resp.Error.IsRetryable {
		logrus.Errorf("client asked for message %s to be retried", resp.MessageId)
		if err := settler.AbandonMessage(ctx, msg, nil); err != nil {
			logrus.Errorf("failed to abandon message: %v", err)
		}
		return
	}

	// Need to understand the overhead of creating a new sender for each response.
	// Since the client is passed down, it is unclear if the sender is on the same tcp connection. Can check with metrics easily.
	sender, err := s.sbClient.NewSender("response", nil)
	if err != nil {
		logrus.Errorf("failed to create sender: %v", err)

		// Clearly need some proper handling here, but this is just a demo
		if err := settler.AbandonMessage(ctx, msg, nil); err != nil {
			logrus.Errorf("failed to abandon message: %v", err)
		}
	}
	shuttleSender := shuttle.NewSender(sender, &shuttle.SenderOptions{Marshaller: &shuttle.DefaultProtoMarshaller{}})

	notificationResponse := &protos.NotficationResponse{
		Error:    resp.Error,
		Response: resp.Response,
	}

	if err := shuttleSender.SendMessage(ctx, notificationResponse, createResponseMsgMiddleware(resp.ClientType, resp.MessageId)); err != nil {
		logrus.Errorf("failed to send response: %v", err)
		if err := settler.AbandonMessage(ctx, msg, nil); err != nil {
			logrus.Errorf("failed to abandon message: %v", err)
		}
	}

	if err := settler.CompleteMessage(ctx, msg, nil); err != nil {
		logrus.Errorf("failed to settle message: %v", err)
	}
}

func createResponseMsgMiddleware(messageType, msgId string) func(*azservicebus.Message) error {
	return func(msg *azservicebus.Message) error {
		msg.MessageID = to.Ptr(msgId)
		msg.ApplicationProperties = map[string]interface{}{
			"type": messageType,
		}
		return nil
	}
}
