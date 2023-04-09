package router

import (
	"fmt"
	"net/http"

	_ "net/http/pprof"

	project "github.com/hacktues-9/tf-api/cmd/projects"
	votes "github.com/hacktues-9/tf-api/cmd/votes"
	database "github.com/hacktues-9/tf-api/pkg/database"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Router struct {
	router *mux.Router
	DB     *gorm.DB
}

func NewRouter(db *gorm.DB) *Router {
	r := mux.NewRouter().PathPrefix("/v1").Subrouter().StrictSlash(true)
	return &Router{r, db}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *Router) GetRouter() *mux.Router {
	return r.router
}

func (r *Router) GetDB() *gorm.DB {
	return r.DB
}

func (r *Router) Init() {
	router := r.GetRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	}).Methods("GET")
	router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
}

func (r *Router) Projects() {
	router := r.GetRouter()
	GetReq := router.PathPrefix("/get").Subrouter().StrictSlash(true)
	GetReq.HandleFunc("/projects", func(writer http.ResponseWriter, request *http.Request) {
		// call function GetProjects from projects package
		project.GetProjects(writer, request, r.GetDB())
	}).Methods("GET")
	GetReq.HandleFunc("/project/{id}", func(writer http.ResponseWriter, request *http.Request) {
		// call function GetProject from projects package
		project.GetProject(writer, request, r.GetDB())
	}).Methods("GET")
	GetReq.HandleFunc("/projects/{category}", func(writer http.ResponseWriter, request *http.Request) {
		// call function GetProjectsByCategory from projects package
		project.GetProjectsByCategory(writer, request, r.GetDB())
	}).Methods("GET")
}

func (r *Router) Votes() {
	router := r.GetRouter()
	PostReq := router.PathPrefix("/post").Subrouter().StrictSlash(true)
	UpdateReq := router.PathPrefix("/update").Subrouter().StrictSlash(true)
	PostReq.HandleFunc("/vote", func(writer http.ResponseWriter, request *http.Request) {
		// call function PostVote from projects package
		votes.PostVote(writer, request, r.GetDB())
	}).Methods("POST")
	UpdateReq.HandleFunc("/verify_vote", func(writer http.ResponseWriter, request *http.Request) {
		// call function VerifyVote from projects package
		votes.VerifyVote(writer, request, r.GetDB())
	}).Methods("PUT")
}

func (r *Router) Database() {
	router := r.GetRouter()
	AdminReq := router.PathPrefix("/admin").Subrouter().StrictSlash(true)
	AdminReq.HandleFunc("/init", func(w http.ResponseWriter, req *http.Request) {
		database.Migrate(r.GetDB())
		// return response with status code 200 and message "Database initialized"
		w.Write([]byte("Database initialized"))
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")
	AdminReq.HandleFunc("/drop", func(w http.ResponseWriter, req *http.Request) {
		database.Drop(r.GetDB())
		// return response with status code 200 and message "Database dropped"
		w.Write([]byte("Database dropped"))
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")
}

func (r *Router) Run() {
	r.Database()
	r.Projects()
	r.Votes()
	r.Init()
	fmt.Println("Routes initialized")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
		return
	}
}
