package usecase

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github/yogabagas/join-app/config"
	"github/yogabagas/join-app/domain/model"
	"github/yogabagas/join-app/domain/repository/cache"
	"github/yogabagas/join-app/domain/repository/sql"
	"github/yogabagas/join-app/domain/service"
	"github/yogabagas/join-app/shared/util"
	"log"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/golang-jwt/jwt"
)

type AuthzServiceImpl struct {
	repo  sql.RepositoryRegistry
	cache cache.Cache
}

type AuthzService interface {
	Login(ctx context.Context, req service.LoginReq) (resp service.LoginResp, err error)
	Logout(ctx context.Context, req service.LogoutReq) error
	HasAuthenticated(ctx context.Context, req service.HasAuthenticatedReq) (resp service.HasAuthenticatedResp, err error)
}

func NewAuthzService(repository sql.RepositoryRegistry, cache cache.Cache) AuthzService {
	return &AuthzServiceImpl{
		repo:  repository,
		cache: cache,
	}
}

func (as *AuthzServiceImpl) Login(ctx context.Context, req service.LoginReq) (resp service.LoginResp, err error) {

	usersRepo := as.repo.UsersRepository()
	jwkRepo := as.repo.JWKRepository()
	credentialsRepo := as.repo.UserCredentialsRepository()

	pwd, err := util.Hash(config.GlobalCfg.PasswordAlg, req.Password)
	if err != nil {
		return resp, err
	}

	user, err := usersRepo.ReadUserByEmail(ctx, &model.ReadUserByEmailReq{
		Email: req.Email,
	})
	if err != nil {
		return resp, err
	}

	if user.UserUID == "" {
		return resp, nil
	}

	crd, err := credentialsRepo.ReadCredentialsByUserUIDAndPassword(ctx, &model.ReadCredentialsByUserUIDAndPasswordReq{
		UserUID:  user.UserUID,
		Password: util.Base64(pwd),
	})
	if err != nil {
		return resp, err
	}

	if !crd.Valid {
		return resp, errors.New("wrong password")
	}

	key, err := jwkRepo.ReadUnexpiredKeyByID(ctx, &model.ReadUnexpiredKeyByIDReq{
		KeyID: user.RoleName,
	})
	if err != nil {
		return resp, err
	}

	obj := &jose.JSONWebKey{}

	privateKeyCache := fmt.Sprintf("jwk::private-key:%s", user.RoleName)

	_ = as.cache.GetObject(ctx, privateKeyCache, &obj)

	if key.ID != obj.KeyID || (key.ID == "" || obj.KeyID == "") {
		if err = as.generateJWKKey(ctx, user.RoleName); err != nil {
			return resp, err
		}

		err = as.cache.GetObject(ctx, privateKeyCache, &obj)
		if err != nil {
			return resp, err
		}
	}

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.SignatureAlgorithm(obj.Algorithm), Key: obj}, nil)
	if err != nil {
		return resp, err
	}

	accessToken, err := as.generateAndSignAccessToken(ctx, &model.GenerateAccessTokenReq{
		KeyID:      user.RoleName,
		UserUID:    user.UserUID,
		RoleUID:    user.RoleUID,
		LastActive: user.LastActive.UTC().Unix(),
		ExpiredAt:  config.GlobalCfg.TokenExpiration,
		Signer:     signer,
	})
	if err != nil {
		return resp, err
	}

	refreshToken, err := as.generateAndSignRefreshToken(ctx, &model.GenerateRefreshTokenReq{
		KeyID:     user.RoleName,
		UserUID:   user.UserUID,
		ExpiredAt: (config.GlobalCfg.TokenExpiration + config.GlobalCfg.RefreshTokenExpiration),
		Signer:    signer,
	})
	if err != nil {
		return resp, err
	}

	cacheKeyLogin := fmt.Sprintf("auth::user-uid:%s", user.UserUID)
	tokenExp := time.Now().UTC().Add(time.Duration(config.GlobalCfg.TokenExpiration) * time.Second).Unix()

	err = as.cache.Set(ctx, cacheKeyLogin, true, int(tokenExp))
	if err != nil {
		log.Fatalln("error set cache auth", err)
	}

	return service.LoginResp{
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
	}, nil
}

func (as *AuthzServiceImpl) Logout(ctx context.Context, req service.LogoutReq) error {

	key := fmt.Sprintf("auth::user-uid:%s", req.UserUID)

	return as.cache.Delete(ctx, key)
}

func (as *AuthzServiceImpl) HasAuthenticated(ctx context.Context, req service.HasAuthenticatedReq) (resp service.HasAuthenticatedResp, err error) {

	cacheKeyLogin := fmt.Sprintf("auth::user-uid:%s", req.Sub)
	if !as.cache.Exist(ctx, cacheKeyLogin) {
		return resp, nil
	}

	return service.HasAuthenticatedResp{
		Valid: true,
	}, nil
}

func (as *AuthzServiceImpl) generateAndSignAccessToken(ctx context.Context, req *model.GenerateAccessTokenReq) (resp *model.GenerateAccessTokenResp, err error) {

	claims := make(jwt.MapClaims)
	claims["sub"] = req.UserUID
	claims["role_uid"] = req.RoleUID
	claims["iat"] = time.Now().UTC().Unix()
	claims["exp"] = time.Now().UTC().Add(time.Duration(req.ExpiredAt) * time.Second).Unix()
	claims["last_active"] = req.LastActive

	data, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}

	obj, err := req.Signer.Sign(data)
	if err != nil {
		return nil, err
	}

	tokenString, err := obj.CompactSerialize()
	if err != nil {
		return nil, err
	}

	return &model.GenerateAccessTokenResp{
		Token: tokenString,
	}, nil
}

func (as *AuthzServiceImpl) generateAndSignRefreshToken(ctx context.Context, req *model.GenerateRefreshTokenReq) (resp *model.GenerateRefreshTokenResp, err error) {

	claims := make(jwt.MapClaims)
	claims["sub"] = req.UserUID
	claims["iat"] = time.Now().UTC().Unix()
	claims["exp"] = time.Now().UTC().Add(time.Duration(req.ExpiredAt) * time.Second).Unix()

	data, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}

	obj, err := req.Signer.Sign(data)
	if err != nil {
		return nil, err
	}

	tokenString, err := obj.CompactSerialize()
	if err != nil {
		return nil, err
	}

	return &model.GenerateRefreshTokenResp{
		Token: tokenString,
	}, nil
}

func (as *AuthzServiceImpl) generateJWKKey(ctx context.Context, keyID string) error {

	jwkRepo := as.repo.JWKRepository()

	key, err := rsa.GenerateKey(rand.Reader, config.GlobalCfg.JWK.Size)
	if err != nil {
		return err
	}

	privateKeyID := keyID

	privateKey := &jose.JSONWebKey{
		Key:       key,
		KeyID:     privateKeyID,
		Algorithm: config.GlobalCfg.JWK.Algorithm,
		Use:       config.GlobalCfg.JWK.Use,
	}

	thumb, err := privateKey.Thumbprint(crypto.SHA256)
	if err != nil {
		return err
	}

	publicKeyID := base64.URLEncoding.EncodeToString(thumb)
	privateKey.KeyID = publicKeyID

	publicKey := jose.JSONWebKey{
		Key:       key.Public(),
		KeyID:     publicKeyID,
		Algorithm: config.GlobalCfg.JWK.Algorithm,
		Use:       config.GlobalCfg.JWK.Use,
	}

	b, err := json.Marshal(publicKey)
	if err != nil {
		return err
	}

	expired := time.Now().Add(time.Duration(config.GlobalCfg.JWK.Expired) * time.Hour).UTC()

	reqJWT := &model.JWK{
		ID:         privateKeyID,
		Key:        string(b),
		PrivateKey: privateKey,
		ExpiredAt:  expired,
	}

	err = jwkRepo.UpsertJWK(ctx, reqJWT)
	if err != nil {
		return err
	}

	privKeyCache := fmt.Sprintf("jwk::private-key:%s", privateKeyID)

	err = as.cache.Set(ctx, privKeyCache, privateKey, int(expired.Unix()))
	if err != nil {
		return err
	}

	return nil
}
