package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_CostGoogleBudget(t *testing.T) {
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
		google = {
		source = "hashicorp/google"
		version = "5.31.1"
		}
	}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewCostGoogleBudgetRule(),
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
		google = {
		source = "hashicorp/google"
		version = "5.31.1"
		}
	}
}

resource "google_billing_budget" "budget" {
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewCostGoogleBudgetRule()

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
