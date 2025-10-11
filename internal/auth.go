package internal

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type JWTManager struct{Secret []byte;Duration time.Duration}


func NewJWT(secret string,d time.Duration)*JWTManager{return &JWTManager{Secret:[]byte(secret),Duration:d}}
func (j *JWTManager) Generate(u *User)(string,error){
    claims:=&Claims{UserID:u.ID,Role:string(u.Role),RegisteredClaims:jwt.RegisteredClaims{ExpiresAt:jwt.NewNumericDate(time.Now().Add(j.Duration)),IssuedAt:jwt.NewNumericDate(time.Now())}}
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(j.Secret)
}
func HashPassword(p string)(string,error){b,_:=bcrypt.GenerateFromPassword([]byte(p),bcrypt.DefaultCost);return string(b),nil}
func ComparePassword(hash,p string)error{return bcrypt.CompareHashAndPassword([]byte(hash),[]byte(p))}
