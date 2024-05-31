package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// CostAwsBudgetRule checks whether ...
type CostAwsBudgetRule struct {
	tflint.DefaultRule
}

// NewCostAwsBudgetRule returns a new rule
func NewCostAwsBudgetRule() *CostAwsBudgetRule {
	return &CostAwsBudgetRule{}
}

// Name returns the rule name
func (r *CostAwsBudgetRule) Name() string {
	return "cost_aws_budget"
}

// Enabled returns whether the rule is enabled by default
func (r *CostAwsBudgetRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *CostAwsBudgetRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *CostAwsBudgetRule) Link() string {
	return "https://search-rug.github.io/iac-cost-patterns/budget/"
}

// Check checks whether ...
func (r *CostAwsBudgetRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent("aws_budgets_budget", &hclext.BodySchema{}, nil)
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
											Type: "aws",
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
