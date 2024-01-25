// Provider

package provider

import (
	"context"
	"os"

    "github.com/networkdude17/passwordstate-api-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure passwordstateProvider satisfies various provider interfaces.
var _ provider.Provider = &passwordstateProvider{}



// passwordstateProviderModel describes the provider data model.
type passwordstateProviderModel struct {
	ApiUrl types.String `tfsdk:"api_url"`
	ApiKey types.String `tfsdk:"api_key"`
}

// passwordstateProvider defines the provider implementation.
type passwordstateProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Function - New Provider
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &passwordstateProvider{
			version: version,
		}
	}
}

// Function - Metadata
func (p *passwordstateProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "passwordstate"
	resp.Version = p.version
}

// Function - Schema
func (p *passwordstateProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				Description: "The PasswordState API URL; for example: https://passwordstate.domain.com/api/passwords",
				Optional:           true,
			},
			"api_key": schema.StringAttribute{
				Description: "The PasswordState API Key; for example: a0000aaa000aa000a0000a0a0a0a0a0a",
				Optional:           true,
				Sensitive: 			true,
			},
		},
	}
}

// Function - Configure
func (p *passwordstateProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Set Config model
	var config passwordstateProviderModel


    diags := req.Config.Get(ctx, &config)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // If practitioner provided a configuration value for any of the
    // attributes, it must be a known value.

	if config.ApiUrl.IsUnknown() {
        resp.Diagnostics.AddAttributeError(
            path.Root("api_url"),
            "Unknown PasswordState API URL",
            "The provider cannot create the PasswordState API client as there is an unknown configuration value for the PasswordState API URL. "+
                "Either target apply the source of the value first, set the value statically in the configuration, or use the PASSWORDSTATE_APIURL environment variable.",
        )
    }

	if config.ApiUrl.IsUnknown() {
        resp.Diagnostics.AddAttributeError(
            path.Root("api_key"),
            "Unknown PasswordState API Key",
            "The provider cannot create the PasswordState API client as there is an unknown configuration value for the PasswordState API Key. "+
                "Either target apply the source of the value first, set the value statically in the configuration, or use the PASSWORDSTATE_APIKEY environment variable.",
        )
    }

    if resp.Diagnostics.HasError() {
        return
    }

    // Default values to environment variables, but override
    // with Terraform configuration value if set.

    api_url := os.Getenv("PASSWORDSTATE_APIURL")
    api_key := os.Getenv("PASSWORDSTATE_APIKEY")

    if !config.ApiUrl.IsNull() {
        api_url = config.ApiUrl.ValueString()
    }

    if !config.ApiKey.IsNull() {
        api_key = config.ApiKey.ValueString()
    }

    // If any of the expected configurations are missing, return
    // errors with provider-specific guidance.

    if api_url == "" {
        resp.Diagnostics.AddAttributeError(
            path.Root("api_url"),
            "Missing PasswordState API URL",
            "The provider cannot create the PasswordState API client as there is a missing or empty value for the PasswordState API URL. "+
                "Set the host value in the configuration or use the PASSWORDSTATE_APIURL environment variable. "+
                "If either is already set, ensure the value is not empty.",
        )
    }

    if api_key == "" {
        resp.Diagnostics.AddAttributeError(
            path.Root("api_key"),
            "Missing PasswordState API Key",
            "The provider cannot create the PasswordState API client as there is a missing or empty value for the PasswordState API Key. "+
                "Set the host value in the configuration or use the PASSWORDSTATE_APIKEY environment variable. "+
                "If either is already set, ensure the value is not empty.",
        )
    }

    if resp.Diagnostics.HasError() {
        return
    }

    // Logging Structure
    ctx = tflog.SetField(ctx, "passwordstate_api_url", api_url)
    ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "passwordstate_api_key", api_key) // Key is masked and will not display in the running logs as plain text

    // Logging
    tflog.Info(ctx, "Starting: Configure PasswordState client")

    // Create a new PasswordState client using the configuration values
    client, err := passwordstateclient.NewClient(&api_url, &api_key)
    if err != nil {
        resp.Diagnostics.AddError(
            "Unable to Create PasswordState API Client",
            "An unexpected error occurred when creating the PasswordState API client. "+
                "If the error is not clear, please contact the provider developers.\n\n"+
                "PasswordState Client Error: "+err.Error(),
        )
        return
    }

    // Make the PasswordState client available during DataSource and Resource
    // type Configure methods.
	resp.DataSourceData = client
	//resp.ResourceData = client

    // Logging
    tflog.Info(ctx, "Completed: Configure PasswordState client", map[string]any{"success": true})
}

// Function - Resources
func (p *passwordstateProvider) Resources(ctx context.Context) []func() resource.Resource {
	//return []func() resource.Resource{
	//	NewPasswordResource,
	//}
	return nil
}

// Function - Data Sources
func (p *passwordstateProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPasswordDataSource,
	}
}
