package tasq

import (
	"encoding/json"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
)

type QAuth struct {
	config      *oauth2.Config
	authCodeURL string
}

var Auth QAuth

func (auth *QAuth) init(cfg *QConfig) error {
	clientSecret, err := ioutil.ReadFile(cfg.Credentials)
	if err != nil {
		return err
	}

	auth.config, err = google.ConfigFromJSON(clientSecret, cfg.Scope)
	if err != nil {
		return err
	}

	auth.authCodeURL = auth.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return nil
}

func encodeToken(token *oauth2.Token) ([]byte, error) {
	return json.Marshal(token)
}

func decodeToken(tokenString []byte) (*oauth2.Token, error) {
	var token *oauth2.Token
	err := json.Unmarshal(tokenString, &token)

	return token, err
}

func (auth *QAuth) getTokenSource(ctx context.Context, tokenString []byte) (oauth2.TokenSource, error) {
	var tokenSource oauth2.TokenSource

	token, err := decodeToken(tokenString)
	if err != nil {
		return tokenSource, err
	}

	return auth.config.TokenSource(ctx, token), nil
}

func (auth *QAuth) GetToken(authCode string) ([]byte, error) {
	var tokenString []byte
	token, err := auth.config.Exchange(context.TODO(), authCode)
	if err != nil {
		return tokenString, err
	}

	tokenString, err = encodeToken(token)
	return tokenString, err
}

func (auth *QAuth) GetAuthCodeURL() string {
	return auth.authCodeURL
}
