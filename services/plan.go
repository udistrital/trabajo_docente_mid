package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/trabajo_docente_mid/models"
	"github.com/udistrital/trabajo_docente_mid/utils"
	request "github.com/udistrital/utils_oas/request"
	requestmanager "github.com/udistrital/utils_oas/requestresponse"
)

// DefinePlanTrabajoDocente ...
func DefinePTD(body map[string]interface{}) requestmanager.APIResponse {
	resultado := map[string]interface{}{}
	resultadoCargas := []map[string]interface{}{}
	var resPlan map[string]interface{}

	for _, carga := range body["carga_plan"].([]interface{}) {
		var resColocacion map[string]interface{}
		var resCarga map[string]interface{}

		espacioFisico := map[string]interface{}{
			"sede_id":     carga.(map[string]interface{})["sede_id"],
			"edificio_id": carga.(map[string]interface{})["edificio_id"],
			"salon_id":    carga.(map[string]interface{})["salon_id"],
		}
		resumenColocacion := map[string]interface{}{
			"colocacion":     carga.(map[string]interface{})["horario"],
			"espacio_fisico": espacioFisico,
		}
		resumenColocacionStr, errCol := json.Marshal(resumenColocacion)
		if errCol != nil {
			panic(errCol)
		}

		colocacionStr, errCol := json.Marshal(carga.(map[string]interface{})["horario"])
		if errCol != nil {
			panic(errCol)
		}
		bodyColocacion := map[string]interface{}{
			"Activo":                         true,
			"ColocacionEspacioAcademico":     utils.GetOrDefault(string(colocacionStr), "NA"),
			"EspacioAcademicoId":             utils.GetOrDefault(carga.(map[string]interface{})["espacio_academico_id"], "NA"),
			"EspacioFisicoId":                utils.GetOrDefault(carga.(map[string]interface{})["salon_id"], -1),
			"PeriodoId":                      carga.(map[string]interface{})["periodo_id"],
			"ResumenColocacionEspacioFisico": utils.GetOrDefault(string(resumenColocacionStr), "NA"),
		}
		bodyCarga := map[string]interface{}{
			"espacio_academico_id": utils.GetOrDefault(carga.(map[string]interface{})["espacio_academico_id"], "NA"),
			"actividad_id":         utils.GetOrDefault(carga.(map[string]interface{})["actividad_id"], "NA"),
			"id":                   utils.GetOrDefault(carga.(map[string]interface{})["id"], "NA"),
			"plan_docente_id":      utils.GetOrDefault(carga.(map[string]interface{})["plan_docente_id"], "NA"),
			"hora_inicio":          utils.GetOrDefault(carga.(map[string]interface{})["hora_inicio"], "NA"),
			"duracion":             utils.GetOrDefault(carga.(map[string]interface{})["duracion"], "NA"),
			"salon_id":             utils.GetOrDefault(carga.(map[string]interface{})["salon_id"], "NA"),
			"activo":               utils.GetOrDefault(carga.(map[string]interface{})["activo"], "NA"),
		}

		if carga.(map[string]interface{})["id"] == nil {
			fmt.Println("ruta creacion ", beego.AppConfig.String("HorarioService")+"colocacion-espacio-academico/")
			if errPostPlacement := request.SendJson("https://"+beego.AppConfig.String("HorarioService")+"colocacion-espacio-academico/",
				"POST", &resColocacion, bodyColocacion); errPostPlacement == nil {
				if resColocacion["Success"].(bool) {
					bodyCarga["colocacion_espacio_academico_id"] = resColocacion["Data"].(map[string]interface{})["_id"]
					fmt.Println("ruta creacion carga ", beego.AppConfig.String("PlanTrabajoDocenteService")+"carga_plan/")
					if errPostCarga := request.SendJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"carga_plan/",
						"POST", &resCarga, bodyCarga); errPostCarga == nil {
						if resCarga["Success"].(bool) {
							resultadoCargas = append(resultadoCargas, map[string]interface{}{"id": resCarga["Data"].(map[string]interface{})["_id"], "creado": true})
						} else {
							resultadoCargas = append(resultadoCargas, map[string]interface{}{"id": carga.(map[string]interface{})["espacio_academico_id"], "creado": false})
						}
					}
				} else {
					resultadoCargas = append(resultadoCargas, map[string]interface{}{"id": carga.(map[string]interface{})["espacio_academico_id"], "creado": false})
				}
			}
		} else if carga.(map[string]interface{})["id"] == "colocacionModuloHorario" {
			bodyCarga["colocacion_espacio_academico_id"] = carga.(map[string]interface{})["colocacion_id"]
			if errPostCarga := request.SendJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"carga_plan/",
				"POST", &resCarga, bodyCarga); errPostCarga == nil {
				if resCarga["Success"].(bool) {
					resultadoCargas = append(resultadoCargas, map[string]interface{}{"id": resCarga["Data"].(map[string]interface{})["_id"], "creado": true})
				} else {
					resultadoCargas = append(resultadoCargas, map[string]interface{}{"id": carga.(map[string]interface{})["espacio_academico_id"], "creado": false})
				}
			}
			if errPutColocacion := request.SendJson("https://"+beego.AppConfig.String("HorarioService")+"colocacion-espacio-academico/"+carga.(map[string]interface{})["colocacion_id"].(string),
				"PUT", &resColocacion, bodyColocacion); errPutColocacion == nil {
			}
		} else {
			var planTrabajoData map[string]interface{}
			if errPlanTrabajo := request.GetJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"carga_plan/"+carga.(map[string]interface{})["id"].(string), &planTrabajoData); errPlanTrabajo == nil {
				if planTrabajoData["Success"].(bool) {
					if colId, colExists := planTrabajoData["Data"].(map[string]interface{})["colocacion_espacio_academico_id"]; colExists {
						if errPutColocacion := request.SendJson("https://"+beego.AppConfig.String("HorarioService")+"colocacion-espacio-academico/"+colId.(string),
							"PUT", &resColocacion, bodyColocacion); errPutColocacion == nil {
							if resColocacion["Success"].(bool) {
								bodyCarga["colocacion_espacio_academico_id"] = resColocacion["Data"].(map[string]interface{})["_id"]
								if errPutCarga := request.SendJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"carga_plan/"+carga.(map[string]interface{})["id"].(string),
									"PUT", &resCarga, bodyCarga); errPutCarga == nil {
									if resCarga["Success"].(bool) {
										resultadoCargas = append(resultadoCargas, map[string]interface{}{"id": resCarga["Data"].(map[string]interface{})["_id"], "actualizado": true})
									} else {
										resultadoCargas = append(resultadoCargas, map[string]interface{}{"id": carga.(map[string]interface{})["espacio_academico_id"], "actualizado": false})
									}
								}
							} else {
								resultadoCargas = append(resultadoCargas, map[string]interface{}{"id": carga.(map[string]interface{})["espacio_academico_id"], "actualizado": false})
							}
						}
					} else {
						if errPutColocacion := request.SendJson("https://"+beego.AppConfig.String("HorarioService")+"colocacion-espacio-academico/",
							"POST", &resColocacion, bodyColocacion); errPutColocacion == nil {
							if resColocacion["Success"].(bool) {
								bodyCarga["colocacion_espacio_academico_id"] = resColocacion["Data"].(map[string]interface{})["_id"]
								if errPutCarga := request.SendJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"carga_plan/"+carga.(map[string]interface{})["id"].(string),
									"PUT", &resCarga, bodyCarga); errPutCarga == nil {
									if resCarga["Success"].(bool) {
										resultadoCargas = append(resultadoCargas, map[string]interface{}{"id": resCarga["Data"].(map[string]interface{})["_id"], "actualizado": true})
									} else {
										resultadoCargas = append(resultadoCargas, map[string]interface{}{"id": carga.(map[string]interface{})["espacio_academico_id"], "actualizado": false})
									}
								}
							} else {
								resultadoCargas = append(resultadoCargas, map[string]interface{}{"id": carga.(map[string]interface{})["espacio_academico_id"], "actualizado": false})
							}
						}
					}

				}
			}
		}
	}

	if body["plan_docente"].(map[string]interface{})["estado_plan"].(string) == "Sin definir" {
		var resEstado map[string]interface{}
		if errEstado := request.GetJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"estado_plan?query=codigo_abreviacion:DEF", &resEstado); errEstado == nil {
			body["plan_docente"].(map[string]interface{})["estado_plan_id"] = resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"]
		}
	} else {
		if _, err := utils.CheckIdString(body["plan_docente"].(map[string]interface{})["estado_plan"].(string)); err == nil {
			body["plan_docente"].(map[string]interface{})["estado_plan_id"] = body["plan_docente"].(map[string]interface{})["estado_plan"].(string)
		}
	}
	if errPutPlan := request.SendJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"plan_docente/"+body["plan_docente"].(map[string]interface{})["id"].(string), "PUT", &resPlan, body["plan_docente"]); errPutPlan == nil {
		if resPlan["Success"].(bool) {
			resultado["plan_actualizado"] = true
		} else {
			resultado["plan_actualizado"] = false
		}
	}

	for _, descartado := range body["descartar"].([]interface{}) {
		_, err := requestmanager.Delete(beego.AppConfig.String("PlanTrabajoDocenteService")+"carga_plan/"+descartado.(map[string]interface{})["id"].(string), requestmanager.ParseResponseFormato1)
		if err == nil {
			resultadoCargas = append(resultadoCargas, map[string]interface{}{"id": descartado.(map[string]interface{})["id"].(string), "desactivado": true})
		} else {
			resultadoCargas = append(resultadoCargas, map[string]interface{}{"id": descartado.(map[string]interface{})["id"].(string), "desactivado": false})
		}
	}

	resultado["carga_plan"] = resultadoCargas

	return requestmanager.APIResponseDTO(true, 200, resultado)
}

// PlanTrabajoDocenteAsignacion ...
func PlanTrabajoDocente(docente, vigencia, vinculacion int64) requestmanager.APIResponse {
	var resPlan map[string]interface{}
	if errPlan := request.GetJson(beego.AppConfig.String("PlanTrabajoDocenteService")+
		fmt.Sprintf("plan_docente?query=activo:true,docente_id:%d,periodo_id:%d&fields=tipo_vinculacion_id,soporte_documental,respuesta,resumen,docente_id,periodo_id,estado_plan_id", docente, vigencia), &resPlan); errPlan == nil {
		if fmt.Sprintf("%v", resPlan["Data"]) != "[]" {
			response := consultarDetallePlan(resPlan["Data"].([]interface{}), vinculacion)
			return requestmanager.APIResponseDTO(true, 200, response)
			/* c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": response} */
		} else {
			return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de preasignaciones")
			/* c.Ctx.Output.SetStatus(404)
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de preasignaciones"} */
		}
	} else {
		logs.Error(errPlan)
		return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de preasignaciones")
		/* c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de preasignaciones"} */
	}
}

func consultarDetallePlan(planes []interface{}, idVinculacion int64) map[string]interface{} {
	memDocente := map[string]interface{}{}
	memVinculacion := []map[string]interface{}{}
	memResumenes := []map[string]interface{}{}
	memEspacios := []interface{}{}
	memEspaciosDetalle := map[string]interface{}{}
	memCarga := []interface{}{}
	memPlanDocente := []string{}
	memEstadoPlan := []string{}
	memEstados := map[string]interface{}{}
	response := map[string]interface{}{}

	var resPeriodo map[string]interface{}
	var resDocente map[string]interface{}
	var resDocumento []map[string]interface{}
	var resVinculacion map[string]interface{}
	var resCarga map[string]interface{}
	var resEstado map[string]interface{}
	var indexSeleccionado int

	if errDocente := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+planes[0].(map[string]interface{})["docente_id"].(string), &resDocente); errDocente == nil {
		memDocente[planes[0].(map[string]interface{})["docente_id"].(string)] = resDocente
		if errDocumento := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=TerceroId.Id:"+planes[0].(map[string]interface{})["docente_id"].(string)+"&fields=Numero", &resDocumento); errDocumento == nil {
			memDocente = map[string]interface{}{
				"id":             planes[0].(map[string]interface{})["docente_id"].(string),
				"nombre":         utils.Capitalize(resDocente["NombreCompleto"].(string)),
				"identificacion": resDocumento[0]["Numero"],
				"nombre1":        utils.Capitalize(resDocente["PrimerNombre"].(string)),
				"apellido1":      utils.Capitalize(resDocente["PrimerApellido"].(string)),
			}
		}
	}

	if errPeriodo := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"periodo/"+fmt.Sprintf("%v", planes[0].(map[string]interface{})["periodo_id"]), &resPeriodo); errPeriodo == nil {
		response["periodo_academico"] = resPeriodo["Data"].(map[string]interface{})["Nombre"].(string)
	}

	for index, plan := range planes {
		var espacioPlan []interface{}
		cargaPlan := []interface{}{}
		if errVinculacion := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro/"+plan.(map[string]interface{})["tipo_vinculacion_id"].(string), &resVinculacion); errVinculacion == nil {
			vinculacion := resVinculacion["Data"].(map[string]interface{})["Nombre"].(string)
			vinculacion = strings.Replace(vinculacion, "DOCENTE DE ", "", 1)
			vinculacion = strings.ToLower(vinculacion)
			memVinculacion = append(memVinculacion, map[string]interface{}{"id": plan.(map[string]interface{})["tipo_vinculacion_id"].(string),
				"nombre": strings.ToUpper(vinculacion[0:1]) + vinculacion[1:]})
		}

		if fmt.Sprintf("%d", idVinculacion) == plan.(map[string]interface{})["tipo_vinculacion_id"].(string) {
			indexSeleccionado = index
		}

		if errCarga := request.GetJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"carga_plan?query=activo:true,plan_docente_id:"+plan.(map[string]interface{})["_id"].(string), &resCarga); errCarga == nil {
			if fmt.Sprintf("%v", resCarga["Data"]) != "[]" {
				for _, carga := range resCarga["Data"].([]interface{}) {
					var horario map[string]interface{}
					var sede []map[string]interface{}
					var edificio map[string]interface{}
					var salon map[string]interface{}
					var resColocacion map[string]interface{}
					var resumenColocacion map[string]interface{}
					var sedeId string
					var edificioId string
					var salonId string

					if colId, colExists := carga.(map[string]interface{})["colocacion_espacio_academico_id"]; colExists {
						if errColocacion := request.GetJson("https://"+beego.AppConfig.String("HorarioService")+"colocacion-espacio-academico/"+colId.(string), &resColocacion); errColocacion == nil {
							if resColocacion["Success"].(bool) {
								json.Unmarshal([]byte(resColocacion["Data"].(map[string]interface{})["ResumenColocacionEspacioFisico"].(string)), &resumenColocacion)
								json.Unmarshal([]byte(resColocacion["Data"].(map[string]interface{})["ColocacionEspacioAcademico"].(string)), &horario)
								sedeId = fmt.Sprintf("%v", resumenColocacion["espacio_fisico"].(map[string]interface{})["sede_id"])
								edificioId = fmt.Sprintf("%v", resumenColocacion["espacio_fisico"].(map[string]interface{})["edificio_id"])
								salonId = fmt.Sprintf("%v", resumenColocacion["espacio_fisico"].(map[string]interface{})["salon_id"])
								cargaDetalle := map[string]interface{}{
									"id":                              carga.(map[string]interface{})["_id"].(string),
									"horario":                         horario,
									"espacio_academico_id":            carga.(map[string]interface{})["espacio_academico_id"].(string),
									"colocacion_espacio_academico_id": carga.(map[string]interface{})["colocacion_espacio_academico_id"].(string),
								}
								if sedeId != "NA" {
									if errSede := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"espacio_fisico?query=Id:"+sedeId+"&fields=Id,Nombre,CodigoAbreviacion", &sede); errSede == nil {
										cargaDetalle["sede"] = sede[0]
									}
								} else {
									cargaDetalle["sede"] = "NA"
								}

								if edificioId != "NA" {
									if errEdificio := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"espacio_fisico/"+edificioId, &edificio); errEdificio == nil {
										cargaDetalle["edificio"] = edificio
									}
								} else {
									cargaDetalle["edificio"] = "NA"
								}

								if salonId != "NA" {
									if errSalon := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"espacio_fisico/"+salonId, &salon); errSalon == nil {
										cargaDetalle["salon"] = salon
									}
								} else {
									cargaDetalle["salon"] = "NA"
								}

								cargaPlan = append(cargaPlan, cargaDetalle)

								if carga.(map[string]interface{})["actividad_id"] != nil {
									cargaPlan[len(cargaPlan)-1].(map[string]interface{})["actividad_id"] = carga.(map[string]interface{})["actividad_id"].(string)
								} else {
									cargaPlan[len(cargaPlan)-1].(map[string]interface{})["espacio_academico_id"] = carga.(map[string]interface{})["espacio_academico_id"].(string)
								}
							}

						}
					}
				}
			}
		}

		var resPreasignacion map[string]interface{}
		if errPreasignacion := request.GetJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"pre_asignacion?query=activo:true,aprobacion_docente:true,aprobacion_proyecto:true,plan_docente_id:"+plan.(map[string]interface{})["_id"].(string), &resPreasignacion); errPreasignacion == nil {
			for _, preasignacion := range resPreasignacion["Data"].([]interface{}) {
				var resEspacioAcademico map[string]interface{}
				if memEspaciosDetalle[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)] == nil {
					if errEspacioAcademico := request.GetJson("https://"+beego.AppConfig.String("EspaciosAcademicosService")+"espacio-academico/"+preasignacion.(map[string]interface{})["espacio_academico_id"].(string), &resEspacioAcademico); errEspacioAcademico == nil {
						modular := false
						if val, ok := resEspacioAcademico["Data"].(map[string]interface{})["espacio_modular"]; ok {
							modular = val.(bool)
						}
						memEspaciosDetalle[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)] = map[string]interface{}{
							"espacio_academico": resEspacioAcademico["Data"].(map[string]interface{})["nombre"].(string),
							"nombre":            resEspacioAcademico["Data"].(map[string]interface{})["nombre"].(string) + " - " + resEspacioAcademico["Data"].(map[string]interface{})["grupo"].(string),
							"grupo":             resEspacioAcademico["Data"].(map[string]interface{})["grupo"],
							"codigo":            resEspacioAcademico["Data"].(map[string]interface{})["codigo"].(string),
							"id":                preasignacion.(map[string]interface{})["espacio_academico_id"].(string),
							"plan_id":           plan.(map[string]interface{})["_id"].(string),
							"proyecto_id":       resEspacioAcademico["Data"].(map[string]interface{})["proyecto_academico_id"],
							"espacio_modular":   modular,
						}
						espacioPlan = append(espacioPlan, memEspaciosDetalle[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)])
					}
				} else {
					espacioPlan = append(espacioPlan, memEspaciosDetalle[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)])
				}

			}
		}

		if plan.(map[string]interface{})["estado_plan_id"] != "Sin definir" {
			if errEstado := request.GetJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"estado_plan/"+plan.(map[string]interface{})["estado_plan_id"].(string), &resEstado); errEstado == nil {
				memEstados[plan.(map[string]interface{})["estado_plan_id"].(string)] = resEstado["Data"].(map[string]interface{})["nombre"].(string)
				memEstadoPlan = append(memEstadoPlan, memEstados[plan.(map[string]interface{})["estado_plan_id"].(string)].(string))
			}
		} else {
			memEstadoPlan = append(memEstadoPlan, plan.(map[string]interface{})["estado_plan_id"].(string))
		}

		resumenJSON := map[string]interface{}{}
		if plan.(map[string]interface{})["resumen"] != nil {
			json.Unmarshal([]byte(plan.(map[string]interface{})["resumen"].(string)), &resumenJSON)
		}

		memResumenes = append(memResumenes, resumenJSON)
		memEspacios = append(memEspacios, espacioPlan)
		memCarga = append(memCarga, cargaPlan)
		memPlanDocente = append(memPlanDocente, plan.(map[string]interface{})["_id"].(string))
	}

	relatedPlans := []string{}
	for _, espacioAcad := range memEspacios[0].([]interface{}) {
		if espacioAcad.(map[string]interface{})["espacio_modular"].(bool) {
			resp, err := requestmanager.Get(beego.AppConfig.String("PlanTrabajoDocenteService")+
				fmt.Sprintf("pre_asignacion?query=activo:true,aprobacion_proyecto:true,aprobacion_docente:true,periodo_id:%v,espacio_academico_id:%v", planes[0].(map[string]interface{})["periodo_id"], espacioAcad.(map[string]interface{})["id"]), requestmanager.ParseResponseFormato1)
			if err == nil {
				for _, preasign := range resp.([]interface{}) {
					if planes[0].(map[string]interface{})["docente_id"] != preasign.(map[string]interface{})["docente_id"] {
						relatedPlans = append(relatedPlans, fmt.Sprintf("?docente=%v&vigencia=%v&vinculacion=%v", preasign.(map[string]interface{})["docente_id"], planes[0].(map[string]interface{})["periodo_id"], preasign.(map[string]interface{})["tipo_vinculacion_id"]))
					}
				}
			}
		}
	}

	relatedPlansSimple := utils.RemoveDuplicated(relatedPlans, func(item interface{}) interface{} {
		return item.(string)
	})

	response["docente"] = memDocente
	response["tipo_vinculacion"] = memVinculacion
	response["carga"] = memCarga
	response["espacios_academicos"] = memEspacios
	response["seleccion"] = indexSeleccionado
	response["plan_docente"] = memPlanDocente
	response["estado_plan"] = memEstadoPlan
	response["vigencia"] = planes[0].(map[string]interface{})["periodo_id"].(string)
	response["resumenes"] = memResumenes
	response["planes_relacionados_query"] = relatedPlansSimple
	// response["actividades"] = memActividades

	return response
}

// CopiarPlanTrabajoDocente ...
func CopiarPTD(docente, vigenciaAnterior, vigencia, vinculacion int64, carga int8) requestmanager.APIResponse {
	response, err := requestmanager.Get(beego.AppConfig.String("PlanTrabajoDocenteService")+
		fmt.Sprintf("plan_docente?query=activo:true,docente_id:%d,periodo_id:%d,tipo_vinculacion_id:%d&fields=_id&limit=1", docente, vigenciaAnterior, vinculacion), requestmanager.ParseResponseFormato1)
	if err != nil {
		logs.Error(err)
		return requestmanager.APIResponseDTO(false, 404, nil, "PlanTrabajoDocenteService (plan_docente): "+err.Error())
		/* badAns, code := requestmanager.MidResponseFormat("PlanTrabajoDocenteService (plan_docente)", "GET", false, map[string]interface{}{
			"response": response,
			"error":    err.Error(),
		})
		c.Ctx.Output.SetStatus(code)
		c.Data["json"] = badAns
		c.ServeJSON()
		return */
	}
	plan_docenteAnterior := []models.PlanDocente{}
	utils.ParseData(response, &plan_docenteAnterior)

	response, err = requestmanager.Get(beego.AppConfig.String("PlanTrabajoDocenteService")+
		fmt.Sprintf("carga_plan?query=activo:true,plan_docente_id:%s&limit=0", plan_docenteAnterior[0].Id), requestmanager.ParseResponseFormato1)
	if err != nil {
		logs.Error(err)
		return requestmanager.APIResponseDTO(false, 404, nil, "PlanTrabajoDocenteService (carga_plan): "+err.Error())
		/* badAns, code := requestmanager.MidResponseFormat("PlanTrabajoDocenteService (carga_plan)", "GET", false, map[string]interface{}{
			"response": response,
			"error":    err.Error(),
		})
		c.Ctx.Output.SetStatus(code)
		c.Data["json"] = badAns
		c.ServeJSON()
		return */
	}
	carga_planAnterior := []models.CargaPlan{}
	utils.ParseData(response, &carga_planAnterior)

	prepareAns := map[string]interface{}{}

	if carga == 1 { // ? Carga -> 1, para Carga Lectiva AKA espacios académicos
		cargaLectiva, errCL := obtenerCargaLectiva(docente, vigencia, vigenciaAnterior, vinculacion, carga_planAnterior)
		if errCL != nil {
			logs.Error(errCL)
			return requestmanager.APIResponseDTO(false, 404, nil, "PlanTrabajoDocenteService (obtenerCargaLectiva): "+errCL.Error())
		}
		prepareAns = cargaLectiva
	} else if carga == 2 { // ? Carga -> 2, para Actividades
		cargaActividades, errCA := obtenerCargaActividades(carga_planAnterior)
		if errCA != nil {
			logs.Error(errCA)
			return requestmanager.APIResponseDTO(false, 404, nil, "PlanTrabajoDocenteService (obtenerCargaActividades): "+errCA.Error())
		}
		prepareAns = cargaActividades
	}

	return requestmanager.APIResponseDTO(true, 200, prepareAns)
}

func obtenerCargaLectiva(docente, vigenciaAnterior, vigencia, vinculacion int64, carga_planAnterior []models.CargaPlan) (map[string]interface{}, error) {
	// * ----------
	// * Consultas sobre las preasignaciones actual y anterior para encontrar similitudes
	//
	preasignacionAnterior, err := consultarEspaciosAcademicosInfoPadre(docente, vigenciaAnterior, vinculacion)
	if err != nil {
		logs.Error(err)
		return nil, fmt.Errorf("func consultarEspaciosAcademicosPadre: %s", err.Error())
		/* badAns, code := requestmanager.MidResponseFormat("func consultarEspaciosAcademicosPadre", "GET", false, err.Error())
		c.Ctx.Output.SetStatus(code)
		c.Data["json"] = badAns
		c.ServeJSON()
		return */
	}
	preasignacionActual, err := consultarEspaciosAcademicosInfoPadre(docente, vigencia, vinculacion)
	if err != nil {
		logs.Error(err)
		return nil, fmt.Errorf("func consultarEspaciosAcademicosPadre: %s", err.Error())
		/* badAns, code := requestmanager.MidResponseFormat("func consultarEspaciosAcademicosPadre", "GET", false, err.Error())
		c.Ctx.Output.SetStatus(code)
		c.Data["json"] = badAns
		c.ServeJSON()
		return */
	}

	igual := utils.JoinEqual(preasignacionAnterior, preasignacionActual, "A", func(valor interface{}) interface{} {
		return valor.(models.EspacioAcademico).Espacio_academico_padre
	})
	preasignIgualAnterior := []models.EspacioAcademico{}
	utils.ParseData(igual, &preasignIgualAnterior)
	//
	// * ----------

	// * ----------
	// * Iteracion sobre preasignación igual y el plan anterior, agrega solo carga academica en base a preasignación actual
	//
	listaCarga := []interface{}{}
	for _, preasignEspacioAcad := range preasignIgualAnterior {
		for _, carga := range carga_planAnterior {
			if carga.Horario == "" && carga.Colocacion_espacio_academico_id != "" {
				var colocacion map[string]interface{}
				if errGetColocacion := request.GetJson("https://"+beego.AppConfig.String("HorarioService")+
					"colocacion-espacio-academico/"+carga.Colocacion_espacio_academico_id, &colocacion); errGetColocacion == nil {
					if colocacion["Success"].(bool) {
						var resumenColocacionJSON map[string]interface{}
						json.Unmarshal([]byte(colocacion["Data"].(map[string]interface{})["ResumenColocacionEspacioFisico"].(string)), &resumenColocacionJSON)
						carga.Horario = colocacion["Data"].(map[string]interface{})["ColocacionEspacioAcademico"].(string)
						carga.Sede_id = resumenColocacionJSON["espacio_fisico"].(map[string]any)["sede_id"].(string)
					}
				} else {
					logs.Error(errGetColocacion)
					return nil, fmt.Errorf("HorarioService (colocacion_espacio_academico): %s", errGetColocacion.Error())
					/* badAns, code := requestmanager.MidResponseFormat("HorarioService (colocacion_espacio_academico)", "GET", false, map[string]interface{}{
						"response": response,
						"error":    errGetColocacion.Error(),
					})
					c.Ctx.Output.SetStatus(code)
					c.Data["json"] = badAns
					c.ServeJSON()
					return */
				}

			}
			if preasignEspacioAcad.Id == carga.Espacio_academico_id {
				encontrado := utils.Find(preasignacionActual, func(valor interface{}) bool {
					return valor.(models.EspacioAcademico).Espacio_academico_padre == preasignEspacioAcad.Espacio_academico_padre
				})
				espacioAcademicoNuevo := models.EspacioAcademico{}
				utils.ParseData(encontrado, &espacioAcademicoNuevo)

				infoEspacio, err := consultarInfoEspacioFisico(carga.Sede_id, carga.Edificio_id, carga.Salon_id)
				if err != nil {
					logs.Error(err)
					return nil, fmt.Errorf("OikosService (espacio_fisico): %s", err.Error())
					/* badAns, code := requestmanager.MidResponseFormat("OikosService (espacio_fisico)", "GET", false, map[string]interface{}{
						"response": response,
						"error":    err.Error(),
					})
					c.Ctx.Output.SetStatus(code)
					c.Data["json"] = badAns
					c.ServeJSON()
					return */
				}

				var horarioJson interface{}
				err = json.Unmarshal([]byte(carga.Horario), &horarioJson)
				if err != nil {
					logs.Error(err)
					return nil, fmt.Errorf("CopiarPlanTrabajoDocente (parse horario): %s", err.Error())
					/* badAns, code := requestmanager.MidResponseFormat("CopiarPlanTrabajoDocente (parse horario)", "GET", false, map[string]interface{}{
						"response": response,
						"error":    err.Error(),
					})
					c.Ctx.Output.SetStatus(code)
					c.Data["json"] = badAns
					c.ServeJSON()
					return */
				}

				listaCarga = append(listaCarga, map[string]interface{}{
					"id":                   nil,
					"sede":                 infoEspacio.(map[string]interface{})["sede"],
					"edificio":             infoEspacio.(map[string]interface{})["edificio"],
					"salon":                infoEspacio.(map[string]interface{})["salon"],
					"espacio_academico_id": espacioAcademicoNuevo.Id,
					"horario":              horarioJson,
				})
			}
		}
	}
	//
	// * ----------

	// * ----------
	// * Resumen de diferencias positivas y negativas de carga lectiva
	//
	norequeridos := utils.SubstractDiff(preasignacionActual, preasignacionAnterior, "B", func(valor interface{}) interface{} {
		return valor.(models.EspacioAcademico).Espacio_academico_padre
	})
	sincarga := utils.SubstractDiff(preasignacionActual, preasignacionAnterior, "A", func(valor interface{}) interface{} {
		return valor.(models.EspacioAcademico).Espacio_academico_padre
	})
	//
	// * ----------

	return map[string]interface{}{
		"carga": listaCarga,
		"espacios_academicos": map[string]interface{}{
			"no_requeridos": norequeridos,
			"sin_carga":     sincarga,
		},
	}, nil
}

func obtenerCargaActividades(carga_planAnterior []models.CargaPlan) (map[string]interface{}, error) {
	listaCarga := []interface{}{}
	for _, carga := range carga_planAnterior {
		if carga.Horario == "" && carga.Colocacion_espacio_academico_id != "" {
			var colocacion map[string]interface{}
			if errGetColocacion := request.GetJson("https://"+beego.AppConfig.String("HorarioService")+
				"colocacion-espacio-academico/"+carga.Colocacion_espacio_academico_id, &colocacion); errGetColocacion == nil {
				if colocacion["Success"].(bool) {
					var resumenColocacionJSON map[string]interface{}
					json.Unmarshal([]byte(colocacion["Data"].(map[string]interface{})["ResumenColocacionEspacioFisico"].(string)), &resumenColocacionJSON)
					carga.Horario = colocacion["Data"].(map[string]interface{})["ColocacionEspacioAcademico"].(string)
					carga.Sede_id = resumenColocacionJSON["espacio_fisico"].(map[string]any)["sede_id"].(string)
				}
			} else {
				logs.Error(errGetColocacion)
				return nil, fmt.Errorf("HorarioService (colocacion_espacio_academico): %s", errGetColocacion.Error())
				/* badAns, code := requestmanager.MidResponseFormat("HorarioService (colocacion_espacio_academico)", "GET", false, map[string]interface{}{
					"response": response,
					"error":    errGetColocacion.Error(),
				})
				c.Ctx.Output.SetStatus(code)
				c.Data["json"] = badAns
				c.ServeJSON()
				return */
			}
		}
		if carga.Actividad_id != "" {
			infoEspacio, err := consultarInfoEspacioFisico(carga.Sede_id, carga.Edificio_id, carga.Salon_id)
			if err != nil {
				logs.Error(err)
				return nil, fmt.Errorf("OikosService (espacio_fisico): %s", err.Error())
				/* badAns, code := requestmanager.MidResponseFormat("OikosService (espacio_fisico)", "GET", false, map[string]interface{}{
					"response": response,
					"error":    err.Error(),
				})
				c.Ctx.Output.SetStatus(code)
				c.Data["json"] = badAns
				c.ServeJSON()
				return */
			}

			var horarioJson interface{}
			err = json.Unmarshal([]byte(carga.Horario), &horarioJson)
			if err != nil {
				logs.Error(err)
				return nil, fmt.Errorf("CopiarPlanTrabajoDocente (parse horario): %s", err.Error())
				/* badAns, code := requestmanager.MidResponseFormat("CopiarPlanTrabajoDocente (parse horario)", "GET", false, map[string]interface{}{
					"response": response,
					"error":    err.Error(),
				})
				c.Ctx.Output.SetStatus(code)
				c.Data["json"] = badAns
				c.ServeJSON()
				return */
			}

			listaCarga = append(listaCarga, map[string]interface{}{
				"id":           nil,
				"sede":         infoEspacio.(map[string]interface{})["sede"],
				"edificio":     infoEspacio.(map[string]interface{})["edificio"],
				"salon":        infoEspacio.(map[string]interface{})["salon"],
				"actividad_id": carga.Actividad_id,
				"horario":      horarioJson,
			})
		}
	}
	return map[string]interface{}{
		"carga": listaCarga,
	}, nil
}

// PlanPreaprobado ...
func ListaPlanPreaprobado(vigencia, proyecto int64) requestmanager.APIResponse {
	rawPlanes := []interface{}{}
	estado := "64c2ca7fd1e67f67f057f3c8" // preaprobado
	resp, err := requestmanager.Get(beego.AppConfig.String("PlanTrabajoDocenteService")+
		fmt.Sprintf("plan_docente?query=activo:true,estado_plan_id:%s,periodo_id:%d&limit=0", estado, vigencia), requestmanager.ParseResponseFormato1)
	if err == nil {
		rawPlanes = append(rawPlanes, resp.([]interface{})...)
	}
	estado = "646fcf784c0bc253c1c720d4" // aprobado
	resp, err = requestmanager.Get(beego.AppConfig.String("PlanTrabajoDocenteService")+
		fmt.Sprintf("plan_docente?query=activo:true,estado_plan_id:%s,periodo_id:%d&limit=0", estado, vigencia), requestmanager.ParseResponseFormato1)
	if err == nil {
		rawPlanes = append(rawPlanes, resp.([]interface{})...)
	}
	estado = "646fcf8a4c0bc253c1c720d6" // no aprobado
	resp, err = requestmanager.Get(beego.AppConfig.String("PlanTrabajoDocenteService")+
		fmt.Sprintf("plan_docente?query=activo:true,estado_plan_id:%s,periodo_id:%d&limit=0", estado, vigencia), requestmanager.ParseResponseFormato1)
	if err == nil {
		rawPlanes = append(rawPlanes, resp.([]interface{})...)
	}

	lista_planes := []models.PlanDocente{}
	utils.ParseData(rawPlanes, &lista_planes)

	planes_proyecto := []models.PlanDocente{}

	for _, plan := range lista_planes {
		_, err := requestmanager.Get("https://"+beego.AppConfig.String("EspaciosAcademicosService")+
			fmt.Sprintf("espacio-academico?query=activo:true,periodo_id:%d,proyecto_academico_id:%d,docente_id:%s&fields=_id&limit=0", vigencia, proyecto, plan.Docente_id), requestmanager.ParseResponseFormato1)
		if err == nil {
			planes_proyecto = append(planes_proyecto, plan)
		}
	}

	if len(planes_proyecto) == 0 {
		return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron planes para el proyecto")
		/* respuesta, statuscode := requestmanager.MidResponseFormat("GetPlanesPreaprobados", "GET", false, nil)
		c.Ctx.Output.SetStatus(statuscode)
		c.Data["json"] = respuesta
		c.ServeJSON()
		return */
	}

	prepareAns := []map[string]interface{}{}
	estadoAprobadoIdActual := "646fcf784c0bc253c1c720d4"

	for _, planProyecto := range planes_proyecto {
		resp, err := requestmanager.Get("http://"+beego.AppConfig.String("TercerosService")+
			fmt.Sprintf("datos_identificacion?query=Activo:true,TerceroId__Id:%v&fields=TerceroId,Numero,TipoDocumentoId&sortby=FechaExpedicion,Id&order=desc&limit=1",
				planProyecto.Docente_id), requestmanager.ParseResonseNoFormat)
		if err != nil {
			logs.Error(err)
			return requestmanager.APIResponseDTO(false, 404, nil, "TercerosService (datos_identificacion): "+err.Error())
			/* badAns, code := requestmanager.MidResponseFormat("TercerosService (datos_identificacion)", "GET", false, map[string]interface{}{
				"response": resp,
				"error":    err.Error(),
			})
			c.Ctx.Output.SetStatus(code)
			c.Data["json"] = badAns
			c.ServeJSON()
			return */
		}
		datos_identificacion := models.DatosIdentificacion{}
		utils.ParseData(resp.([]interface{})[0], &datos_identificacion)

		resp, err = requestmanager.Get("http://"+beego.AppConfig.String("ParametroService")+
			fmt.Sprintf("parametro/%s", planProyecto.Tipo_vinculacion_id), requestmanager.ParseResponseFormato1)
		if err != nil {
			logs.Error(err)
			return requestmanager.APIResponseDTO(false, 404, nil, "ParametroService (parametro): "+err.Error())
			/* badAns, code := requestmanager.MidResponseFormat("ParametroService (parametro)", "GET", false, map[string]interface{}{
				"response": resp,
				"error":    err.Error(),
			})
			c.Ctx.Output.SetStatus(code)
			c.Data["json"] = badAns
			c.ServeJSON()
			return */
		}
		infoVinculacion := models.Parametro{}
		utils.ParseData(resp, &infoVinculacion)

		desactivarSoporte := strings.TrimSpace(planProyecto.Soporte_documental) == "" || planProyecto.Estado_plan_id != estadoAprobadoIdActual

		prepareAns = append(prepareAns, map[string]interface{}{
			"id":                 planProyecto.Id,
			"nombre":             utils.FormatNameTercero(datos_identificacion.TerceroId),
			"identificacion":     datos_identificacion.Numero,
			"tipo_vinculacion":   infoVinculacion.Nombre,
			"periodo_academico":  vigencia,
			"soporte_documental": map[string]interface{}{"value": planProyecto.Soporte_documental, "type": "ver", "disabled": desactivarSoporte},
			"gestion":            map[string]interface{}{"value": nil, "type": "editar", "disabled": false},
			"estado":             planProyecto.Estado_plan_id,
			"tercero_id":         datos_identificacion.TerceroId.Id,
			"vinculacion_id":     planProyecto.Tipo_vinculacion_id,
		})
	}

	return requestmanager.APIResponseDTO(true, 200, prepareAns)
}

// funciones transversales
func consultarEspaciosAcademicosInfoPadre(docente, periodo, vinculacion int64) ([]models.EspacioAcademico, error) {
	espacios := []models.EspacioAcademico{}
	response, err := requestmanager.Get(beego.AppConfig.String("PlanTrabajoDocenteService")+
		fmt.Sprintf("pre_asignacion?query=activo:true,aprobacion_docente:true,aprobacion_proyecto:true,docente_id:%d,periodo_id:%d,tipo_vinculacion_id:%d&fields=espacio_academico_id", docente, periodo, vinculacion),
		requestmanager.ParseResponseFormato1)
	if err != nil {
		return nil, fmt.Errorf("PlanTrabajoDocenteService (pre_asignacion): %s", err.Error())
	}
	preasignaciones := []models.PreAsignacion{}
	utils.ParseData(response, &preasignaciones)
	for _, preasignacion := range preasignaciones {
		response, err := requestmanager.Get("https://"+beego.AppConfig.String("EspaciosAcademicosService")+
			fmt.Sprintf("espacio-academico?query=activo:true,_id:%s&fields=_id,nombre,espacio_academico_padre&limit=1", preasignacion.Espacio_academico_id), requestmanager.ParseResponseFormato1)
		if err != nil {
			return nil, fmt.Errorf("EspaciosAcademicosService (espacio-academico): %s", err.Error())
		}
		espacioacademico := []models.EspacioAcademico{}
		utils.ParseData(response, &espacioacademico)
		espacios = append(espacios, espacioacademico[0])
	}
	return espacios, nil
}
