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
	BASE_URL = "http://localhost" + PORT
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
		middleware.Logging,      // Middleware de logging
		middleware.CompressGzip, // Middleware de compresión GZIP
		//middleware.CatchError,   // Middleware para capturar errores
	)

	// Creamos un servidor
	server := &http.Server{
		Addr:    PORT,          // Puerto en el que escucha el servidor
		Handler: stack(router), // Registramos los middlewares
	}

	// Mostramos un mensaje en consola
	log.Printf("Server is listening at %s ...\n", BASE_URL)
	log.Println("Press Ctrl + C to stop the server")

	// Iniciamos el servidor
	log.Fatal(server.ListenAndServe())
}

func loadRouterX() *http.ServeMux {
	// Creamos los routers
	router := http.NewServeMux()
	// Creamos un sub-enrutador para las rutas que requieren autenticación
	authRouter := http.NewServeMux()

	// Obtenemos el controlador de usuario
	userController := &controllers.UserController{}

	// Registramos las rutas
	router.HandleFunc("/favicon.ico", faviconHandler)
	router.HandleFunc("GET /user/{id}", userController.UserByIDGet)

	authRouter.HandleFunc("GET /", controllers.Home)
	authRouter.HandleFunc("GET /login", userController.LoginGet)
	authRouter.HandleFunc("POST /login", userController.LoginPost)

	authRouter.HandleFunc("GET /register", userController.RegisterGet)
	authRouter.HandleFunc("POST /register", userController.RegisterPost)

	authRouter.HandleFunc("GET /logout", userController.LogoutGet)

	authRouter.HandleFunc("GET /profile", userController.ProfileGet)
	authRouter.HandleFunc("GET /profile/update", userController.UpdateProfileGet)
	authRouter.HandleFunc("POST /profile/update", userController.UpdateProfilePut)

	postController := &controllers.PostController{}

	authRouter.HandleFunc("GET /post/create", postController.CreatePostGet)
	authRouter.HandleFunc("POST /post/create", postController.CreatePostPost)

	authRouter.HandleFunc("GET /post/{id}/update", postController.UpdatePostGet)
	authRouter.HandleFunc("POST /post/{id}/update", postController.UpdatePostPut)

	router.Handle("/", middleware.IsAuthenticated(authRouter))

	return router
}
