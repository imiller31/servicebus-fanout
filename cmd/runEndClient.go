package cmd

import (
	"context"
	"io"

	"github.com/imiller31/servicebus-fanout/protos"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// runEndClientCmd represents the runEndClient command
var runEndClientCmd = &cobra.Command{
	Use:   "run-end-client",
	Short: "runs a client that will register to the notification manager",

	Run: func(cmd *cobra.Command, args []string) {
		grpcAddress, _ := cmd.Flags().GetString("grpc-address")
		clientName, _ := cmd.Flags().GetString("client-name")
		clientType, _ := cmd.Flags().GetString("client-type")

		runEndClient(grpcAddress, clientName, clientType)
	},
}

func runEndClient(grpcAddress, clientName, clientType string) {
	logrus.Infof("running end client with grpc address %s, client name %s", grpcAddress, clientName)

	//Register the client with the notification manager
	conn, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("failed to create grpc connection: %v", err)
	}

	client := protos.NewStreamServiceClient(conn)

	logrus.Infof("registering client %s with type %s", clientName, clientType)

	registrationRequest := &protos.Request{
		ClientName: clientName,
		Type:       clientType,
	}

	requests, err := client.GetRequests(context.Background(), registrationRequest)
	if err != nil {
		logrus.Fatalf("failed to get stream: %v", err)
	}

	logrus.Infof("registered client waiting for requests from notification manager")
	for {
		resp, err := requests.Recv()
		if err == io.EOF || status.Code(err) == 14 {
			logrus.Infof("stream closed with error: %s, reconnecting", err)
			requests, err = client.GetRequests(context.Background(), registrationRequest)
			if err != nil {
				logrus.Fatalf("failed to get stream: %v", err)
			}
			continue
		}
		if err != nil {
			logrus.Fatalf("failed to receive response: %v", err)
		}

		logrus.Infof("received message: %s", resp.Message)
	}

}

func init() {
	rootCmd.AddCommand(runEndClientCmd)

	runEndClientCmd.Flags().StringP("grpc-address", "g", "", "The address of the grpc server")
	runEndClientCmd.Flags().StringP("client-name", "n", "", "The name of the client")
	runEndClientCmd.Flags().StringP("client-type", "t", "", "The type of the client")
}
