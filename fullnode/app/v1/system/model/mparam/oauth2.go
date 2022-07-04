package mparam

import "github.com/cloudslit/cloudslit/fullnode/app/v1/system/model/mmysql"

type EditOauth2 struct {
	ID           int64         `json:"id" binding:"required"`
	Company      string        `json:"name" binding:"required,oneof= github facebook google"`
	ClientId     string        `json:"client_id" binding:"required"`
	ClientSecret string        `json:"client_secret" binding:"required"`
	RedirectUrl  string        `json:"redirect_url" binding:"required"`
	Scopes       mmysql.Scopes `json:"scopes" binding:"required"`
	AuthUrl      string        `json:"auth_url" binding:"required"`
	TokenUrl     string        `json:"token_url" binding:"required"`
}

type AddOauth2 struct {
	Company      string        `json:"name" binding:"required,oneof= github facebook google"`
	ClientId     string        `json:"client_id" binding:"required"`
	ClientSecret string        `json:"client_secret" binding:"required"`
	RedirectUrl  string        `json:"redirect_url" binding:"required"`
	Scopes       mmysql.Scopes `json:"scopes" binding:"required"`
	AuthUrl      string        `json:"auth_url" binding:"required"`
	TokenUrl     string        `json:"token_url" binding:"required"`
}
