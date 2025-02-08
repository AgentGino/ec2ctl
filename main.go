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

// take name, instance type, ami
var createCmd = &cobra.Command{
	Use:   "create [name] [instance-type] [ami]",
	Short: "Create a new EC2 instance",
	Args:  cobra.ExactArgs(3),
	Run:   aws.CreateInstance,
}

var deleteCmd = &cobra.Command{
	Use:   "delete [instance-name]",
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
	initCmd.Flags().Bool("allowLocal", false, "Allow SSH access only from local public IP")
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
