package util

import (
	"errors"
	"fmt"
	"github/yogabagas/join-app/domain/service"
	"net/http"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/golang-jwt/jwt"
)

func GetUserData(token string) (resp service.JWTClaims, err error) {

	jwtToken, err := SplitBearer(token)
	if err != nil {
		return resp, err
	}

	tokenParse, _, err := new(jwt.Parser).ParseUnverified(jwtToken, jwt.MapClaims{})
	if err != nil {
		return resp, err
	}

	if claims, ok := tokenParse.Claims.(jwt.MapClaims); ok {
		if subject, ok := claims["sub"]; ok {
			resp.Sub = subject.(string)
		}
		if role, ok := claims["role_uid"]; ok {
			resp.RoleUID = role.(string)
		}
	}
	return
}

func SplitBearer(token string) (string, error) {
	err := validation.Validate(token,
		validation.Required,
		validation.Match(regexp.MustCompile(`^(s|bearer|Bearer).([a-zA-Z0-9_=]+)\.([a-zA-Z0-9_=]+)\.([a-zA-Z0-9_\-\+\/=]*)`)))
	if err != nil {
		return "", errors.New("unknown jwt format")
	}

	var bearer string

	switch {
	case strings.Contains(token, "bearer"):
		bearer = "bearer "
	case strings.Contains(token, "Bearer"):
		bearer = "Bearer "
	}

	token = strings.TrimPrefix(token, bearer)

	if token == "" {
		return "", errors.New("token is empty")
	}

	return token, nil

}

type UserData struct {
	RoleUUID string `json:"role_uuid"`
	UserUUID string `json:"user_id"`
	Token    string `json:"token"`
}

func (u UserData) GetUserData(r *http.Request) *UserData {
	tokenClaims := r.Context().Value("user_data").(*jwt.Token)
	claims, ok := tokenClaims.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Invalid claims type")
	}
	claimData := make(map[string]interface{})
	for key, value := range claims {
		claimData[key] = value
	}

	u.RoleUUID = claimData["role_uuid"].(string)
	u.UserUUID = claimData["user_uuid"].(string)

	return &u
}
