package controllers

import (
	"module/models"
	"module/shared"
	"net/http"
	"strconv"
)

type PostFormError struct {
	Title         string
	Content       string
	ForSubcribers bool
}

func (pferr PostFormError) hasErrors() bool {
	return pferr.Title != "" || pferr.Content != ""
}

type PostController struct{}

/*
CreatePostGet devuelve la pagina con el formulario para crear un post
*/
func (c *PostController) CreatePostGet(w http.ResponseWriter, r *http.Request) {
	returnView(w, r, "postCreateUpdateForm.html", nil)
}

/*
CreatePostPost crea un post en la base de datos
*/
func (c *PostController) CreatePostPost(w http.ResponseWriter, r *http.Request) {
	// Obtenemos los datos del formulario
	title := r.FormValue("title")
	content := r.FormValue("content")
	forSubcribers := r.FormValue("forSubcribers")

	data := make(map[string]interface{})
	data["Title"] = title
	data["Content"] = content

	// Validamos los datos del formulario
	var formErrors PostFormError
	formErrors.Title = RequiredField(title)
	formErrors.Content = RequiredField(content)
	if formErrors.hasErrors() {
		data["Errors"] = formErrors
		returnView(w, r, "postCreateUpdateForm.html", data)
		return
	}

	user := r.Context().Value(shared.AUTH_USER).(models.User)
	post := models.Post{UserID: user.ID, Title: title, Content: content, ForSubcribers: forSubcribers == "on"}
	err := post.CreatePost()
	if err != nil {
		data["Errors"] = PostFormError{
			Title:         "Check if the title is correct",
			Content:       "Check if the content is correct",
			ForSubcribers: forSubcribers == "on",
		}
		returnView(w, r, "postCreateUpdateForm.html", data)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

/*
UpdatePostGet devuelve la pagina con el formulario de actualización de un post
*/
func (c *PostController) UpdatePostGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
	}

	user := r.Context().Value(shared.AUTH_USER).(models.User)

	post := models.Post{}
	post.GetPostByID(id)

	if post.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	data := make(map[string]interface{})
	data["Post"] = post

	returnView(w, r, "postCreateUpdateForm.html", data)
}

/*
* UpdatePostPost actualiza un post en la base de datos
 */
func (c *PostController) UpdatePostPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
	}

	user := r.Context().Value(shared.AUTH_USER).(models.User)

	post := models.Post{}
	post.GetPostByID(id)

	// Comprobamos si el usuario autenticado es el propietario del post
	if post.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	forSubcribers := r.FormValue("forSubcribers")

	post.Title = title
	post.Content = content
	post.ForSubcribers = forSubcribers == "on"

	data := make(map[string]interface{})
	// Validamos los datos del formulario
	var formErrors PostFormError
	formErrors.Title = RequiredField(title)
	formErrors.Content = RequiredField(content)
	if formErrors.hasErrors() {
		data["Post"] = post
		data["Errors"] = formErrors
		returnView(w, r, "postCreateUpdateForm.html", data)
		return
	}

	// Actualizamos el post
	err = post.UpdatePost()
	if err != nil {
		data["Post"] = post
		data["Errors"] = PostFormError{
			Title:         "Check if the title is correct",
			Content:       "Check if the content is correct",
			ForSubcribers: forSubcribers == "on",
		}
		returnView(w, r, "postCreateUpdateForm.html", data)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

/*
DeletePostGet elimina un post de la base de datos
*/
func (c *PostController) DeletePostGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
	}

	user := r.Context().Value(shared.AUTH_USER).(models.User)

	post := models.Post{}
	post.GetPostByID(id)
	post.DeletePost()

	if post.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
