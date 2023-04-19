package guti

import "fmt"

func AWSResourceARNBuilder(service, region, accountID, resourceType, resourceName string) (string, error) {
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
