package shared

import "html/template"

const BASE_URL = "http://localhost"

// Templates almacena las plantillas cargadas en memoria
var Templates map[string]*template.Template
