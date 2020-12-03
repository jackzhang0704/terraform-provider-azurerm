package monitor

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccMonitorScheduledQueryRulesDataSource_AlertingAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_scheduled_query_rules_alert", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMonitorScheduledQueryRules_AlertingActionConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
				),
			},
		},
	})
}

func TestAccMonitorScheduledQueryRulesDataSource_AlertingActionCrossResource(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_scheduled_query_rules_alert", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMonitorScheduledQueryRules_AlertingActionCrossResourceConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
				),
			},
		},
	})
}

func testAccDataSourceMonitorScheduledQueryRules_AlertingActionConfig(data acceptance.TestData) string {
	ts := time.Now().Format(time.RFC3339)
	template := testAccMonitorScheduledQueryRules_AlertingActionConfigBasic(data, ts)

	return fmt.Sprintf(`
%s

data "azurerm_monitor_scheduled_query_rules_alert" "test" {
  name                = basename(azurerm_monitor_scheduled_query_rules_alert.test.id)
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, template)
}

func testAccDataSourceMonitorScheduledQueryRules_AlertingActionCrossResourceConfig(data acceptance.TestData) string {
	template := testAccMonitorScheduledQueryRules_AlertingActionCrossResourceConfig(data)
	return fmt.Sprintf(`
%s

data "azurerm_monitor_scheduled_query_rules_alert" "test" {
  name                = basename(azurerm_monitor_scheduled_query_rules_alert.test.id)
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, template)
}
