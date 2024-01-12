package service

import "time"

type JWTClaims struct {
	Sub        string    `json:"sub"`
	RoleUID    string    `json:"role_uid"`
	LastActive time.Time `json:"last_active"`
	ExpiredAt  time.Time `json:"expired_at"`
}

type LoginReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type VerifyTokenReq struct {
	Token string
}

type VerifyTokenResp struct {
	Valid      bool      `json:"valid"`
	UserUID    string    `json:"user_uid"`
	RoleUID    string    `json:"role_uid"`
	LastActive time.Time `json:"last_active"`
	ExpiredAt  time.Time `json:"expired_at"`
}

type HasAuthenticatedReq struct {
	Sub string `json:"sub"`
}

type HasAuthenticatedResp struct {
	Valid bool `json:"valid"`
}
