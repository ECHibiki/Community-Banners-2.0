package bannerjwt

import (
  "fmt"
  "time"
  "crypto/rsa"
  "crypto/rand"
  "crypto/x509"
  "encoding/pem"
  "github.com/golang-jwt/jwt"
)

type Claims struct {
	*jwt.StandardClaims
	Name string
	IsDonor bool
	IsMod bool
}

var private_key *rsa.PrivateKey
var public_key *rsa.PublicKey

func init(){
  fmt.Print("")
  sign_key, err := rsa.GenerateKey(rand.Reader , 2048)
  bytes := x509.MarshalPKCS1PrivateKey(sign_key)
  private_pem := pem.EncodeToMemory(&pem.Block{
      Type:  "RSA PRIVATE KEY",
      Bytes: bytes,
  })
  if err != nil{
    panic(err)
  }
  private_key, err = jwt.ParseRSAPrivateKeyFromPEM(private_pem)
  if err != nil{
    panic(err)
  }

  bytes, err = x509.MarshalPKIXPublicKey(&sign_key.PublicKey)
  if err != nil{
    panic(err)
  }
  public_pem := pem.EncodeToMemory(&pem.Block{
      Type:  "RSA PUBLIC KEY",
      Bytes: bytes,
  })
  public_key, err = jwt.ParseRSAPublicKeyFromPEM(public_pem)
  if err != nil{
    panic(err)
  }
}

func IsAuth(token_string string) (string, bool, bool, error){
  token, err := jwt.ParseWithClaims(token_string, &Claims{},
  func(token *jwt.Token) (interface{}, error) {
      return public_key, nil
   },
  )
  if err != nil{
    return "",false, false, err
  }
  claims := token.Claims.(*Claims)
  return claims.Name , claims.IsDonor, claims.IsMod, nil
}

func CreateToken(name string, is_donor bool, is_mod bool) (string, error){
  	token := jwt.New(jwt.GetSigningMethod("RS256"))
  	token.Claims = &Claims{
  		&jwt.StandardClaims{
  			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
  		},
  		name,
      is_donor,
      is_mod,
  	}
  	return token.SignedString(private_key)
  }