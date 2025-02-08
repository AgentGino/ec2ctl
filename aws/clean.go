package aws

import (
	"os"

	logger "ec2ctl/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"
)

func Clean(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadDefaultConfig(cmd.Context())
	if err != nil {
		logger.Error("failed to load configuration: %v", err)
	}

	client := ec2.NewFromConfig(cfg)

	// Get all instances with tags and state running
	instances, err := client.DescribeInstances(cmd.Context(), &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:Project"),
				Values: []string{"ec2ctl"},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"running"},
			},
		},
	})
	if err != nil {
		logger.Error("failed to describe instances: %v", err)
	}

	// Delete all instances
	for _, reservation := range instances.Reservations {
		logger.Info("Found %d instances", len(reservation.Instances))
		for _, instance := range reservation.Instances {
			DeleteInstance(cmd, []string{aws.ToString(instance.Tags[0].Value)})
			logger.Info("Terminated instance %s", aws.ToString(instance.InstanceId))
		}
	}

	// Delete the security group
	client.DeleteSecurityGroup(cmd.Context(), &ec2.DeleteSecurityGroupInput{
		GroupId: aws.String("ec2ctl-sg"),
	})
	if err != nil {
		logger.Error("failed to delete security group: %v", err)
	}

	logger.Success("All instances and security group have been deleted")
	os.Remove(".ec2ctl.json")

}
