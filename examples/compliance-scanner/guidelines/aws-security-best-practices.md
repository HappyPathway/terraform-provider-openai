# AWS Security Best Practices

## Identity and Access Management (IAM)

- Use IAM roles instead of long-term access keys
- Implement least privilege access
- Rotate credentials regularly
- Enable MFA for all users
- Use AWS Organizations for multi-account management

## Network Security

- Use VPCs to isolate resources
- Implement security groups and NACLs
- Enable VPC Flow Logs
- Use AWS PrivateLink for service connections
- Enable GuardDuty for threat detection

## Data Protection

- Encrypt EBS volumes by default
- Enable S3 bucket encryption
- Use KMS for key management
- Enable CloudTrail for API logging
- Regular security assessments with AWS Config

## Compute Security

- Keep AMIs up to date
- Use Systems Manager for patch management
- Enable detailed monitoring
- Use security groups as firewalls
- Regular vulnerability assessments

## Storage Security

- Disable public access to S3 buckets
- Enable versioning for critical buckets
- Use presigned URLs for temporary access
- Regular backup and testing
- Enable access logging

## Database Security

- Use encryption at rest
- Enable automated backups
- Use IAM authentication
- Regular security patching
- Network isolation with security groups

## Monitoring and Logging

- Enable CloudWatch monitoring
- Set up SNS alerts
- Regular log analysis
- Enable AWS Config
- Use Security Hub for centralized security view
