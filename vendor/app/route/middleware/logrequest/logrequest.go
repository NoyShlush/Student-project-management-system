package logrequest

import (
	"app/shared/session"
	"fmt"
	"net/http"
	"time"
)

// Handler will log the HTTP requests
func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := session.Instance(r)
		fmt.Println(time.Now().Format("2006-01-02 03:04:05 PM"), r.RemoteAddr, r.Method, r.URL, "\n --- User info ---\n", "ID:", sess.Values["id"], "eMail:", sess.Values["email"])
		next.ServeHTTP(w, r)
	})
}
