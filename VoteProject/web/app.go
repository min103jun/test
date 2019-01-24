package web

import (
	"fmt"
	"github.com/chainHero/heroes-service/web/controllers"
	"net/http"
)

func Serve(app *controllers.Application) {
	fs := http.FileServer(http.Dir("web/templates"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/home.html", app.HomeHandler)
	http.HandleFunc("/request.html", app.RequestHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/home.html", http.StatusTemporaryRedirect)
	})

	fmt.Println("Listening (http://localhost:8080/) ...")
	http.ListenAndServe(":8080", nil)
}
