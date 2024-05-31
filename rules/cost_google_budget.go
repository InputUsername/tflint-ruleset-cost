package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// CostGoogleBudgetRule checks whether ...
type CostGoogleBudgetRule struct {
	tflint.DefaultRule
}

// NewCostGoogleBudgetRule returns a new rule
func NewCostGoogleBudgetRule() *CostGoogleBudgetRule {
	return &CostGoogleBudgetRule{}
}

// Name returns the rule name
func (r *CostGoogleBudgetRule) Name() string {
	return "cost_google_budget"
}

// Enabled returns whether the rule is enabled by default
func (r *CostGoogleBudgetRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *CostGoogleBudgetRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *CostGoogleBudgetRule) Link() string {
	return "https://search-rug.github.io/iac-cost-patterns/budget/"
}

// Check checks whether ...
func (r *CostGoogleBudgetRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent("google_billing_budget", &hclext.BodySchema{}, nil)
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
											Type: "google",
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

		runner.EmitIssue(r, "no budget defined", tfBlocks.Blocks[0].DefRange)
	}

	return nil
}
