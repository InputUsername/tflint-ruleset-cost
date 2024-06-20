package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// CostAwsObjectStorageLifecycleRule checks whether ...
type CostAwsObjectStorageLifecycleRule struct {
	tflint.DefaultRule
}

// NewCostAwsObjectStorageLifecycleRule returns a new rule
func NewCostAwsObjectStorageLifecycleRule() *CostAwsObjectStorageLifecycleRule {
	return &CostAwsObjectStorageLifecycleRule{}
}

// Name returns the rule name
func (r *CostAwsObjectStorageLifecycleRule) Name() string {
	return "cost_aws_object_storage_lifecycle_rule"
}

// Enabled returns whether the rule is enabled by default
func (r *CostAwsObjectStorageLifecycleRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *CostAwsObjectStorageLifecycleRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *CostAwsObjectStorageLifecycleRule) Link() string {
	return "https://search-rug.github.io/iac-cost-patterns/object-storage-lifecycle-rules/"
}

// Check checks whether ...
func (r *CostAwsObjectStorageLifecycleRule) Check(runner tflint.Runner) error {
	buckets, err := runner.GetResourceContent("aws_s3_bucket", &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: "lifecycle_rule",
				/*Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "enabled"},
					},
					Blocks: []hclext.BlockSchema{
						{Type: "transition"},
					},
				},*/
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	lifecycleConfigs, err := runner.GetResourceContent("aws_s3_bucket_lifecycle_configuration", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: "bucket"},
		},
		/*Blocks: []hclext.BlockSchema{
			{
				Type: "rule",
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "status", Required: true},
					},
				},
			},
		},*/
	}, nil)
	if err != nil {
		return err
	}

	for _, bucket := range buckets.Blocks {
		if len(bucket.Labels) <= 1 {
			continue
		}

		bucketName := bucket.Labels[1]

		hasLifecycleConfig := false
		for _, config := range lifecycleConfigs.Blocks {
			if err := runner.EvaluateExpr(config.Body.Attributes["bucket"].Expr, func(val string) error {
				if val == bucketName {
					hasLifecycleConfig = true
				}
				return nil
			}, nil); err != nil {
				return err
			}

			if hasLifecycleConfig {
				break
			}
		}

		hasLifecycleBlock := len(bucket.Body.Blocks) != 0

		if !hasLifecycleConfig && !hasLifecycleBlock {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("bucket `%s` does not have a lifecycle configuration or lifecycle rule", bucketName),
				bucket.DefRange,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
