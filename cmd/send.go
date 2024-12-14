package cmd

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/go-shuttle/v2"
	"github.com/imiller31/servicebus-fanout/protos"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "send a message to the primary topic",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("send called")
		processorName := args[0]
		connectionString, _ := cmd.Flags().GetString("connection-string")
		primaryTopicName, _ := cmd.Flags().GetString("primary-topic-name")
		msgType, _ := cmd.Flags().GetString("type")
		sendHello(connectionString, primaryTopicName, processorName, msgType)
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringP("connection-string", "c", "", "The connection string for the service bus")
	sendCmd.Flags().StringP("primary-topic-name", "t", "", "The primary topic name")
	sendCmd.Flags().StringP("type", "y", "", "The type of message to send")
}

func sendHello(connectionString, primaryTopicName, processorName, msgType string) {
	// create a servicebus sender
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		logrus.Fatalf("failed to create credential: %v", err)
	}

	client, err := azservicebus.NewClient(connectionString, credential, nil)
	if err != nil {
		logrus.Fatalf("failed to create service bus client: %v", err)
	}

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
		Message:         "Hello from the sender",
		Type:            msgType,
	}

	err = shuttleSender.SendMessage(context.Background(), msg)
	if err != nil {
		logrus.Fatalf("failed to send message: %v", err)
	}
	logrus.Infof("sent message to %s via %s", processorName, primaryTopicName)
}
