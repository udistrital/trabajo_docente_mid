// @APIVersion 1.0.0
// @Title SGA MID - Plan Trabajo Docente
// @Description Microservicio MID del SGA MID que complementa plan de trabajo docente
package routers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/trabajo_docente_mid/controllers"
	"github.com/udistrital/utils_oas/errorhandler"
)

func init() {

	beego.ErrorController(&errorhandler.ErrorHandlerController{})

	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/reporte",
			beego.NSInclude(
				&controllers.ReportesController{},
			),
		),
		beego.NSNamespace("/preasignacion",
			beego.NSInclude(
				&controllers.PreasignacionController{},
			),
		),
		beego.NSNamespace("/asignacion",
			beego.NSInclude(
				&controllers.AsignacionController{},
			),
		),
		beego.NSNamespace("/docente",
			beego.NSInclude(
				&controllers.DocenteController{},
			),
		),
		beego.NSNamespace("/espacio-fisico",
			beego.NSInclude(
				&controllers.EspacioFisicoController{},
			),
		),
		beego.NSNamespace("/espacio-academico",
			beego.NSInclude(
				&controllers.EspacioAcademicoController{},
			),
		),
		beego.NSNamespace("/plan",
			beego.NSInclude(
				&controllers.PlanController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
