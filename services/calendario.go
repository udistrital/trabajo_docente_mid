package services

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/trabajo_docente_mid/models"
	request "github.com/udistrital/utils_oas/request"
	requestmanager "github.com/udistrital/utils_oas/requestresponse"
)

// ListaEventos consulta los eventos en academica
func ListaEventos() requestmanager.APIResponse {
	url := beego.AppConfig.String("AcademicaEspacioAcademicoService") + "eventos"

	var responseXML models.EventosXML
	logs.Info("Consultando endpoint de eventos: ", url)
	if err := request.GetXml(url, &responseXML); err != nil {
		logs.Error("Error en ListaEventos: ", err)
		return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de eventos")
	}

	logs.Info("Cantidad de eventos obtenidos: ", len(responseXML.Eventos))

	response := make([]map[string]interface{}, 0, len(responseXML.Eventos))
	for _, evento := range responseXML.Eventos {
		response = append(response, map[string]interface{}{
			"Descripcion":  evento.Descripcion,
			"CodigoEvento": evento.CodigoEvento,
			"Editable":     evento.Editable,
		})
	}

	return requestmanager.APIResponseDTO(true, 200, response)
}

// ConsultaCalendariosEventos cruza la información de un evento con los proyectos curriculares del coordinador o docente
func ConsultaCalendariosEventos(documento string, codigoEvento string) requestmanager.APIResponse {
	proyectos, err := ObtenerDetalleProyectosCurriculares(documento)
	if err != nil {
		logs.Error("Error al obtener proyectos curriculares del coordinador o docente: ", err)
		return requestmanager.APIResponseDTO(false, 404, nil, "No se encontró código de proyecto curricular para el coordinador o docente")
	}
	logs.Info("Proyectos curriculares encontrados para el coordinador o docente: ", proyectos)

	var resultados []map[string]interface{}
	urlCalendarioBase := beego.AppConfig.String("AcademicaEspacioAcademicoService") + "calendario_eventos/"

	for _, proyectoData := range proyectos {
		proyecto := proyectoData["codigo_carrera"].(string)
		nombreProyecto := proyectoData["nombre_carrera"].(string)

		urlConsulta := urlCalendarioBase + codigoEvento + "/" + proyecto
		logs.Info("Consultando endpoint calendario_eventos: ", urlConsulta)
		var respuestaCalendario models.CalendariosEventosXML
		if err := request.GetXml(urlConsulta, &respuestaCalendario); err == nil {
			logs.Info("Cantidad de calendarios obtenidos para el proyecto ", proyecto, ": ", len(respuestaCalendario.Eventos))
			for _, cal := range respuestaCalendario.Eventos {
				resultados = append(resultados, map[string]interface{}{
					"CodigoEvento":   cal.CodigoEvento,
					"FechaInicio":    cal.FechaInicio,
					"TipoProyecto":   cal.TipoProyecto,
					"Ciclo":          cal.Ciclo,
					"CodigoProyecto": cal.CodigoProyecto,
					"NombreProyecto": nombreProyecto,
					"FechaFin":       cal.FechaFin,
					"CodigoFacultad": cal.CodigoFacultad,
					"Year":           cal.Year,
				})
			}
		} else {
			logs.Error("Error consultando calendario_eventos: ", err)
		}
	}

	logs.Info("Total de resultados unificados de calendarios: ", len(resultados))

	if len(resultados) == 0 {
		return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron calendarios de eventos para este coordinador/docente y código de evento")
	}

	return requestmanager.APIResponseDTO(true, 200, resultados)
}

type proyectosCurricularesXML struct {
	CodigoCarrera      string                  `xml:"codigo_carrera"`
	NombreCarrera      string                  `xml:"nombre_carrera"`
	Coordinadores      []proyectoCurricularXML `xml:"coordinador"`
	CoordinadorUsuario []proyectoCurricularXML `xml:"coordinador_usuario"`
}

type proyectoCurricularXML struct {
	CodigoCarrera string `xml:"codigo_carrera"`
	NombreCarrera string `xml:"nombre_carrera"`
}

type datosCollectionXML struct {
	Datos []datoDocentePlantaXML `xml:"datos"`
}

type datoDocentePlantaXML struct {
	CodigoCarrera string `xml:"codigo_proyecto"`
	NombreCarrera string `xml:"proyecto"`
}

// ObtenerDetalleProyectosCurriculares obtiene los proyectos con código y nombre asociados a un documento (usuario)
func ObtenerDetalleProyectosCurriculares(documento string) ([]map[string]interface{}, error) {
	url := beego.AppConfig.String("AcademicaEspacioAcademicoService") +
		"coordinador_usuario/" + documento

	var responseXML proyectosCurricularesXML
	if err := request.GetXml(url, &responseXML); err != nil {
		return nil, err
	}

	proyectos := []map[string]interface{}{}
	proyectoExiste := map[string]bool{}

	agregarProyecto := func(codigo, nombre string) {
		codigoCarrera := strings.TrimSpace(codigo)
		if codigoCarrera != "" && !proyectoExiste[codigoCarrera] {
			proyectos = append(proyectos, map[string]interface{}{
				"codigo_carrera": codigoCarrera,
				"nombre_carrera": strings.TrimSpace(nombre),
			})
			proyectoExiste[codigoCarrera] = true
		}
	}

	agregarProyecto(responseXML.CodigoCarrera, responseXML.NombreCarrera)

	for _, proyecto := range responseXML.CoordinadorUsuario {
		agregarProyecto(proyecto.CodigoCarrera, proyecto.NombreCarrera)
	}

	for _, proyecto := range responseXML.Coordinadores {
		agregarProyecto(proyecto.CodigoCarrera, proyecto.NombreCarrera)
	}

	if len(proyectos) == 0 {
		urlFallback := beego.AppConfig.String("AcademicaEspacioAcademicoService") + "consulta_datos_docente_planta/" + documento
		var responseFallback datosCollectionXML

		if err := request.GetXml(urlFallback, &responseFallback); err == nil {
			for _, dato := range responseFallback.Datos {
				agregarProyecto(dato.CodigoCarrera, dato.NombreCarrera)
			}
		}

		if len(proyectos) == 0 {
			return nil, fmt.Errorf("no se encontro codigo_carrera para el usuario")
		}
	}

	//ESPACIO PARA EN EL FUTURO IMPLEMENTAR LA BUSQUEDA DE PROYECTOS CURRICULARES DE UN DECANO
	//

	return proyectos, nil
}
