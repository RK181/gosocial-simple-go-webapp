package main

import (
	"module/controllers"
	"net/http"
)

func loadRouter() *http.ServeMux {
	// Creamos los routers
	router := http.NewServeMux()
	// Creamos un sub-enrutador para las rutas que requieren autenticación
	authRouter := http.NewServeMux()

	// Obtenemos el controlador de usuario
	userController := &controllers.UserController{}

	// Registramos las rutas
	router.HandleFunc("GET /user/{id}", userController.UserByIDGet)

	authRouter.HandleFunc("GET /login", userController.LoginGet)
	authRouter.HandleFunc("POST /login", userController.LoginPost)

	authRouter.HandleFunc("GET /register", userController.RegisterGet)
	authRouter.HandleFunc("POST /register", userController.RegisterPost)

	authRouter.HandleFunc("GET /logout", userController.LogoutGet)

	authRouter.HandleFunc("GET /profile", userController.ProfileGet)
	authRouter.HandleFunc("GET /profile/update", userController.UpdateProfileGet)
	authRouter.HandleFunc("POST /profile/update", userController.UpdateProfilePost)

	postController := &controllers.PostController{}

	authRouter.HandleFunc("GET /post/create", postController.CreatePostGet)
	authRouter.HandleFunc("POST /post/create", postController.CreatePostPost)

	authRouter.HandleFunc("GET /post/{id}/update", postController.UpdatePostGet)
	authRouter.HandleFunc("POST /post/{id}/update", postController.UpdatePostPost)

	router.Handle("/", authRouter)

	return router
}
