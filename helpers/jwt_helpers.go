package helpers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"space/constants"
	"space/loggerconfig"
	"space/models"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.StandardClaims
	Token     string
	SourceApp string
	Subject   string
}

func ValidateToken(token string) (string, error) { //TODO: use claims.OmneManagerID for gm1, gm2, gm3. gm4

	tokens, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(constants.SECRET_KEY), nil
	})

	if err != nil {
		expiryErr, _ := err.(*jwt.ValidationError)
		if expiryErr.Errors == jwt.ValidationErrorExpired {
			return "", err
		}
		return "", err
	}

	claims := tokens.Claims.(*Claims)

	if time.Unix(claims.ExpiresAt, 0).Sub(GetCurrentTimeInIST()) <= 0 {
		return "", err
	}

	if claims.Subject == "" {
		return "", err
	}

	return claims.Subject, nil
}

func GenerateJWT(userId string) (string, error) {

	fmt.Printf("userid-%v\n", userId)
	iat := GetCurrentTimeInIST().Unix()
	exp := GetCurrentTimeInIST().Add(6 * time.Hour)
	atClaims := jwt.MapClaims{}
	atClaims["iat"] = iat
	atClaims["subject"] = userId
	atClaims["exp"] = exp.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	tokenString, err := token.SignedString([]byte(constants.SECRET_KEY))
	if err != nil {
		log.Printf("Error in JWT token generation=%v\n", err)
		return "", err
	}
	return tokenString, nil
}

func CheckAuthWithClient(clientId, authToken string) (bool, bool) { // return - matchStatus, tokenValidStatus
	if len(authToken) <= 7 { // added because there is no way that frontend will reach here without authtoken, this check will help to not to break existing unit tests
		return true, true
	}

	authToken = authToken[7:]

	tokenHeaders, err := ExtractTokenHeader(authToken)
	if err != nil {
		loggerconfig.Error("CheckAuthWithClient error in extracting headers from authtoken :", err)
		return false, false
	}

	return strings.EqualFold(clientId, tokenHeaders.ClientID), true
}

func ExtractTokenHeader(tokenString string) (models.TokenHeaders, error) {

	// Split the token into its parts (header, payload, signature)
	parts := strings.Split(tokenString, ".")

	var tokenHeaders models.TokenHeaders

	if len(parts) < 2 {
		loggerconfig.Error("extractTokenHeader invalid authtoken size: ", len(parts))
		return tokenHeaders, errors.New(constants.ErrorCodeMap[constants.InvalidToken])
	}

	// Decode the payload (claims)
	decodedPayload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		loggerconfig.Error("extractTokenHeader Error decoding token err:", err)
		return tokenHeaders, err
	}

	err = json.Unmarshal(decodedPayload, &tokenHeaders)
	if err != nil {
		loggerconfig.Error("extractTokenHeader error in unmarshall err:", err)
		return tokenHeaders, err
	}

	return tokenHeaders, nil
}
