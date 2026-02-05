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
}

// GrupoEspacioAcademico ...
// @Title GrupoEspacioAcademico
// @Description  Lista los grupos de un espacios académico padre por vigencia
// @Param	padre		query 	string	true		"Id del espacio académico padre"
// @Param	vigencia	query 	string	true		"Vigencia del espacio académico"
// @Success 200 {}
// @Failure 404 he request contains an incorrect parameter or no record exist
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
// @Failure 404 he request contains an incorrect parameter or no record exist
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
