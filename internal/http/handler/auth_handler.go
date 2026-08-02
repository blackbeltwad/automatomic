package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"automatomic/internal/auth"
	"automatomic/internal/model"
	"automatomic/internal/repository"
)

type AuthHandler struct {
	ghOAuth *auth.GitHubOAuth
	jwtMgr  *auth.JWTManager
	repo    repository.UserRepository
}


func NewAuthHandler (ghOAuths *auth.GitHubOAuth, jwtMgrs *auth.JWTManager, repos repository.UserRepository) *AuthHandler{
	return &AuthHandler{
		ghOAuth: ghOAuths,
		jwtMgr: jwtMgrs,
		repo: repos,
	}
}

func (h *AuthHandler) HandleGitHubLogin(w http.ResponseWriter, r *http.Request){
	state := "random-state-string"
	url := h.ghOAuth.GetAuthUrl(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) HandleGitHubCallback(w http.ResponseWriter, r *http.Request){
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, `{"error":"missing oauth authorization code"}`, http.StatusBadRequest)
		return
	}
	
	ghUser, err := h.ghOAuth.FetchUserInfo(r.Context(), code)
	if err != nil {
		http.Error(w, `{"error":"failed to authenticate with GitHub"}`, http.StatusInternalServerError)
		return
	}

	user := &model.User{
		GitHubID:  ghUser.ID,
		Username:  ghUser.Login,
		Email:     ghUser.Email,
		AvatarURL: ghUser.AvatarURL,
		Role:      "developer",
	}

	if err := h.repo.UpsertGitHubUser(r.Context(), user); err != nil {
		http.Error(w, `{"error":"failed to save user session"}`, http.StatusInternalServerError)
		return
	}

	scopes := []string{model.ScopePipelineRead, model.ScopePipelineWrite}
	tokenStr, err := h.jwtMgr.Generate(user, scopes)
	if err != nil {
		http.Error(w, `{"error":"token issuance failure"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": tokenStr,
		"token_type":   "Bearer",
		"expires_in":   24 * time.Hour,
		"user":         user,
	})
}