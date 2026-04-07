package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/trabajo_docente_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	requestmanager "github.com/udistrital/utils_oas/requestresponse"
)

// EspacioAcademicoController operations for Asignacion
type EspacioAcademicoController struct {
	beego.Controller
}

// URLMapping ...
func (c *EspacioAcademicoController) URLMapping() {
	c.Mapping("GrupoEspacioAcademico", c.GrupoEspacioAcademico)
	c.Mapping("GrupoEspacioAcademicoPadre", c.GrupoEspacioAcademicoPadre)
	c.Mapping("EspaciosAcademicosProyectoPeriodo", c.EspaciosAcademicosProyectoPeriodo)
	c.Mapping("GruposEspacioPeriodo", c.GruposEspacioPeriodo)
}

// GrupoEspacioAcademico ...
// @Title GrupoEspacioAcademico
// @Description  Lista los grupos de un espacios académico padre por vigencia
// @Param	padre		query 	string	true		"Id del espacio académico padre"
// @Param	vigencia	query 	string	true		"Vigencia del espacio académico"
// @Success 200 {}
// @Failure 400 the request contains an incorrect parameter
// @Failure 404 no record exist
// @router /grupo [get]
func (c *EspacioAcademicoController) GrupoEspacioAcademico() {
	defer errorhandler.HandlePanic(&c.Controller)

	padre := c.GetString("padre")
	vigencia := c.GetString("vigencia")

	if padre == "" || vigencia == "" {
		logs.Error(padre, vigencia)
		c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) con valores no válidos")
		c.Ctx.Output.SetStatus(400)
	} else {
		resultado := services.ListaGruposEspaciosAcademicos(padre, vigencia)
		c.Data["json"] = resultado
		c.Ctx.Output.SetStatus(resultado.Status)
	}

	c.ServeJSON()
}

// GrupoEspacioAcademicoPadre ...
// @Title GrupoEspacioAcademicoPadre
// @Description Lista los grupos de un espacios académico padre
// @Param	padre		query 	string	true		"Id del espacio académico padre"
// @Success 200 {}
// @Failure 400 the request contains an incorrect parameter
// @Failure 404 no record exist
// @router /grupo-padre [get]
func (c *EspacioAcademicoController) GrupoEspacioAcademicoPadre() {
	defer errorhandler.HandlePanic(&c.Controller)

	padre := c.GetString("padre")

	if padre == "" {
		logs.Error(padre)
		c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) con valores no válidos")
		c.Ctx.Output.SetStatus(400)
	} else {
		resultado := services.ListaGruposEspaciosAcademicosPadre(padre)
		c.Data["json"] = resultado
		c.Ctx.Output.SetStatus(resultado.Status)
	}

	c.ServeJSON()
}

// EspaciosAcademicosProyectoPeriodo ...
// @Title EspaciosAcademicosProyectoPeriodo
// @Description Lista espacios academicos por anio, periodo y proyecto curricular
// @Param	anio		query 	string	true		"Anio de consulta"
// @Param	periodo		query 	string	true		"Periodo academico"
// @Param	proyecto	query 	string	false		"Proyecto curricular"
// @Param	documento-coordinador	query 	string	false		"Documento del coordinador"
// @Success 200 {}
// @Failure 400 the request contains an incorrect parameter
// @Failure 404 no record exist
// @router /proyecto-periodo [get]
func (c *EspacioAcademicoController) EspaciosAcademicosProyectoPeriodo() {
	defer errorhandler.HandlePanic(&c.Controller)

	anio := c.GetString("anio")
	periodo := c.GetString("periodo")
	proyecto := c.GetString("proyecto")
	documentoCoordinador := c.GetString("documento-coordinador")
	if documentoCoordinador == "" {
		documentoCoordinador = c.GetString("documento_coordinador")
	}

	if anio == "" || periodo == "" || (proyecto == "" && documentoCoordinador == "") {
		logs.Error(anio, periodo, proyecto, documentoCoordinador)
		c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parametro(s) con valores no validos")
		c.Ctx.Output.SetStatus(400)
	} else {
		resultado := services.ListaEspaciosAcademicosProyectoPeriodo(anio, periodo, proyecto, documentoCoordinador)
		c.Data["json"] = resultado
		c.Ctx.Output.SetStatus(resultado.Status)
	}

	c.ServeJSON()
}

// GruposEspacioPeriodo ...
// @Title GruposEspacioPeriodo
// @Description Lista grupos por anio, periodo y espacio academico
// @Param	anio		query 	string	true		"Anio de consulta"
// @Param	periodo		query 	string	true		"Periodo academico"
// @Param	espacio		query 	string	true		"Espacio academico"
// @Param	espacio_academico_id	query 	string	false		"Id unico del espacio academico"
// @Success 200 {}
// @Failure 400 the request contains an incorrect parameter
// @Failure 404 no record exist
// @router /grupos-periodo [get]
func (c *EspacioAcademicoController) GruposEspacioPeriodo() {
	defer errorhandler.HandlePanic(&c.Controller)

	anio := c.GetString("anio")
	periodo := c.GetString("periodo")
	espacio := c.GetString("espacio")
	if espacio == "" {
		espacio = c.GetString("espacio_academico_id")
	}

	if anio == "" || periodo == "" || espacio == "" {
		logs.Error(anio, periodo, espacio)
		c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parametro(s) con valores no validos")
		c.Ctx.Output.SetStatus(400)
	} else {
		resultado := services.ListaGruposEspacioPeriodo(anio, periodo, espacio)
		c.Data["json"] = resultado
		c.Ctx.Output.SetStatus(resultado.Status)
	}

	c.ServeJSON()
}
