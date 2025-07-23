package dto

type CreateApplicationInput struct {
	Name        string `json:"name" form:"name"`
	Icon        string `json:"icon" form:"icon"`
	RedirectURI string `json:"redirect_uri" form:"redirect_uri"`
	LogoutURL   string `json:"logout_url" form:"logout_url"`
}
