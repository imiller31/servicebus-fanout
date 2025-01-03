package cmd

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/go-shuttle/v2"
	"github.com/imiller31/servicebus-fanout/pkg/streams"
	"github.com/imiller31/servicebus-fanout/protos"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const magicShardNum = 3

var processorName string
var leafTopicName string

// notificationMgrCmd represents the run command
var notificationMgrCmd = &cobra.Command{
	Use:   "notification-manager",
	Short: "the notification manager listens to a subscription on a leaf topic and handles the message lifecycle while forwarding to the appropriate client",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Infof("starting...")
		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}
		connectionString, _ := cmd.Flags().GetString("connection-string")
		primaryTopicName, _ := cmd.Flags().GetString("primary-topic-name")
		processorName, _ = cmd.Flags().GetString("name")
		runNotificationMgr(connectionString, primaryTopicName, processorName)
	},
}

func runNotificationMgr(connectionString, primaryTopicName, processorName string) {
	logrus.Infof("starting run with connection string %s, primary topic name %s, processor name %s", connectionString, primaryTopicName, processorName)

	// create a new context
	ctx := context.Background()
	var err error
	// compute the leaf topic name
	leafTopicName, err = computeLeafTopicName(processorName)
	if err != nil {
		logrus.Fatalf("failed to compute leaf topic name: %v", err)
	}

	// Create a new service bus client
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		logrus.Fatalf("failed to create credential: %w", err)
	}
	// ensure the leaf topic and autofwd subscription exist
	if err := ensureServiceBusAutoForwardTopicExists(ctx, credential, connectionString, primaryTopicName, leafTopicName, processorName); err != nil {
		logrus.Fatalf("failed to ensure service bus auto forward topic exists: %v", err)
	}

	// Create a new service bus client
	client, err := azservicebus.NewClient(connectionString, credential, nil)
	if err != nil {
		logrus.Fatalf("failed to create service bus client: %v", err)
	}

	// create a new stream server to handle client registrations
	streamServer := streams.NewStreamServer(ctx, client)
	go runStreamServer(ctx, streamServer)

	// listen to the leaf topic
	if err := listenToServiceBusAutoForwardTopic(ctx, client, leafTopicName, processorName, streamServer.HandleMessage); err != nil {
		logrus.Fatalf("failed to listen to service bus auto forward topic: %v", err)
	}
}

func runStreamServer(ctx context.Context, streamServer *streams.StreamServer) {
	// dial server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	protos.RegisterStreamServiceServer(grpcServer, streamServer)

	logrus.Infof("starting grpc server")
	grpcServer.Serve(lis)

	logrus.Infof("grpc server stopped")
}

func listenToServiceBusAutoForwardTopic(ctx context.Context, client *azservicebus.Client, leafTopic, processorName string, handlerFunc shuttle.HandlerFunc) error {
	// Create a new receiver on the processor's subscription
	receiver, err := client.NewReceiverForSubscription(leafTopic, processorName, nil)
	if err != nil {
		return fmt.Errorf("failed to create receiver: %w", err)
	}

	// Create a new go-shuttle processor and listen to it
	processor := shuttle.NewProcessor(receiver, handlerFunc, &shuttle.ProcessorOptions{
		MaxConcurrency: 1,
	})

	logrus.Infof("starting processor")
	if err := processor.Start(ctx); err != nil {
		return fmt.Errorf("failed to start processor: %w", err)
	}

	return nil
}

func ensureServiceBusAutoForwardTopicExists(ctx context.Context, creds azcore.TokenCredential, connectionString, rootTopicName, leafName, processorName string) error {
	client, err := admin.NewClient(connectionString, creds, nil)
	if err != nil {
		return fmt.Errorf("failed to create admin service bus client: %w", err)
	}

	if err := ensureTopic(ctx, client, leafName); err != nil {
		return fmt.Errorf("failed to ensure leaf topic: %w", err)
	}

	if err := ensureSubscription(ctx, client, leafName, processorName, nil); err != nil {
		return fmt.Errorf("failed to ensure subscription: %w", err)
	}

	autoFwdOpts := &admin.CreateSubscriptionOptions{
		Properties: &admin.SubscriptionProperties{
			ForwardTo: to.Ptr(fmt.Sprintf("sb://%s/%s", connectionString, leafName)),
		},
	}
	if err := ensureSubscription(ctx, client, rootTopicName, fmt.Sprintf("autofwd-sub-%s", leafName), autoFwdOpts); err != nil {
		return fmt.Errorf("failed to ensure auto-forward subscription: %w", err)
	}

	logrus.Infof("service bus auto forward topic exists")
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
	rootCmd.AddCommand(notificationMgrCmd)

	notificationMgrCmd.Flags().StringP("connection-string", "c", "", "The connection string for the service bus")
	notificationMgrCmd.Flags().StringP("primary-topic-name", "t", "", "The primary topic name")
	notificationMgrCmd.Flags().StringP("name", "n", "", "The name of the process")
}
