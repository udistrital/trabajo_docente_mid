package models

type PlanDocente struct {
	Id                  string `json:"_id,omitempty"`
	Estado_plan_id      string `json:"estado_plan_id,omitempty"`
	Docente_id          string `json:"docente_id,omitempty"`
	Tipo_vinculacion_id string `json:"tipo_vinculacion_id,omitempty"`
	Periodo_id          string `json:"periodo_id,omitempty"`
	Soporte_documental  string `json:"soporte_documental,omitempty"`
	Observacion_docente string `json:"observacion_docente,omitempty"`
	Respuesta           string `json:"respuesta,omitempty"`
	Resumen             string `json:"resumen,omitempty"`
	Activo              bool   `json:"activo,omitempty"`
	Fecha_creacion      string `json:"fecha_creacion,omitempty"`
	Fecha_modificacion  string `json:"fecha_modificacion,omitempty"`
}

type PreAsignacion struct {
	Id                        string `json:"_id,omitempty"`
	Docente_id                string `json:"docente_id,omitempty"`
	Tipo_vinculacion_id       string `json:"tipo_vinculacion_id,omitempty"`
	Espacio_academico_id      string `json:"espacio_academico_id,omitempty"`
	Periodo_id                string `json:"periodo_id,omitempty"`
	Aprobacion_docente        bool   `json:"aprobacion_docente,omitempty"`
	Aprobacion_proyecto       bool   `json:"aprobacion_proyecto,omitempty"`
	Plan_docente_id           string `json:"plan_docente_id,omitempty"`
	Activo                    bool   `json:"activo,omitempty"`
	Proyecto_academico_id     string `json:"proyecto_academico_id,omitempty"`
	Proyecto_academico_nombre string `json:"proyecto_academico_nombre,omitempty"`
	Fecha_creacion            string `json:"fecha_creacion,omitempty"`
	Fecha_modificacion        string `json:"fecha_modificacion,omitempty"`
}

type CargaPlan struct {
	Id                              string  `json:"_id,omitempty"`
	Espacio_academico_id            string  `json:"espacio_academico_id,omitempty"`
	Actividad_id                    string  `json:"actividad_id,omitempty"`
	Plan_docente_id                 string  `json:"plan_docente_id,omitempty"`
	Sede_id                         string  `json:"sede_id,omitempty"`     // Remove from db, only here to avoid edit to much code
	Edificio_id                     string  `json:"edificio_id,omitempty"` // Remove from db
	Salon_id                        string  `json:"salon_id,omitempty"`
	Horario                         string  `json:"horario,omitempty"`                         // Remove from db
	Colocacion_espacio_academico_id string  `json:"colocacion_espacio_academico_id,omitempty"` // those removed is in horarios model
	Hora_inicio                     int     `json:"hora_inicio,omitempty"`
	Duracion                        float64 `json:"duracion,omitempty"`
	Activo                          bool    `json:"activo,omitempty"`
	Fecha_creacion                  string  `json:"fecha_creacion,omitempty"`
	Fecha_modificacion              string  `json:"fecha_modificacion,omitempty"`
}

type Actividad struct {
	Id                 string `json:"_id,omitempty"`
	Nombre             string `json:"nombre,omitempty"`
	Descripcion        string `json:"descripcion,omitempty"`
	Codigo_abreviacion int    `json:"codigo_abreviacion,omitempty"`
	Activo             bool   `json:"activo,omitempty"`
	Fecha_creacion     string `json:"fecha_creacion,omitempty"`
	Fecha_modificacion string `json:"fecha_modificacion,omitempty"`
}

type EstadoPlan struct {
	Id                 string `json:"_id,omitempty"`
	Nombre             string `json:"nombre,omitempty"`
	Descripcion        string `json:"descripcion,omitempty"`
	Codigo_abreviacion int    `json:"codigo_abreviacion,omitempty"`
	Activo             bool   `json:"activo,omitempty"`
	Fecha_creacion     string `json:"fecha_creacion,omitempty"`
	Fecha_modificacion string `json:"fecha_modificacion,omitempty"`
}

type ConsolidadoDocente struct {
	Id                       string `json:"_id,omitempty"`
	Plan_docente_id          string `json:"plan_docente_id,omitempty"`
	Periodo_id               string `json:"periodo_id,omitempty"`
	Proyecto_academico_id    string `json:"proyecto_academico_id,omitempty"`
	Estado_consolidado_id    string `json:"estado_consolidado_id,omitempty"`
	Respuesta_decanatura     string `json:"respuesta_decanatura,omitempty"`
	Consolidado_coordinacion string `json:"consolidado_coordinacion,omitempty"`
	Cumple_normativa         bool   `json:"cumple_normativa,omitempty"`
	Aprobado                 bool   `json:"aprobado,omitempty"`
	Activo                   bool   `json:"activo,omitempty"`
	Fecha_creacion           string `json:"fecha_creacion,omitempty"`
	Fecha_modificacion       string `json:"fecha_modificacion,omitempty"`
}

type EstadoConsolidado struct {
	Id                 string `json:"_id,omitempty"`
	Nombre             string `json:"nombre,omitempty"`
	Descripcion        string `json:"descripcion,omitempty"`
	Codigo_abreviacion int    `json:"codigo_abreviacion,omitempty"`
	Activo             bool   `json:"activo,omitempty"`
	Fecha_creacion     string `json:"fecha_creacion,omitempty"`
	Fecha_modificacion string `json:"fecha_modificacion,omitempty"`
}
