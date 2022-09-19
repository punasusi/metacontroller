package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func CreateAzureADApp(client msgraphsdk.GraphServiceClient, app App) (App, error) {
	newApp := models.NewApplication()
	display_name := get_display_name(app)
	newApp.SetDisplayName(&display_name)
	newApp.SetIdentifierUris([]string{""})
	sign_in_audience := os.Getenv("SIGN_IN_AUDIENCE")
	newApp.SetSignInAudience(&sign_in_audience)
	is_fallback_public_client := true
	newApp.SetIsFallbackPublicClient(&is_fallback_public_client)
	appApi := compile_api(app)
	newApp.SetApi(&appApi)
	claim := get_claims(app)
	newApp.SetOptionalClaims(&claim)
	resources := get_resources(app)
	newApp.SetRequiredResourceAccess(resources)
	app = setup_uris(app)
	//TODO Add Roles when read from JSON

	registeredapp, err := client.Applications().Post(context.Background(), newApp, nil)
	if err != nil {
		fmt.Printf("Error creating the app: %v\n", err)
		fmt.Print(err)
		return app, errors.New("error creating the app")
	}
	app.AppName = *registeredapp.GetAppId()
	return app, nil
}

func compile_api(app App) models.ApiApplication {
	var api_access_token_version int32 = 2
	var admin_consent_description = "Defines the default application access' scope."
	var admin_consent_display_name = "Application access"
	var oauth_permission_enabled = true
	var scopetype = "Admin"
	var default_api_scope = os.Getenv("DEFAULT_API_SCOPE")
	default_scope_uuid := get_display_name(app)
	// api oath scope
	scope := models.NewPermissionScope()
	scope.SetAdminConsentDescription(&admin_consent_description)
	scope.SetAdminConsentDisplayName(&admin_consent_display_name)
	scope.SetIsEnabled(&oauth_permission_enabled)
	scope.SetType(&scopetype)
	scope.SetUserConsentDescription(nil)
	scope.SetUserConsentDisplayName(nil)
	scope.SetValue(&default_api_scope)
	scope.SetId(&default_scope_uuid)
	//APP API
	newApi := models.NewApiApplication()
	newApi.SetRequestedAccessTokenVersion(&api_access_token_version)
	newApi.SetOauth2PermissionScopes([]models.PermissionScopeable{scope})
	return *newApi
}

func get_display_name(app App) string {
	var suffix string
	switch app.Env {
	case "dv":
		suffix = "DEV"
	case "in":
		suffix = "INT"
	case "qa":
		suffix = "QUA"
	case "ua":
		suffix = "UAT"
	case "st":
		suffix = "STG"
	case "pr":
		suffix = "PRD"
	default:
		suffix = ""
	}

	return "SP " + app.ArdFullName + " " + suffix
}

func get_claims(app App) models.OptionalClaims {
	var tempfalsebool = false
	var claim_sid_name = "sid"
	var claim_groups_name = "sid"
	var groupclaim []string
	groupclaim = append(groupclaim, "netbios_domain_and_sam_account_name")

	// app optional claims
	claimsid := []models.OptionalClaimable{}
	claim1 := models.NewOptionalClaim()
	claim1.SetName(&claim_sid_name)
	claim1.SetSource(nil)
	claim1.SetEssential(&tempfalsebool)
	claim2 := models.NewOptionalClaim()
	claim2.SetName(&claim_groups_name)
	claim2.SetSource(nil)
	claim2.SetEssential(&tempfalsebool)
	claim2.SetAdditionalProperties(groupclaim[:])
	claimsid = append(claimsid, claim1, claim2)
	// claimsid
	claim := models.NewOptionalClaims()
	claim.SetAccessToken(claimsid)
	return *claim
}

func get_resources(app App) []models.RequiredResourceAccessable {

	var resources []models.RequiredResourceAccessable
	var resource_access []models.ResourceAccessable
	var resource_access_content models.ResourceAccessable
	var resource models.RequiredResourceAccessable
	var resource_app_id = os.Getenv("RESOURCE_APP_ID")
	var resource_access_type = os.Getenv("RESOURCE_ACCESS_TYPE")
	var resource_access_id = os.Getenv("RESOURCE_ACCESS_ID")
	// resources
	resource = models.NewRequiredResourceAccess()
	resource.SetResourceAppId(&resource_app_id)
	resource_access_content.SetType(&resource_access_type)
	resource_access_content.SetId(&resource_access_id)
	resource_access = append(resource_access, resource_access_content)
	resource.SetResourceAccess(resource_access)
	resources = append(resources, resource)
	for _, access := range app.RequiredResourceAccess {
		//TODO read this from input json!
		var delegated_resource_access []models.ResourceAccessable
		var delegated_resource_access_content models.ResourceAccessable
		var delegated_resource models.RequiredResourceAccessable
		var delegated_resource_app_id = ""
		var delegated_resource_access_type = ""
		var delegated_resource_access_id = ""
		// resources
		delegated_resource = models.NewRequiredResourceAccess()
		delegated_resource.SetResourceAppId(&delegated_resource_app_id)
		delegated_resource_access_content.SetType(&delegated_resource_access_type)
		delegated_resource_access_content.SetId(&delegated_resource_access_id)
		delegated_resource_access = append(delegated_resource_access, resource_access_content)
		delegated_resource.SetResourceAccess(delegated_resource_access)
		fmt.Println(access)
		resources = append(resources, delegated_resource)
	}
	return resources

}

func setup_uris(app App) App {
	switch app.Kind {
	case "API":
		//TODO add homepage, logout, redirect url/i
	case "APP":
		//TODO add redirect uri
	default:
	}
	return app
}
