package main

import (
	"log"

	"github.com/darshanj/terraform-provider-azurerm-sqlplugin/mssql"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	// remove date and time stamp from log output as the plugin SDK already adds its own
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: mssql.Provider,
	})
}
