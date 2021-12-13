package ginjwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
)

// GinJwtToken is a method to generate new token with expiry
// with dynamic payload passed in arguments
func GinJwtToken(setClaims map[string]interface{}) (map[string]interface{}, map[string]interface{}) {

	// intializing middleware
	mw := MwInitializer()

	// Create the token
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	// extracting claims in form of map
	claims := token.Claims.(jwt.MapClaims)

	// extracting expire time
	expire := mw.TimeFunc().Add(mw.Timeout)

	// setting claims
	for key, val := range setClaims {
		claims[key] = val
	}
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = mw.TimeFunc().Unix()

	mapd := map[string]interface{}{"token": "", "expire": ""}

	//passing map with extra information
	extraInfo := map[string]interface{}{
		"start":  mw.TimeFunc().Unix(),
		"end":    expire.Unix(),
		"expire": config.Conf.JWT.JWTExpireTime,
	}

	// signing token
	tokenString, err := token.SignedString(mw.Key)
	if err != nil {
		return mapd, extraInfo
	}

	// passing map with all information
	mapd = map[string]interface{}{
		"error":  false,
		"token":  tokenString,
		"expire": expire.Format(time.RFC3339),
	}

	return mapd, extraInfo
}
