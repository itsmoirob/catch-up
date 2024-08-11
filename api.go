package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/medias/{id}", makeHTTPHandleFunc(s.handleMediaWithId))
	router.HandleFunc("/medias", makeHTTPHandleFunc(s.handleMedia))

	log.Println("JSON API server running on port: ", s.listenAddr)

	// server := http.Server{
	// 	Addr:    s.listenAddr,
	// 	Handler: Logging(router),
	// }

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleMediaWithId(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetMediaById(w, r)
	case "DELETE":
		return s.handleDeleteMedia(w, r)
	case "PUT":
		return s.handleUpdateMedia(w, r)
	default:
		return fmt.Errorf("method no allowed %s", r.Method)
	}
}

func (s *APIServer) handleMedia(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetMedias(w, r)
	case "POST":
		return s.handleCreateMedia(w, r)
	default:
		return fmt.Errorf("method no allowed %s", r.Method)
	}
}

func (s *APIServer) handleGetMediaById(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIdStrToInt(r)
	if err != nil {
		return err
	}

	media, err := s.store.GetMediaById(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, media)
}

func (s *APIServer) handleDeleteMedia(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIdStrToInt(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteMedia(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleUpdateMedia(w http.ResponseWriter, r *http.Request) error {
	updateReq := new(UpdateMediaRequest)
	if err := json.NewDecoder(r.Body).Decode(updateReq); err != nil {
		return err
	}
	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, updateReq)
}

func (s *APIServer) handleGetMedias(w http.ResponseWriter, r *http.Request) error {
	media, err := s.store.GetMedias()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, media)
}

func (s *APIServer) handleCreateMedia(w http.ResponseWriter, r *http.Request) error {
	mediaReq := new(CreateMediaRequest)
	if err := json.NewDecoder(r.Body).Decode(mediaReq); err != nil {
		return err
	}

	media := NewMedia(mediaReq.Title, mediaReq.Form)
	if err := s.store.CreateMedia(media); err != nil {
		return nil
	}

	return WriteJSON(w, http.StatusOK, media)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path, time.Since(start))
	})
}

func parseIdStrToInt(r *http.Request) (int, error) {
	idString := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idString)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idString)
	}

	return id, nil
}
