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

type PostController struct{}

func (c *PostController) CreatePostGet(w http.ResponseWriter, r *http.Request) {
	returnView(w, r, "postCreateUpdateForm.html", nil)
}

func (c *PostController) CreatePostPost(w http.ResponseWriter, r *http.Request) {
	// Obtenemos los datos del formulario
	title := r.FormValue("title")
	content := r.FormValue("content")
	forSubcribers := r.FormValue("forSubcribers")

	data := make(map[string]interface{})
	// TO-DO: Validar los datos del formulario

	user := r.Context().Value(shared.AUTH_USER).(models.User)
	post := models.Post{UserID: user.ID, Title: title, Content: content, ForSubcribers: forSubcribers == "on"}
	err := post.CreatePost()
	if err != nil {
		data["Title"] = title
		data["Content"] = content
		data["Errors"] = PostFormError{
			Title:         "Check if the title is correct",
			Content:       "Check if the content is correct",
			ForSubcribers: forSubcribers == "on",
		}
		tmpl := shared.Templates["postCreateUpdateForm.html"]
		tmpl.Execute(w, data)
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

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

func (c *PostController) UpdatePostPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
	}

	user := r.Context().Value(shared.AUTH_USER).(models.User)

	post := models.Post{}
	post.GetPostByID(id)

	if post.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	forSubcribers := r.FormValue("forSubcribers")

	// TO-DO: Validar los datos del formulario

	post.Title = title
	post.Content = content
	post.ForSubcribers = forSubcribers == "on"

	err = post.UpdatePost()
	if err != nil {
		data := make(map[string]interface{})
		data["Post"] = post
		data["Errors"] = PostFormError{
			Title:         "Check if the title is correct",
			Content:       "Check if the content is correct",
			ForSubcribers: forSubcribers == "on",
		}

		returnView(w, r, "postCreateUpdateForm.html", data)
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

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
