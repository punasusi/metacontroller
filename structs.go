package main

import meta "k8s.io/apimachinery/pkg/apis/meta/v1"

type Spec struct {
	// TODO add here
	APIDelegatedApps string `json:"api_delegated_apps"`
	APISuffix        string `json:"api_suffix"`
	AppKind          string `json:"app_kind"`
	Appname          string `json:"appname"`
	Ardfullname      string `json:"ardfullname"`
	Ardid            int    `json:"ardid"`
	Env              string `json:"env"`
	KeyVaultID       string `json:"key_vault_id"`
	SolutionID       string `json:"solution_id"`
}
type Controller struct {
	APIVersion      string `json:"apiVersion"`
	Kind            string `json:"kind"`
	Spec            Spec   `json:"spec"`
	meta.ObjectMeta `json:"metadata"`
}

type SyncRequest struct {
	Parent   Controller          `json:"parent"`
	Children SyncRequestChildren `json:"children"`
	Finalize bool                `json:"finalizing"`
}
type SyncRequestChildren struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata"`
}
type ControllerStatus struct {
	Replicas  int `json:"replicas"`
	Succeeded int `json:"succeeded"`
}
type SyncResponse struct {
	Status   ControllerStatus `json:"status"`
	Children []Registration   `json:"children"`
}
type Registration struct {
	APIVersion string               `json:"apiVersion"`
	Kind       string               `json:"kind"`
	Metadata   RegistrationMetadata `json:"metadata"`
	Spec       RegistrationSpec     `json:"spec"`
}
type RegistrationMetadata struct {
	Name string `json:"name"`
}
type RegistrationSpec struct {
	// TODO add here
	Identname string `json:"identname"`
	UUID      string `json:"uuid"`
}
type AppReg struct {
	// TODO add here
	APIVersion string
	Kind       string
	Env        string
	Appname    string
	Version    string
	Reg        Registration
	Final      bool
}
