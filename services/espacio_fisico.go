package services

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	request "github.com/udistrital/utils_oas/request"
	requestmanager "github.com/udistrital/utils_oas/requestresponse"
)

// EspacioFisicoDependencia ...
func ArbolEspaciosFisicosDependencia(dependencia int64) requestmanager.APIResponse {
	_ = dependencia

	baseOikosURL := "http://" + beego.AppConfig.String("OikosService")

	sedes, err := requestmanager.Get(baseOikosURL+
		"espacio_fisico?query=Activo:true,TipoEspacioFisicoId__Id:38&limit=0&sortby=Nombre&order=asc",
		requestmanager.ParseResonseNoFormat)
	if err != nil {
		logs.Error(err)
		return requestmanager.APIResponseDTO(false, 404, nil, "OikosService (sedes: TipoEspacioFisicoId 38): "+err.Error())
	}
	sedes = filterDistinctByField(sedes, "Nombre")

	edificios, err := requestmanager.Get(baseOikosURL+
		"espacio_fisico?query=Activo:true,TipoEspacioFisicoId__Id:39&limit=0&sortby=Nombre&order=asc",
		requestmanager.ParseResonseNoFormat)
	if err != nil {
		logs.Error(err)
		return requestmanager.APIResponseDTO(false, 404, nil, "OikosService (edificios: TipoEspacioFisicoId 39): "+err.Error())
	}
	edificios = filterDistinctByField(edificios, "Nombre")

	salones, err := requestmanager.Get(baseOikosURL+
		"espacio_fisico?query=Activo:true,TipoEspacioFisicoId__Id:2&limit=0&sortby=Nombre&order=asc",
		requestmanager.ParseResonseNoFormat)
	if err != nil {
		logs.Error(err)
		return requestmanager.APIResponseDTO(false, 404, nil, "OikosService (salones: TipoEspacioFisicoId 2): "+err.Error())
	}
	salones = filterDistinctByField(salones, "Nombre")

	return requestmanager.APIResponseDTO(true, 200, map[string]interface{}{
		"Salones":    salones,
		"Edificios":  edificios,
		"Sedes":      sedes,
		"PorAsignar": false,
	})

	/*
		// Lógica previa basada en jerarquías de asignación. Se conserva comentada mientras se corrigen los datos en Oikos.
		inBog, _ := time.LoadLocation("America/Bogota")
		horaes := time.Now().In(inBog).Format(time.RFC3339)

		resp, err := requestmanager.Get("http://"+beego.AppConfig.String("ProyectoAcademicoService")+
			fmt.Sprintf("proyecto_academico_institucion/%d", dependencia), requestmanager.ParseResonseNoFormat)
		if err != nil {
			logs.Error(err)
			return requestmanager.APIResponseDTO(false, 404, nil, "ProyectoAcademicoService (proyecto_academico_institucion): "+err.Error())
		}

		dependencia = int64(resp.(map[string]interface{})["DependenciaId"].(float64))
		if dependencia <= 0 {
			err = fmt.Errorf("no valid Id: %d > 0 = false", dependencia)
			logs.Error(err)
			return requestmanager.APIResponseDTO(false, 404, nil, "GetEspaciosFisicosDependencia (param: dependencia): "+err.Error())
		}

		Salones := map[string][]map[string]interface{}{}
		Edificios := map[string][]map[string]interface{}{}
		Sedes := []map[string]interface{}{}

		resp, err = requestmanager.Get("http://"+beego.AppConfig.String("OikosService")+
			fmt.Sprintf("asignacion_espacio_fisico_dependencia?query=Activo:true,DependenciaId:%d,FechaInicio__lte:%v,FechaFin__gte:%v&fields=EspacioFisicoId&limit=0",
				dependencia, horaes, horaes), requestmanager.ParseResonseNoFormat)
		if err != nil {
			resp, err = requestmanager.Get("http://"+beego.AppConfig.String("OikosService")+"espacio_fisico?query=Nombre:POR%20ASIGNAR,TipoEspacioFisicoId__Id:2", requestmanager.ParseResonseNoFormat)
			if err != nil {
				logs.Error(err)
				return requestmanager.APIResponseDTO(false, 404, nil, "OikosService (espacio_fisico): "+err.Error())
			}
			Idstr := fmt.Sprintf("%v", resp.([]interface{})[0].(map[string]interface{})["Id"])
			Opcion := map[string]interface{}{
				"Id":     resp.([]interface{})[0].(map[string]interface{})["Id"],
				"Nombre": resp.([]interface{})[0].(map[string]interface{})["Nombre"],
			}
			Salones[Idstr] = append(Salones[Idstr], Opcion)
			Edificios[Idstr] = append(Edificios[Idstr], Opcion)
			Sedes = append(Sedes, Opcion)

			return requestmanager.APIResponseDTO(true, 200, map[string]interface{}{
				"Salones":    Salones,
				"Edificios":  Edificios,
				"Sedes":      Sedes,
				"PorAsignar": true,
			})
		}

		for _, EspacioFisico := range resp.([]interface{}) {
			resp, err := requestmanager.Get("http://"+beego.AppConfig.String("OikosService")+
				fmt.Sprintf("espacio_fisico_padre?query=HijoId:%v", EspacioFisico.(map[string]interface{})["EspacioFisicoId"].(map[string]interface{})["Id"]), requestmanager.ParseResonseNoFormat)
			if err == nil {
				tipoEspacio := resp.([]interface{})[0].(map[string]interface{})["PadreId"].(map[string]interface{})["TipoEspacioFisicoId"].(map[string]interface{})["Id"].(float64)
				PadreSalon := fmt.Sprintf("%v", resp.([]interface{})[0].(map[string]interface{})["PadreId"].(map[string]interface{})["Id"])
				for tipoEspacio != 39 {
					resp, err := requestmanager.Get("http://"+beego.AppConfig.String("OikosService")+
						fmt.Sprintf("espacio_fisico_padre?query=HijoId:%v", PadreSalon), requestmanager.ParseResonseNoFormat)
					if err == nil {
						PadreSalon = fmt.Sprintf("%v", resp.([]interface{})[0].(map[string]interface{})["PadreId"].(map[string]interface{})["Id"])
						tipoEspacio = resp.([]interface{})[0].(map[string]interface{})["PadreId"].(map[string]interface{})["TipoEspacioFisicoId"].(map[string]interface{})["Id"].(float64)
					}
				}

				if _, ok := Salones[PadreSalon]; !ok {
					Salones[PadreSalon] = []map[string]interface{}{}
				}
				Salones[PadreSalon] = append(Salones[PadreSalon], map[string]interface{}{
					"Id":                resp.([]interface{})[0].(map[string]interface{})["HijoId"].(map[string]interface{})["Id"],
					"Nombre":            resp.([]interface{})[0].(map[string]interface{})["HijoId"].(map[string]interface{})["Nombre"],
					"Descripcion":       resp.([]interface{})[0].(map[string]interface{})["HijoId"].(map[string]interface{})["Descripcion"],
					"CodigoAbreviacion": resp.([]interface{})[0].(map[string]interface{})["HijoId"].(map[string]interface{})["CodigoAbreviacion"],
				})

			}
		}

		for PadreSalon := range Salones {
			resp, err := requestmanager.Get("http://"+beego.AppConfig.String("OikosService")+
				fmt.Sprintf("espacio_fisico_padre?query=HijoId:%v", PadreSalon), requestmanager.ParseResonseNoFormat)
			if err == nil {
				PadreEdificio := fmt.Sprintf("%v", resp.([]interface{})[0].(map[string]interface{})["PadreId"].(map[string]interface{})["Id"])
				if _, ok := Edificios[PadreEdificio]; !ok {
					Edificios[PadreEdificio] = []map[string]interface{}{}
				}
				Edificios[PadreEdificio] = append(Edificios[PadreEdificio], map[string]interface{}{
					"Id":                resp.([]interface{})[0].(map[string]interface{})["HijoId"].(map[string]interface{})["Id"],
					"Nombre":            resp.([]interface{})[0].(map[string]interface{})["HijoId"].(map[string]interface{})["Nombre"],
					"Descripcion":       resp.([]interface{})[0].(map[string]interface{})["HijoId"].(map[string]interface{})["Descripcion"],
					"CodigoAbreviacion": resp.([]interface{})[0].(map[string]interface{})["HijoId"].(map[string]interface{})["CodigoAbreviacion"],
				})
			}
		}

		for PadreEficio := range Edificios {
			resp, err := requestmanager.Get("http://"+beego.AppConfig.String("OikosService")+
				fmt.Sprintf("espacio_fisico_padre?query=HijoId:%v", PadreEficio), requestmanager.ParseResonseNoFormat)
			if err == nil {
				Sedes = append(Sedes, map[string]interface{}{
					"Id":                resp.([]interface{})[0].(map[string]interface{})["HijoId"].(map[string]interface{})["Id"],
					"Nombre":            resp.([]interface{})[0].(map[string]interface{})["HijoId"].(map[string]interface{})["Nombre"],
					"Descripcion":       resp.([]interface{})[0].(map[string]interface{})["HijoId"].(map[string]interface{})["Descripcion"],
					"CodigoAbreviacion": resp.([]interface{})[0].(map[string]interface{})["HijoId"].(map[string]interface{})["CodigoAbreviacion"],
				})
			}
		}

		return requestmanager.APIResponseDTO(true, 200, map[string]interface{}{
			"Salones":   Salones,
			"Edificios": Edificios,
			"Sedes":     Sedes,
		})
	*/

}

// DisponibilidadEspacioFisico ...
func OcupacionEspacioFisico(salon, vigencia, plan string) requestmanager.APIResponse {
	var planTrabajoDocente map[string]interface{}
	var cargaPlan map[string]interface{}
	var colocacion map[string]interface{}
	var cargas []map[string]interface{}
	var errorGetAll bool

	if errGetPlan := request.GetJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"plan_docente?query=activo:true,periodo_id:"+vigencia+"&fields=_id", &planTrabajoDocente); errGetPlan == nil {
		if fmt.Sprintf("%v", planTrabajoDocente["Data"]) != "[]" {
			planes := planTrabajoDocente["Data"].([]interface{})
			for _, plan := range planes {
				if errGetCargas := request.GetJson(beego.AppConfig.String("PlanTrabajoDocenteService")+"carga_plan?query=activo:true,salon_id:"+salon+",plan_docente_id:"+plan.(map[string]interface{})["_id"].(string)+"&fields=horario,plan_docente_id,colocacion_espacio_academico_id", &cargaPlan); errGetCargas == nil {
					if fmt.Sprintf("%v", cargaPlan["Data"]) != "[]" {
						for _, carga := range cargaPlan["Data"].([]interface{}) {
							if carga.(map[string]interface{})["plan_docente_id"] != plan {
								if colId, colExists := carga.(map[string]interface{})["colocacion_espacio_academico_id"]; colExists {
									var horarioJSON map[string]interface{}
									if errGetColocacion := request.GetJson("https://"+beego.AppConfig.String("HorarioService")+"colocacion-espacio-academico/"+colId.(string), &colocacion); errGetColocacion == nil {
										if colocacion["Success"].(bool) {
											json.Unmarshal([]byte(colocacion["Data"].(map[string]interface{})["ColocacionEspacioAcademico"].(string)), &horarioJSON)
											cargas = append(cargas, map[string]interface{}{
												"finalPosition": horarioJSON["finalPosition"],
												"horas":         horarioJSON["horas"],
												"id":            carga.(map[string]interface{})["_id"],
											})
										}
									}
								}
							}
						}
					}
				}
			}
		} else {
			return requestmanager.APIResponseDTO(false, 404, nil, "No hay planes de trabajo docente para la vigencia seleccionada")
		}
	}

	if errorGetAll {
		return requestmanager.APIResponseDTO(false, 404, nil, "No hay planes de trabajo docente para la vigencia seleccionada")
	} else {
		return requestmanager.APIResponseDTO(true, 200, cargas)
	}
}

// filtro para que no se repitan los espacios físicos con el mismo nombre
func filterDistinctByField(data interface{}, field string) interface{} {
	list, ok := data.([]interface{})
	if !ok {
		return data
	}

	seen := make(map[string]struct{})
	unique := make([]interface{}, 0, len(list))

	for _, item := range list {
		row, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		key := fmt.Sprintf("%v", row[field])
		if key == "" {
			continue
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		unique = append(unique, item)
	}

	return unique
}
