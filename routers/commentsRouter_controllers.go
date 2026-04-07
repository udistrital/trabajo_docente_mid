package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:AsignacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:AsignacionController"],
		beego.ControllerComments{
			Method:           "Asignacion",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:AsignacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:AsignacionController"],
		beego.ControllerComments{
			Method:           "AsignacionDocente",
			Router:           "/docente",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:DocenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:DocenteController"],
		beego.ControllerComments{
			Method:           "DocumentoDocenteVinculacion",
			Router:           "/documento",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:DocenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:DocenteController"],
		beego.ControllerComments{
			Method:           "NombreDocenteVinculacion",
			Router:           "/nombre",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:EspacioAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:EspacioAcademicoController"],
		beego.ControllerComments{
			Method:           "GrupoEspacioAcademico",
			Router:           "/grupo",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:EspacioAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:EspacioAcademicoController"],
		beego.ControllerComments{
			Method:           "GrupoEspacioAcademicoPadre",
			Router:           "/grupo-padre",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:EspacioFisicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:EspacioFisicoController"],
		beego.ControllerComments{
			Method:           "EspacioFisicoDependencia",
			Router:           "/dependencia",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:EspacioFisicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:EspacioFisicoController"],
		beego.ControllerComments{
			Method:           "DisponibilidadEspacioFisico",
			Router:           "/disponibilidad",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:PlanController"] = append(beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:PlanController"],
		beego.ControllerComments{
			Method:           "DefinePlanTrabajoDocente",
			Router:           "/",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:PlanController"] = append(beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:PlanController"],
		beego.ControllerComments{
			Method:           "PlanTrabajoDocenteAsignacion",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:PlanController"] = append(beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:PlanController"],
		beego.ControllerComments{
			Method:           "CopiarPlanTrabajoDocente",
			Router:           "/copiar",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:PlanController"] = append(beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:PlanController"],
		beego.ControllerComments{
			Method:           "PlanPreaprobado",
			Router:           "/preaprobado",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:PreasignacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:PreasignacionController"],
		beego.ControllerComments{
			Method:           "Preasignacion",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:PreasignacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:PreasignacionController"],
		beego.ControllerComments{
			Method:           "DeletePreasignacion",
			Router:           "/:preasignacion_id",
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:PreasignacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:PreasignacionController"],
		beego.ControllerComments{
			Method:           "Aprobar",
			Router:           "/aprobar",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:PreasignacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:PreasignacionController"],
		beego.ControllerComments{
			Method:           "PreasignacionDocente",
			Router:           "/docente",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:ReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:ReportesController"],
		beego.ControllerComments{
			Method:           "ReporteCargaLectiva",
			Router:           "/plan-trabajo-docente",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:ReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/trabajo_docente_mid/controllers:ReportesController"],
		beego.ControllerComments{
			Method:           "ReporteVerificacionCumplimientoPTD",
			Router:           "/verificacion-cumplimiento-ptd",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
