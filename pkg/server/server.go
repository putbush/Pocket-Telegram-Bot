package server

import (
	"context"
	"errors"
	"github.com/zhashkevych/go-pocket-sdk"
	"net/http"
	"pocketer_bot/pkg/storage"
	"strconv"
)

type AuthorizationServer struct {
	server          *http.Server
	pocketClient    *pocket.Client
	tokenRepository storage.TokenRepository
	redirectURL     string
}

func NewAuthorizationServer(pocketClient *pocket.Client, tokenRepository storage.TokenRepository, redirectURL string) *AuthorizationServer {
	return &AuthorizationServer{pocketClient: pocketClient, tokenRepository: tokenRepository, redirectURL: redirectURL}
}

func (s *AuthorizationServer) Start() error {
	s.server = &http.Server{
		Addr:    ":80",
		Handler: s,
	}

	return s.server.ListenAndServe()
}

func (s *AuthorizationServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chatIDParam := r.URL.Query().Get("chat_id")
	if chatIDParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatID, err := strconv.ParseInt(chatIDParam, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = s.createAccessToken(r.Context(), chatID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", s.redirectURL)
	w.WriteHeader(http.StatusMovedPermanently)
}

func (s *AuthorizationServer) createAccessToken(ctx context.Context, chatID int64) error {
	requestToken, err := s.tokenRepository.Get(chatID, storage.RequestTokens)
	if err != nil {
		return errors.New("failed to request token")
	}

	authResponse, err := s.pocketClient.Authorize(context.Background(), requestToken)
	if err != nil {
		return errors.New("failed to auth at Pocket")
	}

	err = s.tokenRepository.Save(chatID, authResponse.AccessToken, storage.AccessTokens)
	if err != nil {
		return errors.New("failed to save access-token to storage")
	}

	return nil
}
