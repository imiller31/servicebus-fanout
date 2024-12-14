package cmd

import (
	"context"

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
}

func NewStreamServer(ctx context.Context) *StreamServer {
	return &StreamServer{
		serverCtx:         ctx,
		RegisteredStreams: NewStreamManager(),
	}
}

func (s StreamServer) GetRegisteredStream(streamType streamType) grpc.ServerStreamingServer[protos.Response] {
	// downstream client could have multiple streams open, so we need to select one randomly to send the response to
	if stream, ok := s.RegisteredStreams.GetRandomStream(streamType); ok {
		return stream
	}
	return nil
}

func (s StreamServer) GetRequests(request *protos.Request, stream grpc.ServerStreamingServer[protos.Response]) error {
	s.RegisteredStreams.SetStream(streamType(request.Type), clientName(request.ClientName), stream)

	logrus.Infof("registered stream for client %s, keeping stream open", request.ClientName)
	for {
		select {
		case <-stream.Context().Done():
			logrus.Infof("stream for client %s closed", request.ClientName)
			s.RegisteredStreams.DeleteStream(streamType(request.Type), clientName(request.ClientName))
			logrus.Infof("deleted stream for client %s", request.ClientName)
			return nil
		case <-s.serverCtx.Done():
			logrus.Errorf("server context done")
			return nil
		}
	}
}

func (s StreamServer) handleMessage(ctx context.Context, settler shuttle.MessageSettler, msg *azservicebus.ReceivedMessage) {
	logrus.Debugf("received message: %s", msg)
	var receivedMsg protos.ServiceBusMessage
	if err := proto.Unmarshal(msg.Body, &receivedMsg); err != nil {
		logrus.Errorf("failed to unmarshal message: %v", err)
		return
	}

	stream := s.GetRegisteredStream(streamType(receivedMsg.Type))
	if stream == nil {
		logrus.Errorf("no stream registered for message type %s", receivedMsg.Type)
	} else {
		//TODO: this needs to be a bi-directional stream so the client can inform us when it has finished processing the message
		//TODO: we also need to have some mechanism to ensure the correct handler gets the correpsonding response
		//TODO: we could use the message ID to correlate the response with the request, and have the handler register a channel to receive the response for the given messade ID
		if err := stream.Send(&protos.Response{Message: receivedMsg.Message}); err != nil {
			logrus.Errorf("failed to send message: %v", err)
		}
	}

	if err := settler.CompleteMessage(ctx, msg, nil); err != nil {
		logrus.Errorf("failed to settle message: %v", err)
	}
}
