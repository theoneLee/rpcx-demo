package model

type AuthRequest struct {
	UserName string
	Password string
}

type AuthResponse struct {
	Token string
}

type SayRequest string
type SayResponse string
