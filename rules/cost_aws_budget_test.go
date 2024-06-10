package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_CostAwsBudget(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "issue found",
			Content: `
terraform {
	required_providers {
		aws {
			source = "hashicorp/aws"
			version = "5.51.1"
		}
	}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewCostAwsBudgetRule(),
					Message: "no budget defined",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 10},
					},
				},
			},
		},
		{
			Name: "no issue found",
			Content: `
terraform {
	required_providers {
		aws {
			source = "hashicorp/aws"
			version = "5.51.1"
		}
	}
}

resource "aws_budgets_budget" "budget" {
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewCostAwsBudgetRule()

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
