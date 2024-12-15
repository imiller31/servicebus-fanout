package cmd

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/go-shuttle/v2"
	"github.com/google/uuid"
	"github.com/imiller31/servicebus-fanout/protos"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
)

var l bool

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "send a message to the primary topic",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("send called")
		processorName := args[0]
		connectionString, _ := cmd.Flags().GetString("connection-string")
		primaryTopicName, _ := cmd.Flags().GetString("primary-topic-name")
		msgType, _ := cmd.Flags().GetString("type")
		l, _ = cmd.Flags().GetBool("listen-indefinitely")
		sendHello(connectionString, primaryTopicName, processorName, msgType)
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringP("connection-string", "c", "", "The connection string for the service bus")
	sendCmd.Flags().StringP("primary-topic-name", "t", "", "The primary topic name")
	sendCmd.Flags().StringP("type", "y", "", "The type of message to send")
	sendCmd.Flags().BoolP("listen-indefinitely", "l", false, "Listen until the connection is closed")
}

func sendHello(connectionString, primaryTopicName, processorName, msgType string) {
	// create a servicebus sender
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		logrus.Fatalf("failed to create credential: %v", err)
	}

	adminClient, err := admin.NewClient(connectionString, credential, nil)
	if err != nil {
		logrus.Fatalf("failed to create service bus admin client: %v", err)
	}

	// ensure the response topic exists
	if err := ensureTopic(context.Background(), adminClient, "response"); err != nil {
		logrus.Fatalf("failed to ensure response topic exists: %v", err)
	}

	// ensure the response subscription exists
	if err := ensureSubscription(context.Background(), adminClient, "response", msgType, nil); err != nil {
		logrus.Fatalf("failed to ensure response subscription exists: %v", err)
	}

	// create the singleton client
	client, err := azservicebus.NewClient(connectionString, credential, nil)
	if err != nil {
		logrus.Fatalf("failed to create service bus client: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	responseChan := make(chan *protos.NotficationResponse)
	handler := responseHandler{responseChan: responseChan}
	go func() {
		defer wg.Done()
		if err := listenToResponseTopic(context.Background(), client, "response", msgType, handler); err != nil {
			logrus.Fatalf("failed to listen to response topic: %v", err)
		}
	}()

	sender, err := client.NewSender(primaryTopicName, nil)
	if err != nil {
		logrus.Fatalf("failed to create sender: %v", err)
	}

	shuttleSender := shuttle.NewSender(sender, &shuttle.SenderOptions{Marshaller: &shuttle.DefaultProtoMarshaller{}})

	leafTopicName, err := computeLeafTopicName(processorName)
	if err != nil {
		logrus.Fatalf("failed to compute leaf topic name: %v", err)
	}

	msg := &protos.ServiceBusMessage{
		TargetLeaf:      leafTopicName,
		TargetProcessor: processorName,
		Message:         fmt.Sprintf("Hello from the sender at %s", time.Now()),
		MessageId:       uuid.Must(uuid.NewV6()).String(),
		Type:            msgType,
	}

	err = shuttleSender.SendMessage(context.Background(), msg)
	if err != nil {
		logrus.Fatalf("failed to send message: %v", err)
	}
	logrus.Infof("sent message to %s via %s with msgId: %s, and type: %s", msg.TargetProcessor, primaryTopicName, msg.MessageId, msg.Type)

	wg.Wait()
	logrus.Infof("successfully shut down")
}

func listenToResponseTopic(
	ctx context.Context,
	client *azservicebus.Client,
	topicName, subscriptionName string,
	handler responseHandler) error {
	receiver, err := client.NewReceiverForSubscription(topicName, subscriptionName, nil)
	if err != nil {
		return fmt.Errorf("failed to create receiver: %w", err)
	}

	receiverCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	shuttleReceiver := shuttle.NewProcessor(receiver, handler.handleResponses, nil)

	go shuttleReceiver.Start(receiverCtx)

	for {
		select {
		case response := <-handler.responseChan:
			logrus.Infof("received response for msgId: %s, of type: %s, with error: %s", response.MessageId, response.Type, response.Error)
			if !l {
				break
			}
		}
		if !l {
			break
		}
	}

	logrus.Infof("shutting down")
	return nil
}

type responseHandler struct {
	responseChan chan *protos.NotficationResponse
}

func (r responseHandler) handleResponses(ctx context.Context, settler shuttle.MessageSettler, msg *azservicebus.ReceivedMessage) {
	var response protos.NotficationResponse
	if err := proto.Unmarshal(msg.Body, &response); err != nil {
		logrus.Errorf("failed to unmarshal response: %v", err)
		return
	}

	r.responseChan <- &response

	if err := settler.CompleteMessage(ctx, msg, nil); err != nil {
		logrus.Errorf("failed to settle message: %v", err)
	}
}
