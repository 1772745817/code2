package main

import (
	"fmt"
	"net/http"
	"time"
)

const loginPage = `
	<html>
		<body>
			<h2>Login</h2>
				<form method="POST" action="/login">
					Username: <input type="text"  name="username" />
					<input type="submit" value="Login" />
				</form>
		</body>
	</html> 
	`

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write([]byte(loginPage))
	} else if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")

		if username != "" {
			expiration := time.Now().Add(24 * time.Hour)

			cookie := http.Cookie{
				Name:     "username",
				Value:    username,
				Expires:  expiration,
				HttpOnly: true,
			}
			//响应标头 : Set-Cookie: username=deng; Expires=Mon, 14 Oct 2024 03:19:31 GMT; HttpOnly                     0
			http.SetCookie(w, &cookie)

			w.Write([]byte("Login Successful. Welcome " + username + "!\n"))
		} else {
			w.Write([]byte("username cannot be empty !\n"))
		}
	}

}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil {
		w.Write([]byte("you are not logged , please login first"))
		return
	}
	username := cookie.Value

	w.Write([]byte("welcome back! ," + username))
}

func logOutHandler(w http.ResponseWriter, r *http.Request) {

	cookie := http.Cookie{
		Name:   "username",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
	w.Write([]byte("you have been logged out"))
}
func main() {
	http.HandleFunc("/login", loginPageHandler)
	http.HandleFunc("/welcome", welcomeHandler)
	http.HandleFunc("/logout", logOutHandler)

	fmt.Println("sevel running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
