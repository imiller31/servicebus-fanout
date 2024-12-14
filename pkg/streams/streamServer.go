package streams

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/go-shuttle/v2"
	"github.com/imiller31/servicebus-fanout/protos"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

var _ protos.StreamServiceServer = StreamServer{}

type StreamServer struct {
	protos.UnimplementedStreamServiceServer

	serverCtx context.Context

	RegisteredStreams *StreamManager

	Responses *ResponseManager
}

func NewStreamServer(ctx context.Context) *StreamServer {
	return &StreamServer{
		serverCtx:         ctx,
		RegisteredStreams: NewStreamManager(),
		Responses:         NewResponseManager(),
	}
}

func (s StreamServer) GetRegisteredStream(streamType streamType) grpc.BidiStreamingServer[protos.ClientRequest, protos.ServiceBusMessage] {
	// downstream client could have multiple streams open, so we need to select one randomly to send the response to
	if stream, ok := s.RegisteredStreams.GetRandomStream(streamType); ok {
		return stream
	}
	return nil
}

func (s StreamServer) GetMessages(stream grpc.BidiStreamingServer[protos.ClientRequest, protos.ServiceBusMessage]) error {

	req, err := stream.Recv()
	if err != nil {
		logrus.Errorf("failed to receive message: %v", err)
		return err
	}

	if !req.IsRegistration {
		logrus.Errorf("received message that is not a registration request")
		return errors.New("received message that is not a registration request")
	}

	s.RegisteredStreams.SetStream(streamType(req.ClientType), clientName(req.ClientName), stream)
	logrus.Infof("registered stream for client %s of type %s, keeping stream open", req.ClientName, req.ClientType)
	for {
		req, err := stream.Recv()
		if err != nil {
			logrus.Errorf("stream for client %s closed with error: %s", req.ClientName, err)
			s.RegisteredStreams.DeleteStream(streamType(req.ClientType), clientName(req.ClientName))
			logrus.Infof("deleted stream for client %s", req.ClientName)
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
	var receivedMsg protos.ServiceBusMessage
	if err := proto.Unmarshal(msg.Body, &receivedMsg); err != nil {
		logrus.Errorf("failed to unmarshal message: %v", err)
		return
	}

	stream := s.GetRegisteredStream(streamType(receivedMsg.Type))
	if stream == nil {
		logrus.Errorf("no stream registered for message type %s", receivedMsg.Type)
		if err := settler.CompleteMessage(ctx, msg, nil); err != nil {
			logrus.Errorf("failed to settle message: %v", err)
		}
		return
	}
	// we need to register a response channel for this message id
	responseChannel := make(chan *protos.ClientRequest)
	defer close(responseChannel)

	s.Responses.Set(receivedMsg.MessageId, responseChannel)
	logrus.Tracef("registered response channel for message id %s", receivedMsg.MessageId)

	// now we can send forward the message to the client for some processing
	if err := stream.Send(&receivedMsg); err != nil {
		logrus.Errorf("failed to send message to client: %v", err)
	}

	logrus.Infof("sent message %s to client for processing", receivedMsg.MessageId)

	// wait for a response from the client
	resp := <-responseChannel
	logrus.Infof("received response from client: %s", resp)

	// remove the response channel from the manager
	s.Responses.Delete(receivedMsg.MessageId)

	if resp.Error != nil && resp.Error.IsRetryable {
		logrus.Errorf("client returned a retryable error: %s", resp.Error)
		if err := settler.AbandonMessage(ctx, msg, nil); err != nil {
			logrus.Errorf("failed to abandon message: %v", err)
		}
		return
	}

	//TODO: send response to the response topic

	if err := settler.CompleteMessage(ctx, msg, nil); err != nil {
		logrus.Errorf("failed to settle message: %v", err)
	}
}
