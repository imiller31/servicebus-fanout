package cmd

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/sirupsen/logrus"
)

func ensureTopic(ctx context.Context, client *admin.Client, topicName string) error {
	topicResp, err := client.GetTopic(ctx, topicName, nil)
	if err != nil {
		return fmt.Errorf("failed to get topic: %w", err)
	}

	if topicResp == nil {
		logrus.Infof("topic %s does not exist, creating it", topicName)
		_, err = client.CreateTopic(ctx, topicName, nil)
		if err != nil {
			return fmt.Errorf("failed to create leaf topic: %w", err)
		}
	}
	return nil
}

func ensureSubscription(ctx context.Context, client *admin.Client, topicName, subscriptionName string, opts *admin.CreateSubscriptionOptions) error {
	subResp, err := client.GetSubscription(ctx, topicName, subscriptionName, nil)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}
	if subResp == nil {
		logrus.Infof("subscription %s does not exist, creating it", subscriptionName)
		_, err = client.CreateSubscription(ctx, topicName, subscriptionName, opts)
		if err != nil {
			return fmt.Errorf("failed to create subscription: %w", err)
		}
	}
	return nil
}
