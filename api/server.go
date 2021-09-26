package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Item struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Server struct {
	*mux.Router
	shoppintItems []Item
}

func NewServer() *Server {
	s := &Server{
		Router:        mux.NewRouter(),
		shoppintItems: []Item{},
	}
	s.setupRoutes()
	return s
}
func (s *Server) setupRoutes() {
	s.HandleFunc("/page", s.creatShoppingItem()).Methods("POST")
	s.HandleFunc("/page", s.helloWorld()).Methods("GET")
}
func (s *Server) helloWorld() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// id :=  mux.Vars(r)["id"]
		rw.Header().Set("Content-Type", "application/html")
		rw.Write([]byte("Hello World"))
	}
}
func (s *Server) creatShoppingItem() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var i Item
		if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		i.ID = uuid.New()
		s.shoppintItems = append(s.shoppintItems, i)
		rw.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(rw).Encode(i); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}