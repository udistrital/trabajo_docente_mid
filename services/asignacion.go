package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/trabajo_docente_mid/utils"
	request "github.com/udistrital/utils_oas/request"
	requestmanager "github.com/udistrital/utils_oas/requestresponse"
)

// Asignacion ...
func ListaAsignacion(vigencia string) requestmanager.APIResponse {
	var resPreasignaciones map[string]interface{}

	if errPreasignacion := request.GetJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"pre_asignacion?query=activo:true,aprobacion_docente:true,aprobacion_proyecto:true,"+
		"periodo_id:"+vigencia+"&fields=docente_id,tipo_vinculacion_id,plan_docente_id,periodo_id", &resPreasignaciones); errPreasignacion == nil {
		if fmt.Sprintf("%v", resPreasignaciones["Data"]) != "[]" {
			response := consultarDetalleAsignacion(resPreasignaciones["Data"].([]interface{}), false)
			return requestmanager.APIResponseDTO(true, 200, response)
			/* c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": response} */
		} else {
			return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de preasignaciones")
			/* c.Ctx.Output.SetStatus(404)
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de preasignaciones"} */
		}
	} else {
		logs.Error(errPreasignacion)
		return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de preasignaciones")
		/* c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de preasignaciones"} */
	}
}

func consultarDetalleAsignacion(asignaciones []interface{}, forTeacher bool) []map[string]interface{} {
	memEstados := map[string]interface{}{}
	memPeriodo := map[string]interface{}{}
	memDocente := map[string]interface{}{}
	memDocumento := map[string]interface{}{}
	memVinculacion := map[string]interface{}{}
	response := []map[string]interface{}{}

	var resPeriodo map[string]interface{}
	var resDocente map[string]interface{}
	var resDocumento []map[string]interface{}
	var resVinculacion map[string]interface{}
	var resEstado map[string]interface{}

	for _, asignacion := range asignaciones {
		if errDocente := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+asignacion.(map[string]interface{})["docente_id"].(string), &resDocente); errDocente == nil {
			memDocente[asignacion.(map[string]interface{})["docente_id"].(string)] = resDocente
			if errDocumento := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=TerceroId.Id:"+asignacion.(map[string]interface{})["docente_id"].(string)+"&fields=Numero", &resDocumento); errDocumento == nil {
				memDocumento[asignacion.(map[string]interface{})["docente_id"].(string)] = resDocumento[0]["Numero"]
			}
		}

		if errVinculacion := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro/"+asignacion.(map[string]interface{})["tipo_vinculacion_id"].(string), &resVinculacion); errVinculacion == nil {
			vinculacion := resVinculacion["Data"].(map[string]interface{})["Nombre"].(string)
			vinculacion = strings.Replace(vinculacion, "DOCENTE DE ", "", 1)
			vinculacion = strings.ToLower(vinculacion)
			memVinculacion[asignacion.(map[string]interface{})["tipo_vinculacion_id"].(string)] = strings.ToUpper(vinculacion[0:1]) + vinculacion[1:]
		}

		if memPeriodo[asignacion.(map[string]interface{})["periodo_id"].(string)] == nil {
			if errPeriodo := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"periodo/"+fmt.Sprintf("%v", asignacion.(map[string]interface{})["periodo_id"]), &resPeriodo); errPeriodo == nil {
				memPeriodo[asignacion.(map[string]interface{})["periodo_id"].(string)] = resPeriodo["Data"].(map[string]interface{})["Nombre"].(string)
			}
		}

		var resPlan map[string]interface{}
		var idDocumental interface{}
		var tieneObservaciones bool = false
		if memEstados[asignacion.(map[string]interface{})["plan_docente_id"].(string)] == nil {

			estadoPlan := "Sin definir"
			plan_id := ""
			if errPlan := request.GetJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"plan_docente/"+fmt.Sprintf("%v", asignacion.(map[string]interface{})["plan_docente_id"]), &resPlan); errPlan == nil {
				idEstado := resPlan["Data"].(map[string]interface{})["estado_plan_id"].(string)
				plan_id = resPlan["Data"].(map[string]interface{})["_id"].(string)
				if idEstado == "Sin definir" {
					memEstados[asignacion.(map[string]interface{})["plan_docente_id"].(string)] = resPlan["Data"].(map[string]interface{})["estado_plan_id"].(string)
				} else {
					if errEstado := request.GetJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"estado_plan/"+idEstado, &resEstado); errEstado == nil {
						memEstados[asignacion.(map[string]interface{})["plan_docente_id"].(string)] = resEstado["Data"].(map[string]interface{})["nombre"].(string)
						estadoPlan = resEstado["Data"].(map[string]interface{})["codigo_abreviacion"].(string)
					}
				}
				if resPlan["Data"].(map[string]interface{})["soporte_documental"] != nil {
					idDocumental = resPlan["Data"].(map[string]interface{})["soporte_documental"]
				}
				// Verificar si el plan tiene observaciones
				if resumen, ok := resPlan["Data"].(map[string]interface{})["resumen"].(string); ok {
					tieneObservaciones = verificarSiTieneObservaciones(resumen)
				}
			}

			desactivarEnviar := false
			tipoGestion := "ver"

			if forTeacher {
				switch estadoPlan {
				case "ENV_COO":
					tipoGestion = "editar"
					desactivarEnviar = false
				case "N_APR":
					tipoGestion = "editar"
					desactivarEnviar = false
				default:
					tipoGestion = "ver"
					desactivarEnviar = true
				}
			} else {
				tipoGestion = "editar"
				switch estadoPlan {
				case "Sin definir":
					desactivarEnviar = true
				case "ENV_COO":
					desactivarEnviar = true
				case "ENV_DOC":
					desactivarEnviar = true
				case "PAPR":
					desactivarEnviar = true
				case "APR":
					desactivarEnviar = true
				default:
					tipoGestion = "editar"
					desactivarEnviar = false
				}
			}

			response = append(response, map[string]interface{}{
				"plan_docente_id":     plan_id,
				"id":                  asignacion.(map[string]interface{})["_id"],
				"docente_id":          asignacion.(map[string]interface{})["docente_id"].(string),
				"docente":             utils.Capitalize(memDocente[asignacion.(map[string]interface{})["docente_id"].(string)].(map[string]interface{})["NombreCompleto"].(string)),
				"tipo_vinculacion_id": asignacion.(map[string]interface{})["tipo_vinculacion_id"].(string),
				"tipo_vinculacion":    memVinculacion[asignacion.(map[string]interface{})["tipo_vinculacion_id"].(string)],
				"identificacion":      memDocumento[asignacion.(map[string]interface{})["docente_id"].(string)],
				"periodo_academico":   memPeriodo[asignacion.(map[string]interface{})["periodo_id"].(string)],
				"periodo_id":          asignacion.(map[string]interface{})["periodo_id"].(string),
				"estado":              memEstados[asignacion.(map[string]interface{})["plan_docente_id"].(string)],
				"codigo_estado":       estadoPlan,
				"tiene_observaciones": tieneObservaciones,
				"soporte_documental":  map[string]interface{}{"value": idDocumental, "type": "ver", "disabled": idDocumental == nil || estadoPlan != "APR"},
				"enviar":              map[string]interface{}{"value": nil, "type": "enviar", "disabled": desactivarEnviar},
				"gestion":             map[string]interface{}{"value": nil, "type": tipoGestion, "disabled": false}})
		}

	}
	return response
}

// AsignacionDocente ...
func ListaAsignacionDocente(docente, vigencia string) requestmanager.APIResponse {
	var resPreasignaciones map[string]interface{}

	if errPreasignacion := request.GetJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"pre_asignacion?query=activo:true,aprobacion_docente:true,aprobacion_proyecto:true,docente_id:"+docente+",periodo_id:"+vigencia+"&fields=docente_id,tipo_vinculacion_id,plan_docente_id,periodo_id", &resPreasignaciones); errPreasignacion == nil {
		if fmt.Sprintf("%v", resPreasignaciones["Data"]) != "[]" {
			response := consultarDetalleAsignacion(resPreasignaciones["Data"].([]interface{}), true)
			return requestmanager.APIResponseDTO(true, 200, response)
			/* c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": response} */
		} else {
			return requestmanager.APIResponseDTO(false, 404, "No se encontraron registros de preasignaciones")
			/* c.Ctx.Output.SetStatus(404)
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de preasignaciones"} */
		}
	} else {
		logs.Error(errPreasignacion)
		return requestmanager.APIResponseDTO(false, 404, "No se encontraron registros de preasignaciones")
		/* c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de preasignaciones"} */
	}
}

// verificarSiTieneObservaciones verifica si el campo de observaciones tiene contenido
func verificarSiTieneObservaciones(resumenJSON string) bool {
	if resumenJSON == "" {
		return false
	}

	var resumen map[string]interface{}
	if err := json.Unmarshal([]byte(resumenJSON), &resumen); err != nil {
		return false
	}

	observacion, exists := resumen["observacion"]
	if !exists {
		return false
	}

	// Verificar si observacion no está vacía
	if obsTrimmed, ok := observacion.(string); ok {
		return strings.TrimSpace(obsTrimmed) != ""
	}

	return false
}
