package frontend

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/square/go-jose.v2"
)

// AppJWTClaims is the jwt token claims for this application
type AppJWTClaims struct {
	Role     int32              `json:"role,omitempty"`
	Auctions jwt.StandardClaims `json:"auctions,omitempty"`
	Login    jwt.StandardClaims `json:"login,omitempty"`
}

var sharedKey = []byte("%D*G-KaPdSgUkXp2s5v8y/B?E(H+MbQe")

// EncJwt encrypts the given claims and returns a jwt token
func EncJWT(claims AppJWTClaims) (JWTToken string, err error) {

	// Instantiate an encrypter using AES128-GCM with AES-GCM key wrap.
	encrypter, err := jose.NewEncrypter(jose.A256GCM, jose.Recipient{
		Algorithm: jose.DIRECT,
		Key:       sharedKey}, nil)

	if err != nil {
		return "", err
	}

	data, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	encObject, err := encrypter.Encrypt(data)
	if err != nil {
		return "", err
	}

	token, err := encObject.CompactSerialize()
	if err != nil {
		return "", err
	}

	return token, nil
}

// DecJWT decrypts the given jwt token and return the claims
func DecJWT(JWTToken string) (claims AppJWTClaims, err error) {
	encObject, err := jose.ParseEncrypted(JWTToken)
	if err != nil {
		return AppJWTClaims{}, err
	}

	claimsJson, err := encObject.Decrypt(sharedKey)
	if err != nil {
		return AppJWTClaims{}, err
	}

	claims = AppJWTClaims{}
	err = json.Unmarshal(claimsJson, &claims)
	if err != nil {
		return claims, err
	}

	return claims, nil
}

// OneTimeEnc creates a hash for the given password text.
func OneTimeEnc(password string) (string, error) {
	enc := sha256.New()
	_, err := enc.Write([]byte(password))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(enc.Sum(nil)), nil
}
