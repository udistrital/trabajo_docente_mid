package test

import (
	"fmt"
	"strings"
	"testing"

	"bou.ke/monkey"
	"github.com/udistrital/trabajo_docente_mid/models"
	"github.com/udistrital/trabajo_docente_mid/services"
	"github.com/udistrital/utils_oas/request"
	requestmanager "github.com/udistrital/utils_oas/requestresponse"
)

const separadorSlash = "=================================================================================="

func TestListaAsignacion(t *testing.T) {
	t.Log(separadorSlash)
	t.Log("Inicio TestListaAsignacion")
	t.Log(separadorSlash)

	t.Run("Caso 1: Consulta exitosa con datos válidos", func(t *testing.T) {
		guardReq := monkey.Patch(request.GetJson, func(url string, target interface{}) error {
			if strings.Contains(url, "pre_asignacion?query=activo:true,aprobacion_docente:true,aprobacion_proyecto:true") {
				if data, ok := target.(*map[string]interface{}); ok {
					*data = map[string]interface{}{
						"Data": []interface{}{
							map[string]interface{}{
								"_id":                 "asig-1",
								"docente_id":          "123",
								"tipo_vinculacion_id": "1",
								"plan_docente_id":     "456",
								"periodo_id":          "2024",
							},
						},
					}
				}
				return nil
			}
			return fmt.Errorf("URL no mockeada: %s", url)
		})
		defer guardReq.Unpatch()

		guardDetalle := monkey.Patch(services.ConsultarDetalleAsignacion, func(asignaciones []interface{}, forTeacher bool) []map[string]interface{} {
			return []map[string]interface{}{
				{
					"plan_docente_id":     "456",
					"id":                  "asig-1",
					"docente_id":          "123",
					"docente":             "Juan Perez",
					"tipo_vinculacion_id": "1",
					"tipo_vinculacion":    "Tiempo completo",
					"identificacion":      "123456",
					"periodo_academico":   "2024-1",
					"periodo_id":          "2024",
					"estado":              "Aprobado",
					"codigo_estado":       "APR",
					"tiene_observaciones": false,
					"soporte_documental":  map[string]interface{}{"value": nil, "type": "ver", "disabled": true},
					"enviar":              map[string]interface{}{"value": nil, "type": "enviar", "disabled": true},
					"gestion":             map[string]interface{}{"value": nil, "type": "ver", "disabled": false},
				},
			}
		})
		defer guardDetalle.Unpatch()

		result := services.ListaAsignacion("2024")

		if !result.Success {
			t.Fatalf("Se esperaba Success=true, pero se obtuvo: %v", result.Success)
		}

		if result.Status != 200 {
			t.Errorf("Se esperaba Status=200, pero se obtuvo: %d", result.Status)
		}

		if result.Data == nil {
			t.Errorf("Se esperaba Data != nil, pero se obtuvo nil")
		}

		_, ok := result.Data.([]map[string]interface{})
		if !ok {
			t.Errorf("Se esperaba Data de tipo []map[string]interface{}, pero se obtuvo: %T", result.Data)
		}

		t.Log("Test de caso exitoso para ListaAsignacion ejecutado correctamente")
	})

	t.Run("Caso 2: Sin datos en la consulta", func(t *testing.T) {
		guardReq := monkey.Patch(request.GetJson, func(url string, target interface{}) error {
			if strings.Contains(url, "pre_asignacion?query=activo:true,aprobacion_docente:true,aprobacion_proyecto:true") {
				if data, ok := target.(*map[string]interface{}); ok {
					*data = map[string]interface{}{"Data": []interface{}{}}
				}
				return nil
			}
			return fmt.Errorf("URL no mockeada: %s", url)
		})
		defer guardReq.Unpatch()

		result := services.ListaAsignacion("2024")

		if result.Success {
			t.Fatalf("Se esperaba Success=false, pero se obtuvo: %v", result.Success)
		}

		if result.Status != 404 {
			t.Errorf("Se esperaba Status=404, pero se obtuvo: %d", result.Status)
		}

		if result.Data != nil {
			t.Errorf("Se esperaba Data=nil, pero se obtuvo: %v", result.Data)
		}

		msg, _ := result.Message.(string)
		if !strings.Contains(msg, "No se encontraron registros") {
			t.Errorf("Se esperaba mensaje con 'No se encontraron registros', pero se obtuvo: %s", msg)
		}

		t.Log("Test de caso sin datos para ListaAsignacion ejecutado correctamente")
	})

	t.Run("Caso 3: Error en la consulta", func(t *testing.T) {
		guardReq := monkey.Patch(request.GetJson, func(url string, target interface{}) error {
			return fmt.Errorf("error simulado en request")
		})
		defer guardReq.Unpatch()

		result := services.ListaAsignacion("2024")

		if result.Success {
			t.Fatalf("Se esperaba Success=false, pero se obtuvo: %v", result.Success)
		}

		if result.Status != 404 {
			t.Errorf("Se esperaba Status=404, pero se obtuvo: %d", result.Status)
		}

		msg, _ := result.Message.(string)
		if !strings.Contains(msg, "No se encontraron registros") {
			t.Errorf("Se esperaba mensaje con 'No se encontraron registros', pero se obtuvo: %s", msg)
		}

		t.Log("Test de caso error para ListaAsignacion ejecutado correctamente")
	})

	t.Log(separadorSlash)
	t.Log("Fin TestListaAsignacion")
	t.Log(separadorSlash)
}

func TestListaAsignacionDocente(t *testing.T) {
	t.Log(separadorSlash)
	t.Log("Inicio TestListaAsignacionDocente")
	t.Log(separadorSlash)

	t.Run("Caso 1: Consulta exitosa del docente", func(t *testing.T) {
		guardReq := monkey.Patch(request.GetJson, func(url string, target interface{}) error {
			if strings.Contains(url, "pre_asignacion?query=activo:true,aprobacion_docente:true,aprobacion_proyecto:true,docente_id:") {
				if data, ok := target.(*map[string]interface{}); ok {
					*data = map[string]interface{}{
						"Data": []interface{}{
							map[string]interface{}{
								"_id":                 "asig-1",
								"docente_id":          "123",
								"tipo_vinculacion_id": "1",
								"plan_docente_id":     "456",
								"periodo_id":          "2024",
							},
						},
					}
				}
				return nil
			}
			return fmt.Errorf("URL no mockeada: %s", url)
		})
		defer guardReq.Unpatch()

		guardDetalle := monkey.Patch(services.ConsultarDetalleAsignacion, func(asignaciones []interface{}, forTeacher bool) []map[string]interface{} {
			return []map[string]interface{}{
				{
					"plan_docente_id":     "456",
					"id":                  "asig-1",
					"docente_id":          "123",
					"docente":             "Juan Perez",
					"tipo_vinculacion_id": "1",
					"tipo_vinculacion":    "Tiempo completo",
					"identificacion":      "123456",
					"periodo_academico":   "2024-1",
					"periodo_id":          "2024",
					"estado":              "Aprobado",
					"codigo_estado":       "APR",
					"tiene_observaciones": false,
					"soporte_documental":  map[string]interface{}{"value": nil, "type": "ver", "disabled": true},
					"enviar":              map[string]interface{}{"value": nil, "type": "enviar", "disabled": false},
					"gestion":             map[string]interface{}{"value": nil, "type": "editar", "disabled": false},
				},
			}
		})
		defer guardDetalle.Unpatch()

		result := services.ListaAsignacionDocente("123", "2024")

		if !result.Success {
			t.Fatalf("Se esperaba Success=true, pero se obtuvo: %v", result.Success)
		}

		if result.Status != 200 {
			t.Errorf("Se esperaba Status=200, pero se obtuvo: %d", result.Status)
		}

		t.Log("Test de caso exitoso para ListaAsignacionDocente ejecutado correctamente")
	})

	t.Run("Caso 2: Sin datos del docente", func(t *testing.T) {
		guardReq := monkey.Patch(request.GetJson, func(url string, target interface{}) error {
			if strings.Contains(url, "pre_asignacion?query=activo:true,aprobacion_docente:true,aprobacion_proyecto:true,docente_id:") {
				if data, ok := target.(*map[string]interface{}); ok {
					*data = map[string]interface{}{"Data": []interface{}{}}
				}
				return nil
			}
			return fmt.Errorf("URL no mockeada: %s", url)
		})
		defer guardReq.Unpatch()

		result := services.ListaAsignacionDocente("123", "2024")

		if result.Success {
			t.Fatalf("Se esperaba Success=false, pero se obtuvo: %v", result.Success)
		}

		if result.Status != 404 {
			t.Errorf("Se esperaba Status=404, pero se obtuvo: %d", result.Status)
		}

		t.Log("Test de caso sin datos para ListaAsignacionDocente ejecutado correctamente")
	})

	t.Log(separadorSlash)
	t.Log("Fin TestListaAsignacionDocente")
	t.Log(separadorSlash)
}
