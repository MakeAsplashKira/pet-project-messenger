package models

type Tokens struct {
	AccessToken  string `json:"acess_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
