package mssql

import (
	"fmt"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
)
var (
	armAuthorizer      autorest.Authorizer
)

// OAuthGrantType specifies which grant type to use.
type OAuthGrantType int
const (
	// OAuthGrantTypeServicePrincipal for client credentials flow
	OAuthGrantTypeServicePrincipal OAuthGrantType = iota
)
// GrantType returns what grant type has been configured.
func grantType() OAuthGrantType {
	return OAuthGrantTypeServicePrincipal
}
// GetResourceManagementAuthorizer gets an OAuthTokenAuthorizer for Azure Resource Manager
func GetResourceManagementAuthorizer(config Config) (autorest.Authorizer, error) {
	if armAuthorizer != nil {
		return armAuthorizer, nil
	}

	var a autorest.Authorizer
	var err error

	a, err = getAuthorizerForResource(grantType(), config)

	if err == nil {
		// cache
		armAuthorizer = a
	} else {
		// clear cache
		armAuthorizer = nil
	}
	return armAuthorizer, err
}


func getAuthorizerForResource(grantType OAuthGrantType, config Config) (autorest.Authorizer, error) {
	var a autorest.Authorizer
	var err error

	switch grantType {

	case OAuthGrantTypeServicePrincipal:
		oauthConfig, err := adal.NewOAuthConfig(
			config.Environment().ActiveDirectoryEndpoint, config.TenantID)
		if err != nil {
			return nil, err
		}

		token, err := adal.NewServicePrincipalToken(
			*oauthConfig, config.ClientID, config.ClientSecret, config.Environment().ResourceManagerEndpoint)
		if err != nil {
			return nil, err
		}
		a = autorest.NewBearerAuthorizer(token)

	default:
		return a, fmt.Errorf("invalid grant type specified")
	}

	return a, err
}