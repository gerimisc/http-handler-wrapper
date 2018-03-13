### http-handler-wrapper
Performing authorization using http handler wrappers

To better understand the handler, let's start from the main:
```
func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "welcome")
	})

	http.Handle("/secret", mustAuth(&secretHandler{}))

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

```

The line of interest is `http.Handle("/secret", mustAuth(&secretHandler{}))`

Next, what does mustAuth function takes in and return?
```
func mustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}
```
The mustAuth function takes in a `http.Handler` and returning another `http.Handler` and in this case, it is the `secretHandler{}`. It will cause the execution to go through authHandler first and if true, it will run the `secretHandler{}`

#What does authHandler check for?
```
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
```

If a cookie called "auth" is set, it will call the next handler, which is `secretHandler` which already satisfies the http.Handler Interface. It will return "Secret is 'secret'" as you can deduce from the handler below:
```
func (s *secretHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Secret is 'secret'")
}
```

This will be the basis of wrapping handlers within handlers to perform auth checks.

For more info:
https://medium.com/@matryer/the-http-handler-wrapper-technique-in-golang-updated-bc7fbcffa702
