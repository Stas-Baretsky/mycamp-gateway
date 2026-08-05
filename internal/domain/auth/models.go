package auth

type RegisterParams struct {
	Email    string
	Password string
}

type LoginParams struct {
	Email    string
	Password string
	AppID    int32
}
