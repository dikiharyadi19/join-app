package middlewares

import (
	"context"
	"errors"
	"github/yogabagas/join-app/adapter/controller"
	"github/yogabagas/join-app/config"
	"github/yogabagas/join-app/domain/service"
	"github/yogabagas/join-app/registry"
	"github/yogabagas/join-app/shared/constant"
	"github/yogabagas/join-app/transport/rest/handler/response"
	"net/http"
	"strings"
)

type MiddlewareImpl struct {
	appController controller.AppController
}

type Middleware interface {
	AuthenticationMiddleware(next http.Handler) http.Handler
	CORSHandle(next http.Handler) http.Handler
}

func NewMiddleware(r registry.Registry) Middleware {
	return &MiddlewareImpl{
		appController: r.NewAppController(),
	}
}

// AuthenticationMiddleware validates the JWT token.
func (mi *MiddlewareImpl) AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if !mi.isWhitelist(r.URL.Path, r.Method) {
			res := response.NewJSONResponse()

			token := r.Header.Get("Authorization")

			if token == "" {
				res.SetError(response.ErrUnauthorized).SetMessage(errors.New("authorization header is required").Error()).Send(w)
				return
			}

			newCtx, valid := mi.parseJwt(ctx, token)
			if !valid {
				res.SetError(response.ErrUnauthorized).SetMessage(errors.New("invalid authorized token, please re-authenticate").Error()).Send(w)
				return
			}

			ctx = newCtx
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (mi *MiddlewareImpl) isWhitelist(endpoint, method string) bool {
	mapAPI := make(map[string][]string)

	for _, v := range config.GlobalCfg.Whitelist.APIs {

		if strings.ContainsAny(v.Endpoint, "*") {
			if v.Endpoint[:strings.Index(v.Endpoint, "*")] == endpoint[:strings.Index(v.Endpoint, "*")] {
				v.Endpoint = endpoint
			}
		}

		mapAPI[v.Endpoint] = append(mapAPI[v.Endpoint], v.Methods...)
	}

	if methods, ok := mapAPI[endpoint]; ok {
		for _, m := range methods {
			if method == m {
				return true
			}
		}
	}

	return false
}

func (mi *MiddlewareImpl) parseJwt(ctx context.Context, token string) (context.Context, bool) {

	jwtSvc := mi.appController.JWKController
	authzSvc := mi.appController.AuthzController

	resp, err := jwtSvc.VerifyJWT(ctx, service.VerifyTokenReq{
		Token: token,
	})
	if err != nil {
		return ctx, false
	}

	claims := service.JWTClaims{
		Sub:        resp.UserUID,
		RoleUID:    resp.RoleUID,
		LastActive: resp.LastActive,
		ExpiredAt:  resp.ExpiredAt,
	}

	auth, _ := authzSvc.HasAuthenticated(ctx, service.HasAuthenticatedReq{
		Sub: claims.Sub,
	})

	if !auth.Valid {
		return ctx, false
	}

	return context.WithValue(ctx, constant.Claim, claims), true
}
