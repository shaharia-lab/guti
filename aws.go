package guti

import "fmt"

// AWSGenerateARN generates an Amazon Resource Name (ARN) for the specified AWS resource.
// The function takes a service name, region, AWS account ID, resource type, and resource name as input parameters.
// Depending on the service and resource type, the function generates an ARN using the following format:
//
// S3 bucket: arn:aws:s3:::{bucket-name}
// S3 object: arn:aws:s3:::{bucket-name}/{region}
// EC2 instance: arn:aws:ec2:{region}:{account-id}:instance:{instance-id}
// RDS database: arn:aws:rds:{region}:{account-id}:db:{database-name}
// ECR repository: arn:aws:ecr:{region}:{account-id}:repository/{repository-name}
//
// If an invalid service or resource type is provided, the function returns an empty string and an error.
func AWSGenerateARN(service, region, accountID, resourceType, resourceName string) (string, error) {
	switch service {
	case "s3":
		if resourceType == "bucket" {
			return fmt.Sprintf("arn:aws:%s:::%s", service, resourceName), nil
		} else if resourceType == "object" {
			return fmt.Sprintf("arn:aws:%s:::%s/%s", service, resourceName, region), nil
		} else {
			return "", fmt.Errorf("invalid resource type for S3")
		}
	case "ec2":
		if resourceType == "instance" {
			return fmt.Sprintf("arn:aws:%s:%s:%s:%s:%s", service, region, accountID, "instance", resourceName), nil
		} else {
			return "", fmt.Errorf("invalid resource type for EC2")
		}
	case "rds":
		if resourceType == "db" {
			return fmt.Sprintf("arn:aws:%s:%s:%s:%s:%s", service, region, accountID, "db", resourceName), nil
		} else {
			return "", fmt.Errorf("invalid resource type for RDS")
		}
	case "ecr":
		if resourceType == "repository" {
			return fmt.Sprintf("arn:aws:%s:%s:%s:%s/%s", service, region, accountID, "repository", resourceName), nil
		} else {
			return "", fmt.Errorf("invalid resource type for ECR")
		}
	default:
		return "", fmt.Errorf("invalid service")
	}
}
