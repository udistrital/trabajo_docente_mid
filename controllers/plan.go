package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/trabajo_docente_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	requestmanager "github.com/udistrital/utils_oas/requestresponse"
)

// PlanController operations for Asignacion
type PlanController struct {
	beego.Controller
}

// URLMapping ...
func (c *PlanController) URLMapping() {
	c.Mapping("DefinePlanTrabajoDocente", c.DefinePlanTrabajoDocente)
	c.Mapping("PlanTrabajoDocenteAsignacion", c.PlanTrabajoDocenteAsignacion)
	c.Mapping("CopiarPlanTrabajoDocente", c.CopiarPlanTrabajoDocente)
	c.Mapping("PlanPreaprobado", c.PlanPreaprobado)
}

// DefinePlanTrabajoDocente ...
// @Title DefinePlanTrabajoDocente
// @Description Actualiza la información de los planes de trabajo
// @Param   body        body    {}  true        "body plan a actualizar"
// @Success 200 {}
// @Failure 404 the request contains an incorrect parameter or no record exist
// @router / [put]
func (c *PlanController) DefinePlanTrabajoDocente() {
	defer errorhandler.HandlePanic(&c.Controller)

	var body map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	if err != nil {
		logs.Error(err)
		c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) no válido(s) o faltante(s)")
		c.Ctx.Output.SetStatus(400)
	}

	params := []string{"carga_plan", "plan_docente", "descartar"}
	errParam := false
	for _, param := range params {
		if _, ok := body[param]; !ok {
			errParam = true
			logs.Error("No existe el parametro %s", param)
			break
		}
	}
	if errParam {
		c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) no válido(s) o faltante(s)")
		c.Ctx.Output.SetStatus(400)
	} else {
		resultado := services.DefinePTD(body)
		c.Data["json"] = resultado
		c.Ctx.Output.SetStatus(resultado.Status)
	}

	c.ServeJSON()
}

// PlanTrabajoDocenteAsignacion ...
// @Title PlanTrabajoDocenteAsignacion
// @Description Traer la información de las asignaciones de un docente en la vigencia determinada
// @Param	docente		query 	int	true		"Id docente"
// @Param	vigencia	query 	int	true		"Vigencia de las asignaciones"
// @Param	vinculacion	query 	int	true		"Id vinculacion"
// @Success 200 {}
// @Failure 400 The request contains an incorrect data type or an invalid parameter
// @Failure 404 he request contains an incorrect parameter or no record exist
// @router / [get]
func (c *PlanController) PlanTrabajoDocenteAsignacion() {
	defer errorhandler.HandlePanic(&c.Controller) // not mine

	docente, errdoc := c.GetInt64("docente")
	vigencia, errvig := c.GetInt64("vigencia")
	vinculacion, errvin := c.GetInt64("vinculacion")

	if errdoc != nil || errvig != nil || errvin != nil {
		logs.Error(errdoc, errvig, errvin)
		c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) no válido(s) o faltante(s)")
		c.Ctx.Output.SetStatus(400)
	} else if docente <= 0 || vigencia <= 0 || vinculacion <= 0 {
		logs.Error(docente, vigencia, vinculacion)
		c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) con valores no válidos")
		c.Ctx.Output.SetStatus(400)
	} else {
		resultado := services.PlanTrabajoDocente(docente, vigencia, vinculacion)
		c.Data["json"] = resultado
		c.Ctx.Output.SetStatus(resultado.Status)
	}

	c.ServeJSON()
}

// CopiarPlanTrabajoDocente ...
// @Title CopiarPlanTrabajoDocente
// @Description Copia de plan de trabajo docente de una vigencia anterior
// @Param	docente				query 	int	true		"Id docente"
// @Param	vigenciaAnterior	query 	int	true		"Vigencia de la cual se pretende hacer copia"
// @Param	vigencia			query 	int	true		"Vigencia actual para encontrar diferencias"
// @Param	vinculacion			query 	int	true		"Id vinculacion"
// @Param	carga				query 	int	true		"Carga Lectiva 1, Actividades 2"
// @Success 200 {}
// @Failure 400 The request contains an incorrect data type or an invalid parameter
// @Failure 404 he request contains an incorrect parameter or no record exist
// @router /copiar [get]
func (c *PlanController) CopiarPlanTrabajoDocente() {
	defer errorhandler.HandlePanic(&c.Controller)

	docente, errdoc := c.GetInt64("docente")
	vigenciaAnterior, errvigA := c.GetInt64("vigenciaAnterior")
	vigencia, errvig := c.GetInt64("vigencia")
	vinculacion, errvin := c.GetInt64("vinculacion")
	carga, errcar := c.GetInt8("carga")

	if errdoc != nil || errvigA != nil || errvig != nil || errvin != nil || errcar != nil {
		logs.Error(errdoc, errvigA, errvig, errvin, errcar)
		c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) no válido(s) o faltante(s)")
		c.Ctx.Output.SetStatus(400)
	} else if docente <= 0 || vigenciaAnterior <= 0 || vigencia <= 0 || vinculacion <= 0 || carga != 1 && carga != 2 {
		logs.Error(docente, vigenciaAnterior, vigencia, vinculacion, carga)
		c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) con valores no válidos")
		c.Ctx.Output.SetStatus(400)
	} else {
		resultado := services.CopiarPTD(docente, vigenciaAnterior, vigencia, vinculacion, carga)
		c.Data["json"] = resultado
		c.Ctx.Output.SetStatus(resultado.Status)
	}

	c.ServeJSON()
}

// PlanPreaprobado ...
// @Title PlanPreaprobado
// @Description Listar planes que han sido aprobados en asignar ptd
// @Param	vigencia		query 	int	true	"Id periodo"
// @Param	proyecto		query 	int	true	"Id proyecto"
// @Success 200 {}
// @Failure 404 he request contains an incorrect parameter or no record exist
// @router /preaprobado [get]
func (c *PlanController) PlanPreaprobado() {
	defer errorhandler.HandlePanic(&c.Controller)

	vigencia, errvig := c.GetInt64("vigencia")
	proyecto, errpro := c.GetInt64("proyecto")

	if errvig != nil || errpro != nil {
		logs.Error(errvig, errpro)
		c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) no válido(s) o faltante(s)")
		c.Ctx.Output.SetStatus(400)
	} else if vigencia <= 0 || proyecto <= 0 {
		logs.Error(vigencia, proyecto)
		c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) con valores no válidos")
		c.Ctx.Output.SetStatus(400)
	} else {
		resultado := services.ListaPlanPreaprobado(vigencia, proyecto)
		c.Data["json"] = resultado
		c.Ctx.Output.SetStatus(resultado.Status)
	}

	c.ServeJSON()
}
