package main

import (
	"github.com/inputusername/tflint-ruleset-cost/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "cost",
			Version: "0.1.0",
			Rules: []tflint.Rule{
				rules.NewCostAwsBudgetRule(),
				rules.NewCostGoogleBudgetRule(),
				rules.NewCostAzureBudgetRule(),
			},
		},
	})
}
