package services

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/trabajo_docente_mid/helpers"
	"github.com/udistrital/trabajo_docente_mid/models"
	"github.com/udistrital/trabajo_docente_mid/utils"
	request "github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
	requestmanager "github.com/udistrital/utils_oas/requestresponse"
)

func obtenerIdRelacionado(valor interface{}) string {
	if valor == nil {
		return ""
	}

	switch v := valor.(type) {
	case string:
		return v
	case map[string]interface{}:
		if id, ok := v["_id"]; ok && id != nil {
			return fmt.Sprintf("%v", id)
		}
	}

	return ""
}

// Preasignacion ...
func ListaPreasignacion(vigencia string) requestmanager.APIResponse {
	var resPreasignaciones map[string]interface{}

	if errPreasignacion := request.GetJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"pre_asignacion?query=activo:true,periodo_id:"+vigencia, &resPreasignaciones); errPreasignacion == nil {
		if fmt.Sprintf("%v", resPreasignaciones["Data"]) != "[]" {
			response := consultarDetallePreasignacion(resPreasignaciones["Data"].([]interface{}))

			for _, preasignacion := range response {
				preasignacion["aprobacion_docente"].(map[string]interface{})["disabled"] = true
				preasignacion["aprobacion_proyecto"].(map[string]interface{})["disabled"] = true
			}
			return requestmanager.APIResponseDTO(true, 200, response)
		} else {
			return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de preasignaciones")
		}
	} else {
		logs.Error(errPreasignacion)
		return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de preasignaciones")
	}
}

func consultarDetallePreasignacion(preasignaciones []interface{}) []map[string]interface{} {
	memEspacios := map[string]interface{}{}
	memPeriodo := map[string]interface{}{}
	memDocente := map[string]interface{}{}
	response := []map[string]interface{}{}
	var resPeriodo map[string]interface{}
	var resDocente map[string]interface{}

	for _, preasignacion := range preasignaciones {
		if errDocente := request.GetJson(beego.AppConfig.String("TercerosService")+"tercero/"+preasignacion.(map[string]interface{})["docente_id"].(string), &resDocente); errDocente == nil {
			memDocente[preasignacion.(map[string]interface{})["docente_id"].(string)] = resDocente
		}

		if memEspacios[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)] == nil {
			// Consultar detalle del curso usando el nuevo endpoint que consolida información de Oracle
			cursoDetalle := DetalleCursoId(preasignacion.(map[string]interface{})["espacio_academico_id"].(string))
			if cursoDetalle.Success && cursoDetalle.Data != nil {
				cursoData := cursoDetalle.Data.(map[string]interface{})
				memEspacios[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)] = map[string]interface{}{
					"espacio_academico":         cursoData["EspacioAcademico"].(string),
					"grupo":                     cursoData["grupo"].(string),
					"codigo":                    cursoData["CodigoEspacioAcademico"].(string),
					"proyecto":                  cursoData["ProyectoAcademico"].(string),
					"codigo_proyecto_academico": cursoData["CodigoProyectoAcademico"].(string),
					"nivel":                     cursoData["Nivel"].(string),
				}
			}
		}

		if memPeriodo[preasignacion.(map[string]interface{})["periodo_id"].(string)] == nil {
			if errPeriodo := request.GetJson(beego.AppConfig.String("ParametroService")+"periodo/"+fmt.Sprintf("%v", preasignacion.(map[string]interface{})["periodo_id"]), &resPeriodo); errPeriodo == nil {
				memPeriodo[preasignacion.(map[string]interface{})["periodo_id"].(string)] = resPeriodo["Data"].(map[string]interface{})["Nombre"].(string)
			}
		}

		response = append(response, map[string]interface{}{
			"id":                        preasignacion.(map[string]interface{})["_id"],
			"docente_id":                preasignacion.(map[string]interface{})["docente_id"].(string),
			"docente":                   utils.Capitalize(memDocente[preasignacion.(map[string]interface{})["docente_id"].(string)].(map[string]interface{})["NombreCompleto"].(string)),
			"tipo_vinculacion_id":       preasignacion.(map[string]interface{})["tipo_vinculacion_id"].(string),
			"espacio_academico":         memEspacios[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)].(map[string]interface{})["espacio_academico"],
			"espacio_academico_id":      preasignacion.(map[string]interface{})["espacio_academico_id"].(string),
			"espacio_academico_padre":   memEspacios[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)].(map[string]interface{})["codigo"],
			"grupo":                     memEspacios[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)].(map[string]interface{})["grupo"],
			"proyecto":                  memEspacios[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)].(map[string]interface{})["proyecto"],
			"codigo_proyecto_academico": memEspacios[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)].(map[string]interface{})["codigo_proyecto_academico"],
			"nivel":                     memEspacios[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)].(map[string]interface{})["nivel"],
			"codigo":                    memEspacios[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)].(map[string]interface{})["codigo"],
			"periodo":                   memPeriodo[preasignacion.(map[string]interface{})["periodo_id"].(string)],
			"periodo_id":                preasignacion.(map[string]interface{})["periodo_id"].(string),
			"aprobacion_docente":        map[string]interface{}{"value": preasignacion.(map[string]interface{})["aprobacion_docente"].(bool), "disabled": false},
			"aprobacion_proyecto":       map[string]interface{}{"value": preasignacion.(map[string]interface{})["aprobacion_proyecto"].(bool), "disabled": false},
			"editar":                    map[string]interface{}{"value": nil, "type": "editar", "disabled": false},
			"enviar":                    map[string]interface{}{"value": nil, "type": "enviar", "disabled": preasignacion.(map[string]interface{})["aprobacion_proyecto"].(bool)},
			"borrar":                    map[string]interface{}{"value": nil, "type": "borrar", "disabled": false},
		})
	}
	return response
}

// PreasignacionDocente ...
func ListaPreasignacionDocente(docente, vigencia string) requestmanager.APIResponse {
	var resPreasignaciones map[string]interface{}

	if errPreasignacion := request.GetJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"pre_asignacion?query=aprobacion_proyecto:true,activo:true,periodo_id:"+vigencia+",docente_id:"+docente, &resPreasignaciones); errPreasignacion == nil {
		if fmt.Sprintf("%v", resPreasignaciones["Data"]) != "[]" {
			response := consultarDetallePreasignacion(resPreasignaciones["Data"].([]interface{}))

			for _, preasignacion := range response {
				preasignacion["aprobacion_proyecto"].(map[string]interface{})["disabled"] = true
				preasignacion["aprobacion_docente"].(map[string]interface{})["disabled"] = preasignacion["aprobacion_docente"].(map[string]interface{})["value"]
			}

			return requestmanager.APIResponseDTO(true, 200, response)
		} else {
			return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros para el docente")
		}
	} else {
		logs.Error(errPreasignacion)
		return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de docentes")
	}
}

// Aprobar ...
func DefinePreasignacion(body map[string]interface{}) requestmanager.APIResponse {
	var PreasignacionPut map[string]interface{}
	var EspacioPut map[string]interface{}
	resultado := []map[string]interface{}{}

	var preasignacionPut map[string]interface{}

	// Preasignaciones aceptadas
	if body["docente"].(bool) {
		preasignacionPut = map[string]interface{}{"aprobacion_docente": true}
	} else {
		preasignacionPut = map[string]interface{}{"aprobacion_proyecto": true}
	}

	for _, preasignacion := range body["preasignaciones"].([]interface{}) {
		if errAprobacion := request.SendJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"pre_asignacion/"+fmt.Sprintf("%v", preasignacion.(map[string]interface{})["Id"]), "PUT", &PreasignacionPut, preasignacionPut); errAprobacion == nil {
			// Actualización de espacio academico hijo con docente cuando es aprobado por el docente
			if body["docente"] == true {
				// Trae el espacio academico hijo para posterior actualización con el docente asigando
				var EspacioAcademicoHijo map[string]interface{}
				if errEspacios := request.GetJson(beego.AppConfig.String("EspaciosAcademicosService")+"espacio-academico/"+fmt.Sprintf("%v", PreasignacionPut["Data"].(map[string]interface{})["espacio_academico_id"]), &EspacioAcademicoHijo); errEspacios == nil {
					if espacioData, ok := EspacioAcademicoHijo["Data"].(map[string]interface{}); ok && espacioData != nil {
						EspacioAcademicoHijoPut := espacioData
						espacioAcademicoPadreID := obtenerIdRelacionado(EspacioAcademicoHijoPut["espacio_academico_padre"])
						estadoAprobacionID := obtenerIdRelacionado(EspacioAcademicoHijoPut["estado_aprobacion_id"])

						if esp_mod, ok := EspacioAcademicoHijoPut["espacio_modular"]; ok {
							if esp_mod.(bool) {

								resp, err := requestmanager.Get(beego.AppConfig.String("PlanTrabajoDocenteService")+
									fmt.Sprintf("pre_asignacion?query=activo:true,espacio_academico_id:%s,periodo_id:%s,aprobacion_docente:true,aprobacion_proyecto:true", PreasignacionPut["Data"].(map[string]interface{})["espacio_academico_id"], PreasignacionPut["Data"].(map[string]interface{})["periodo_id"]), requestmanager.ParseResponseFormato1)
								if err == nil {
									preasign_list := []models.PreAsignacion{}
									utils.ParseData(resp, &preasign_list)
									listDocents := []int{}
									for _, preasign := range preasign_list {
										id, _ := strconv.Atoi(preasign.Docente_id)
										listDocents = append(listDocents, id)
									}
									EspacioAcademicoHijoPut["lista_modular_docentes"] = listDocents
								}

								if espacioAcademicoPadreID != "" {
									EspacioAcademicoHijoPut["espacio_academico_padre"] = espacioAcademicoPadreID
								}
								if estadoAprobacionID != "" {
									EspacioAcademicoHijoPut["estado_aprobacion_id"] = estadoAprobacionID
								}
							} else {
								EspacioAcademicoHijoPut["docente_id"], _ = strconv.Atoi(PreasignacionPut["Data"].(map[string]interface{})["docente_id"].(string))
								if espacioAcademicoPadreID != "" {
									EspacioAcademicoHijoPut["espacio_academico_padre"] = espacioAcademicoPadreID
								}
								if estadoAprobacionID != "" {
									EspacioAcademicoHijoPut["estado_aprobacion_id"] = estadoAprobacionID
								}
							}
						} else {
							EspacioAcademicoHijoPut["docente_id"], _ = strconv.Atoi(PreasignacionPut["Data"].(map[string]interface{})["docente_id"].(string))
							if espacioAcademicoPadreID != "" {
								EspacioAcademicoHijoPut["espacio_academico_padre"] = espacioAcademicoPadreID
							}
							if estadoAprobacionID != "" {
								EspacioAcademicoHijoPut["estado_aprobacion_id"] = estadoAprobacionID
							}
						}
						// Put al espacio academico hijo con el docente asignado cuando se aprueba la preasignacion
						if errPutEspacio := request.SendJson(beego.AppConfig.String("EspaciosAcademicosService")+"espacio-academico/"+fmt.Sprintf("%v", PreasignacionPut["Data"].(map[string]interface{})["espacio_academico_id"]), "PUT", &EspacioPut, EspacioAcademicoHijoPut); errPutEspacio == nil {
						} else {
							resultado = append(resultado, map[string]interface{}{"Id": preasignacion.(map[string]interface{})["Id"], "actualizado": false})
						}

						//------------------------------------------Finalización Actualización------------------------------------------------------
					} else {
						logs.Warn("No fue posible actualizar espacio academico en aprobacion docente; Data ausente o invalido para preasignacion %v", preasignacion.(map[string]interface{})["Id"])
					}
				} else {
					logs.Error(errEspacios)
					logs.Warn("No fue posible consultar espacio academico hijo para preasignacion %v", preasignacion.(map[string]interface{})["Id"])
				}

			}

			if body["docente"].(bool) && PreasignacionPut["Data"].(map[string]interface{})["plan_docente_id"] == nil {
				var planDocenteGet map[string]interface{}
				if errGetPlan := request.GetJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"plan_docente?query=docente_id:"+fmt.Sprintf("%v", PreasignacionPut["Data"].(map[string]interface{})["docente_id"])+",periodo_id:"+fmt.Sprintf("%v", PreasignacionPut["Data"].(map[string]interface{})["periodo_id"])+",tipo_vinculacion_id:"+fmt.Sprintf("%v", PreasignacionPut["Data"].(map[string]interface{})["tipo_vinculacion_id"]), &planDocenteGet); errGetPlan == nil {
					if resultado != nil {
						if fmt.Sprintf("%v", planDocenteGet["Data"]) != "[]" {
							idPlanDocente := planDocenteGet["Data"].([]interface{})[0].(map[string]interface{})["_id"].(string)
							preasignacionPut = map[string]interface{}{"plan_docente_id": idPlanDocente}

							if errAprobacion := request.SendJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"pre_asignacion/"+fmt.Sprintf("%v", preasignacion.(map[string]interface{})["Id"]), "PUT", &PreasignacionPut, preasignacionPut); errAprobacion == nil {
								resultado = append(resultado, map[string]interface{}{"Id": PreasignacionPut["Data"].(map[string]interface{})["_id"], "actualizado": true, "plan_trabajo": true})
							}
						} else {
							planDocente := map[string]interface{}{
								"estado_plan_id":      "Sin definir",
								"docente_id":          PreasignacionPut["Data"].(map[string]interface{})["docente_id"],
								"tipo_vinculacion_id": PreasignacionPut["Data"].(map[string]interface{})["tipo_vinculacion_id"],
								"periodo_id":          PreasignacionPut["Data"].(map[string]interface{})["periodo_id"],
								"activo":              true,
							}

							var planDocentePost map[string]interface{}
							if errPlan := request.SendJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"plan_docente", "POST", &planDocentePost, planDocente); errPlan == nil {
								idPlanDocente := planDocentePost["Data"].(map[string]interface{})["_id"].(string)
								preasignacionPut = map[string]interface{}{"plan_docente_id": idPlanDocente}

								if errAprobacion := request.SendJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"pre_asignacion/"+fmt.Sprintf("%v", preasignacion.(map[string]interface{})["Id"]), "PUT", &PreasignacionPut, preasignacionPut); errAprobacion == nil {
									resultado = append(resultado, map[string]interface{}{"Id": PreasignacionPut["Data"].(map[string]interface{})["_id"], "actualizado": true, "plan_trabajo": true})
								}
							}
						}
					}
				}
			} else {
				resultado = append(resultado, map[string]interface{}{"Id": PreasignacionPut["Data"].(map[string]interface{})["_id"], "actualizado": true})
			}
		} else {
			resultado = append(resultado, map[string]interface{}{"Id": preasignacion.(map[string]interface{})["Id"], "actualizado": false})
		}
	}

	// Preasignaciones negadas
	if body["docente"].(bool) {
		preasignacionPut = map[string]interface{}{"aprobacion_docente": false}
	} else {
		preasignacionPut = map[string]interface{}{"aprobacion_proyecto": false}
	}

	for _, preasignacion := range body["no-preasignaciones"].([]interface{}) {
		if errAprobacion := request.SendJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"pre_asignacion/"+fmt.Sprintf("%v", preasignacion.(map[string]interface{})["Id"]), "PUT", &PreasignacionPut, preasignacionPut); errAprobacion == nil {
			resultado = append(resultado, map[string]interface{}{"Id": PreasignacionPut["Data"].(map[string]interface{})["_id"], "actualizado": true})
		} else {
			resultado = append(resultado, map[string]interface{}{"Id": preasignacion.(map[string]interface{})["Id"], "actualizado": false})
		}
	}

	return requestmanager.APIResponseDTO(true, 200, resultado)
}

func DeletePreasignacion(preAsignacionId string) requestmanager.APIResponse {
	urlPreasignacion := beego.AppConfig.String("PlanTrabajoDocenteService") + "pre_asignacion/" + preAsignacionId
	var preAsignacion map[string]interface{}
	if err := request.GetJson(urlPreasignacion, &preAsignacion); err != nil {
		return requestresponse.APIResponseDTO(false, 404, nil, "Error en el servicio plan docente"+err.Error())
	}

	espacioAcademicoId := preAsignacion["Data"].(map[string]interface{})["espacio_academico_id"].(string)
	docenteId := preAsignacion["Data"].(map[string]interface{})["docente_id"].(string)

	urlColocaciones := beego.AppConfig.String("PlanTrabajoDocenteService") + "carga_plan?query=activo:true,espacio_academico_id:" + espacioAcademicoId

	var colocacionesRes map[string]interface{}
	if err := request.GetJson(urlColocaciones, &colocacionesRes); err != nil {
		return requestresponse.APIResponseDTO(false, 404, nil, "Error en el servicio plan docente"+err.Error())
	}

	if len(colocacionesRes["Data"].([]interface{})) > 0 {
		return requestmanager.APIResponseDTO(false, 200, nil, "tiene colocaciones")
	}

	_, err := helpers.DesactivarPreAsignacion(preAsignacionId)
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, err.Error())
	}

	if planDocenteId, exists := preAsignacion["Data"].(map[string]interface{})["plan_docente_id"]; exists || planDocenteId != nil {
		planDocenteId := preAsignacion["Data"].(map[string]interface{})["plan_docente_id"].(string)
		_, err := helpers.CambiarEstadoDePlanDocente(planDocenteId, "DEF") //DEF es el codigo de abreviacion de Definido
		if err != nil {
			return requestresponse.APIResponseDTO(false, 500, nil, err.Error())
		}
	}

	_, err = helpers.DesasignarDocenteDeEspacioAcademico(espacioAcademicoId, docenteId)
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, err.Error())
	}

	return requestmanager.APIResponseDTO(true, 200, nil, "eliminado correctamente")
}
