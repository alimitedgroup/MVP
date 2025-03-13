package service

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"log"
	"sync"
	"time"

	common "github.com/alimitedgroup/MVP/srv/authenticator/authCommon"
	servicecmd "github.com/alimitedgroup/MVP/srv/authenticator/service/cmd"
	serviceobject "github.com/alimitedgroup/MVP/srv/authenticator/service/object"
	serviceportout "github.com/alimitedgroup/MVP/srv/authenticator/service/portOut"
	serviceresponse "github.com/alimitedgroup/MVP/srv/authenticator/service/response"
	serviceauthenticator "github.com/alimitedgroup/MVP/srv/authenticator/service/strategy"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
)

type AuthServiceParams struct {
	fx.In
	CheckKeyPairExistancePort serviceportout.ICheckKeyPairExistance
	GetPemPrivateKeyPort      serviceportout.IGetPemPrivateKeyPort
	GetPemPublicKeyPort       serviceportout.IGetPemPublicKeyPort
	StorePemKeyPairPort       serviceportout.IStorePemKeyPair
	PublishPort               serviceportout.IPublishPort
	AuthenticatorStrategy     serviceauthenticator.IAuthenticateUserStrategy
}

type AuthService struct {
	checkKeyPairExistancePort serviceportout.ICheckKeyPairExistance
	getPemPrivateKeyPort      serviceportout.IGetPemPrivateKeyPort
	getPemPublicKeyPort       serviceportout.IGetPemPublicKeyPort
	storePemKeyPairPort       serviceportout.IStorePemKeyPair
	publishPort               serviceportout.IPublishPort
	authenticatorStrategy     serviceauthenticator.IAuthenticateUserStrategy
	mutex                     sync.Mutex
}

func NewAuthService(p AuthServiceParams) *AuthService {
	return &AuthService{checkKeyPairExistancePort: p.CheckKeyPairExistancePort, getPemPrivateKeyPort: p.GetPemPrivateKeyPort, getPemPublicKeyPort: p.GetPemPublicKeyPort, storePemKeyPairPort: p.StorePemKeyPairPort, publishPort: p.PublishPort, authenticatorStrategy: p.AuthenticatorStrategy}
}

func (as *AuthService) generatePemKey() (*[]byte, *[]byte, error) {
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	puk := pk.Public()
	privKeyBytes, err := x509.MarshalECPrivateKey(pk)
	if err != nil {
		return nil, nil, err
	}
	privateKeyFile := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privKeyBytes})
	pubkeybytes, err := x509.MarshalPKIXPublicKey(puk)
	if err != nil {
		return nil, nil, err
	}
	publicKeyFile := pem.EncodeToMemory(&pem.Block{Type: "EC PUBLIC KEY", Bytes: pubkeybytes})
	return &privateKeyFile, &publicKeyFile, nil
}

func (as *AuthService) StorePemKeyPair(cmd *servicecmd.StorePemKeyPairCmd) error {
	storeRespose := as.storePemKeyPairPort.StorePemKeyPair(cmd)
	if storeRespose.GetError() != nil {
		return storeRespose.GetError()
	}
	return nil
}

func (as *AuthService) getPrivateKeyFromPem(prk *[]byte) (*ecdsa.PrivateKey, error) {
	decodedKey, _ := pem.Decode(*prk)
	if decodedKey == nil {
		return nil, common.ErrNoPrivateKey
	}
	prkDecoded, errprk := x509.ParseECPrivateKey(decodedKey.Bytes)
	if errprk != nil {
		return nil, common.ErrNoPrivateKey
	}
	return prkDecoded, nil
}

func (as *AuthService) generateToken(username string, role string) (string, error) {
	as.mutex.Lock()
	defer as.mutex.Unlock()
	response := as.checkKeyPairExistancePort.CheckKeyPairExistance(servicecmd.NewCheckPemKeyPairExistence())
	if response.GetError() != nil {
		prk, puk, err := as.generatePemKey()
		if err != nil {
			return "", err
		}
		err = as.StorePemKeyPair(servicecmd.NewStorePemKeyPairCmd(prk, puk))
		if err != nil {
			return "", err
		}
		issuer := as.getPemPublicKeyPort.GetPemPublicKey(servicecmd.NewGetPemPublicKeyCmd())
		err = as.publishPort.PublishKey(servicecmd.NewPublishPublicKeyCmd(puk, issuer.GetIssuer())).GetError()
		if err != nil {
			log.Fatal("Cannot publish key, turning off service")
		}
	}
	storePrk := as.getPemPrivateKeyPort.GetPemPrivateKey(servicecmd.NewGetPemPrivateKeyCmd())
	if storePrk.GetError() != nil {
		return "", nil
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sub":  username,
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iss":  storePrk.GetIssuer(),
	})
	signKey, err := as.getPrivateKeyFromPem(storePrk.GetPemPrivateKey())
	if err != nil {
		return "", err
	}
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", common.ErrNoToken
	}
	return tokenString, nil
}

func (as *AuthService) GetToken(cmd *servicecmd.GetTokenCmd) *serviceresponse.GetTokenResponse {
	role, err := as.authenticatorStrategy.Authenticate(*serviceobject.NewUserData(cmd.GetUsername()))
	if err != nil {
		return serviceresponse.NewGetTokenResponse("", common.ErrUserNotLegit)
	}
	token, err := as.generateToken(cmd.GetUsername(), role)
	if err != nil {
		return serviceresponse.NewGetTokenResponse("", common.ErrNoToken)
	}
	return serviceresponse.NewGetTokenResponse(token, nil)
}
