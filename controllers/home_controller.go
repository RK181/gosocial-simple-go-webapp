package controllers

import (
	"module/middlewares"
	"module/models"
	"module/shared"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	posts, _ := models.NewPost().GetAllPosts()
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
	if data == nil {
		data = make(map[string]interface{})
	}
	//data["isAuth"] = (r.Context().Value(middlewares.AUTH_USER) != nil)
	_, data["isAuth"] = (r.Context().Value(middlewares.AUTH_USER).(models.User))

	tmpl := shared.Templates[tmplName]
	tmpl.Execute(w, data)
}
