package services

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/trabajo_docente_mid/utils"
	request "github.com/udistrital/utils_oas/request"
	requestmanager "github.com/udistrital/utils_oas/requestresponse"
)

// GrupoEspacioAcademico ...
func ListaGruposEspaciosAcademicos(padre, vigencia string) requestmanager.APIResponse {
	var response []interface{}
	queryParams := "query=activo:true,espacio_academico_padre:" + padre +
		",periodo_id:" + vigencia
	if resSpaces, errSpace := getAcademicSpacesByQuery(queryParams); errSpace == nil {
		if resSpaces != nil {
			spaces := resSpaces.([]interface{})
			for _, space := range spaces {
				spaceMap := space.(map[string]interface{})
				if spaceMap["espacio_modular"] == true || fmt.Sprintf("%v", spaceMap["docente_id"]) == "0" {
					var resProject []interface{}
					queryParams = "query=Id:" +
						fmt.Sprintf("%v", spaceMap["proyecto_academico_id"]) +
						"&fields=Nombre,Id,NivelFormacionId"
					if errProject := getAcademicProjectByQuery(queryParams, &resProject); errProject == nil {
						if resProject[0].(map[string]interface{})["Id"] != nil {
							response = append(response, map[string]interface{}{
								"Id":                spaceMap["_id"],
								"Nombre":            spaceMap["nombre"],
								"ProyectoAcademico": resProject[0].(map[string]interface{})["Nombre"],
								"Nivel":             resProject[0].(map[string]interface{})["NivelFormacionId"].(map[string]interface{})["Nombre"],
								"grupo":             spaceMap["grupo"],
							})
						}
					}
				}
			}
			return requestmanager.APIResponseDTO(true, 200, response)
			/* c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": response} */
		} else {
			return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de espacios académicos")
			/* c.Ctx.Output.SetStatus(404)
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron espacios académicos 1"} */
		}
	} else {
		logs.Error(errSpace)
		return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de espacios académicos")
		/* c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron espacios académicos"} */
	}
}

func getAcademicSpacesByQuery(query string) (any, error) {
	var resSpaces interface{}
	urlAcademicSpaces := beego.AppConfig.String("EspaciosAcademicosService") +
		"espacio-academico?" + query
	if errSpace := request.GetJson(urlAcademicSpaces, &resSpaces); errSpace == nil {
		if resSpaces.(map[string]interface{})["Data"] != nil {
			return resSpaces.(map[string]interface{})["Data"], nil
		} else {
			return nil, fmt.Errorf("EspaciosAcademicosService No se encuentran espacios académicos")
		}
	} else {
		return nil, errSpace
	}
}

func getAcademicProjectByQuery(query string, resProject *[]any) error {
	urlAcademicProject := beego.AppConfig.String("ProyectoAcademicoService") +
		"proyecto_academico_institucion?" + query

	if errProject := request.GetJson(urlAcademicProject, &resProject); errProject == nil {
		return nil
	} else {
		return errProject
	}
}

// GrupoEspacioAcademicoPadre ...
func ListaGruposEspaciosAcademicosPadre(padre string) requestmanager.APIResponse {
	if response, errGroupsSpace := getAcademicSpaces2AssignPeriodByParent(padre); errGroupsSpace == nil {
		return requestmanager.APIResponseDTO(true, 200, response)
		/* c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": response} */
	} else {
		logs.Error(errGroupsSpace)
		return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de espacios académicos")
		/* c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron espacios académicos"} */
	}
}

func getAcademicSpaces2AssignPeriodByParent(parent string) (any, error) {
	var response []any
	queryParams := "query=_id:" + parent +
		"&fields=nombre,grupo,proyecto_academico_id"
	if resSpaces, errSpace := getAcademicSpacesByQuery(queryParams); errSpace == nil {
		spaces := resSpaces.([]any)
		for _, space := range spaces {
			groups := utils.SplitTrimSpace(fmt.Sprintf("%v", space.(map[string]interface{})["grupo"]),
				",")
			var resProject []any
			queryParams = "query=Id:" +
				fmt.Sprintf("%v", space.(map[string]any)["proyecto_academico_id"]) +
				"&fields=Nombre,Id,NivelFormacionId"
			if errProject := getAcademicProjectByQuery(queryParams, &resProject); errProject == nil {
				projectData := resProject[0].(map[string]any)
				if projectData["Id"] != nil {
					response = append(response, map[string]interface{}{
						"Id":                space.(map[string]interface{})["_id"],
						"Nombre":            space.(map[string]interface{})["nombre"],
						"ProyectoAcademico": projectData["Nombre"],
						"Nivel":             projectData["NivelFormacionId"].(map[string]interface{})["NivelFormacionPadreId"].(map[string]interface{})["Nombre"],
						"Subnivel":          projectData["NivelFormacionId"].(map[string]interface{})["Nombre"],
						"Grupos":            groups,
					})
				}
			}
		}
		return response, nil
	} else {
		return nil, errSpace
	}
}
