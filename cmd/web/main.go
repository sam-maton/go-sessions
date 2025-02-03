package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	sessions     = make(map[string]bool)
	sessionMutex = &sync.Mutex{}
	users        = map[string]string{
		"test@email.com": "password123",
	}
)

func main() {
	addr := flag.String("addr", ":4000", "http address")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		render(w, "index")
	})

	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		render(w, "login")
	})

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		email := r.PostForm.Get("email")
		password := r.PostForm.Get("password")

		pw := users[email]

		if pw != password {
			fmt.Println("No password")
		}

		fmt.Printf("Login user with email %s and password %s", email, password)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	mux.HandleFunc("GET /protected", func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		render(w, "protected")
	})

	mux.HandleFunc("GET /logout", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

func render(w http.ResponseWriter, page string) {
	t, _ := template.New(page).ParseFiles("ui/html/base.html", "ui/html/partials/nav.html", fmt.Sprintf("ui/html/pages/%s.html", page))
	err := t.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err)
	}
}
