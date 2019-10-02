package web

import (
	"fmt"
	"gitlab.com/systemz/tasktab/model"
	"net/http"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			//w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		//w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	// We then get the name of the user from our cache, where we set the session token
	res := model.Redis.Get(sessionToken)
	//response, err := cache.Do("GET", sessionToken)
	//if err != nil {
	//	// If there is an error fetching from cache, return an internal server error status
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}

	//if response == nil {
	if len(res.String()) < 1 {
		// If the session token is not present in cache, return an unauthorized error
		//w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	// Finally, return the welcome message to the user
	w.Write([]byte(fmt.Sprintf("Welcome %s!", res.Val())))
}
