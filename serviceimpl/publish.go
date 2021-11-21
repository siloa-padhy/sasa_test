package serviceimpl

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"main.go/dbservice"
)

func Publishstatusmessage(projectID, topicID string, transaction dbservice.Upitransaction, username dbservice.Userapimap) error {
	fmt.Println(username.User_name)
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		fmt.Println("Error is :", err)
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	t := client.Topic(topicID)
	// fmt.Println("create new topic", t)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte("STATUS : " + transaction.Status + "TxnStatusCode :" + transaction.Transaction_status_code),
	})
	id, err := result.Get(ctx)
	if err != nil {
		fmt.Println("Error in getting msg id  ") // Getting Error in getting id of the published messages
		return err
	}
	fmt.Println("Published a message; msg ID:", id)
	return nil
}
