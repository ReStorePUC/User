package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strings"
)

const contentType = "application/x-www-form-urlencoded"

type KongConfig struct {
	Host            string `yaml:"host"`
	ConsumerRequest string `yaml:"consumer_request"`
	JwtRequest      string `yaml:"jwt_request"`
}

type Kong struct {
	cfg *KongConfig
}

type Claims struct {
	Iss string `json:"iss"`
	jwt.StandardClaims
}

type jwtResponse struct {
	ConsumerId string `json:"consumer_id"`
	Key        string `json:"key"`
	Secret     string `json:"secret"`
}

func NewKong(cfg *KongConfig) *Kong {
	k := &Kong{
		cfg: cfg,
	}
	return k
}

// CreateCustomer Creates a customer on kong POST to `http://konghost:8001/consumer`
func (k *Kong) CreateCustomer(email string) error {
	log := zap.NewNop()

	urlRequest := k.cfg.Host + k.cfg.ConsumerRequest
	data := url.Values{}
	data.Set("username", email)

	r, err := http.NewRequest(http.MethodPost, urlRequest, strings.NewReader(data.Encode()))
	if err != nil {
		log.Error(
			"error creating kong request",
			zap.Error(err),
		)
		return err
	}

	r.Header.Set("Content-Type", contentType)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Error(
			"error making kong request",
			zap.Error(err),
		)
		return err
	}

	if resp.StatusCode != 201 {
		log.Error(
			"error creating consumer",
			zap.Any("status_code", resp.StatusCode),
		)
		return errors.New("error creating consumer")
	}

	return nil
}

// CreateCredentials Creates a JWT keys to a user POST to `http://konghost:8001/consumers/%s/jwt`
func (k *Kong) CreateCredentials(email string) (string, error) {
	log := zap.NewNop()

	urlReq := k.cfg.Host + fmt.Sprintf(k.cfg.JwtRequest, email)
	r, err := http.NewRequest(http.MethodPost, urlReq, nil)
	if err != nil {
		log.Error(
			"error creating credentials request",
			zap.Error(err),
		)
		return "", err
	}

	r.Header.Set("Content-Type", contentType)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Error(
			"error making credentials request",
			zap.Error(err),
		)
		return "", err
	}

	if resp.StatusCode != 201 {
		log.Error(
			"error creating consumer",
			zap.Any("status_code", resp.StatusCode),
		)
		return "", errors.New("error creating credentials")
	}

	defer resp.Body.Close()
	jwtR := jwtResponse{}

	err = json.NewDecoder(resp.Body).Decode(&jwtR)
	if err != nil {
		log.Error(
			"error decoding credentials request",
			zap.Error(err),
		)
		return "", err
	}

	jwtCode, err := k.createJWT(&jwtR)

	return jwtCode, nil
}

func (k *Kong) createJWT(response *jwtResponse) (string, error) {
	claims := &Claims{
		Iss: response.Key,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(response.Secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
