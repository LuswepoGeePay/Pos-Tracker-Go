package models

type PocketBaseCredentials struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

type PocketBaseAuthResponse struct {
	Token  string                 `json:"token"`
	Record map[string]interface{} `json:"record"`
}
