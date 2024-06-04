package synq

import (
	"context"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"net/url"
)

type TokenSource interface {
	oauth2.TokenSource
	credentials.PerRPCCredentials
}

func LongLivedTokenSource(longLivedToken string, apiEndpoint *url.URL) (TokenSource, error) {
	initialToken, err := obtainToken(apiEndpoint, longLivedToken)
	if err != nil {
		return nil, err
	}

	return oauth.TokenSource{TokenSource: oauth2.ReuseTokenSource(initialToken, &tokenSource{apiEndpoint: apiEndpoint, longLivedToken: longLivedToken})}, nil
}

type tokenSource struct {
	longLivedToken string
	apiEndpoint    *url.URL
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	return obtainToken(t.apiEndpoint, t.longLivedToken)
}

func obtainToken(apiEndpoint *url.URL, longLivedToken string) (*oauth2.Token, error) {

	tokenURL, _ := url.Parse(apiEndpoint.String())
	tokenURL.Path = "/oauth2/token"
	conf := oauth2.Config{
		Endpoint: oauth2.Endpoint{
			TokenURL:  tokenURL.String(),
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
	return conf.PasswordCredentialsToken(context.Background(), "synq", longLivedToken)
}
