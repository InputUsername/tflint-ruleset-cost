package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_CostAwsObjectStorageLifecycleRule(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "issue found",
			Content: `
resource "aws_s3_bucket" "bucket0" {
}`,
			Expected: helper.Issues{
				{
					Rule:    NewCostAwsObjectStorageLifecycleRule(),
					Message: "bucket `bucket0` does not have a lifecycle configuration or lifecycle rule",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 35},
					},
				},
			},
		},
		{
			Name: "lifecycle configuration defined",
			Content: `
resource "aws_s3_bucket" "bucket1" {
}

resource "aws_s3_bucket_lifecycle_configuration" "bucket1-config" {
  bucket = "bucket1"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "lifecycle rule defined",
			Content: `
resource "aws_s3_bucket" "bucket2" {
  lifecycle_rule {
  }
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewCostAwsObjectStorageLifecycleRule()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": test.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, test.Expected, runner.Issues)
		})
	}
}
