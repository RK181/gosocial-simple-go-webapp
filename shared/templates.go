package shared

import (
	"html/template"
	"module/models"
	"net/http"
)

type contextKey string

const (
	AUTH_USER       contextKey = "authUser"      // Clave para almacenar el usuario autenticado en el contexto
	AUTH_USER_TOKEN string     = "authUserToken" // Clave para almacenar el token de sesión en las cookies
)

// Templates almacena las plantillas cargadas en memoria
var Templates map[string]*template.Template

// Devuelve la vista con la plantilla, información de autenticación y datos adicionales
func ReturnView(w http.ResponseWriter, r *http.Request, tmplName string, data map[string]interface{}) {
	// Añadimos la información de autenticación
	if data == nil {
		data = make(map[string]interface{})
	}
	data["AuthUser"], data["isAuth"] = (r.Context().Value(AUTH_USER).(models.User))

	tmpl := Templates[tmplName]
	tmpl.Execute(w, data)
}
