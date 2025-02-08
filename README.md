# EC2CTL

EC2CTL is a CLI tool for managing AWS EC2 instances, designed as an alternative to Multipass CLI.

## Installation

1. Install Go (version 1.20 or higher)
2. Clone this repository
3. Build the project:
   ```bash
   go build -o ec2ctl
   ```
4. Move the binary to your PATH:
   ```bash
   sudo mv ec2ctl /usr/local/bin/
   ```

## Usage

```bash
# List all EC2 instances
ec2ctl list

# Create a new EC2 instance
ec2ctl create

# Delete an EC2 instance
ec2ctl delete [instance-id]
```

## Configuration

The tool uses the default AWS credentials and configuration. Make sure you have:
1. AWS CLI installed and configured
2. Proper IAM permissions for EC2 management

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

MIT