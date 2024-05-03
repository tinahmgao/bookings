package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tinahmgao/bookings/pkg/config"
	"github.com/tinahmgao/bookings/pkg/handlers"
	"github.com/tinahmgao/bookings/pkg/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber=":8080"
var app config.AppConfig
var session *scs.SessionManager

// mian is the main application function
func main() {

	// Change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false
	
	repo:= handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)
	
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		// Do nothing.Temparary fix the extra request on refresh the home page
	})

	fmt.Printf("Starting application on port %s\n", portNumber)
	// http.ListenAndServe(portNumber, nil)
	srv := &http.Server{
		Addr: portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err) 		
}