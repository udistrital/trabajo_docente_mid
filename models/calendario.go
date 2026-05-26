package models

type EventosXML struct {
	Eventos []EventoXML `xml:"evento"`
}

type EventoXML struct {
	Descripcion  string `xml:"descripcion"`
	CodigoEvento string `xml:"codigo_evento"`
	Editable     string `xml:"editable"`
}

type CalendariosEventosXML struct {
	Eventos []CalendarioEventoXML `xml:"evento"`
}

type CalendarioEventoXML struct {
	CodigoEvento   string `xml:"codigo_evento"`
	FechaInicio    string `xml:"fecha_inicio"`
	TipoProyecto   string `xml:"tipo_proyecto"`
	Ciclo          string `xml:"ciclo"`
	CodigoProyecto string `xml:"codigo_proyecto"`
	FechaFin       string `xml:"fecha_fin"`
	CodigoFacultad string `xml:"codigo_facultad"`
	Year           string `xml:"anio"`
}

type FacultadesDecanoXML struct {
	Decanos []FacultadDecanoXML `xml:"decano"`
}

type FacultadDecanoXML struct {
	FechaDesde     string `xml:"fecha_desde"`
	CodigoFacultad string `xml:"codigo_facultad"`
	Nombre         string `xml:"nombre"`
	FechaHasta     string `xml:"fecha_hasta"`
	Facultad       string `xml:"facultad"`
}

type ProyectosFacultadXML struct {
	Proyectos []ProyectoFacultadXML `xml:"proyecto"`
}

type ProyectoFacultadXML struct {
	NombreProyectoCurricular string `xml:"nombre_proyecto_curricular"`
	CodigoProyectoCurricular string `xml:"codigo_proyecto_curricular"`
}
