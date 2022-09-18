package main

import (
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	a "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	log "github.com/sirupsen/logrus"
)

func RegisterApp(app App) (App, error) {
	client, err := createClient()
	failedApp := &App{}
	if err != nil {
		fmt.Printf("Error creating credentials: %v\n", err)
		return *failedApp, errors.New("failed to create client")
	}

	app, err = createAzureADApp(client, app)
	if err != nil {
		return app, err
	}
	app, _ = createAzureADAppPass(client, app)
	app, _ = createAzureADAppSP(client, app)
	app, _ = createAzureADgroups(client, app)
	app, _ = createAzureKV(client, app)

	log.Debug(string(app.AppName))
	return app, nil
}

func createClient() (msgraphsdk.GraphServiceClient, error) {

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	failedauth, _ := a.NewAzureIdentityAuthenticationProvider(cred)
	failedadapter, _ := msgraphsdk.NewGraphRequestAdapter(failedauth)
	failed := msgraphsdk.NewGraphServiceClient(failedadapter)
	if err != nil {
		message := "error creating credentials"
		return *failed, errors.New(message)
	}

	auth, err := a.NewAzureIdentityAuthenticationProviderWithScopes(cred, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		message := "error authentication provider"
		return *failed, errors.New(message)
	}
	adapter, err := msgraphsdk.NewGraphRequestAdapter(auth)
	if err != nil {
		message := "error creating adapter"
		return *failed, errors.New(message)
	}
	client := msgraphsdk.NewGraphServiceClient(adapter)
	return *client, nil
}

func createAzureADApp(client msgraphsdk.GraphServiceClient, app App) (App, error) {
	newApp, err := CreateAzureADApp(client, app)
	if err != nil {
		log.Error(err)
		return app, errors.New("error creating the app")
	}
	return newApp, nil
}
func createAzureADAppPass(client msgraphsdk.GraphServiceClient, app App) (App, error) {
	//TODO
	return app, nil
}
func createAzureADAppSP(client msgraphsdk.GraphServiceClient, app App) (App, error) {
	//TODO
	return app, nil
}
func createAzureADgroups(client msgraphsdk.GraphServiceClient, app App) (App, error) {
	//TODO
	return app, nil
}
func createAzureKV(client msgraphsdk.GraphServiceClient, app App) (App, error) {
	//TODO
	return app, nil
}
