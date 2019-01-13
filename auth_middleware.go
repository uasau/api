package main

import (
	"net/http"
	"strings"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/keratin/authn-go/authn"
)

// JWTAuthMiddleware handles requests authorising them with a JWT
type JWTAuthMiddleware struct {
	logger     hclog.Logger
	next       http.Handler
	getSubject func(string) (string, error)
}

// NewJWTAuthMiddleware create a new JWT middleware
func NewJWTAuthMiddleware(authNServer, audience string, logger hclog.Logger, next http.Handler) (*JWTAuthMiddleware, error) {
	err := authn.Configure(authn.Config{
		Issuer:   authNServer,
		Audience: audience,
	})

	if err != nil {
		return nil, err
	}

	return &JWTAuthMiddleware{
		logger,
		next,
		func(jwt string) (string, error) {
			return authn.SubjectFrom(jwt)
		},
	}, nil
}

// Handle the request and if JWT is value call the next handler
func (j *JWTAuthMiddleware) Handle(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	j.logger.Info("Validate request from", "host", r.Host, "path", r.URL.Path, "query", r.URL.RawQuery)

	jwtHeader := r.Header.Get("Authorization")
	if jwtHeader == "" {
		j.logger.Error("No authorization header")
		http.Error(rw, "No auth header", http.StatusUnauthorized)
		return
	}

	authParts := strings.Split(jwtHeader, " ")
	if len(authParts) < 2 || authParts[0] != "jwt" {
		j.logger.Error("Invalid authorization header")
		http.Error(rw, "No auth token", http.StatusUnauthorized)
		return
	}

	//validate the token
	_, err := j.getSubject(authParts[1])
	if err != nil {
		j.logger.Error("Invalid JWT", "error", err, "jwt", authParts[1])
		http.Error(rw, "Invalid jwt", http.StatusUnauthorized)
		return
	}

	j.next.ServeHTTP(rw, r)
}