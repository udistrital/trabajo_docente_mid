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

func TestEnviarAsignacionCoordinador(t *testing.T) {
	t.Log(separadorSlash)
	t.Log("Inicio TestEnviarAsignacionCoordinador")
	t.Log(separadorSlash)

	t.Run("Caso 1: Envío exitoso de asignación", func(t *testing.T) {
		guardReq := monkey.Patch(request.GetJson, func(url string, target interface{}) error {
			switch {
			case strings.Contains(url, "plan_docente/456"):
				if data, ok := target.(*map[string]interface{}); ok {
					*data = map[string]interface{}{
						"Success": true,
						"Data": map[string]interface{}{
							"_id":                 "456",
							"docente_id":          "123",
							"periodo_id":          "2024",
							"tipo_vinculacion_id": "1",
							"estado_plan_id":      "ENV_DOC",
						},
					}
				}
				return nil
			case strings.Contains(url, "estado_plan?query=codigo_abreviacion:ENV_DOC"):
				if data, ok := target.(*map[string]interface{}); ok {
					*data = map[string]interface{}{
						"Data": []interface{}{
							map[string]interface{}{
								"_id":                "envdocid",
								"codigo_abreviacion": "ENV_DOC",
								"nombre":             "Enviado a Docente",
							},
						},
					}
				}
				return nil
			case strings.Contains(url, "pre_asignacion?query=activo:true,aprobacion_docente:true,aprobacion_proyecto:true,plan_docente_id:"):
				if data, ok := target.(*map[string]interface{}); ok {
					*data = map[string]interface{}{
						"Data": []interface{}{
							map[string]interface{}{
								"espacio_academico_id": "esp-1",
								"_id":                  "preasig-1",
								"activo":               true,
							},
						},
					}
				}
				return nil
			case strings.Contains(url, "carga_plan?query=activo:true,plan_docente_id:"):
				if data, ok := target.(*map[string]interface{}); ok {
					*data = map[string]interface{}{
						"Data": []interface{}{
							map[string]interface{}{
								"Id":                              "carga-1",
								"Espacio_academico_id":            "esp-1",
								"Actividad_id":                    "act-1",
								"Plan_docente_id":                 "456",
								"Hora_inicio":                     8,
								"Duracion":                        2,
								"Salon_id":                        "sal-1",
								"Sede_id":                         "sed-1",
								"Edificio_id":                     "edif-1",
								"Activo":                          true,
								"Colocacion_espacio_academico_id": "",
							},
						},
					}
				}
				return nil
			}
			return nil
		})
		defer guardReq.Unpatch()

		guardConvertir := monkey.Patch(services.ConvertirCargaPlanParaDefinePTD, func(carga models.CargaPlan) (map[string]interface{}, error) {
			return map[string]interface{}{
				"id":                   carga.Id,
				"espacio_academico_id": carga.Espacio_academico_id,
				"actividad_id":         carga.Actividad_id,
				"plan_docente_id":      carga.Plan_docente_id,
				"hora_inicio":          carga.Hora_inicio,
				"duracion":             carga.Duracion,
				"salon_id":             carga.Salon_id,
				"sede_id":              carga.Sede_id,
				"edificio_id":          carga.Edificio_id,
				"activo":               carga.Activo,
				"horario": map[string]interface{}{
					"hora_inicio": carga.Hora_inicio,
					"duracion":    carga.Duracion,
				},
			}, nil
		})
		defer guardConvertir.Unpatch()

		guardDefine := monkey.Patch(services.DefinePTD, func(body map[string]interface{}) requestmanager.APIResponse {
			return requestmanager.APIResponseDTO(true, 200, map[string]interface{}{"resultado": "exitoso"})
		})
		defer guardDefine.Unpatch()

		result := services.EnviarAsignacionCoordinador(map[string]interface{}{"plan_docente_id": "456"})

		if !result.Success {
			t.Fatalf("Se esperaba Success=true, pero se obtuvo: %v", result.Success)
		}

		if result.Status != 200 {
			t.Errorf("Se esperaba Status=200, pero se obtuvo: %d", result.Status)
		}

		t.Log("Test de caso exitoso para EnviarAsignacionCoordinador ejecutado correctamente")
	})

	t.Run("Caso 2: plan_docente_id inválido", func(t *testing.T) {
		result := services.EnviarAsignacionCoordinador(map[string]interface{}{"plan_docente_id": ""})

		if result.Success {
			t.Fatalf("Se esperaba Success=false, pero se obtuvo: %v", result.Success)
		}

		if result.Status != 400 {
			t.Errorf("Se esperaba Status=400, pero se obtuvo: %d", result.Status)
		}

		msg, _ := result.Message.(string)
		if !strings.Contains(msg, "plan_docente_id inválido") {
			t.Errorf("Se esperaba mensaje con 'plan_docente_id inválido', pero se obtuvo: %s", msg)
		}

		t.Log("Test de parámetro inválido para EnviarAsignacionCoordinador ejecutado correctamente")
	})

	t.Run("Caso 3: Error al consultar plan docente", func(t *testing.T) {
		defer monkey.Unpatch(request.GetJson)
		monkey.Patch(request.GetJson, func(url string, target interface{}) error {
			return fmt.Errorf("error simulado en request")
		})

		result := services.EnviarAsignacionCoordinador(map[string]interface{}{"plan_docente_id": "456"})

		if result.Success {
			t.Fatalf("Se esperaba Success=false, pero se obtuvo: %v", result.Success)
		}

		if result.Status != 404 {
			t.Errorf("Se esperaba Status=404, pero se obtuvo: %d", result.Status)
		}

		msg, _ := result.Message.(string)
		if !strings.Contains(msg, "No se pudo consultar el plan docente") {
			t.Errorf("Se esperaba mensaje con 'No se pudo consultar el plan docente', pero se obtuvo: %s", msg)
		}

		t.Log("Test de error en plan docente para EnviarAsignacionCoordinador ejecutado correctamente")
	})

	t.Log(separadorSlash)
	t.Log("Fin TestEnviarAsignacionCoordinador")
	t.Log(separadorSlash)
}
