package sys

import (
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/gowool/middleware/keyauth"
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/pkg/jwt"
)

var jwtParser = jwtv4.NewParser(jwtv4.WithValidMethods([]string{jwtv4.SigningMethodPS256.Alg()}))

func parseJWT(cfg *jwt.Config, key string) (UserClaims, error) {
	var claims UserClaims

	_, err := jwtParser.ParseWithClaims(key, &claims, func(token *jwtv4.Token) (any, error) {
		return cfg.GetPublicKey(), nil
	})

	return claims, err
}

func GetClaims(c wool.Ctx) UserClaims {
	value, _ := c.Get("claims").(UserClaims)
	return value
}

func JWTMiddleware(cfg *jwt.Config, checkOTP bool) wool.Middleware {
	authCfg := &keyauth.Config{Validator: func(c wool.Ctx, key string, source keyauth.ExtractorSource) (bool, error) {
		if source != keyauth.ExtractorSourceHeader {
			return false, jwtv4.ErrInvalidKey
		}

		claims, err := parseJWT(cfg, key)
		if err != nil {
			return false, err
		}

		c.Set("claims", claims)

		if checkOTP && !claims.OTP {
			return false, jwtv4.ErrInvalidKey
		}

		return true, nil
	}}

	return keyauth.Middleware(authCfg)
}
