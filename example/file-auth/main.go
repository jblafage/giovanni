package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/blob/containers"
)

func main() {
	log.Printf("[DEBUG] Started..")

	// NOTE: fill this in
	storageAccountName := "example"
	storageAccountKey := "example"

	log.Printf("[DEBUG] Building Client..")
	client, err := buildClient(storageAccountName, storageAccountKey)
	if err != nil {
		panic(fmt.Errorf("Error building client: %s", err))
	}

	ctx := context.TODO()
	containerName := "armauth"
	input := containers.CreateInput{
		AccessLevel: containers.Private,
		MetaData: map[string]string{
			"hello": "world",
		},
	}
	log.Printf("[DEBUG] Creating Container..")
	if _, err := client.ContainersClient.Create(ctx, storageAccountName, containerName, input); err != nil {
		panic(fmt.Errorf("Error creating container: %s", err))
	}

	log.Printf("[DEBUG] Retrieving Container..")
	container, err := client.ContainersClient.GetProperties(ctx, storageAccountName, containerName)
	if err != nil {
		panic(fmt.Errorf("Error reading properties for container: %s", err))
	}

	log.Printf("[DEBUG] MetaData: %+v", container.MetaData)
}

type Client struct {
	ContainersClient containers.Client
}

func buildClient(accountName, accountKey string) (*Client, error) {
	env, err := authentication.DetermineEnvironment(os.Getenv("ARM_ENVIRONMENT"))
	if err != nil {
		return nil, err
	}

	storageAuth := autorest.NewSharedKeyLiteAuthorizer(accountName, accountKey)
	containersClient := containers.New()
	containersClient.Client.Authorizer = storageAuth

	result := &Client{
		ContainersClient: containersClient,
	}

	return result, nil
}
