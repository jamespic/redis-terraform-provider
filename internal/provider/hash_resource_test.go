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

func TestAccHashResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccHashResourceConfig("tf-acc-hash", map[string]string{
					"field1": "val1",
					"field2": "val2",
				}),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"redis_hash.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("tf-acc-hash"),
					),
					statecheck.ExpectKnownValue(
						"redis_hash.test",
						tfjsonpath.New("fields").AtMapKey("field1"),
						knownvalue.StringExact("val1"),
					),
					statecheck.ExpectKnownValue(
						"redis_hash.test",
						tfjsonpath.New("fields").AtMapKey("field2"),
						knownvalue.StringExact("val2"),
					),
				},
			},
			// ImportState
			{
				ResourceName:      "redis_hash.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update: change a value, add a field, remove a field
			{
				Config: testAccHashResourceConfig("tf-acc-hash", map[string]string{
					"field1": "updated",
					"field3": "val3",
				}),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"redis_hash.test",
						tfjsonpath.New("fields").AtMapKey("field1"),
						knownvalue.StringExact("updated"),
					),
					statecheck.ExpectKnownValue(
						"redis_hash.test",
						tfjsonpath.New("fields").AtMapKey("field3"),
						knownvalue.StringExact("val3"),
					),
				},
			},
		},
	})
}

func testAccHashResourceConfig(key string, fields map[string]string) string {
	pairs := ""
	for k, v := range fields {
		pairs += fmt.Sprintf("    %q = %q\n", k, v)
	}
	return fmt.Sprintf(`
resource "redis_hash" "test" {
  key    = %q
  fields = {
%s  }
}
`, key, pairs)
}
