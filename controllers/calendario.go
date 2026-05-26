package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/trabajo_docente_mid/services"
	requestmanager "github.com/udistrital/utils_oas/requestresponse"
)

// CalendarioController operations for Calendario
type CalendarioController struct {
	beego.Controller
}

// URLMapping ...
func (c *CalendarioController) URLMapping() {
	c.Mapping("GetEventos", c.GetEventos)
	c.Mapping("GetCalendariosEventos", c.GetCalendariosEventos)
	c.Mapping("GetProyectosFacultadDecano", c.GetProyectosFacultadDecano)
}

// GetEventos ...
// @Title GetEventos
// @Description Lista los eventos académicos
// @Success 200 {object} []map[string]interface{}
// @Failure 404 no record exist
// @router /eventos [get]
func (c *CalendarioController) GetEventos() {
	//defer errorhandler.HandlePanic(&c.Controller)
	logs.Info("[CalendarioController] GetEventos - endpoint /eventos llamado")

	resultado := services.ListaEventos()
	c.Data["json"] = resultado
	// El código de negocio va en el body; se usa 200 para evitar que Beego
	// intercepte 4xx/5xx e intente renderizar templates de error inexistentes.
	c.Ctx.Output.SetStatus(200)

	c.ServeJSON()
}

// GetCalendariosEventos ...
// @Title GetCalendariosEventos
// @Description Realiza el cruce de eventos del calendario con los proyectos curriculares a cargo de un docente/coordinador
// @Param   documento   query   string  true        "Documento del coordinador o docente"
// @Param   codigo_evento           query   string  true        "Codigo del evento del calendario"
// @Success 200 {object} []map[string]interface{}
// @Failure 404 no record exist
// @router /calendario_eventos [get]
func (c *CalendarioController) GetCalendariosEventos() {
	//defer errorhandler.HandlePanic(&c.Controller)

	documento := c.GetString("documento")
	codigoEvento := c.GetString("codigo_evento")
	proyecto := c.GetString("proyecto")

	if documento == "" || codigoEvento == "" {
		logs.Error("documento o codigo_evento vacio")
		c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) con valores no válidos")
	} else {
		resultado := services.ConsultaCalendariosEventos(documento, codigoEvento, proyecto)
		c.Data["json"] = resultado
	}
	// El código de negocio va en el body; se usa 200 para evitar que Beego
	// intercepte 4xx/5xx e intente renderizar templates de error inexistentes.
	c.Ctx.Output.SetStatus(200)

	c.ServeJSON()
}

// GetProyectosFacultadDecano ...
// @Title GetProyectosFacultadDecano
// @Description Consulta la facultad activa del decano y retorna los proyectos de PREGRADO y POSGRADO
// @Param   documento   query   string  true        "Documento del decano"
// @Success 200 {object} []map[string]interface{}
// @Failure 404 no record exist
// @router /proyectos_facultad_decano [get]
func (c *CalendarioController) GetProyectosFacultadDecano() {
	documento := c.GetString("documento")

	if documento == "" {
		logs.Error("documento vacio")
		c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) con valores no válidos")
	} else {
		resultado := services.ConsultaProyectosFacultadDecano(documento)
		c.Data["json"] = resultado
	}

	c.Ctx.Output.SetStatus(200)
	c.ServeJSON()
}
