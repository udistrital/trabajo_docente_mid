package services

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/trabajo_docente_mid/utils"
	request "github.com/udistrital/utils_oas/request"
	requestmanager "github.com/udistrital/utils_oas/requestresponse"
)

type espaciosAcademicosXML struct {
	Espacios []espacioAcademicoXML `xml:"espacio_academico"`
}

type espacioAcademicoXML struct {
	CodigoCarrera  string `xml:"codigo_carrera"`
	NombreCarrera  string `xml:"nombre_carrera"`
	Id             string `xml:"_id"`
	Nombre         string `xml:"nombre"`
	EspacioModular bool   `xml:"espacio_modular"`
}

type gruposEspacioPeriodoXML struct {
	Grupos []grupoEspacioPeriodoXML `xml:"grupo"`
}

type grupoEspacioPeriodoXML struct {
	Grupo             string `xml:"grupo"`
	ProyectoAcademico string `xml:"proyecto_academico"`
	Id                string `xml:"id"`
	Nivel             string `xml:"nivel"`
	Nombre            string `xml:"nombre"`
}

type detalleCursoIdXML struct {
	Id                      string `xml:"id"`
	Nivel                   string `xml:"nivel"`
	CodigoProyectoAcademico string `xml:"codigoProyectoAcademico"`
	ProyectoAcademico       string `xml:"proyectoAcademico"`
	CodigoEspacioAcademico  string `xml:"codigoEspacioAcademico"`
	EspacioAcademico        string `xml:"espacioAcademico"`
	Grupo                   string `xml:"grupo"`
}

type informacionCursoXML struct {
	Detalle detalleCursoIdXML `xml:"detalle"`
}

type coordinadorUsuarioXML struct {
	CodigoCarrera      string                 `xml:"codigo_carrera"`
	Coordinadores      []coordinadorUsuarioID `xml:"coordinador"`
	CoordinadorUsuario []coordinadorUsuarioID `xml:"coordinador_usuario"`
}

type coordinadorUsuarioID struct {
	CodigoCarrera string `xml:"codigo_carrera"`
}

type informacionHorariosXML struct {
	Horarios []horarioXML `xml:"horario"`
}

type horarioXML struct {
	IdEdificio             string `xml:"id_edificio"`
	Periodo                string `xml:"periodo"`
	ActivoEspacioAcademico string `xml:"activo_espacio_academico"`
	HoraInicio             string `xml:"hora_inicio"`
	CantidadHoras          string `xml:"cantidad_horas"`
	Grupo                  string `xml:"grupo"`
	IdHorario              string `xml:"id_horario"`
	ActivoHorario          string `xml:"activo_horario"`
	IdSalon                string `xml:"id_salon"`
	IdEspacioAcademico     string `xml:"id_espacio_academico"`
	IdSede                 string `xml:"id_sede"`
	IdEspacioFisico        string `xml:"id_espacio_fisico"`
	HoraFin                string `xml:"hora_fin"`
	NombreEspacioAcademico string `xml:"nombre_espacio_academico"`
	DiaSemana              string `xml:"dia_semana"`
}

func obtenerProyectosCurricularesCoordinador(documento string) ([]string, error) {
	url := "http://" + beego.AppConfig.String("AcademicaEspacioAcademicoService") +
		"coordinador_usuario/" + documento

	var responseXML coordinadorUsuarioXML
	if err := request.GetXml(url, &responseXML); err != nil {
		return nil, err
	}

	proyectos := []string{}
	proyectoExiste := map[string]bool{}

	agregarProyecto := func(codigo string) {
		codigoCarrera := strings.TrimSpace(codigo)
		if codigoCarrera != "" && !proyectoExiste[codigoCarrera] {
			proyectos = append(proyectos, codigoCarrera)
			proyectoExiste[codigoCarrera] = true
		}
	}

	agregarProyecto(responseXML.CodigoCarrera)

	for _, coordinador := range responseXML.CoordinadorUsuario {
		agregarProyecto(coordinador.CodigoCarrera)
	}

	for _, coordinador := range responseXML.Coordinadores {
		agregarProyecto(coordinador.CodigoCarrera)
	}

	if len(proyectos) == 0 {
		return nil, fmt.Errorf("no se encontro codigo_carrera para el coordinador")
	}

	return proyectos, nil
}

// ListaEspaciosAcademicosProyectoPeriodo consulta los espacios academicos
// en academica_pruebas para un anio, periodo y proyecto curricular.
func ListaEspaciosAcademicosProyectoPeriodo(anio, periodo, proyecto, documentoCoordinador string) requestmanager.APIResponse {
	proyectoConsulta := strings.TrimSpace(proyecto)
	proyectosConsulta := []string{}

	if proyectoConsulta == "" {
		if strings.TrimSpace(documentoCoordinador) == "" {
			return requestmanager.APIResponseDTO(false, 400, nil, "Error: Parametro(s) con valores no validos")
		}

		var err error
		proyectosConsulta, err = obtenerProyectosCurricularesCoordinador(documentoCoordinador)
		if err != nil {
			logs.Error(err)
			return requestmanager.APIResponseDTO(false, 404, nil, "No se encontro proyecto curricular del coordinador")
		}
	} else {
		proyectosConsulta = append(proyectosConsulta, proyectoConsulta)
	}

	response := []map[string]interface{}{}
	espaciosIds := map[string]bool{}

	for _, proyectoItem := range proyectosConsulta {
		url := "http://" + beego.AppConfig.String("AcademicaEspacioAcademicoService") +
			"espacios_academicos_proyecto_periodo/" + anio + "/" + periodo + "/" + proyectoItem

		var responseXML espaciosAcademicosXML
		if err := request.GetXml(url, &responseXML); err != nil {
			logs.Error(err)
			continue
		}

		for _, espacio := range responseXML.Espacios {
			if espaciosIds[espacio.Id] {
				continue
			}

			espaciosIds[espacio.Id] = true
			response = append(response, map[string]interface{}{
				"_id":             espacio.Id,
				"codigo":          espacio.Id,
				"nombre":          espacio.Nombre,
				"espacio_modular": espacio.EspacioModular,
				"codigo_carrera":  espacio.CodigoCarrera,
				"nombre_carrera":  espacio.NombreCarrera,
			})
		}
	}

	if len(response) == 0 {
		return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de espacios academicos")
	}

	return requestmanager.APIResponseDTO(true, 200, response)
}

// ListaGruposEspacioPeriodo consulta grupos de un espacio academico
// por anio y periodo en academica_pruebas.
func ListaGruposEspacioPeriodo(anio, periodo, espacio string) requestmanager.APIResponse {
	url := "http://" + beego.AppConfig.String("AcademicaEspacioAcademicoService") +
		"grupos_espacio_periodo/" + anio + "/" + periodo + "/" + espacio

	var responseXML gruposEspacioPeriodoXML
	if err := request.GetXml(url, &responseXML); err != nil {
		logs.Error(err)
		return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de grupos")
	}

	response := make([]map[string]interface{}, 0, len(responseXML.Grupos))
	for _, grupo := range responseXML.Grupos {
		response = append(response, map[string]interface{}{
			"Id":                grupo.Id,
			"Nombre":            grupo.Nombre,
			"ProyectoAcademico": grupo.ProyectoAcademico,
			"Nivel":             grupo.Nivel,
			"grupo":             grupo.Grupo,
		})
	}

	return requestmanager.APIResponseDTO(true, 200, response)
}

// DetalleCursoId consulta el detalle de un curso por CUR_ID
func DetalleCursoId(id string) requestmanager.APIResponse {
	url := "http://" + beego.AppConfig.String("AcademicaEspacioAcademicoService") +
		"informacion_curso/" + id

	var responseXML informacionCursoXML
	if err := request.GetXml(url, &responseXML); err != nil {
		logs.Error(err)
		return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de cursos")
	}

	detalle := responseXML.Detalle
	if strings.TrimSpace(detalle.Id) == "" {
		return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de cursos")
	}

	response := map[string]interface{}{
		"id":                      detalle.Id,
		"Nivel":                   detalle.Nivel,
		"CodigoProyectoAcademico": detalle.CodigoProyectoAcademico,
		"ProyectoAcademico":       detalle.ProyectoAcademico,
		"CodigoEspacioAcademico":  detalle.CodigoEspacioAcademico,
		"EspacioAcademico":        detalle.EspacioAcademico,
		"grupo":                   detalle.Grupo,
	}

	return requestmanager.APIResponseDTO(true, 200, response)
}

func parseNumber(value string) float64 {
	parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		return 0
	}
	return parsed
}

func formatHour(value float64) string {
	totalMinutes := int(math.Round(value * 60))
	hour := totalMinutes / 60
	minute := totalMinutes % 60
	return fmt.Sprintf("%02d:%02d", hour, minute)
}

func parseBool(value string) bool {
	normalized := strings.ToLower(strings.TrimSpace(value))
	return normalized == "true" || normalized == "1"
}

// InformacionHorarios consulta el endpoint de academica y transforma la respuesta
// al formato requerido por colocacion-espacio-academico para el cliente.
func InformacionHorarios(anio, periodo, asignaturaID, grupoID string) requestmanager.APIResponse {
	url := "http://" + beego.AppConfig.String("AcademicaEspacioAcademicoService") +
		"informacion_horarios/" + anio + "/" + periodo + "/" + asignaturaID + "/" + grupoID

	var responseXML informacionHorariosXML
	if err := request.GetXml(url, &responseXML); err != nil {
		logs.Error(err)
		return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron horarios para los parametros consultados")
	}

	response := make([]map[string]interface{}, 0, len(responseXML.Horarios))
	for _, horario := range responseXML.Horarios {
		if !parseBool(horario.ActivoHorario) || !parseBool(horario.ActivoEspacioAcademico) {
			continue
		}

		diaSemana := parseNumber(horario.DiaSemana)
		horaInicio := parseNumber(horario.HoraInicio)
		cantidadHoras := parseNumber(horario.CantidadHoras)
		horaFin := parseNumber(horario.HoraFin)
		if horaFin == 0 {
			horaFin = horaInicio + cantidadHoras
		}

		posX := (diaSemana - 1) * 110
		posY := ((horaInicio*60 - 360) / 15) * 22.5

		colocacion := map[string]interface{}{
			"dragPosition": map[string]interface{}{
				"x": posX,
				"y": posY,
			},
			"estado": 2,
			"finalPosition": map[string]interface{}{
				"x": posX,
				"y": posY,
			},
			"horaFormato": fmt.Sprintf("%s - %s", formatHour(horaInicio), formatHour(horaFin)),
			"horas":       cantidadHoras,
			"prevPosition": map[string]interface{}{
				"x": posX,
				"y": posY,
			},
			"tipo": 1,
		}

		resumenColocacion := map[string]interface{}{
			"colocacion": colocacion,
			"espacio_fisico": map[string]interface{}{
				"edificio_id": horario.IdEdificio,
				"salon_id":    horario.IdSalon,
				"sede_id":     horario.IdSede,
				"sede": map[string]interface{}{
					"Id":                horario.IdSede,
					"CodigoAbreviacion": horario.IdSede,
					"Nombre":            horario.IdSede,
				},
				"edificio": map[string]interface{}{
					"Id":     horario.IdEdificio,
					"Nombre": horario.IdEdificio,
				},
				"salon": map[string]interface{}{
					"Id":     horario.IdSalon,
					"Nombre": horario.IdSalon,
				},
			},
		}

		espacioAcademico := map[string]interface{}{
			"_id":                     horario.IdEspacioAcademico,
			"activo":                  parseBool(horario.ActivoEspacioAcademico),
			"espacio_academico_padre": nil,
			"grupo":                   horario.Grupo,
			"nombre":                  horario.NombreEspacioAcademico,
		}

		response = append(response, map[string]interface{}{
			"_id":                            horario.IdHorario,
			"EspacioAcademicoId":             horario.IdEspacioAcademico,
			"EspacioFisicoId":                horario.IdEspacioFisico,
			"ColocacionEspacioAcademico":     colocacion,
			"ResumenColocacionEspacioFisico": resumenColocacion,
			"Periodo":                        horario.Periodo,
			"Activo":                         parseBool(horario.ActivoHorario),
			"EspacioAcademico":               espacioAcademico,
		})
	}

	return requestmanager.APIResponseDTO(true, 200, response)
}

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
