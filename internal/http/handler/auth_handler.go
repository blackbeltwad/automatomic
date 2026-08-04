package handler

import (
	
	"net/http"
	"log"

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
		log.Printf("[OAuth Error] Failed to exchange code for GitHub token: %v\n", err)
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
		log.Printf("[OAuth Error] Database Upsert Failed: %v\n", err)
		http.Error(w, `{"error":"failed to save user session"}`, http.StatusInternalServerError)
		return
	}

	scopes := []string{model.ScopePipelineRead, model.ScopePipelineWrite}
	tokenStr, err := h.jwtMgr.Generate(user, scopes)
	if err != nil {
		log.Printf("[OAuth Error] JWT Generation Failed: %v\n", err)
		http.Error(w, `{"error":"token issuance failure"}`, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
        Name:     "access_token",
        Value:    tokenStr,
        Path:     "/",
        HttpOnly: true,
        Secure:   false, // set to true in production with HTTPS
        SameSite: http.SameSiteLaxMode,
        MaxAge:   86400, // 24 hours
    })

    http.Redirect(w, r, "http://localhost:3000/app/dashboard", http.StatusTemporaryRedirect)
}