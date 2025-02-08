package main

import (
	"ec2ctl/aws"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ec2ctl",
	Short: "EC2CTL is a CLI tool for managing AWS EC2 instances",
	Long: `A simple and efficient alternative to Multipass CLI for managing EC2 instances.
Complete documentation is available at https://github.com/yourusername/ec2ctl`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all EC2 instances",
	Run:   aws.ListInstances,
}

var createCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new EC2 instance",
	Args:  cobra.ExactArgs(1),
	Run:   aws.CreateInstance,
	Example: `  # Create instance with default settings
  ec2ctl create myinstance

  # Create instance with specific type and AMI
  ec2ctl create myinstance --instance-type t2.micro --ami ami-0440d3b780d96b29d`,
}

var deleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete an EC2 instance",
	Args:  cobra.ExactArgs(1),
	Run:   aws.DeleteInstance,
}

var sshCmd = &cobra.Command{
	Use:   "ssh [name]",
	Short: "SSH into an EC2 instance",
	Args:  cobra.ExactArgs(1),
	Run:   aws.SSHInstance,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize ec2ctl by creating a security group",
	Long:  `Creates a security group that will be used by all instances created by ec2ctl`,
	Run:   aws.Init,
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up all instances and security groups",
	Run:   aws.Clean,
}

func init() {
	// Add flags for createCmd
	createCmd.Flags().String("instance-type", "t2.micro", "Instance type")
	createCmd.Flags().String("ami", "ami-0440d3b780d96b29d", "AMI ID")

	// Add flag for initCmd
	initCmd.Flags().Bool("allow-local", false, "Allow SSH access only from local public IP")

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(sshCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(cleanCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
