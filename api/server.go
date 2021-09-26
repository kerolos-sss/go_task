package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project/scraper"

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

type ScrapePageRequestBody struct {
	Url string `json:"url"`
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
	s.HandleFunc("/pageDetails", s.scrapePageDetailsAndCount()).Methods("POST")
	s.HandleFunc("/page", s.creatShoppingItem()).Methods("POST")
	s.HandleFunc("/page", s.helloWorld()).Methods("GET")
	fs := http.FileServer(http.Dir("./static"))
	s.Handle("/static/", http.StripPrefix("/static/", fs))
}
func (s *Server) helloWorld() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// id :=  mux.Vars(r)["id"]
		rw.Header().Set("Content-Type", "application/html")
		rw.Write([]byte("Hello World"))
	}
}
func (s *Server) scrapePageDetailsAndCount() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var body ScrapePageRequestBody
		fmt.Println("BODY #s####")
		fmt.Println(r.Body)

		if r.Header.Get("Content-Type") == "application/json" {
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(rw, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			body.Url = ""
			fmt.Println("###d##")
			r.ParseMultipartForm(2000)
			fmt.Println(r.Form)
			// b, _ := ioutil.ReadAll(r.Body)
			// fmt.Println(string(b))
			urlArr := r.Form["url"]

			if len(urlArr) > 0 {
				body.Url = urlArr[0]
			}
			fmt.Println(urlArr)

		}

		counts := scraper.GetPageDetailsAndCounts(body.Url)

		rw.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(rw).Encode(counts); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
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
