package main

import (
	"flag"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	addr := flag.String("addr", ":4000", "http address")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		t, _ := template.New("home").ParseFiles("ui/html/base.html", "ui/html/partials/nav.html", "ui/html/pages/index.html")
		err := t.ExecuteTemplate(w, "base", nil)
		if err != nil {
			log.Println(err)
		}
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.New("login").ParseFiles("ui/html/base.html", "ui/html/partials/nav.html", "ui/html/pages/login.html")
		err := t.ExecuteTemplate(w, "base", nil)
		if err != nil {
			log.Println(err)
		}
	})

	mux.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.New("login").ParseFiles("ui/html/base.html", "ui/html/partials/nav.html", "ui/html/pages/protected.html")
		err := t.ExecuteTemplate(w, "base", nil)
		if err != nil {
			log.Println(err)
		}
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusOK)
	})

	srv := &http.Server{
		Addr:         *addr,
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	slog.Info("starting server on http://localhost" + *addr)
	err := srv.ListenAndServe()
	slog.Error(err.Error())
	os.Exit(1)
}
