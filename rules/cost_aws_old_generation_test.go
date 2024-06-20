package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_CostAwsOldGenerationRule(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "old instance type defined",
			Content: `
resource "aws_instance" "instance" {
  instance_type = "t2.micro"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewCostAwsOldGenerationRule(),
					Message: "older generation instances might have worse performance and cost more than newer generations",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 29},
					},
				},
			},
		},
		{
			Name: "old volume type defined",
			Content: `
resource "aws_instance" "instance" {
  root_block_device {
    volume_type = "gp2"
  }
}

resource "aws_ebs_volume" "volume" {
  type = "gp2"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewCostAwsOldGenerationRule(),
					Message: "older generation volumes might have worse performance and cost more than newer generations",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 5},
						End:      hcl.Pos{Line: 4, Column: 24},
					},
				},
				{
					Rule:    NewCostAwsOldGenerationRule(),
					Message: "older generation volumes might have worse performance and cost more than newer generations",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 9, Column: 3},
						End:      hcl.Pos{Line: 9, Column: 15},
					},
				},
			},
		},
		{
			Name: "no issues found",
			Content: `
resource "aws_instance" "instance" {
  instance_type = "t3.micro"
  root_block_device {
    volume_type = "gp3"
  }
}

resource "aws_ebs_volume" "volume" {
  type = "gp3"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewCostAwsOldGenerationRule()

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
