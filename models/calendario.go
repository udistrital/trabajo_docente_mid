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
