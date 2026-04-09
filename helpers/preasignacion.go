package helpers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
)

func DesactivarPreAsignacion(preAsignacionId string) (map[string]interface{}, error) {
	var preAsignacion map[string]interface{}
	urlPreAsignacion := beego.AppConfig.String("PlanTrabajoDocenteService") + "pre_asignacion/" + preAsignacionId
	if err := request.GetJson(urlPreAsignacion, &preAsignacion); err != nil {
		return nil, fmt.Errorf("error en el servicio plan docente: %v", err)
	}

	preAsignacion = preAsignacion["Data"].(map[string]interface{})
	preAsignacion["activo"] = false

	urlPreasignacionPut := beego.AppConfig.String("PlanTrabajoDocenteService") + "pre_asignacion/" + preAsignacion["_id"].(string)
	var preAsignacionPut map[string]interface{}
	if err := request.SendJson(urlPreasignacionPut, "PUT", &preAsignacionPut, preAsignacion); err != nil {
		return nil, fmt.Errorf("error en el servicio plan docente: %v", err)
	}

	return preAsignacionPut["Data"].(map[string]interface{}), nil
}

// CambiarEstadoDePlanDocente
//
// EstadoACambiar: Es el codigo de abreviacion del estado_plan
func CambiarEstadoDePlanDocente(planDocenteId, estadoACambiar string) (map[string]interface{}, error) {
	var planDocente map[string]interface{}
	urlPlanDocente := beego.AppConfig.String("PlanTrabajoDocenteService") + "plan_docente/" + planDocenteId
	if err := request.GetJson(urlPlanDocente, &planDocente); err != nil {
		return nil, fmt.Errorf("error en el servicio plan docente: %v", err)
	}

	var estadoPlan map[string]interface{}
	urlEstadoPlan := beego.AppConfig.String("PlanTrabajoDocenteService") + "estado_plan?query=codigo_abreviacion:" + estadoACambiar

	if err := request.GetJson(urlEstadoPlan, &estadoPlan); err != nil {
		return nil, fmt.Errorf("error en el servicio plan docente assas: %v", err)
	}

	//el nuevo estado plan id que tendra el plan docente, segun el parametro estadoACambiar
	estadoPlanId := estadoPlan["Data"].([]interface{})[0].(map[string]interface{})["_id"].(string)

	planDocente = planDocente["Data"].(map[string]interface{})
	planDocente["estado_plan_id"] = estadoPlanId

	urlPlanDocentePut := beego.AppConfig.String("PlanTrabajoDocenteService") + "plan_docente/" + planDocenteId
	var planDocentePut map[string]interface{}
	if err := request.SendJson(urlPlanDocentePut, "PUT", &planDocentePut, planDocente); err != nil {
		return nil, fmt.Errorf("error en el servicio plan docente: %v", err)
	}

	return planDocentePut["Data"].(map[string]interface{}), nil
}

func DesasignarDocenteDeEspacioAcademico(espacioAcademicoId, docenteId string) (map[string]interface{}, error) {
	var espacioAcademico map[string]interface{}
	urlPlanDocente := beego.AppConfig.String("EspaciosAcademicosService") + "espacio-academico?query=_id:" + espacioAcademicoId
	if err := request.GetJson(urlPlanDocente, &espacioAcademico); err != nil {
		return nil, fmt.Errorf("error en el servicio espacios academicos: %v", err)
	}

	data, ok := espacioAcademico["Data"].([]interface{})
	if !ok || len(data) == 0 {
		logs.Warn("no se encontro espacio academico para desasignar docente, espacio_academico_id=%s", espacioAcademicoId)
		return nil, nil
	}

	espacioAcademico, ok = data[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("respuesta invalida del servicio espacios academicos para _id %s", espacioAcademicoId)
	}

	if espacioAcademico["espacio_modular"] == true {
		listaModularDocentes := espacioAcademico["lista_modular_docentes"].([]interface{})
		nuevaListaModularDocentes := []interface{}{}

		for _, id := range listaModularDocentes {
			idStr := fmt.Sprintf("%v", id)
			if idStr != docenteId {
				nuevaListaModularDocentes = append(nuevaListaModularDocentes, id)
			}
		}

		espacioAcademico["lista_modular_docentes"] = nuevaListaModularDocentes
	} else {
		espacioAcademico["docente_id"] = 0
	}

	urlEspacioAcademicoPut := beego.AppConfig.String("EspaciosAcademicosService") + "espacio-academico/" + espacioAcademicoId
	var espacioAcademicoPut map[string]interface{}
	if err := request.SendJson(urlEspacioAcademicoPut, "PUT", &espacioAcademicoPut, espacioAcademico); err != nil {
		return nil, fmt.Errorf("error en el servicio plan docente: %v", err)
	}

	return espacioAcademicoPut["Data"].(map[string]interface{}), nil
}
