package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"log"
	"os"
	"os/user"
)

func main() {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		fmt.Println("Error connecting to aws with config")
	}

	iamClient := iam.NewFromConfig(cfg)

	ctx := context.TODO()

	cred, err := iamClient.ListAccessKeys(ctx, nil)
	if err != nil {
		fmt.Printf("Failed to get access keys - %s\n", err)
	}

	userName := *cred.AccessKeyMetadata[0].UserName

	if len(cred.AccessKeyMetadata) == 2 {
		keyToBeDeleted := getOldestAccessKey(cred)
		deleteOldestAccessKey(ctx, iamClient, keyToBeDeleted, userName)
	}

	accessKey, secretAccessKey := createNewAccessKey(ctx, iamClient, userName)
	if accessKey != "error" {
		writeToFile(accessKey, secretAccessKey)
		fmt.Println("Credentials Updated!")
	} else {
		fmt.Println("Something went wrong creating new access keys")
	}
}

func writeToFile(accessKey string, secretAccessKey string) {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Error getting user home directory %s", err)
	}
	path := usr.HomeDir + "/.aws/credentials"

	awsCredentials := "[default]\n" +
		"aws_access_key_id=" + accessKey + "\n" +
		"aws_secret_access_key=" + secretAccessKey

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		log.Fatalf("Failed creating file: %s", err)
	}

	dataWriter := bufio.NewWriter(file)

	_, err = dataWriter.WriteString(awsCredentials + "\n")
	if err != nil {
		fmt.Println("Error writing credentials to file")
	}

	err = dataWriter.Flush()
	if err != nil {
		fmt.Println(err)
	}

	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func deleteOldestAccessKey(ctx context.Context, iamClient *iam.Client, keyToBeDeleted string, userName string) {
	input := &iam.DeleteAccessKeyInput{
		AccessKeyId: aws.String(keyToBeDeleted),
		UserName:    aws.String(userName),
	}

	_, err := iamClient.DeleteAccessKey(ctx, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func createNewAccessKey(ctx context.Context, iamClient *iam.Client, userName string) (string, string) {
	input := &iam.CreateAccessKeyInput{
		UserName: aws.String(userName),
	}

	result, err := iamClient.CreateAccessKey(ctx, input)
	if err != nil {
		fmt.Println(err.Error())
		return "error", "error"
	}
	accessKey := *result.AccessKey.AccessKeyId
	secretAccessKey := *result.AccessKey.SecretAccessKey
	return accessKey, secretAccessKey
}

func getOldestAccessKey(cred *iam.ListAccessKeysOutput) string {
	accessKeyData := cred.AccessKeyMetadata
	if accessKeyData[0].CreateDate.Before(*accessKeyData[1].CreateDate) {
		return *accessKeyData[0].AccessKeyId
	}
	return *accessKeyData[1].AccessKeyId
}
