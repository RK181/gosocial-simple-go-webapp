package shared

import (
	"html/template"
	"module/models"
	"net/http"
)

type contextKey string

const (
	BASE_URL                   = "http://localhost"
	AUTH_USER       contextKey = "authUser"
	AUTH_USER_TOKEN string     = "authUserToken"
)

// Templates almacena las plantillas cargadas en memoria
var Templates map[string]*template.Template

func ReturnView(w http.ResponseWriter, r *http.Request, tmplName string, data map[string]interface{}) {
	// Añadimos la información de autenticación
	if data == nil {
		data = make(map[string]interface{})
	}
	_, data["isAuth"] = (r.Context().Value(AUTH_USER).(models.User))

	tmpl := Templates[tmplName]
	tmpl.Execute(w, data)
}
