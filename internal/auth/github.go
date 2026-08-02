
package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

type GitHubUser struct{
	ID int64 `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
	AvatarURL string `json:"avatar_url"`

}

type GitHubOAuth struct{
	config *oauth2.Config
}

func NewGitHubOAuth(clientID string, clientSecret string, redirection_url string) *GitHubOAuth{
	return &GitHubOAuth{
		config: &oauth2.Config {
			ClientID: clientID,
			ClientSecret: clientSecret,
			Endpoint: githuboauth.Endpoint,
			RedirectURL: redirection_url,
			Scopes:       []string{"user:email", "read:user"},
		},
	}
}

func (g *GitHubOAuth) GetAuthUrl(state string) string {
	return g.config.AuthCodeURL(state)
}

func (g *GitHubOAuth) FetchUserInfo(ctx context.Context, code string) (*GitHubUser, error){
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("oauth exchange failed: %w", err)
	}
	client := g.config.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed fetching github user profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("github api returned non-200 status: %d", resp.StatusCode)
	}

	var ghUser GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&ghUser); err != nil {
		return nil, fmt.Errorf("failed decoding github profile: %w", err)
	}

	return &ghUser, nil
}