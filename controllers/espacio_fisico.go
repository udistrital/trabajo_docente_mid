package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/trabajo_docente_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	requestmanager "github.com/udistrital/utils_oas/requestresponse"
)

// EspacioFisicoController operations for Asignacion
type EspacioFisicoController struct {
	beego.Controller
}

// URLMapping ...
func (c *EspacioFisicoController) URLMapping() {
	c.Mapping("EspacioFisicoDependencia", c.EspacioFisicoDependencia)
	c.Mapping("DisponibilidadEspacioFisico", c.DisponibilidadEspacioFisico)
}

// EspacioFisicoDependencia ...
// @Title EspacioFisicoDependencia
// @Description Lista opciones espacios físicos asignados a una dependencia
// @Param	dependencia		query	int	true	"Id dependencia"
// @Success 200 {}
// @Failure 404 he request contains an incorrect parameter or no record exist
// @router /dependencia [get]
func (c *EspacioFisicoController) EspacioFisicoDependencia() {
	defer errorhandler.HandlePanic(&c.Controller)

	dependencia, errdep := c.GetInt64("dependencia")
	if errdep != nil {
		logs.Error(errdep)
		c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) no válido(s) o faltante(s)")
		c.Ctx.Output.SetStatus(400)
	} else if dependencia <= 0 {
		logs.Error(dependencia)
		c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) con valores no válidos")
		c.Ctx.Output.SetStatus(400)
	} else {
		resultado := services.ArbolEspaciosFisicosDependencia(dependencia)
		c.Data["json"] = resultado
		c.Ctx.Output.SetStatus(resultado.Status)
	}

	c.ServeJSON()
}

// DisponibilidadEspacioFisico ...
// @Title DisponibilidadEspacioFisico
// @Description Consulta la disponibilidad de un espacio fisico
// @Param	salon 		query 	string	true		"Salon de las asignaciones"
// @Param	vigencia 	query 	string	true		"Vigencia de las asignaciones"
// @Param	plan 		query 	string	true		"Id del plan de trabajo"
// @Success 200 {}
// @Failure 404 he request contains an incorrect parameter or no record exist
// @router /disponibilidad [get]
func (c *EspacioFisicoController) DisponibilidadEspacioFisico() {
	defer errorhandler.HandlePanic(&c.Controller)

	salon := c.GetString("salon")
	vigencia := c.GetString("vigencia")
	plan := c.GetString("plan")

	if salon == "" || vigencia == "" || plan == "" {
		logs.Error(salon, vigencia, plan)
		c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) con valores no válidos")
		c.Ctx.Output.SetStatus(400)
	} else {
		resultado := services.OcupacionEspacioFisico(salon, vigencia, plan)
		c.Data["json"] = resultado
		c.Ctx.Output.SetStatus(resultado.Status)
	}

	c.ServeJSON()
}
