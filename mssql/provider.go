package mssql

import (
	sql "github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"time"
)

func Provider() terraform.ResourceProvider {
	dataSources := make(map[string]*schema.Resource)

	resources := map[string]*schema.Resource{
		"azurerm_mssql_user": resourceArmMsSqlUser(),
	}

	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"subscription_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_SUBSCRIPTION_ID", ""),
				Description: "The Subscription ID which should be used.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_ID", ""),
				Description: "The Client ID which should be used.",
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_TENANT_ID", ""),
				Description: "The Tenant ID which should be used.",
			},
		},
		DataSourcesMap: dataSources,
		ResourcesMap:   resources,
	}
	p.ConfigureFunc = providerConfigure(p)
	return p
}

func resourceArmMsSqlUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMsSqlUserCreateUpdate,
		Read:   resourceArmMsSqlUserRead,
		Update: resourceArmMsSqlUserCreateUpdate,
		Delete: resourceArmMsSqlUserDelete,
		//Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
		//	_, err := parse.MsSqlDatabaseID(id)
		//	return err
		//}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				// TODO: ValidateFunc: azure.ValidateMsSqlDatabaseName,
			},

			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"database_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				// TODO: ValidateFunc: azure.ValidateMsSqlDatabaseName,
			},
		},
	}
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		config := &Config{
			SubscriptionID:     d.Get("subscription_id").(string),
			ClientID:           d.Get("client_id").(string),
			ClientSecret:       d.Get("client_secret").(string),
			TenantID:           d.Get("tenant_id").(string),
		}
		client := getDbClient(d, *config)
		return client,nil
	}
}

func getDbClient(d *schema.ResourceData,config Config) sql.DatabasesClient {
	subscriptionID := d.Get("subscription_id").(string)
	dbClient := sql.NewDatabasesClient(subscriptionID)
	a, _ := GetResourceManagementAuthorizer(config)
	dbClient.Authorizer = a
	dbClient.AddToUserAgent(config.UserAgent())
	return dbClient
}

func resourceArmMsSqlUserCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceArmMsSqlUserDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceArmMsSqlUserRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}