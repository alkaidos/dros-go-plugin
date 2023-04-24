package user_proxy_vo

type AppAuthenticationDTO struct {
	Token   string `json:"token"`
	AppCode string `json:"appCode"`
	AppName string `json:"appName"`
}
