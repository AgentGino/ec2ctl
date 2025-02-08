package aws

import (
	logger "ec2ctl/logger"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"
)

func ListInstances(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadDefaultConfig(cmd.Context())
	if err != nil {
		logger.Error("failed to load configuration, %v", err)
	}

	client := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:Project"),
				Values: []string{"ec2ctl"},
			},
		},
	}

	result, err := client.DescribeInstances(cmd.Context(), input)
	if err != nil {
		logger.Error("failed to describe instances, %v", err)
	}

	logger.Info("%-20s %-20s %-20s %-20s\n", "Name", "Instance ID", "State", "Type")
	logger.Info(strings.Repeat("-", 80))

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instanceName := ""
			for _, tag := range instance.Tags {
				if aws.ToString(tag.Key) == "Name" {
					instanceName = aws.ToString(tag.Value)
					break
				}
			}
			logger.Info("%-20s %-20s %-20s %-20s\n",
				instanceName,
				aws.ToString(instance.InstanceId),
				string(instance.State.Name),
				string(instance.InstanceType))
		}
	}
}
