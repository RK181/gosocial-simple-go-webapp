package controllers

import (
	"module/models"
	"module/shared"
	"net/http"
)

type Post struct {
	ID              int `storm:"id,increment"`
	UserID          int `storm:"index"`
	Title           string
	Content         string
	ForSupscriptors bool
}

func Home(w http.ResponseWriter, r *http.Request) {
	user, isAuth := r.Context().Value(shared.AUTH_USER).(models.User)
	subscriptions := []models.UserUserSubscription{}
	if isAuth {
		// Obtener las suscripciones del usuario
		subscriptions, _ = models.NewUserUserSubscription().GetSubscriptionsByUserID(user.ID)
	}
	// Obtener los posts
	posts, _ := models.NewPost().GetAllPosts(subscriptions)
	data := make(map[string]interface{})
	data["Posts"] = posts

	returnView(w, r, "home.html", data)
}

func FavIcon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "relative/path/to/favicon.ico")
}

/*
Devuelve la vista con los datos proporcionados y la información de autenticación
*/
func returnView(w http.ResponseWriter, r *http.Request, tmplName string, data map[string]interface{}) {
	// Añadimos la información de autenticación
	shared.ReturnView(w, r, tmplName, data)
}
