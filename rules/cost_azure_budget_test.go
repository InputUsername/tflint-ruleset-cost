package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_CostAzurermBudget(t *testing.T) {
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
		azurerm = {
			source = "hashicorp/azurerm"
			version = "3.107.0"
		}
	}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewCostAzureBudgetRule(),
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
		azurerm = {
			source = "hashicorp/azurerm"
			version = "3.107.0"
		}
	}
}

resource "azurerm_consumption_budget_subscription" "consumption_budget" {
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewCostAzureBudgetRule()

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
