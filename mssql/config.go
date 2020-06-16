package mssql

import (
	"fmt"
	"github.com/Azure/go-autorest/autorest/azure"
)

type Config struct {
	ClientID               string
	ClientSecret           string
	TenantID               string
	SubscriptionID         string
	ActiveDirectoryEndpoint string
	environment            *azure.Environment
}

func (c Config) UserAgent() string {
	return "msql-azurerm-sqlplugin"
}

func (c Config) Environment() *azure.Environment {
	if c.environment != nil {
		return c.environment
	}
	cloudName := "AzurePublicCloud"
	env, err := azure.EnvironmentFromName(cloudName)
	if err != nil {
		// TODO: move to initialization of var
		panic(fmt.Sprintf(
			"invalid cloud name '%s' specified, cannot continue\n", cloudName))
	}
	c.environment = &env
	return c.environment
}