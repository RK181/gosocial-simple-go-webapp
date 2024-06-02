package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"module/controllers"
	middleware "module/middlewares"
	"module/models"
	"module/shared"
	"net/http"
)

// Server configuration
const (
	// Ruta de la plantilla base
	BASE_TEMPLATE_PATH = "./views/base.html"
	// Puerto en el que escuchará el servidor
	PORT = ":8080"
	// URL base del servidor
	BASE_URL = "https://localhost" + PORT
)

// Configuración de las rutas de las plantillas
const (
	layoutsDir   = "views/layouts"
	templatesDir = "views"
	extension    = "/*.html"
)

// Incrustamos las plantillas en el binario
var (
	//go:embed views/* views/layouts/*
	files embed.FS
)

// LoadTemplates carga las plantillas en memoria
func LoadTemplates() error {
	if shared.Templates == nil {
		shared.Templates = make(map[string]*template.Template)
	}
	tmplFiles, err := fs.ReadDir(files, templatesDir)
	if err != nil {
		return err
	}
	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}
		pt, err := template.ParseFS(files, templatesDir+"/"+tmpl.Name(), layoutsDir+extension)
		if err != nil {
			return err
		}
		shared.Templates[tmpl.Name()] = pt
	}
	return nil
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/favicon.ico")
}

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	models.InitDatabase()
	// Cargamos las plantillas
	err := LoadTemplates()
	if err != nil {
		log.Fatal(err)
	}

	// Creamos un enrutador y registramos las rutas con middlewares especificos
	router := loadRouterX()

	// Registramos los middlewares generales
	stack := middleware.CreateStack(
		middleware.SecureHeaders,
		middleware.Logging,         // Middleware de logging
		middleware.CompressGzip,    // Middleware de compresión GZIP
		middleware.HandleErrorPage, // Middleware para capturar errores
	)

	// Creamos un servidor
	server := &http.Server{
		Addr:    PORT,          // Puerto en el que escucha el servidor
		Handler: stack(router), // Registramos los middlewares

	}

	// Mostramos un mensaje en consola
	log.Printf("Server is listening at %s ...\n", "https://localhost"+PORT)
	log.Println("Press Ctrl + C to stop the server")

	// Iniciamos el servidor
	log.Fatal(server.ListenAndServeTLS("localhost.crt", "localhost.key"))
}

func loadRouterX() *http.ServeMux {
	// Creamos los routers
	router := http.NewServeMux()
	// Creamos un sub-enrutador para las rutas que comprueban la información de autenticación
	routerWithAuthInfo := http.NewServeMux()
	// Creamos un sub-enrutador para las rutas que requieren autenticación
	routerRequireAuth := http.NewServeMux()

	// Obtenemos el controlador de usuario
	userController := &controllers.UserController{}

	// Registramos las rutas
	router.HandleFunc("/favicon.ico", faviconHandler)

	routerWithAuthInfo.HandleFunc("GET /user/{id}", userController.UserByIDGet)

	routerWithAuthInfo.HandleFunc("GET /home", controllers.Home)
	routerWithAuthInfo.HandleFunc("GET /login", userController.LoginGet)
	routerWithAuthInfo.HandleFunc("POST /login", userController.LoginPost)

	routerWithAuthInfo.HandleFunc("GET /register", userController.RegisterGet)
	routerWithAuthInfo.HandleFunc("POST /register", userController.RegisterPost)

	routerRequireAuth.HandleFunc("GET /logout", userController.LogoutGet)

	routerRequireAuth.HandleFunc("GET /profile", userController.ProfileGet)
	routerRequireAuth.HandleFunc("GET /profile/update", userController.UpdateProfileGet)
	routerRequireAuth.HandleFunc("POST /profile/update", userController.UpdateProfilePut)

	postController := &controllers.PostController{}

	routerRequireAuth.HandleFunc("GET /post/create", postController.CreatePostGet)
	routerRequireAuth.HandleFunc("POST /post/create", postController.CreatePostPost)

	routerRequireAuth.HandleFunc("GET /post/{id}/update", postController.UpdatePostGet)
	routerRequireAuth.HandleFunc("POST /post/{id}/update", postController.UpdatePostPut)

	routerWithAuthInfo.Handle("/", middleware.RequireAuth(routerRequireAuth))
	router.Handle("/", middleware.FetchAuthInfo(routerWithAuthInfo))

	return router
}
