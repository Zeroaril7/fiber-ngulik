package jwtrsa

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

type GenerateInputJWT struct {
	PrivateKey string
	Claims     map[string]interface{}
	TimeExpire time.Duration
}

func updatePrivateKey(signKeyStr string) error {
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(signKeyStr))
	if err != nil {
		return errors.New("error in pkg jwtrsa, when updatePrivateKey at jwt.ParseRSAPRivateKeyFromPEM")
	}

	privateKey = signKey

	return nil
}

func updatePublicKey(signKeyStr string) error {
	signKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(signKeyStr))
	if err != nil {
		return errors.New("error in pkg jwtrsa, when updatePublicKey at jwt.ParseRSAPublicKeyFromPEM")
	}

	publicKey = signKey

	return nil
}

func GetPrivateKey(path string) (*rsa.PrivateKey, error) {
	err := updatePrivateKey(path)

	if err != nil {
		return nil, errors.New("error in pkg jwtrsa, when updatePrivateKey at GetPrivateKey")
	}

	return privateKey, nil
}

func GetPublicKey(path string) (*rsa.PublicKey, error) {
	err := updatePublicKey(path)

	if err != nil {
		return nil, errors.New("error in pkg jwtrsa, when updatePublicKey at GetPublicKey")
	}

	return publicKey, nil
}

func GenerateJWT(input GenerateInputJWT) (tokenStr string, expiresIn time.Time, err error) {
	err = updatePrivateKey(input.PrivateKey)
	if err != nil {
		return tokenStr, expiresIn, errors.New("error in pkg jwtrsa, when updatePrivateKey at GenerateJWT")
	}

	token := jwt.New(jwt.GetSigningMethod("RS256"))

	tokenClaims := token.Claims.(jwt.MapClaims)

	for k, v := range input.Claims {
		tokenClaims[k] = v
	}

	exp := time.Now().Add(input.TimeExpire)
	tokenClaims["exp"] = exp.Unix()

	tkn, err := token.SignedString(privateKey)
	return tkn, exp, err
}
