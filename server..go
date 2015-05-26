//server.go
package main

import (
	"bitbucket.org/ollietrex/sleepyfish/controller"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	ApiAuthAuthenticate = "/api/v1/auth/authenticate"
	ApiAuthRegister     = "/api/v1/auth/register"
	ApiAuthLogout       = "/api/v1/auth/logout"
	ApiCount            = "/api/v1/count"
	ApiSleeps           = "/api/v1/sleeps"
	ApiSleepsWithId     = "/api/v1/sleeps/{id:[0-9]+}"
	SiteContentCss      = "/content/css"
	SiteContentJs       = "/content/js"
)

func Handlers() *mux.Router {
	r := mux.NewRouter()
	//Home page
	r.HandleFunc("/", controller.HomeIndex).Methods("GET")
	//Static serving
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	r.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts/", http.FileServer(http.Dir("static/components/bootstrap/fonts/"))))
	//Site content optimisation
	r.HandleFunc(SiteContentJs, controller.ContentJsHandler).Methods("GET")
	r.HandleFunc(SiteContentCss, controller.ContentCssHandler).Methods("GET")
	//Account
	r.HandleFunc(ApiAuthAuthenticate, controller.AccountLoginHandler).Methods("POST")
	r.HandleFunc(ApiAuthRegister, controller.AccountRegisterHandler).Methods("POST")
	r.HandleFunc(ApiAuthLogout, controller.AccountLogoutHandler).Methods("GET")
	r.HandleFunc(ApiCount, controller.CountIndexHandler).Methods("GET")
	//Sleep
	r.HandleFunc(ApiSleeps, controller.SleepCreateHandler).Methods("POST")
	r.HandleFunc(ApiSleepsWithId, controller.SleepUpdateHandler).Methods("PUT")
	r.HandleFunc(ApiSleepsWithId, controller.SleepDeleteHandler).Methods("DELETE")
	r.HandleFunc(ApiSleepsWithId, controller.SleepGetHandler).Methods("GET")
	r.HandleFunc(ApiSleeps, controller.SleepSearchHandler).Methods("GET")

	return r
}

func main() {
	//Set the server
	err := http.ListenAndServe(":8000", Handlers())
	//Output an error if it does not run
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
