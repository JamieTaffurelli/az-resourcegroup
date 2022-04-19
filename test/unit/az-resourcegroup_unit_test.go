package test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/elgs/gojq"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func getJsonMap(m map[string]interface{}, key string) map[string]interface{} {
	raw := m[key]
	sub, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}
	return sub
}

func TestUT_CanNotDeleteLockIsDeployed(t *testing.T) {
	t.Parallel()

	tfOptions := &terraform.Options{
		TerraformDir: "../../",
		Vars: map[string]interface{}{
			"resource_group_name": "testrg",
			"location":            "francecentral",
			"tags":                map[string]string{},
		},
		NoColor: true,
	}

	tfPlanOutput := "terraform.tfplan"
	terraform.Init(t, tfOptions)
	terraform.RunTerraformCommand(t, tfOptions, terraform.FormatArgs(tfOptions, "plan", "-out="+tfPlanOutput)...)

	tfOptionsEmpty := &terraform.Options{}
	planJSON, err := terraform.RunTerraformCommandAndGetStdoutE(t, tfOptions, terraform.FormatArgs(tfOptionsEmpty, "show", "-json", tfPlanOutput)...)

	parser, err := gojq.NewStringQuery(planJSON)
	if err != nil {
		t.Fatal(err)
	}

	lockType, err := parser.Query("configuration.root_module.resources.[0].expressions.lock_level.constant_value")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, lockType, "CanNotDelete", "A CanNotDelete lock should be deployed to the resource group")
}
