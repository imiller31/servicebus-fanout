package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/go-shuttle/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const magicShardNum = 20

var processorName string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Infof("starting...")
		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}
		connectionString, _ := cmd.Flags().GetString("connection-string")
		primaryTopicName, _ := cmd.Flags().GetString("primary-topic-name")
		processorName, _ = cmd.Flags().GetString("name")
		run(connectionString, primaryTopicName, processorName)
	},
}

func run(connectionString, primaryTopicName, processorName string) {
	logrus.Infof("starting run with connection string %s, primary topic name %s, processor name %s", connectionString, primaryTopicName, processorName)

	// create a new context
	ctx := context.Background()

	// compute the leaf topic name
	leafTopicName, err := computeLeafTopicName(processorName)
	if err != nil {
		logrus.Fatalf("failed to compute leaf topic name: %v", err)
	}

	// ensure the leaf topic and autofwd subscription exist
	if err := ensureServiceBusAutoForwardTopicExists(ctx, connectionString, primaryTopicName, leafTopicName, processorName); err != nil {
		logrus.Fatalf("failed to ensure service bus auto forward topic exists: %v", err)
	}

	// listen to the leaf topic
	if err := listenToServiceBusAutoForwardTopic(ctx, connectionString, leafTopicName, processorName); err != nil {
		logrus.Fatalf("failed to listen to service bus auto forward topic: %v", err)
	}
}

func listenToServiceBusAutoForwardTopic(ctx context.Context, namespace, leafTopic, processorName string) error {
	// Create a new service bus client
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return fmt.Errorf("failed to create credential: %w", err)
	}

	client, err := azservicebus.NewClient(namespace, credential, nil)
	if err != nil {
		return fmt.Errorf("failed to create service bus client: %w", err)
	}

	// Create a new reciever on the processor's subscription
	reciever, err := client.NewReceiverForSubscription(leafTopic, processorName, nil)
	if err != nil {
		return fmt.Errorf("failed to create receiver: %w", err)
	}

	// Create a new go-shuttle processor and listen to it
	processor := shuttle.NewProcessor(reciever, handleMessage, nil)
	if err := processor.Start(ctx); err != nil {
		return fmt.Errorf("failed to start processor: %w", err)
	}

	return nil
}

func handleMessage(ctx context.Context, settler shuttle.MessageSettler, msg *azservicebus.ReceivedMessage) {
	// check if the message is for this processor
	logrus.Debugf("message received: %s", string(msg.Body))
	var recievedMsg Message
	if err := json.Unmarshal(msg.Body, &recievedMsg); err != nil {
		logrus.Errorf("failed to unmarshal message: %v", err)
		return
	}

	logrus.Infof("message received: %s", recievedMsg)

	if recievedMsg.To != processorName {
		// this message is not for this processor, settle it
		logrus.Infof("message is not for this processor %s, settling it: %s", processorName, recievedMsg.To)
	} else {
		//otherwise, wait 5 seconds and say hello
		time.Sleep(5 * time.Second)
		logrus.Infof("hello from processor %s", processorName)
	}

	if err := settler.CompleteMessage(ctx, msg, nil); err != nil {
		logrus.Errorf("failed to settle message: %s", err)
	}
	return
}

func ensureServiceBusAutoForwardTopicExists(ctx context.Context, namespace, primaryTopicName, leafName, processorName string) error {
	// Create a new service bus client
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return fmt.Errorf("failed to create credential: %w", err)
	}
	logrus.Debugf("creds: %s", credential)
	client, err := admin.NewClient(namespace, credential, nil)
	if err != nil {
		return fmt.Errorf("failed to create admin service bus client: %w", err)
	}

	// Create a new leaf topic if it doesn't exist
	topicResp, err := client.GetTopic(ctx, leafName, nil)
	if err != nil {
		return fmt.Errorf("failed to get leaf topic: %w", err)
	}

	if topicResp == nil {
		logrus.Infof("leaf topic %s does not exist, creating it", leafName)
		_, err = client.CreateTopic(ctx, leafName, nil)
		if err != nil {
			return fmt.Errorf("failed to create leaf topic: %w", err)
		}
	}

	// Create a new subscription on the leaf topic if it doesn't exist
	leafSubResp, err := client.GetSubscription(ctx, leafName, processorName, nil)
	if err != nil {
		return fmt.Errorf("failed to get leaf subscription: %w", err)
	}
	if leafSubResp == nil {
		logrus.Infof("leaf subscription %s does not exist, creating it", processorName)
		_, err = client.CreateSubscription(ctx, leafName, processorName, nil)
		if err != nil {
			return fmt.Errorf("failed to create leaf subscription: %w", err)
		}
	}

	// Create a new subscription to autoforward messages from the primary topic to the new topic
	fwdSubResp, err := client.GetSubscription(ctx, primaryTopicName, fmt.Sprintf("autofwd-sub-%s", leafName), nil)
	if err != nil {
		return fmt.Errorf("failed to get auto-forward subscription: %w", err)
	}
	if fwdSubResp == nil {
		logrus.Infof("auto-forward subscription from %s to %s does not exist, creating it", primaryTopicName, leafName)
		_, err = client.CreateSubscription(
			ctx,
			primaryTopicName,
			fmt.Sprintf("autofwd-sub-%s", leafName),
			&admin.CreateSubscriptionOptions{
				Properties: &admin.SubscriptionProperties{
					ForwardTo: to.Ptr(fmt.Sprintf("sb://%s/%s", namespace, leafName)),
				},
			})
		if err != nil {
			return fmt.Errorf("failed to create auto-forward subscription: %w", err)
		}
	}
	return nil
}

func computeLeafTopicName(processorName string) (string, error) {
	logrus.Infof("computing leaf topic name for processor %s", processorName)
	// split the name by the dash, take the last element and convert it to an int
	nameParts := strings.Split(processorName, "-")
	suffix, err := strconv.Atoi(nameParts[len(nameParts)-1])
	if err != nil {
		return "", fmt.Errorf("failed to convert suffix to int: %w", err)
	}

	moddedSuffix := suffix % magicShardNum

	return fmt.Sprintf("leaf-%d", moddedSuffix), nil
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringP("connection-string", "c", "", "The connection string for the service bus")
	runCmd.Flags().StringP("primary-topic-name", "t", "", "The primary topic name")
	runCmd.Flags().StringP("name", "n", "", "The name of the process")
}
