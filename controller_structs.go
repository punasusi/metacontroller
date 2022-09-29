package main

import meta "k8s.io/apimachinery/pkg/apis/meta/v1"

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
	Finalize bool             `json:"finalizing"`
	Final    bool             `json:"finalized"`
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
	Identname string `json:"identname"`
	UUID      string `json:"uuid"`
}
