package aws

import (
	"ec2ctl/utils"
	"encoding/json"
	"fmt"
	"os"
	"time"

	logger "ec2ctl/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"
)

func CreateInstance(cmd *cobra.Command, args []string) {
	// Get the name from positional argument
	name := args[0]

	// Get flags
	instanceType, _ := cmd.Flags().GetString("instance-type")
	ami, _ := cmd.Flags().GetString("ami")

	logger.Info("Creating instance with name: %s, instance type: %s, AMI: %s",
		name, instanceType, ami)

	if instanceType == "" {
		instanceType = "t2.micro"
	}

	if ami == "" {
		ami = "ami-0440d3b780d96b29d"
	}

	cfg, err := config.LoadDefaultConfig(cmd.Context())
	if err != nil {
		logger.Error("failed to load configuration, %v", err)
	}

	client := ec2.NewFromConfig(cfg)

	// Read security group ID from config file
	configBytes, err := os.ReadFile(".ec2ctl.json")
	if err != nil {

		Init(cmd, []string{})

	}

	logger.Info("Generating SSH key pair")
	// Generate SSH key pair
	privateKey, publicKey, err := utils.GenerateSSHKeyPair()
	if err != nil {
		logger.Error("failed to generate SSH key pair, %v", err)
	}

	// Create SSH key pair in AWS
	keyName := fmt.Sprintf("ec2ctl-%s", name)
	createKeyPairInput := &ec2.ImportKeyPairInput{
		KeyName:           aws.String(keyName),
		PublicKeyMaterial: []byte(publicKey),
	}

	_, err = client.ImportKeyPair(cmd.Context(), createKeyPairInput)
	if err != nil {
		logger.Error("failed to create SSH key pair in AWS, %v", err)
	}

	logger.Success("SSH key pair created")
	// Save private key to file
	privateKeyPath := fmt.Sprintf("%s.pem", name)
	err = os.WriteFile(privateKeyPath, []byte(privateKey), 0600)
	if err != nil {
		logger.Error("failed to save private key, %v", err)
	}

	var config struct {
		SecurityGroupID string `json:"security_group_id"`
	}
	if err := json.Unmarshal(configBytes, &config); err != nil {
		logger.Error("failed to parse config file: %v", err)
	}

	logger.Info("Creating instance")
	// Remove the security group creation code and use the existing security group
	input := &ec2.RunInstancesInput{
		ImageId:      aws.String(ami),
		InstanceType: types.InstanceType(instanceType),
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
		KeyName:      aws.String(keyName),
		SecurityGroupIds: []string{
			config.SecurityGroupID,
		},
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstance,
				Tags: []types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(name),
					},
					{
						Key:   aws.String("Project"),
						Value: aws.String("ec2ctl"),
					},
				},
			},
		},
	}

	result, err := client.RunInstances(cmd.Context(), input)
	if err != nil {
		logger.Error("failed to create instance, %v", err)
	}

	instanceId := aws.ToString(result.Instances[0].InstanceId)

	logger.Info("Waiting for instance to be running")
	// Wait for the instance to be running
	waiter := ec2.NewInstanceRunningWaiter(client)
	err = waiter.Wait(cmd.Context(), &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceId},
	}, 5*time.Minute)
	if err != nil {
		logger.Error("failed to wait for instance to be running, %v", err)
	}

	logger.Success("Instance is running")

	// Get the instance's public IP address
	describeInput := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceId},
	}
	describeResult, err := client.DescribeInstances(cmd.Context(), describeInput)
	if err != nil {
		logger.Error("failed to describe instance, %v", err)
	}

	publicIp := aws.ToString(describeResult.Reservations[0].Instances[0].PublicIpAddress)

	logger.Info("Generating SSH config entry")
	// Generate SSH config entry
	sshConfigEntry := fmt.Sprintf(`
Host %s
	HostName %s
	User ec2-user
	IdentityFile %s
`, name, publicIp, privateKeyPath)

	// Write SSH config entry to file
	err = os.WriteFile(fmt.Sprintf("%s.config", name), []byte(sshConfigEntry), 0600)
	if err != nil {
		logger.Error("failed to write SSH config entry, %v", err)
	}

	logger.Success("Created instance %s with SSH key %s\nSSH config saved to %s.config\n", instanceId, keyName, name)
}
