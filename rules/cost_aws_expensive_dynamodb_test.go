package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_CostAwsExpensiveDynamoDbRule(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "billing mode not set to PAY_PER_REQUEST",
			Content: `
resource "aws_dynamodb_table" "table0" {
}

resource "aws_dynamodb_table" "table1" {
  billing_mode = "PROVISIONED"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewCostAwsExpensiveDynamoDbRule(),
					Message: "billing mode is not set to PAY_PER_REQUEST which may be expensive",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 39},
					},
				},
				{
					Rule:    NewCostAwsExpensiveDynamoDbRule(),
					Message: "billing mode is not set to PAY_PER_REQUEST which may be expensive",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 3},
						End:      hcl.Pos{Line: 6, Column: 31},
					},
				},
			},
		},
		{
			Name: "high read/write capacity",
			Content: `
resource "aws_dynamodb_table" "table0" {
  read_capacity  = 20
  write_capacity = 20
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewCostAwsExpensiveDynamoDbRule(),
					Message: "billing mode is not set to PAY_PER_REQUEST which may be expensive",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 39},
					},
				},
				{
					Rule:    NewCostAwsExpensiveDynamoDbRule(),
					Message: "high read capacity might lead to higher cost",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 22},
					},
				},
				{
					Rule:    NewCostAwsExpensiveDynamoDbRule(),
					Message: "high write capacity might lead to higher cost",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 3},
						End:      hcl.Pos{Line: 4, Column: 22},
					},
				},
			},
		},
		{
			Name: "global secondary indices used",
			Content: `
resource "aws_dynamodb_table" "table0" {
  billing_mode = "PAY_PER_REQUEST"

  global_secondary_index {
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewCostAwsExpensiveDynamoDbRule(),
					Message: "global secondary indices are expensive",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 3},
						End:      hcl.Pos{Line: 5, Column: 25},
					},
				},
			},
		},
		{
			Name: "no issues found",
			Content: `
resource "aws_dynamodb_table" "table0" {
  billing_mode = "PAY_PER_REQUEST"
}

resource "aws_dynamodb_table" "table1" {
  read_capacity  = 1
  write_capacity = 1
}`,
			Expected: helper.Issues{
				// We expect an issue because billing mode is not PAY_PER_REQUEST.
				// There shouldn't be any issue for high r/w and global secondary indices though.
				{
					Rule:    NewCostAwsExpensiveDynamoDbRule(),
					Message: "billing mode is not set to PAY_PER_REQUEST which may be expensive",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 1},
						End:      hcl.Pos{Line: 6, Column: 39},
					},
				},
			},
		},
	}

	rule := NewCostAwsExpensiveDynamoDbRule()

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
