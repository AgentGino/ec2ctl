package aws

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	logger "ec2ctl/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"
)

func Init(cmd *cobra.Command, args []string) {
	allowLocal, _ := cmd.Flags().GetBool("allow-local")

	cfg, err := config.LoadDefaultConfig(cmd.Context())
	if err != nil {
		logger.Error("failed to load configuration: %v", err)
	}

	client := ec2.NewFromConfig(cfg)

	// Create security group
	sgInput := &ec2.CreateSecurityGroupInput{
		GroupName:   aws.String("ec2ctl-sg"),
		Description: aws.String("Security group for ec2ctl instances"),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeSecurityGroup,
				Tags: []types.Tag{
					{
						Key:   aws.String("Project"),
						Value: aws.String("ec2ctl"),
					},
				},
			},
		},
	}

	sgResult, err := client.CreateSecurityGroup(cmd.Context(), sgInput)
	if err != nil {
		logger.Error("failed to create security group: %v", err)
	}

	var cidrIP string
	if allowLocal {
		// Get public IP
		resp, err := http.Get("https://api.ipify.org")
		if err != nil {
			logger.Error("failed to get public IP: %v", err)
		}
		defer resp.Body.Close()

		ip, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error("failed to read public IP: %v", err)
		}
		cidrIP = fmt.Sprintf("%s/32", string(ip))
	} else {
		cidrIP = "0.0.0.0/0"
	}

	// Add inbound rules for all ports
	_, err = client.AuthorizeSecurityGroupIngress(cmd.Context(), &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: sgResult.GroupId,
		IpPermissions: []types.IpPermission{
			{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int32(1),
				ToPort:     aws.Int32(65535),
				IpRanges: []types.IpRange{
					{
						CidrIp:      aws.String(cidrIP),
						Description: aws.String("All ports access"),
					},
				},
			},
		},
	})
	if err != nil {
		logger.Error("failed to authorize security group ingress: %v", err)
	}

	// Save security group ID to a config file
	config := struct {
		SecurityGroupID string `json:"security_group_id"`
	}{
		SecurityGroupID: aws.ToString(sgResult.GroupId),
	}

	configBytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		logger.Error("failed to marshal config: %v", err)
	}

	err = os.WriteFile(".ec2ctl.json", configBytes, 0600)
	if err != nil {
		logger.Error("failed to write config file: %v", err)
	}

	logger.Success("Initialized ec2ctl with security group %s\n", aws.ToString(sgResult.GroupId))
	if allowLocal {
		logger.Info("Allowing SSH access only from %s\n", cidrIP)
	} else {
		logger.Info("Allowing SSH access from anywhere\n")
	}
}
