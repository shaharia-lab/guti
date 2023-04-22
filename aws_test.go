package guti

import (
	"testing"
)

func TestAWSGenerateARN(t *testing.T) {
	tests := []struct {
		service      string
		region       string
		accountID    string
		resourceType string
		resourceName string
		expectedARN  string
		expectedErr  error
	}{
		// S3 bucket test case
		{
			service:      "s3",
			resourceType: "bucket",
			resourceName: "my-bucket",
			expectedARN:  "arn:aws:s3:::my-bucket",
			expectedErr:  nil,
		},
		// EC2 instance test case
		{
			service:      "ec2",
			region:       "us-west-2",
			accountID:    "123456789012",
			resourceType: "instance",
			resourceName: "i-1234567890abcdef0",
			expectedARN:  "arn:aws:ec2:us-west-2:123456789012:instance:i-1234567890abcdef0",
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		actualARN, actualErr := AWSGenerateARN(tt.service, tt.region, tt.accountID, tt.resourceType, tt.resourceName)

		if actualARN != tt.expectedARN {
			t.Errorf("unexpected ARN; want %s, got %s", tt.expectedARN, actualARN)
		}

		if actualErr != tt.expectedErr {
			t.Errorf("unexpected error; want %v, got %v", tt.expectedErr, actualErr)
		}
	}
}
