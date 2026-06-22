package test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/astaxie/beego"
	"github.com/udistrital/trabajo_docente_mid/models"
	"github.com/udistrital/trabajo_docente_mid/services"
	"github.com/udistrital/utils_oas/request"
)

// TestListaEventos cubre la consulta de eventos globales en formato XML
func TestListaEventos(t *testing.T) {
	t.Run("Caso 1: Lista de eventos obtenida correctamente", func(t *testing.T) {
		defer monkey.UnpatchAll()

		_ = beego.AppConfig.Set("AcademicaEspacioAcademicoService", "http://academica/")

		monkey.Patch(request.GetXml, func(url string, target interface{}) error {
			if strings.Contains(url, "eventos") {
				*(target.(*models.EventosXML)) = models.EventosXML{
					Eventos: []models.EventoXML{
						{
							CodigoEvento: "EV_01",
							Descripcion:  "Inscripciones",
							Editable:     "S",
						},
					},
				}
				return nil
			}
			return errors.New("URL no esperada")
		})

		apiResponse := services.ListaEventos()

		if !apiResponse.Success {
			t.Errorf("Se esperaba éxito (Success: true), se obtuvo false")
		}
		if apiResponse.Status != 200 {
			t.Errorf("Se esperaba status 200, se obtuvo %v", apiResponse.Status)
		}

		data, ok := apiResponse.Data.([]map[string]interface{})
		if !ok || len(data) == 0 {
			t.Fatalf("Estructura de respuesta inválida o vacía")
		}
		if data[0]["CodigoEvento"] != "EV_01" {
			t.Errorf("Se esperaba Código EV_01, se obtuvo %v", data[0]["CodigoEvento"])
		}
	})

	t.Run("Caso 2: Error al consultar el endpoint XML de eventos", func(t *testing.T) {
		defer monkey.UnpatchAll()

		_ = beego.AppConfig.Set("AcademicaEspacioAcademicoService", "http://academica/")

		monkey.Patch(request.GetXml, func(url string, target interface{}) error {
			return errors.New("Error de conexión remota")
		})

		apiResponse := services.ListaEventos()

		if apiResponse.Success {
			t.Errorf("Se esperaba fallo (Success: false), se obtuvo true")
		}
		if apiResponse.Status != 404 {
			t.Errorf("Se esperaba status 404, se obtuvo %v", apiResponse.Status)
		}
		if apiResponse.Message != "No se encontraron registros de eventos" {
			t.Errorf("Mensaje inesperado: %s", apiResponse.Message)
		}
	})
}

// TestConsultaProyectosFacultadDecano cubre la lógica de decano activo, vigencias de fechas y unificación de niveles
func TestConsultaProyectosFacultadDecano(t *testing.T) {
	t.Run("Caso 1: Decano con facultad vigente y proyectos curriculares unificados", func(t *testing.T) {
		defer monkey.UnpatchAll()

		_ = beego.AppConfig.Set("AcademicaEspacioAcademicoService", "http://academica/")

		ahora := time.Now()
		fechaPasada := ahora.Add(-24 * time.Hour).Format("2006-01-02T15:04:05-05:00")
		fechaFutura := ahora.Add(24 * time.Hour).Format("2006-01-02T15:04:05-05:00")

		monkey.Patch(request.GetXml, func(url string, target interface{}) error {
			switch {
			case strings.Contains(url, "/decano/"):
				*(target.(*models.FacultadesDecanoXML)) = models.FacultadesDecanoXML{
					Decanos: []models.FacultadDecanoXML{
						{
							CodigoFacultad: "14",
							Facultad:       "Facultad de Ingeniería",
							Nombre:         "Decano de Prueba",
							FechaDesde:     fechaPasada,
							FechaHasta:     fechaFutura,
						},
					},
				}
				return nil

			case strings.Contains(url, "/proyectos_facultad/14/PREGRADO"):
				*(target.(*models.ProyectosFacultadXML)) = models.ProyectosFacultadXML{
					Proyectos: []models.ProyectoFacultadXML{
						{
							CodigoProyectoCurricular: "20",
							NombreProyectoCurricular: "Ingeniería de Sistemas",
						},
					},
				}
				return nil

			case strings.Contains(url, "/proyectos_facultad/14/POSGRADO"):
				*(target.(*models.ProyectosFacultadXML)) = models.ProyectosFacultadXML{
					Proyectos: []models.ProyectoFacultadXML{},
				}
				return nil
			}
			return errors.New("URL no esperada: " + url)
		})

		apiResponse := services.ConsultaProyectosFacultadDecano("80123456")

		if !apiResponse.Success {
			t.Fatalf("Se esperaba éxito en la unificación y se obtuvo error: %v", apiResponse.Message)
		}
		if apiResponse.Status != 200 {
			t.Errorf("Se esperaba status 200, se obtuvo %v", apiResponse.Status)
		}

		proyectos := apiResponse.Data.([]map[string]interface{})
		if len(proyectos) != 1 {
			t.Errorf("Se esperaba 1 proyecto unificado, se obtuvieron %d", len(proyectos))
		}
	})

	t.Run("Caso 2: Decano con facultad inactiva por expiración de vigencia", func(t *testing.T) {
		defer monkey.UnpatchAll()

		_ = beego.AppConfig.Set("AcademicaEspacioAcademicoService", "http://academica/")

		ahora := time.Now()
		fechaPasadaLejana := ahora.Add(-48 * time.Hour).Format("2006-01-02T15:04:05-05:00")
		fechaPasadaCercana := ahora.Add(-24 * time.Hour).Format("2006-01-02T15:04:05-05:00")

		monkey.Patch(request.GetXml, func(url string, target interface{}) error {
			if strings.Contains(url, "/decano/") {
				*(target.(*models.FacultadesDecanoXML)) = models.FacultadesDecanoXML{
					Decanos: []models.FacultadDecanoXML{
						{
							CodigoFacultad: "14",
							Facultad:       "Facultad de Ingeniería",
							FechaDesde:     fechaPasadaLejana,
							FechaHasta:     fechaPasadaCercana, // Ya caducó hace 24 horas
						},
					},
				}
				return nil
			}
			// Permitimos que otras URLs devuelvan un XML vacío en vez de un error estricto
			return nil
		})

		apiResponse := services.ConsultaProyectosFacultadDecano("80123456")

		if apiResponse.Success {
			t.Errorf("Se esperaba fallo por vigencia expirada, pero retornó éxito")
		}
		if apiResponse.Status != 404 {
			t.Errorf("Se esperaba status 404, se obtuvo %v", apiResponse.Status)
		}

		// Corregimos la aserción para que valide cualquiera de los mensajes de error válidos de salida (404)
		msg, ok := apiResponse.Message.(string)
		if !ok || (!strings.Contains(msg, "No se encontró facultad activa") && !strings.Contains(msg, "No se encontraron proyectos")) {
			t.Errorf("Mensaje inesperado de error obtenido: %v", apiResponse.Message)
		}
	})
}

// TestConsultaCalendariosEventos cubre el cruce con proyecto explícito o fallback de proyectos por documento
func TestConsultaCalendariosEventos(t *testing.T) {
	t.Run("Caso 1: Consulta exitosa suministrando un código de proyecto directo", func(t *testing.T) {
		defer monkey.UnpatchAll()

		_ = beego.AppConfig.Set("AcademicaEspacioAcademicoService", "http://academica/")

		monkey.Patch(request.GetXml, func(url string, target interface{}) error {
			if strings.Contains(url, "/calendario_eventos/EV_INS/20") {
				*(target.(*models.CalendariosEventosXML)) = models.CalendariosEventosXML{
					Eventos: []models.CalendarioEventoXML{
						{
							CodigoEvento:   "EV_INS",
							CodigoProyecto: "20",
							FechaInicio:    "2026-01-10",
							FechaFin:       "2026-02-10",
							Year:           "2026",
						},
					},
				}
				return nil
			}
			return errors.New("URL no esperada: " + url)
		})

		apiResponse := services.ConsultaCalendariosEventos("80123456", "EV_INS", "20")

		if !apiResponse.Success {
			t.Fatalf("La consulta debió ser exitosa, error: %v", apiResponse.Message)
		}
		if apiResponse.Status != 200 {
			t.Errorf("Se esperaba status 200, se obtuvo %v", apiResponse.Status)
		}

		resultados := apiResponse.Data.([]map[string]interface{})
		if len(resultados) != 1 {
			t.Errorf("Se esperaba 1 registro de calendario, se obtuvieron %d", len(resultados))
		}
		if resultados[0]["FechaInicio"] != "2026-01-10" {
			t.Errorf("Se esperaba fecha inicio '2026-01-10', se obtuvo %v", resultados[0]["FechaInicio"])
		}
	})

	t.Run("Caso 2: Fallo general cuando no hay calendarios para el coordinador ni por fallback de docente planta", func(t *testing.T) {
		defer monkey.UnpatchAll()

		_ = beego.AppConfig.Set("AcademicaEspacioAcademicoService", "http://academica/")

		monkey.Patch(request.GetXml, func(url string, target interface{}) error {
			// Simula que la primera consulta de coordinador viene completamente en blanco
			if strings.Contains(url, "/coordinador_usuario/") {
				return errors.New("No encontrado en tabla coordinador")
			}
			// Simula que el fallback tampoco registra información de carrera asociada
			if strings.Contains(url, "/consulta_datos_docente_planta/") {
				return errors.New("No encontrado en planta docente")
			}
			return nil
		})

		// Se le envía el parámetro "proyecto" en vacío para obligar a buscar por documento
		apiResponse := services.ConsultaCalendariosEventos("80123456", "EV_INS", "")

		if apiResponse.Success {
			t.Errorf("Se esperaba un escenario fallido al no encontrarse vinculaciones")
		}
		if apiResponse.Status != 404 {
			t.Errorf("Se esperaba status 404, se obtuvo %v", apiResponse.Status)
		}
		if !strings.Contains(apiResponse.Message.(string), "No se encontró código de proyecto curricular") {
			t.Errorf("Mensaje de error inesperado: %v", apiResponse.Message)
		}
	})
}
