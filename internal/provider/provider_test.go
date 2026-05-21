// Copyright (c) 2026 James Pickering
// SPDX-License-Identifier: MIT

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories is used to instantiate the provider during
// acceptance testing. Set REDIS_ADDR (default: localhost:6379) to point at a
// running Redis instance before executing acceptance tests.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"redis": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// Acceptance tests require a running Redis instance.
	// Configure the address via the REDIS_ADDR environment variable (default: localhost:6379).
}
