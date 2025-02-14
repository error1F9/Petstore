package service

type LoginOut struct {
	AccessToken string `json:"access_token"`
	Err         error
}

type LoginIn struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
