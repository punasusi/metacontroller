package main

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// TODO add here
func RegisterApp(app AppReg) (id uuid.UUID) {
	log.Debug(string(app.Appname))
	log.Debug(string(app.APIVersion))
	log.Debug(string(app.Env))
	return uuid.New()
}
