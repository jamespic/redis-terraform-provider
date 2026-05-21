// Copyright (c) 2026 James Pickering
// SPDX-License-Identifier: MIT

package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccSetResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccSetResourceConfig("tf-acc-set", []string{"alpha", "beta", "gamma"}),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"redis_set.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("tf-acc-set"),
					),
					statecheck.ExpectKnownValue(
						"redis_set.test",
						tfjsonpath.New("members"),
						knownvalue.SetSizeExact(3),
					),
					statecheck.ExpectKnownValue(
						"redis_set.test",
						tfjsonpath.New("members"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("alpha"),
							knownvalue.StringExact("beta"),
							knownvalue.StringExact("gamma"),
						}),
					),
				},
			},
			// ImportState
			{
				ResourceName:      "redis_set.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update: remove one member, add another
			{
				Config: testAccSetResourceConfig("tf-acc-set", []string{"beta", "gamma", "delta"}),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"redis_set.test",
						tfjsonpath.New("members"),
						knownvalue.SetSizeExact(3),
					),
					statecheck.ExpectKnownValue(
						"redis_set.test",
						tfjsonpath.New("members"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("beta"),
							knownvalue.StringExact("gamma"),
							knownvalue.StringExact("delta"),
						}),
					),
				},
			},
		},
	})
}

func testAccSetResourceConfig(key string, members []string) string {
	quoted := make([]string, len(members))
	for i, m := range members {
		quoted[i] = fmt.Sprintf("%q", m)
	}
	return fmt.Sprintf(`
resource "redis_set" "test" {
  key     = %q
  members = [%s]
}
`, key, strings.Join(quoted, ", "))
}
