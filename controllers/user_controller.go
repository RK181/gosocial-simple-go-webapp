package controllers

import (
	"log"
	"module/models"
	"module/shared"
	"module/utils"
	"net/http"
	"strconv"
)

type UserFormError struct {
	Username    string
	Email       string
	Password    string
	Phone       string
	Description string
}

func (uferr UserFormError) hasErrors() bool {
	return uferr.Username != "" || uferr.Email != "" || uferr.Password != "" || uferr.Phone != "" || uferr.Description != ""
}

type UserController struct{}

func (c *UserController) LogoutGet(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	cookie := http.Cookie{
		Name:     shared.AUTH_USER_TOKEN,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
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

	data := make(map[string]interface{})
	data["Email"] = email
	data["Password"] = password

	// Validamos los datos del formulario
	var formErrors UserFormError
	formErrors.Email = RequiredField(email)
	formErrors.Email = IsValidEmail(email)
	formErrors.Password = RequiredField(password)

	if formErrors.hasErrors() {
		data["Errors"] = formErrors
		returnView(w, r, "register.html", data)
		return
	}

	user := models.User{Email: email, Password: []byte(password)}
	err := user.LoginUser()
	if err != nil {
		data := make(map[string]interface{})
		data["Errors"] = UserFormError{
			Email:    "Check if the email is correct",
			Password: "Check if the password is correct",
		}
		returnView(w, r, "login.html", data)
		return
	}

	cookie := http.Cookie{
		Name:     shared.AUTH_USER_TOKEN,
		Value:    utils.Encode64(user.SessionToken),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (c *UserController) RegisterGet(w http.ResponseWriter, r *http.Request) {
	returnView(w, r, "register.html", nil)
}

func (c *UserController) RegisterPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	description := r.FormValue("description")

	data := make(map[string]interface{})
	data["Username"] = username
	data["Password"] = password
	data["Email"] = email
	data["Phone"] = phone
	data["Description"] = description
	// Validate the form
	var formErrors UserFormError
	formErrors.Username = RequiredField(username)
	formErrors.Email = IsValidEmail(email)
	formErrors.Password = VarifyPassword(password)
	formErrors.Phone = RequiredField(phone)
	if UserAlredyExists(email) {
		formErrors.Email = "Email already exists"
	}

	if formErrors.hasErrors() {
		data["Errors"] = formErrors
		returnView(w, r, "register.html", data)
		return
	}

	user := models.User{
		Username:    username,
		Password:    []byte(password),
		Email:       email,
		PhoneNumber: phone,
		Description: description,
	}

	// Register the user
	err := user.RegisterUser()
	if err != nil {
		returnView(w, r, "register.html", data)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

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

	// Obtenemos el usuario autenticado
	authUser := r.Context().Value(shared.AUTH_USER).(models.User)
	// Comprobamos si el usuario autenticado esta suscrito al usuario
	subscriber := models.UserUserSubscription{}
	isSubscribed := false
	// Comprobamos si el usuario esta autenticado
	if authUser.ID > 0 {
		isSubscribed = subscriber.CheckSubscriptionByUserID(user.ID, authUser.ID)
	}
	// Obtenemos los posts del usuario por ID y si esta suscrito
	posts, err := models.NewPost().GetPostsByUserID(user.ID, isSubscribed)
	if err != nil {
		http.Error(w, "Posts not found", http.StatusNotFound)
		return
	}

	data := make(map[string]interface{})
	data["User"] = user
	data["Posts"] = posts
	data["isSubscribed"] = isSubscribed

	returnView(w, r, "user.html", data)
}

func (c *UserController) ProfileGet(w http.ResponseWriter, r *http.Request) {
	// Obtenemos el usuario autenticado
	user := r.Context().Value(shared.AUTH_USER).(models.User)
	// Obtenemos los posts del usuario
	posts, err := models.NewPost().GetPostsByUserID(user.ID, true)
	if err != nil {
		log.Println(err.Error())
	}

	data := make(map[string]interface{})
	data["User"] = user
	data["Posts"] = posts

	returnView(w, r, "profile.html", data)
}

func (c *UserController) UpdateProfileGet(w http.ResponseWriter, r *http.Request) {
	// Obtenemos el usuario autenticado
	user := r.Context().Value(shared.AUTH_USER).(models.User)

	data := make(map[string]interface{})
	data["User"] = user

	returnView(w, r, "profileUpdateForm.html", data)
}

func (c *UserController) UpdateProfilePost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	description := r.FormValue("description")

	data := make(map[string]interface{})
	data["Username"] = username
	data["Email"] = email
	data["Phone"] = phone
	data["Description"] = description
	// Validate the form
	var formErrors UserFormError
	formErrors.Username = RequiredField(username)
	formErrors.Email = IsValidEmail(email)
	formErrors.Phone = RequiredField(phone)
	if UserAlredyExists(email) {
		formErrors.Email = "Email already exists"
	}

	if formErrors.hasErrors() {
		data["Errors"] = formErrors
		returnView(w, r, "register.html", data)
		return
	}

	user := r.Context().Value(shared.AUTH_USER).(models.User)
	user.Username = username
	user.Email = email
	user.PhoneNumber = phone
	user.Description = description

	// Register the user
	err := user.UpdateUser()
	if err != nil {
		http.Redirect(w, r, "/profile/update", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (c *UserController) UpdatePrivacyGet(w http.ResponseWriter, r *http.Request) {
	// Obtenemos el usuario autenticado
	user := r.Context().Value(shared.AUTH_USER).(models.User)

	data := make(map[string]interface{})
	data["User"] = user

	returnView(w, r, "privacyUpdateForm.html", data)
}

func (c *UserController) UpdatePrivacyPost(w http.ResponseWriter, r *http.Request) {
	showEmail := r.FormValue("showEmail")
	showPhoneNumber := r.FormValue("showPhoneNumber")

	// TO-DO: Validar los datos del formulario

	user := r.Context().Value(shared.AUTH_USER).(models.User)
	user.ShowEmail = showEmail == "on"
	user.ShowPhoneNumber = showPhoneNumber == "on"

	// Register the user
	err := user.UpdateUser()
	if err != nil {
		http.Redirect(w, r, "/privacy/update", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (c *UserController) Subscribe(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(shared.AUTH_USER).(models.User)
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	subscriber := models.UserUserSubscription{}
	err = subscriber.SubscribeToUserByID(user.ID, id)
	if err != nil {
		http.Error(w, "Error subscribing", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/user/"+r.PathValue("id"), http.StatusSeeOther)
}

func (c *UserController) UnSubscribe(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(shared.AUTH_USER).(models.User)
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	subscriber := models.UserUserSubscription{}
	err = subscriber.UnsubscribeToUserByID(user.ID, id)
	if err != nil {
		http.Error(w, "Error unsubscribing", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/user/"+r.PathValue("id"), http.StatusSeeOther)
}
