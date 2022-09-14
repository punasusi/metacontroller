package main

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var (
	subscriptionID    string
	location          string
	resourceGroupName string
)

// TODO add here
func RegisterApp(app AppReg) (id uuid.UUID) {
	location = os.Getenv("RESOURCE_LOCATION")
	resourceGroupName = os.Getenv("RESOURCE_GROUP_NAME")
	subscriptionID = os.Getenv("SUBSCRIPTION_ID")
	if len(subscriptionID) == 0 {
		log.Fatal("AZURE_SUBSCRIPTION_ID is not set.")
	}
	conn, err := connectionAzure()
	if err != nil {
		log.Fatalf("cannot connect to Azure:%+v", err)
	}
	ctx := context.Background()

	exits, err := checkExistenceResourceGroup(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("resources group exist:", exits)

	resourceGroup, err := getResourceGroup(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("get resources group:", *resourceGroup.ID)

	log.Debug(string(app.Appname))
	log.Debug(string(app.APIVersion))
	log.Debug(string(app.Env))
	return uuid.New()
}

func connectionAzure() (azcore.TokenCredential, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}
	return cred, nil
}

func checkExistenceResourceGroup(ctx context.Context, cred azcore.TokenCredential) (bool, error) {
	resourceGroupClient, err := armresources.NewResourceGroupsClient(subscriptionID, cred, nil)
	if err != nil {
		return false, err
	}

	boolResp, err := resourceGroupClient.CheckExistence(ctx, resourceGroupName, nil)
	if err != nil {
		return false, err
	}
	return boolResp.Success, nil
}
func getResourceGroup(ctx context.Context, cred azcore.TokenCredential) (*armresources.ResourceGroup, error) {
	resourceGroupClient, err := armresources.NewResourceGroupsClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	resourceGroupResp, err := resourceGroupClient.Get(ctx, resourceGroupName, nil)
	if err != nil {
		return nil, err
	}
	return &resourceGroupResp.ResourceGroup, nil
}
