// Data Source - PasswordState - Password

package provider

import (
	"context"
	"fmt"

	"github.com/networkdude17/passwordstate-api-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var ( 
	_ datasource.DataSource 				= &PasswordDataSource{}
	_ datasource.DataSourceWithConfigure 	= &PasswordDataSource{}
)

// Function - New Data Source
func NewPasswordDataSource() datasource.DataSource {
	return &PasswordDataSource{}
}

// PasswordDataSource defines the data source implementation.
type PasswordDataSource struct {
	client *passwordstateclient.Client
}

// PasswordModel describes the data source data model.
type PasswordDataSourceModel struct {
	PasswordID 	types.Int64 	`tfsdk:"passwordid"`
	UserName   	types.String 	`tfsdk:"username"`
	Password    types.String 	`tfsdk:"password"`
}

// Function - Metadata
func (d *PasswordDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password"
}

// Function - Schema
func (d *PasswordDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Retrieve a password and username from PasswordState.",
		Attributes: map[string]schema.Attribute{
			"passwordid": schema.Int64Attribute{
				Description: "The Password ID (PID) of the object.",
				Required:    		 true,
			},
			"username": schema.StringAttribute{
				Description: "The username of the account.",
				Computed:            true,
			},
			"password": schema.StringAttribute{
				Description: "The password of the account.",
				Computed:            true,
				Sensitive: 			 true,
			},
		},
	}
}

// Function - Configure
func (d *PasswordDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	// Set Client
	client, ok := req.ProviderData.(*passwordstateclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *passwordstateclient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	// Set Data Client
	d.client = client

}

// Function - Read
func (d *PasswordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Retrieve values from data
	var state PasswordDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

    // Logging Structure
    ctx = tflog.SetField(ctx, "passwordstate_pid", state.PasswordID.ValueInt64())

    // Logging
    tflog.Info(ctx, "Starting: Data request to the PasswordState client")

	// API Call
	passwords, err := d.client.GetPassword(state.PasswordID.ValueInt64())
	if err != nil {
	  resp.Diagnostics.AddError(
		"Unable to Read PasswordState Password",
		err.Error(),
	  )
	  return
	}

	// Map response body to model
	for _, password := range passwords {
		PasswordState := PasswordDataSourceModel{
			PasswordID:		types.Int64Value(int64(password.PasswordID)),
			UserName:       types.StringValue(password.UserName),
			Password:      	types.StringValue(password.Password),
		}

		// Set returned data
		state = PasswordState
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

    // Logging
    tflog.Info(ctx, "Completed: Data request to the PasswordState client", map[string]any{"success": true})
}
