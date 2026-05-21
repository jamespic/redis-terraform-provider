// Copyright (c) 2026 James Pickering
// SPDX-License-Identifier: MIT

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccStringResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccStringResourceConfig("tf-acc-string", "hello"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"redis_string.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("tf-acc-string"),
					),
					statecheck.ExpectKnownValue(
						"redis_string.test",
						tfjsonpath.New("key"),
						knownvalue.StringExact("tf-acc-string"),
					),
					statecheck.ExpectKnownValue(
						"redis_string.test",
						tfjsonpath.New("value"),
						knownvalue.StringExact("hello"),
					),
				},
			},
			// ImportState
			{
				ResourceName:      "redis_string.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update value
			{
				Config: testAccStringResourceConfig("tf-acc-string", "world"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"redis_string.test",
						tfjsonpath.New("value"),
						knownvalue.StringExact("world"),
					),
				},
			},
		},
	})
}

func testAccStringResourceConfig(key, value string) string {
	return fmt.Sprintf(`
resource "redis_string" "test" {
  key   = %q
  value = %q
}
`, key, value)
}
