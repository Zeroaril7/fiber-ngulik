package models

type AccessTokenClaims struct {
	Aud      string `claim:"aud"`
	Username string `claim:"username"`
	Role     string `claim:"role"`
	Exp      int    `claim:"exp"`
}
