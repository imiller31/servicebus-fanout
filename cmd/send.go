package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "send a message to the primary topic",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("send called")
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringP("connection-string", "c", "", "The connection string for the service bus")
	sendCmd.Flags().StringP("primary-topic-name", "t", "", "The primary topic name")
}

func sendHello() {
	// create
}
