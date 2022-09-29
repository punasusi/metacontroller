package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func loglevel() log.Level {
	switch getEnv("LOG_LEVEL", "trace") {
	case "debug":
		return log.DebugLevel
	case "trace":
		return log.TraceLevel
	default:
		return log.InfoLevel

	}
}

func main() {
	// load .env file
	// godotenv.Load(".env")

	// if err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }
	fmt.Println("Starting?")
	log.SetLevel(loglevel())
	log.Info("starting")
	gin.DefaultWriter = log.New().Out
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/sync", syncHandler)
	r.POST("/final", finalizerHandler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	r.Run(fmt.Sprintf(":%s", port))
}

// TODO add finalizer to delete
func syncHandler(c *gin.Context) {
	response := &SyncResponse{}
	body, err := io.ReadAll(c.Request.Body)
	log.Trace(string(body))
	if err != nil {
		log.Error("JSON body could not be retrieved")
		c.String(http.StatusBadRequest, "JSON body could not be retrieved")
		return
	}
	request := &SyncRequest{}
	err = json.Unmarshal(body, request)
	if err != nil {
		log.Error("JSON could not be unmarshalled")
		c.String(http.StatusBadRequest, "JSON could not be unmarshalled")
		return
	}

	identVersion := request.Parent.ObjectMeta.ResourceVersion
	for _, reg := range request.Children.ResourceVersion {
		response.Status.Replicas++
		response.Status.Succeeded++
		log.Trace(reg)
	}

	app := &App{}
	app.RequiredResourceAccess = append(app.RequiredResourceAccess, request.Parent.Spec.APIDelegatedApps)
	app.ApiSuffix = request.Parent.Spec.APISuffix
	app.Kind = request.Parent.Kind
	app.ArdFullName = request.Parent.Spec.Appname
	app.Ardid = int32(request.Parent.Spec.Ardid)
	app.Env = request.Parent.Spec.Env
	app.KeyVaultID = request.Parent.Spec.KeyVaultID
	app.SolutionID = request.Parent.Spec.SolutionID
	app.AppRoles = append(app.AppRoles, request.Parent.Spec.AppRoles)
	app.AppTags = append(app.AppTags, request.Parent.Spec.AppTags)
	app.Final = request.Finalize

	app.Version = identVersion
	appdump := spew.Sdump(app)
	log.Trace(appdump)
	appuuid, err := RegisterApp(*app)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	regspec := &RegistrationSpec{request.Parent.Spec.Appname, appuuid.ArdFullName}
	regname := &RegistrationMetadata{request.Parent.Spec.Appname}

	reg := &Registration{
		APIVersion: "security.punasusi.com/v1alpha1",
		Kind:       "Identreg",
		Metadata:   *regname,
		Spec:       *regspec,
	}
	response.Children = append(response.Children, *reg)
	log.Trace(appuuid.ArdFullName)
	c.JSON(http.StatusOK, response)
}

func finalizerHandler(c *gin.Context) {
	response := &SyncResponse{}
	response.Final = true
	c.JSON(http.StatusOK, response)
}
