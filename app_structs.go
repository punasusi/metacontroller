package main

type Spec struct {
	APIDelegatedApps string `json:"api_delegated_apps"`
	APISuffix        string `json:"api_suffix"`
	AppKind          string `json:"app_kind"`
	Appname          string `json:"appname"`
	Ardfullname      string `json:"ardfullname"`
	Ardid            int    `json:"ardid"`
	Env              string `json:"env"`
	KeyVaultID       string `json:"key_vault_id"`
	SolutionID       string `json:"solution_id"`
	AppRoles         string `json:"roles"`
	AppTags          string `json:"tags"`
}

type App struct {
	// TODO add KV
	Ardid                  int32
	AppName                string
	ApiSuffix              string
	Kind                   string
	Env                    string
	ArdFullName            string
	AppAPI                 string
	OptionalClaims         string
	RequiredResourceAccess []string
	SinglePageApps         []string
	Web_apps               []string
	AppRoles               []string
	AppTags                []string
	KeyVaultID             string
	SolutionID             string
	Reg                    Registration
	Final                  bool
	Version                string
}
