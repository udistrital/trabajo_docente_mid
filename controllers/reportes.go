package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/trabajo_docente_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	requestmanager "github.com/udistrital/utils_oas/requestresponse"
)

// ReportesController operations for Reportes
type ReportesController struct {
	beego.Controller
}

// URLMapping ...
func (c *ReportesController) URLMapping() {
	c.Mapping("ReporteCargaLectiva", c.ReporteCargaLectiva)
	c.Mapping("ReporteVerificacionCumplimientoPTD", c.ReporteVerificacionCumplimientoPTD)
}

// ReporteCargaLectiva ...
// @Title ReporteCargaLectiva
// @Description Generar reporte excel de carga lectiva para docente
// @Param 	docente 		query 	int true	 "Id de docente"
// @Param 	vinculacion 	query	int true	 "Id vinculacion"
// @Param 	periodo 		query	int true	 "Id periodo academico"
// @Param 	carga 			query	string true	 "Tipo carga: C) Carga lectiva, A) Actividades"
// @Success 200 Report Generation successful
// @Failure 400 The request contains an incorrect data type or an invalid parameter
// @Failure 404 he request contains an incorrect parameter or no record exist
// @router /plan-trabajo-docente [get]
func (c *ReportesController) ReporteCargaLectiva() {
	defer errorhandler.HandlePanic(&c.Controller)

	docente, errdoc := c.GetInt64("docente")
	vinculacion, errvin := c.GetInt64("vinculacion")
	periodo, errper := c.GetInt64("periodo")
	carga := c.GetString("carga")

	if errdoc != nil || errvin != nil || errper != nil || carga == "" {
		logs.Error(errdoc, errvin, errper, carga)
		c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) no válido(s) o faltante(s)")
		c.Ctx.Output.SetStatus(400)
	} else if docente <= 0 || vinculacion <= 0 || periodo <= 0 || carga != "C" && carga != "A" && carga != "CA" {
		logs.Error(docente, vinculacion, periodo, carga)
		c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) con valores no válidos")
		c.Ctx.Output.SetStatus(400)
	} else {
		resultado := services.RepCargaLectiva(docente, vinculacion, periodo, carga)
		c.Data["json"] = resultado
		c.Ctx.Output.SetStatus(resultado.Status)
	}

	c.ServeJSON()
}

// ReporteVerificacionCumplimientoPTD ...
// @Title ReporteVerificacionCumplimientoPTD
// @Description Generar reporte excel de verificacion cumplimiento PTD
// @Param 	vigencia 		query 	int true	 "Id periodo academico"
// @Param 	proyecto 		query 	int false	 "Id proyecto academico"
// @Success 200 Report Generation successful
// @Failure 400 The request contains an incorrect data type or an invalid parameter
// @Failure 404 he request contains an incorrect parameter or no record exist
// @router /verificacion-cumplimiento-ptd [get]
func (c *ReportesController) ReporteVerificacionCumplimientoPTD() {
	defer errorhandler.HandlePanic(&c.Controller)

	vigencia, errvig := c.GetInt64("vigencia")
	proyecto := c.GetString("proyecto") // tomado como opcional

	if errvig != nil {
		logs.Error(errvig)
		c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) no válido(s) o faltante(s)")
		c.Ctx.Output.SetStatus(400)
	} else if vigencia <= 0 {
		logs.Error(vigencia, proyecto)
		c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) con valores no válidos")
	} else {
		resultado := services.RepCumplimiento(vigencia, proyecto)
		c.Data["json"] = resultado
		c.Ctx.Output.SetStatus(resultado.Status)
	}

	c.ServeJSON()
}
