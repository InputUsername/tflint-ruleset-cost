package rules

import (
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// CostAwsOldGenerationRule checks whether ...
type CostAwsOldGenerationRule struct {
	tflint.DefaultRule
}

// NewCostAwsOldGenerationRule returns a new rule
func NewCostAwsOldGenerationRule() *CostAwsOldGenerationRule {
	return &CostAwsOldGenerationRule{}
}

// Name returns the rule name
func (r *CostAwsOldGenerationRule) Name() string {
	return "cost_aws_old_generation"
}

// Enabled returns whether the rule is enabled by default
func (r *CostAwsOldGenerationRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *CostAwsOldGenerationRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *CostAwsOldGenerationRule) Link() string {
	return "https://search-rug.github.io/iac-cost-patterns/old-generation/"
}

func isOldInstanceType(instanceType string) (bool, error) {
	return regexp.MatchString("t2|m4", instanceType)
}

func isOldVolumeType(volumeType string) (bool, error) {
	return regexp.MatchString("gp2", volumeType)
}

// Check checks whether ...
func (r *CostAwsOldGenerationRule) Check(runner tflint.Runner) error {
	instances, err := runner.GetResourceContent("aws_instance", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: "instance_type"},
		},
		Blocks: []hclext.BlockSchema{
			{
				Type: "root_block_device",
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "volume_type"},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, instance := range instances.Blocks {
		if instanceType, exists := instance.Body.Attributes["instance_type"]; exists {
			if err := runner.EvaluateExpr(instanceType.Expr, func(val string) error {
				isOldInstance, err := isOldInstanceType(val)
				if err != nil {
					return err
				}
				if isOldInstance {
					return runner.EmitIssue(r, "older generation instances might have worse performance and cost more than newer generations", instanceType.Range)
				}
				return nil
			}, nil); err != nil {
				return err
			}
		}

		for _, rootBlockDevice := range instance.Body.Blocks {
			if volumeType, exists := rootBlockDevice.Body.Attributes["volume_type"]; exists {
				if err := runner.EvaluateExpr(volumeType.Expr, func(val string) error {
					isOldVolume, err := isOldVolumeType(val)
					if err != nil {
						return err
					}
					if isOldVolume {
						return runner.EmitIssue(r, "older generation volumes might have worse performance and cost more than newer generations", volumeType.Range)
					}
					return nil
				}, nil); err != nil {
					return err
				}
			}
		}
	}

	volumes, err := runner.GetResourceContent("aws_ebs_volume", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: "type"},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, volume := range volumes.Blocks {
		if volumeType, exists := volume.Body.Attributes["type"]; exists {
			if err := runner.EvaluateExpr(volumeType.Expr, func(val string) error {
				isOldVolume, err := isOldVolumeType(val)
				if err != nil {
					return err
				}
				if isOldVolume {
					return runner.EmitIssue(r, "older generation volumes might have worse performance and cost more than newer generations", volumeType.Range)
				}
				return nil
			}, nil); err != nil {
				return err
			}
		}
	}

	return nil
}
