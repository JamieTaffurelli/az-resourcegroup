package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgs/gojq"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func showTerraformPlanAsJSON(t *testing.T) *gojq.JQ {
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
	if err != nil {
		t.Fatal(err)
	}

	parser, err := gojq.NewStringQuery(planJSON)
	if err != nil {
		t.Fatal(err)
	}

	return parser
}
func TestUT_CanNotDeleteLockIsDeployed(t *testing.T) {
	t.Parallel()

	terraformJSON := showTerraformPlanAsJSON(t)

	lockType, err := terraformJSON.Query("configuration.root_module.resources.[0].expressions.lock_level.constant_value")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, lockType, "CanNotDelete", "A CanNotDelete lock should be deployed to the resource group")
}

func TestUT_PreventDeletionIfContainsResourcesIsEnabled(t *testing.T) {
	t.Parallel()

	terraformJSON := showTerraformPlanAsJSON(t)

	preventDelete, err := terraformJSON.Query("configuration.provider_config.azurerm.expressions.features.[0].resource_group.[0].prevent_deletion_if_contains_resources.constant_value")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, preventDelete, true, "Non-empty Resource Groups should be prvented from being deleted")
}
