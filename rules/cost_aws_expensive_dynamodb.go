package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// CostAwsExpensiveDynamoDbRule checks whether ...
type CostAwsExpensiveDynamoDbRule struct {
	tflint.DefaultRule
}

// NewCostAwsExpensiveDynamoDbRule returns a new rule
func NewCostAwsExpensiveDynamoDbRule() *CostAwsExpensiveDynamoDbRule {
	return &CostAwsExpensiveDynamoDbRule{}
}

// Name returns the rule name
func (r *CostAwsExpensiveDynamoDbRule) Name() string {
	return "cost_aws_expensive_dynamodb"
}

// Enabled returns whether the rule is enabled by default
func (r *CostAwsExpensiveDynamoDbRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *CostAwsExpensiveDynamoDbRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *CostAwsExpensiveDynamoDbRule) Link() string {
	return "https://search-rug.github.io/iac-cost-patterns/budget/"
}

// Check checks whether ...
func (r *CostAwsExpensiveDynamoDbRule) Check(runner tflint.Runner) error {
	dynamoDbTables, err := runner.GetResourceContent("aws_dynamodb_table", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: "billing_mode"},
			{Name: "read_capacity"},
			{Name: "write_capacity"},
		},
		Blocks: []hclext.BlockSchema{
			{
				Type: "global_secondary_index",
				Body: &hclext.BodySchema{},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, table := range dynamoDbTables.Blocks {
		// billing mode defaults to provisioned
		payPerRequest := false

		billingMode, billingModeDefined := table.Body.Attributes["billing_mode"]

		if billingModeDefined {
			if err := runner.EvaluateExpr(billingMode.Expr, func(val string) error {
				payPerRequest = val == "PAY_PER_REQUEST"
				return nil
			}, nil); err != nil {
				return err
			}
		}

		if !payPerRequest {
			issueRange := table.DefRange
			if billingModeDefined {
				issueRange = billingMode.Range
			}

			if err := runner.EmitIssue(
				r,
				"billing mode is not set to PAY_PER_REQUEST which may be expensive",
				issueRange,
			); err != nil {
				return err
			}
		}

		if !payPerRequest {
			// read/write capacity only work for PROVISIONED mode, so we only need to check them if we are not using PAY_PER_REQUEST

			readCapacity, readCapDefined := table.Body.Attributes["read_capacity"]
			if readCapDefined {
				if err := runner.EvaluateExpr(readCapacity.Expr, func(val int) error {
					if val > 1 {
						return runner.EmitIssue(r, "high read capacity might lead to higher cost", readCapacity.Range)
					}
					return nil
				}, nil); err != nil {
					return err
				}
			}

			writeCapacity, writeCapDefined := table.Body.Attributes["write_capacity"]
			if writeCapDefined {
				if err := runner.EvaluateExpr(writeCapacity.Expr, func(val int) error {
					if val > 1 {
						return runner.EmitIssue(r, "high write capacity might lead to higher cost", writeCapacity.Range)
					}
					return nil
				}, nil); err != nil {
					return err
				}
			}
		}

		for _, globalSecondaryIndex := range table.Body.Blocks {
			if err := runner.EmitIssue(r, "global secondary indices are expensive", globalSecondaryIndex.DefRange); err != nil {
				return err
			}
		}
	}

	return nil
}
