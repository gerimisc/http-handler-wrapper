package main

import (
	"fmt"
	"log"
	"net/http"
)

type authHandler struct {
	next http.Handler
}

type secretHandler struct{}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		// not authenticated
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	if err != nil {
		// some other error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, success. call the next handler
	h.next.ServeHTTP(w, r)
}

func mustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func (s *secretHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Secret is 'secret'")
}

// comparison
func authHandlerV2(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("before")
		h.ServeHTTP(w, r)
		log.Println("after")
	})
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "welcome")
	})

	http.Handle("/secret", mustAuth(&secretHandler{}))
	http.Handle("/secret2", authHandlerV2(&secretHandler{}))

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
