package guti

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateARN(t *testing.T) {
	tests := []struct {
		service       string
		region        string
		accountID     string
		resourceType  string
		resourceName  string
		expectedARN   string
		expectedError bool
	}{
		// S3 bucket test cases
		{
			service:       "s3",
			resourceType:  "bucket",
			resourceName:  "my-bucket",
			expectedARN:   "arn:aws:s3:::my-bucket",
			expectedError: false,
		},
		{
			service:       "s3",
			resourceType:  "bucket",
			resourceName:  "",
			expectedARN:   "",
			expectedError: true,
		},
		// S3 object test cases
		{
			service:       "s3",
			region:        "us-east-1",
			resourceType:  "object",
			resourceName:  "my-bucket",
			expectedARN:   "arn:aws:s3:::my-bucket/us-east-1",
			expectedError: false,
		},
		{
			service:       "s3",
			resourceType:  "object",
			resourceName:  "my-bucket/my-object",
			expectedARN:   "",
			expectedError: true,
		},
		// EC2 instance test cases
		{
			service:       "ec2",
			region:        "us-east-1",
			accountID:     "123456789012",
			resourceType:  "instance",
			resourceName:  "i-1234567890abcdef0",
			expectedARN:   "arn:aws:ec2:us-east-1:123456789012:instance:i-1234567890abcdef0",
			expectedError: false,
		},
		{
			service:       "ec2",
			resourceType:  "instance",
			resourceName:  "i-1234567890abcdef0",
			expectedARN:   "",
			expectedError: true,
		},
		// RDS DB instance test cases
		{
			service:       "rds",
			region:        "us-east-1",
			accountID:     "123456789012",
			resourceType:  "db",
			resourceName:  "my-db-instance",
			expectedARN:   "arn:aws:rds:us-east-1:123456789012:db:my-db-instance",
			expectedError: false,
		},
		{
			service:       "rds",
			resourceType:  "db",
			resourceName:  "my-db-instance",
			expectedARN:   "",
			expectedError: true,
		},
		// ECR repository test cases
		{
			service:       "ecr",
			region:        "us-east-1",
			accountID:     "123456789012",
			resourceType:  "repository",
			resourceName:  "my-repo",
			expectedARN:   "arn:aws:ecr:us-east-1:123456789012:repository/my-repo",
			expectedError: false,
		},
		{
			service:       "ecr",
			resourceType:  "repository",
			resourceName:  "my-repo",
			expectedARN:   "",
			expectedError: true,
		},
		// Invalid service test case
		{
			service:       "invalid",
			resourceType:  "bucket",
			resourceName:  "my-bucket",
			expectedARN:   "",
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("test for %s", tc.service), func(t *testing.T) {
			arn, err := AWSResourceARNBuilder(tc.service, tc.region, tc.accountID, tc.resourceType, tc.resourceName)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedARN, arn)
			}
		})
	}
}
