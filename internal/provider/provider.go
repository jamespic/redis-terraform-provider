// Copyright (c) 2026 James Pickering
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"crypto/tls"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/redis/go-redis/v9"
)

// Ensure RedisProvider satisfies the provider interface.
var _ provider.Provider = &RedisProvider{}

// RedisProvider defines the provider implementation.
type RedisProvider struct {
	version string
}

// RedisProviderModel describes the provider data model.
type RedisProviderModel struct {
	Addr                  types.String `tfsdk:"addr"`
	Password              types.String `tfsdk:"password"`
	Username              types.String `tfsdk:"username"`
	DB                    types.Int64  `tfsdk:"db"`
	TLS                   types.Bool   `tfsdk:"tls"`
	TLSInsecureSkipVerify types.Bool   `tfsdk:"tls_insecure_skip_verify"`
}

// New returns a function that creates a new RedisProvider instance.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &RedisProvider{version: version}
	}
}

func (p *RedisProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "redis"
	resp.Version = p.version
}

func (p *RedisProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provider for populating a Redis instance with data. Supports string key-value pairs, hashes, and sets.",
		Attributes: map[string]schema.Attribute{
			"addr": schema.StringAttribute{
				MarkdownDescription: "Redis server address in `host:port` form. Falls back to the `REDIS_ADDR` environment variable, then `localhost:6379`.",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Redis password. Falls back to the `REDIS_PASSWORD` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Redis username (Redis 6+ ACL). Falls back to the `REDIS_USERNAME` environment variable.",
				Optional:            true,
			},
			"db": schema.Int64Attribute{
				MarkdownDescription: "Redis database index (0–15). Defaults to `0`.",
				Optional:            true,
			},
			"tls": schema.BoolAttribute{
				MarkdownDescription: "Enable TLS for the connection. Falls back to the `REDIS_TLS` environment variable (`true`/`false`). Defaults to `false`.",
				Optional:            true,
			},
			"tls_insecure_skip_verify": schema.BoolAttribute{
				MarkdownDescription: "Disable TLS certificate verification. Only takes effect when `tls` is `true`. Falls back to the `REDIS_TLS_INSECURE_SKIP_VERIFY` environment variable (`true`/`false`). **Not recommended in production.** Defaults to `false`.",
				Optional:            true,
			},
		},
	}
}

func (p *RedisProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data RedisProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addr := "localhost:6379"
	if v := os.Getenv("REDIS_ADDR"); v != "" {
		addr = v
	}
	if !data.Addr.IsNull() && !data.Addr.IsUnknown() {
		addr = data.Addr.ValueString()
	}

	password := os.Getenv("REDIS_PASSWORD")
	if !data.Password.IsNull() && !data.Password.IsUnknown() {
		password = data.Password.ValueString()
	}

	username := os.Getenv("REDIS_USERNAME")
	if !data.Username.IsNull() && !data.Username.IsUnknown() {
		username = data.Username.ValueString()
	}

	db := 0
	if !data.DB.IsNull() && !data.DB.IsUnknown() {
		db = int(data.DB.ValueInt64())
	}

	useTLS := os.Getenv("REDIS_TLS") == "true"
	if !data.TLS.IsNull() && !data.TLS.IsUnknown() {
		useTLS = data.TLS.ValueBool()
	}

	skipVerify := os.Getenv("REDIS_TLS_INSECURE_SKIP_VERIFY") == "true"
	if !data.TLSInsecureSkipVerify.IsNull() && !data.TLSInsecureSkipVerify.IsUnknown() {
		skipVerify = data.TLSInsecureSkipVerify.ValueBool()
	}

	opts := &redis.Options{
		Addr:     addr,
		Password: password,
		Username: username,
		DB:       db,
	}

	if useTLS {
		opts.TLSConfig = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: skipVerify, //nolint:gosec // controlled by explicit provider option
		}
	}

	resp.DataSourceData = redis.NewClient(opts)
	resp.ResourceData = redis.NewClient(opts)
}

func (p *RedisProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewStringResource,
		NewHashResource,
		NewSetResource,
	}
}

func (p *RedisProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
