package controllers

import (
	"log"
	"module/middlewares"
	"module/models"
	"module/utils"
	"net/http"
	"strconv"
)

type UserFormError struct {
	Username string
	Email    string
	Password string
	Phone    string
}

type UserController struct{}

func (c *UserController) LogoutGet(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	cookie := http.Cookie{
		Name:   middlewares.AUTH_USER_TOKEN,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// User login GET page
func (c *UserController) LoginGet(w http.ResponseWriter, r *http.Request) {
	returnView(w, r, "login.html", nil)
}

// User login POST form
func (c *UserController) LoginPost(w http.ResponseWriter, r *http.Request) {
	// Obtenemos los datos del formulario
	email := r.FormValue("email")
	password := r.FormValue("password")

	// TO-DO: Validar los datos del formulario

	user := models.User{Email: email, Password: []byte(password)}
	err := user.LoginUser()
	if err != nil {
		data := make(map[string]interface{})
		data["Email"] = email
		data["Password"] = password
		data["Errors"] = UserFormError{
			Email:    "Check if the email is correct",
			Password: "Check if the password is correct",
		}
		returnView(w, r, "login.html", data)
		return
	}

	cookie := http.Cookie{
		Name:     middlewares.AUTH_USER_TOKEN,
		Value:    utils.Encode64(user.SessionToken),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

// User register GET form page
func (c *UserController) RegisterGet(w http.ResponseWriter, r *http.Request) {
	returnView(w, r, "register.html", nil)
}

// User register POST form
func (c *UserController) RegisterPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	phone := r.FormValue("phone")

	// TO-DO: Validar los datos del formulario

	user := models.User{
		Username:    username,
		Password:    []byte(password),
		Email:       email,
		PhoneNumber: phone,
	}
	// Register the user
	err := user.RegisterUser()
	if err != nil {
		returnView(w, r, "register.html", nil)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// User by ID GET page
func (c *UserController) UserByIDGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
	}
	user := models.User{}
	err = user.GetUserByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	posts, err := models.NewPost().GetPostsByUserID(id)
	if err != nil {
		http.Error(w, "Posts not found", http.StatusNotFound)
		return
	}

	data := make(map[string]interface{})
	data["User"] = user
	data["Posts"] = posts

	returnView(w, r, "user.html", data)
}

// User profile GET page
func (c *UserController) ProfileGet(w http.ResponseWriter, r *http.Request) {
	// Obtenemos el usuario autenticado
	user := r.Context().Value(middlewares.AUTH_USER).(models.User)
	// Obtenemos los posts del usuario
	posts, err := models.NewPost().GetPostsByUserID(user.ID)
	if err != nil {
		log.Println(err.Error())
	}

	data := make(map[string]interface{})
	data["User"] = user
	data["Posts"] = posts

	returnView(w, r, "profile.html", data)
}

// User profile update GET form page
func (c *UserController) UpdateProfileGet(w http.ResponseWriter, r *http.Request) {
	// Obtenemos el usuario autenticado
	user := r.Context().Value(middlewares.AUTH_USER).(models.User)

	data := make(map[string]interface{})
	data["User"] = user

	returnView(w, r, "profileUpdateForm.html", data)
}

// User profile update PUT form
func (c *UserController) UpdateProfilePut(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	phone := r.FormValue("phone")

	// TO-DO: Validar los datos del formulario

	user := r.Context().Value(middlewares.AUTH_USER).(models.User)
	user.Username = username
	user.Email = email
	user.PhoneNumber = phone

	// Register the user
	err := user.UpdateUser()
	if err != nil {
		http.Redirect(w, r, "/profile/update", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
