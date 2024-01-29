package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Service interface {
	Send(routingKey string, data []byte) error
}

type HTTPServer struct {
	Service Service
	token   string
}

func NewHttpServer(service Service, authToken string) *HTTPServer {
	if len(authToken) < 12 {
		log.Fatalln("Необходимо указать токен аутентификации длиной не менее 12 символов")
	}
	return &HTTPServer{service, authToken}
}

func (s *HTTPServer) Start(address string) {
	http.HandleFunc("/", s.handleRequest)
	fmt.Printf("Сервер запущен: %s...\n", address)

	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalln("Ошибка при запуске сервера: ", err)
	}
}

func (s *HTTPServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	// Получаем routingKey из URL
	routingKey := strings.TrimPrefix(r.URL.Path, "/")

	if r.Method != http.MethodPost {
		s.getJSONError(w, errors.New("метод не разрешен"), http.StatusMethodNotAllowed)
		return
	}

	authorizationHeader := r.Header.Get("Authorization")
	actualToken := strings.TrimPrefix(authorizationHeader, "Token ")
	if actualToken != s.token {
		s.getJSONError(w, errors.New("неверный токен"), http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.getJSONError(w, errors.New("ошибка чтения тела запроса"), http.StatusInternalServerError)
	}

	if err := s.Service.Send(routingKey, body); err != nil {
		s.getJSONError(w, err, http.StatusInternalServerError)
	} else {
		log.Printf("Publish. Routing Key: %s", routingKey)
		w.WriteHeader(http.StatusCreated)
	}
}

func (s *HTTPServer) getJSONError(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	jsonData, _ := json.Marshal(map[string]string{"detail": fmt.Sprintf("%s", err)})
	w.Write(jsonData)
}
