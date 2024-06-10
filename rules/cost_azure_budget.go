package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// CostAwsBudgetRule checks whether ...
type CostAzureBudgetRule struct {
	tflint.DefaultRule
}

// NewCostAwsBudgetRule returns a new rule
func NewCostAzureBudgetRule() *CostAzureBudgetRule {
	return &CostAzureBudgetRule{}
}

// Name returns the rule name
func (r *CostAzureBudgetRule) Name() string {
	return "cost_azure_budget"
}

// Enabled returns whether the rule is enabled by default
func (r *CostAzureBudgetRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *CostAzureBudgetRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *CostAzureBudgetRule) Link() string {
	return "https://search-rug.github.io/iac-cost-patterns/budget/"
}

// Check checks whether ...
func (r *CostAzureBudgetRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent("azurerm_consumption_budget_subscription", &hclext.BodySchema{}, nil)
	if err != nil {
		return err
	}

	if len(resources.Blocks) == 0 {
		tfBlocks, err := runner.GetModuleContent(&hclext.BodySchema{
			Blocks: []hclext.BlockSchema{
				{
					Type: "terraform",
					Body: &hclext.BodySchema{
						Blocks: []hclext.BlockSchema{
							{
								Type: "required_providers",
								Body: &hclext.BodySchema{
									Blocks: []hclext.BlockSchema{
										{
											Type: "azurerm",
											Body: &hclext.BodySchema{},
										},
									},
								},
							},
						},
					},
				},
			},
		}, nil)
		if err != nil {
			return err
		}

		if len(tfBlocks.Blocks) < 1 {
			logger.Warn("no terraform block defined")

			return nil
		}

		err = runner.EmitIssue(r, "no budget defined", tfBlocks.Blocks[0].DefRange)
		if err != nil {
			return err
		}
	}

	return nil
}
