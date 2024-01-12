package model

import (
	"time"

	"github.com/go-jose/go-jose/v3"
)

type JWK struct {
	ID         string
	Key        interface{}
	PrivateKey *jose.JSONWebKey
	ExpiredAt  time.Time
}

type ReadUnexpiredKeyByIDReq struct {
	KeyID string
}

type ReadUnexpiredKeyByIDResp struct {
	ID        string
	Key       interface{}
	ExpiredAt time.Time
}

type ReadUnexpiredKeyReq struct {
	Time time.Time
}

type ReadUnexpiredKeyResp struct {
	ID        string
	Key       interface{}
	ExpiredAt time.Time
}
