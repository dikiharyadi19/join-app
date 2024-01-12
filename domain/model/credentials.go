package model

import "time"

type UserCredential struct {
	UserUID   string
	Username  string
	Password  string
	IsActive  bool
	CreatedAt time.Time
}

type ReadCredentialsByUserUIDAndPasswordReq struct {
	UserUID  string
	Password string
}

type ReadCredentialsByUserUIDAndPasswordResp struct {
	Valid bool
}
