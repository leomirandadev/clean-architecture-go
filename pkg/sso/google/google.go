package google

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleSSO interface {
	GetUserData(ctx context.Context, state, code string) (*GoogleResp, error)
	GenerateURL() string
	IsRandomStrValid(state string) bool
}

func New(opts Options) GoogleSSO {
	return &googleImpl{
		randomStr: opts.RandomStr,
		oauthConfig: &oauth2.Config{
			RedirectURL:  opts.RedirectURL,
			ClientID:     opts.ClientID,
			ClientSecret: opts.ClientSecret,
			// Scopes: https://developers.google.com/identity/protocols/oauth2/scopes?hl=pt-br
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

type Options struct {
	RedirectURL  string `mapstructure:"REDIRECT_URL"`
	ClientID     string `mapstructure:"CLIENT_ID"`
	ClientSecret string `mapstructure:"CLIENT_SECRET"`
	RandomStr    string `mapstructure:"RANDOM_STR"`
}

type googleImpl struct {
	randomStr   string
	oauthConfig *oauth2.Config
}

func (g googleImpl) GetUserData(ctx context.Context, state, code string) (*GoogleResp, error) {
	if state != g.randomStr {
		slog.Debug("state missmatch")
		return nil, errors.New("state missmatch")
	}

	token, err := g.oauthConfig.Exchange(ctx, code)
	if err != nil {
		slog.Debug("exchange oauth fails", "err", err)
		return nil, err
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		slog.Debug("request api fails", "err", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Debug("status not ok", "status", resp.StatusCode)
		return nil, errors.New("status not ok")
	}

	var result GoogleResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Debug("parse result fails", "err", err)
		return nil, err
	}

	return &result, nil
}

// TODO in the future we can use one string per user
func (g googleImpl) GenerateURL() string {
	return g.oauthConfig.AuthCodeURL(g.randomStr)
}

// TODO in the future we can use one string per user
func (g googleImpl) IsRandomStrValid(state string) bool {
	return state == g.randomStr
}

type GoogleResp struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
}
