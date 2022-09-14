package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func loglevel() log.Level {
	switch getEnv("LOG_LEVEL", "info") {
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
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	log.SetLevel(loglevel())
	gin.DefaultWriter = log.New().Out
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/sync", syncHandler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8088"
	}
	r.Run(fmt.Sprintf(":%s", port))
}

func syncHandler(c *gin.Context) {
	response := &SyncResponse{}
	body, err := io.ReadAll(c.Request.Body)
	log.Debug(string(body))
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
	log.Info(request.Parent.APIVersion)

	identVersion := request.Parent.ObjectMeta.ResourceVersion
	log.Trace(request.Parent.Spec.APIDelegatedApps)
	for _, reg := range request.Children.ResourceVersion {
		response.Status.Replicas++
		response.Status.Succeeded++
		log.Trace(reg)
	}

	app := &AppReg{}
	// TODO add here
	app.APIVersion = request.Parent.APIVersion
	app.Appname = request.Parent.Spec.Appname
	app.Env = request.Parent.Spec.Env
	app.Final = request.Finalize
	app.Kind = request.Parent.Kind
	app.Version = identVersion

	log.Trace(app)
	appuuid := RegisterApp(*app)
	regspec := &RegistrationSpec{request.Parent.Spec.Appname, appuuid.String()}
	regname := &RegistrationMetadata{request.Parent.Spec.Appname}

	reg := &Registration{
		APIVersion: "security.punasusi.com/v1alpha1",
		Kind:       "Identreg",
		Metadata:   *regname,
		Spec:       *regspec,
	}
	response.Children = append(response.Children, *reg)
	log.Trace(appuuid.String())
	c.JSON(http.StatusOK, response)
}
