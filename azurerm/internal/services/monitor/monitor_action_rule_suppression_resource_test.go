package monitor

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/monitor/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccMonitorActionRuleSuppression_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_rule_suppression", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMonitorActionRuleSuppressionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorActionRuleSuppression_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorActionRuleSuppressionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccMonitorActionRuleSuppression_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_rule_suppression", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMonitorActionRuleSuppressionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorActionRuleSuppression_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorActionRuleSuppressionExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccMonitorActionRuleSuppression_requiresImport),
		},
	})
}

func TestAccMonitorActionRuleSuppression_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_rule_suppression", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMonitorActionRuleSuppressionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorActionRuleSuppression_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorActionRuleSuppressionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccMonitorActionRuleSuppression_updateSuppressionConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_rule_suppression", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMonitorActionRuleSuppressionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorActionRuleSuppression_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorActionRuleSuppressionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccMonitorActionRuleSuppression_dailyRecurrence(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorActionRuleSuppressionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccMonitorActionRuleSuppression_monthlyRecurrence(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorActionRuleSuppressionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccMonitorActionRuleSuppression_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorActionRuleSuppressionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccMonitorActionRuleSuppression_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorActionRuleSuppressionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckMonitorActionRuleSuppressionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.ActionRulesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can not found Monitor ActionRule: %s", resourceName)
		}
		id, err := parse.ActionRuleID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.GetByName(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: monitor action_rule %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on Monitor ActionRulesClient: %+v", err)
		}
		return nil
	}
}

func testCheckMonitorActionRuleSuppressionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.ActionRulesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_monitor_action_rule_suppression" {
			continue
		}
		id, err := parse.ActionRuleID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.GetByName(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Monitor ActionRulesClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccMonitorActionRuleSuppression_basic(data acceptance.TestData) string {
	template := testAccMonitorActionRuleSuppression_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_rule_suppression" "test" {
  name                = "acctest-moniter-%d"
  resource_group_name = azurerm_resource_group.test.name

  suppression {
    recurrence_type = "Always"
  }
}
`, template, data.RandomInteger)
}

func testAccMonitorActionRuleSuppression_requiresImport(data acceptance.TestData) string {
	template := testAccMonitorActionRuleSuppression_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_rule_suppression" "import" {
  name                = azurerm_monitor_action_rule_suppression.test.name
  resource_group_name = azurerm_monitor_action_rule_suppression.test.resource_group_name

  suppression {
    recurrence_type = azurerm_monitor_action_rule_suppression.test.suppression.0.recurrence_type
  }
}
`, template)
}

func testAccMonitorActionRuleSuppression_complete(data acceptance.TestData) string {
	template := testAccMonitorActionRuleSuppression_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_rule_suppression" "test" {
  name                = "acctest-moniter-%d"
  resource_group_name = azurerm_resource_group.test.name
  enabled             = false
  description         = "actionRule-test"

  scope {
    type         = "ResourceGroup"
    resource_ids = [azurerm_resource_group.test.id]
  }

  suppression {
    recurrence_type = "Weekly"

    schedule {
      start_date_utc = "2019-01-01T01:02:03Z"
      end_date_utc   = "2019-01-03T15:02:07Z"

      recurrence_weekly = ["Sunday", "Monday", "Friday", "Saturday"]
    }
  }

  condition {
    alert_context {
      operator = "Contains"
      values   = ["context1", "context2"]
    }

    alert_rule_id {
      operator = "Contains"
      values   = ["ruleId1", "ruleId2"]
    }

    description {
      operator = "DoesNotContain"
      values   = ["description1", "description2"]
    }

    monitor {
      operator = "NotEquals"
      values   = ["Fired"]
    }

    monitor_service {
      operator = "Equals"
      values   = ["Data Box Edge", "Data Box Gateway", "Resource Health"]
    }

    severity {
      operator = "Equals"
      values   = ["Sev0", "Sev1", "Sev2"]
    }

    target_resource_type {
      operator = "Equals"
      values   = ["Microsoft.Compute/VirtualMachines", "microsoft.batch/batchaccounts"]
    }
  }

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccMonitorActionRuleSuppression_dailyRecurrence(data acceptance.TestData) string {
	template := testAccMonitorActionRuleSuppression_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_rule_suppression" "test" {
  name                = "acctest-moniter-%d"
  resource_group_name = azurerm_resource_group.test.name

  scope {
    type         = "ResourceGroup"
    resource_ids = [azurerm_resource_group.test.id]
  }

  suppression {
    recurrence_type = "Daily"

    schedule {
      start_date_utc = "2019-01-01T01:02:03Z"
      end_date_utc   = "2019-01-03T15:02:07Z"
    }
  }
}
`, template, data.RandomInteger)
}

func testAccMonitorActionRuleSuppression_monthlyRecurrence(data acceptance.TestData) string {
	template := testAccMonitorActionRuleSuppression_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_rule_suppression" "test" {
  name                = "acctest-moniter-%d"
  resource_group_name = azurerm_resource_group.test.name

  scope {
    type         = "ResourceGroup"
    resource_ids = [azurerm_resource_group.test.id]
  }

  suppression {
    recurrence_type = "Monthly"

    schedule {
      start_date_utc     = "2019-01-01T01:02:03Z"
      end_date_utc       = "2019-01-03T15:02:07Z"
      recurrence_monthly = [1, 2, 15, 30, 31]
    }
  }
}
`, template, data.RandomInteger)
}

func testAccMonitorActionRuleSuppression_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-monitor-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
