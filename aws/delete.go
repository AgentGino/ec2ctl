package aws

import (
	"fmt"
	"os"

	logger "ec2ctl/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"
)

func DeleteInstance(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadDefaultConfig(cmd.Context())
	if err != nil {
		logger.Error("failed to load configuration, %v", err)
	}

	client := ec2.NewFromConfig(cfg)

	logger.Info("Getting instance details")
	// Get instance details to find the key name
	describeInput := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []string{args[0]},
			},
		},
	}
	describeResult, err := client.DescribeInstances(cmd.Context(), describeInput)
	if err != nil {
		logger.Error("failed to describe instance, %v", err)
	}

	keyName := ""
	instanceId := ""
	for _, reservation := range describeResult.Reservations {
		for _, instance := range reservation.Instances {
			keyName = aws.ToString(instance.KeyName)
			instanceId = aws.ToString(instance.InstanceId)
			break
		}
	}

	if keyName == "" {
		logger.Error("no key name found for instance %s", args[0])
	}

	logger.Info("Deleting SSH key pair from AWS")
	// Delete the SSH key pair from AWS
	deleteKeyPairInput := &ec2.DeleteKeyPairInput{
		KeyName: aws.String(keyName),
	}
	_, err = client.DeleteKeyPair(cmd.Context(), deleteKeyPairInput)
	if err != nil {
		logger.Error("failed to delete SSH key pair in AWS, %v", err)
	}

	logger.Info("Deleting private key file")
	// Delete the private key file
	privateKeyPath := fmt.Sprintf("%s.pem", args[0])
	err = os.Remove(privateKeyPath)
	if err != nil {
		logger.Error("failed to delete private key file, %v", err)
	}

	logger.Info("Terminating instance")
	// Terminate the instance
	input := &ec2.TerminateInstancesInput{
		InstanceIds: []string{instanceId},
	}

	logger.Info("Deleting config file")
	// Delete the config file
	err = os.Remove(fmt.Sprintf("%s.config", args[0]))
	if err != nil {
		logger.Error("failed to delete config file, %v", err)
	}

	logger.Info("Terminating instance")
	result, err := client.TerminateInstances(cmd.Context(), input)
	if err != nil {
		logger.Error("failed to terminate instance, %v", err)
	}

	logger.Success("Terminated instance %s and deleted SSH key %s\n", aws.ToString(result.TerminatingInstances[0].InstanceId), keyName)
}
