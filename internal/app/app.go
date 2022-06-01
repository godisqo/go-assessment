package app

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"go-assessment/internal/config"
	"io"
	"log"
	"net/http"
)

type User struct {
	Id        string `json:"id"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
}

type App struct {
	Config *config.Config
}

// New creates a new App
func New(cfg config.Config) App {

	app := App{
		Config: &cfg,
	}

	return app
}

// router function handles assignment of routes to handlers
// define your paths and middleware here
func (a *App) router() http.Handler {
	r := mux.NewRouter()

	// health check end point, initialized without validation middleware
	r.HandleFunc("/health-check", a.HealthCheckHandler).Methods(http.MethodGet)

	// user endpoints
	r.HandleFunc("/user", a.CreateUserHandler).Methods(http.MethodPost)
	r.HandleFunc("/user/{userId}", a.UpdateUserHandler).Methods(http.MethodPatch)
	r.HandleFunc("/user/{userId}", a.GetUserHandler).Methods(http.MethodGet)
	r.HandleFunc("/user/{userId}", a.DeleteUserHandler).Methods(http.MethodDelete)

	return r
}

// HealthCheckHandler should be used to check the health of the application
func (a *App) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (a *App) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:password@tcp(0.0.0.0:1444)/badass_db")
	if err != nil {
		panic(err.(any))
	}

	var userData User

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(b, &userData)
	if err != nil {
		panic(err.(any))
	}

	_, err = db.Exec("INSERT INTO badass_users(first_name, last_name) VALUES(?,?)", userData.Firstname, userData.Lastname)
	if err != nil {
		panic(err.(any))
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *App) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:password@tcp(0.0.0.0:1444)/badass_db")
	if err != nil {
		panic(err.(any))
	}

	var userData User
	params := mux.Vars(r)
	userData.Id = params["userId"]

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(b, &userData)
	if err != nil {
		panic(err.(any))
	}

	_, err = db.Exec("UPDATE badass_users SET first_name=?, last_name=? WHERE id = ?", userData.Firstname, userData.Lastname, userData.Id)
	if err != nil {
		panic(err.(any))
	}

	w.WriteHeader(http.StatusOK)
}

func (a *App) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// implement Get User Handler
	//w.Header().Set()
	//w.WriteHeader()
	//_, err = w.Write()
}

func (a *App) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:password@tcp(0.0.0.0:1444)/badass_db")
	if err != nil {
		panic(err.(any))
	}

	var userData User
	params := mux.Vars(r)
	userData.Id = params["userId"]

	_, err = db.Exec("DELETE FROM badass_users WHERE id = ?", userData.Id)

	if err != nil {
		panic(err.(any))
	}

	w.WriteHeader(http.StatusNoContent)
}
