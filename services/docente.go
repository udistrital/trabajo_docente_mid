package services

import (
	"fmt"
	"sync"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/trabajo_docente_mid/utils"
	request "github.com/udistrital/utils_oas/request"
	requestmanager "github.com/udistrital/utils_oas/requestresponse"
)

// DocumentoDocenteVinculacion ...
func ListaDocentesxDocumentoVinculacion(documento string, vinculacion int64) requestmanager.APIResponse {
	resVinculacion := []interface{}{}
	resDocumento := []interface{}{}
	response := []interface{}{}

	if errDocumento := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=Activo:true,Numero:"+documento+"&fields=TerceroId", &resDocumento); errDocumento == nil {
		if fmt.Sprintf("%v", resDocumento) != "[map[]]" {
			for _, documentoGet := range resDocumento {
				if errVinculacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"vinculacion?query=Activo:true,"+
					"TipoVinculacionId:"+fmt.Sprintf("%d", vinculacion)+",TerceroPrincipalId.Id:"+fmt.Sprintf("%v", documentoGet.(map[string]interface{})["TerceroId"].(map[string]interface{})["Id"])+"&fields=TerceroPrincipalId", &resVinculacion); errVinculacion == nil {
					if fmt.Sprintf("%v", resVinculacion) != "[map[]]" {
						response = append(response, map[string]interface{}{
							"Nombre":    utils.Capitalize(resVinculacion[0].(map[string]interface{})["TerceroPrincipalId"].(map[string]interface{})["NombreCompleto"].(string)),
							"Documento": documento,
							"Id":        resVinculacion[0].(map[string]interface{})["TerceroPrincipalId"].(map[string]interface{})["Id"]})
						/* c.Ctx.Output.SetStatus(200)
						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": response} */
					} else {
						return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de docente")
						/* c.Ctx.Output.SetStatus(404)
						c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de docente"} */
					}
				}
			}
		} else {
			return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de docente")
			/* c.Ctx.Output.SetStatus(404)
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de docente"} */
		}
	} else {
		logs.Error(errDocumento)
		return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de docentes")
		/* c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de docentes"} */
	}

	return requestmanager.APIResponseDTO(true, 200, response)
}

// NombreDocenteVinculacion ...
func ListaDocentesxNombreVinculacion(nombre string, vinculacion int64) requestmanager.APIResponse {
	resVinculacion := []interface{}{}
	resDocumento := []interface{}{}
	response := []interface{}{}

	if errVinculacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"vinculacion?limit=0&query=TipoVinculacionId__in:"+fmt.Sprintf("%d", vinculacion)+
		",Activo:true,TerceroPrincipalId.NombreCompleto__icontains:"+nombre+"&fields=TerceroPrincipalId", &resVinculacion); errVinculacion == nil {
		if fmt.Sprintf("%v", resVinculacion) != "[map[]]" {
			var tercerosIds string
			for _, vinculacion := range resVinculacion {
				tercerosIds += fmt.Sprintf("%v", vinculacion.(map[string]interface{})["TerceroPrincipalId"].(map[string]interface{})["Id"]) + "|"
			}
			tercerosIds = tercerosIds[:len(tercerosIds)-1]

			if errDocumento := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=Activo:true,TerceroId__in:"+tercerosIds+"&fields=Numero,TerceroId", &resDocumento); errDocumento == nil {
				for _, vinculacion := range resVinculacion {
					for indexDocumento, documento := range resDocumento {

						if vinculacion.(map[string]interface{})["TerceroPrincipalId"].(map[string]interface{})["Id"] == documento.(map[string]interface{})["TerceroId"].(map[string]interface{})["Id"] {
							response = append(response, map[string]interface{}{
								"Nombre":    utils.Capitalize(vinculacion.(map[string]interface{})["TerceroPrincipalId"].(map[string]interface{})["NombreCompleto"].(string)),
								"Documento": resDocumento[0].(map[string]interface{})["Numero"],
								"Id":        vinculacion.(map[string]interface{})["TerceroPrincipalId"].(map[string]interface{})["Id"]})
							resDocumento = append(resDocumento[:indexDocumento], resDocumento[indexDocumento+1:]...)
							break
						}
					}
				}
			} else {
				logs.Error(errDocumento)
				return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de docentes")
				/* c.Ctx.Output.SetStatus(404)
				c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de docentes"} */
			}
			return requestmanager.APIResponseDTO(true, 200, response)
			/* c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": response} */
		} else {
			return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de docentes")
			/* c.Ctx.Output.SetStatus(404)
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de docentes"} */
		}
	} else {
		logs.Error(errVinculacion)
		return requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de docentes")
		/* c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de docentes"} */
	}
}

// Busca el docente por documento y sus vinculaciones asociadas de tipo docente
func BuscarDocentesPorDocumentoConVinculaciones(documento string) (response requestmanager.APIResponse) {
	// 0. Inicialización de variables
	var docentes []map[string]interface{}
	response = requestmanager.APIResponseDTO(false, 404, docentes, "No se encontraron registros de docentes")
	idsVinculacionTipoDocente := []int64{293, 294, 296, 297, 298, 299}

	// 1. Busca los terceros relacionados al documento
	resTercero := []interface{}{}
	errResTercero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=Activo:true,Numero:"+documento, &resTercero)
	if errResTercero != nil {
		logs.Error(errResTercero)
		response = requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de docentes")
		return response
	}
	if fmt.Sprint(resTercero) == "[map[]]" {
		response = requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de docentes")
		return response
	}

	for _, tercero := range resTercero {
		// 1.1 Obtener el id del tercero
		terceroId := fmt.Sprintf("%v", tercero.(map[string]interface{})["TerceroId"].(map[string]interface{})["Id"])

		// 2. Buscar las vinculaciones del tercero por id del tercero
		resVinculacion := []interface{}{}
		errResVinculacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"vinculacion?query=Activo:true,tercero_principal_id:"+terceroId+"&fields=TipoVinculacionId", &resVinculacion)
		if errResVinculacion != nil {
			logs.Error(errResVinculacion)
			continue // Continuar con el siguiente tercero
		}
		if fmt.Sprint(resVinculacion) == "[map[]]" {
			continue // Continuar con el siguiente tercero
		}

		// 3. Buscar el tipo de vinculacion docente
		var docenteVinculacionesIds []int64
		for _, vinculacion := range resVinculacion {
			vinculacionId := fmt.Sprintf("%v", vinculacion.(map[string]interface{})["TipoVinculacionId"])

			// 4. Verificar si la vinculacion es de tipo docente
			for _, idVinculacionTipoDocente := range idsVinculacionTipoDocente {
				if vinculacionId == fmt.Sprint(idVinculacionTipoDocente) {
					docenteVinculacionesIds = append(docenteVinculacionesIds, idVinculacionTipoDocente)
				}
			}
		}

		// 5. Si el tercero tiene vinculaciones de tipo docente
		if len(docenteVinculacionesIds) > 0 {
			docente := map[string]interface{}{
				"Nombre":        utils.Capitalize(tercero.(map[string]interface{})["TerceroId"].(map[string]interface{})["NombreCompleto"].(string)),
				"Documento":     documento,
				"Id":            tercero.(map[string]interface{})["TerceroId"].(map[string]interface{})["Id"],
				"Vinculaciones": docenteVinculacionesIds,
			}
			docentes = append(docentes, docente)
		}
	}

	// Si se encontraron docentes con vinculaciones, actualiza la respuesta
	if len(docentes) > 0 {
		response = requestmanager.APIResponseDTO(true, 200, docentes, "Registros de docentes encontrados")
	}

	return response
}

// BuscarDocentesPorNombreConVinculaciones busca docentes por nombre y sus vinculaciones de tipo docente
func BuscarDocentesPorNombreConVinculaciones(nombre string) (response requestmanager.APIResponse) {
	// Inicialización de variables
	var docentes []map[string]interface{}
	response = requestmanager.APIResponseDTO(false, 404, docentes, "No se encontraron registros de docentes")
	idsVinculacionTipoDocente := []int64{293, 294, 296, 297, 298, 299}

	// 1. Busca los terceros relacionados al nombre
	resTercero := []interface{}{}
	errResTercero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=Activo:true,TerceroId__NombreCompleto__icontains:"+nombre+"&limit=0", &resTercero)
	if errResTercero != nil {
		logs.Error(errResTercero)
		response = requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de docentes")
		return response
	}
	if fmt.Sprint(resTercero) == "[map[]]" {
		response = requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de docentes")
		return response
	}

	for index, tercero := range resTercero {
		fmt.Println(index)
		// 1.1 Obtener el id del tercero
		terceroId := fmt.Sprintf("%v", tercero.(map[string]interface{})["TerceroId"].(map[string]interface{})["Id"])

		// 2. Buscar las vinculaciones del tercero por id del tercero
		resVinculacion := []interface{}{}
		errResVinculacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"vinculacion?query=Activo:true,tercero_principal_id:"+terceroId+"&fields=TipoVinculacionId", &resVinculacion)
		if errResVinculacion != nil {
			logs.Error(errResVinculacion)
			continue // Continuar con el siguiente tercero
		}
		if fmt.Sprint(resVinculacion) == "[map[]]" {
			continue // Continuar con el siguiente tercero
		}

		// 3. Buscar el tipo de vinculacion docente
		var docenteVinculacionesIds []int64
		for _, vinculacion := range resVinculacion {
			vinculacionId := fmt.Sprintf("%v", vinculacion.(map[string]interface{})["TipoVinculacionId"])

			// 4. Verificar si la vinculacion es de tipo docente
			for _, idVinculacionTipoDocente := range idsVinculacionTipoDocente {
				if vinculacionId == fmt.Sprint(idVinculacionTipoDocente) {
					docenteVinculacionesIds = append(docenteVinculacionesIds, idVinculacionTipoDocente)
				}
			}
		}

		// 5. Si el tercero tiene vinculaciones de tipo docente
		if len(docenteVinculacionesIds) > 0 {
			docente := map[string]interface{}{
				"Nombre":        utils.Capitalize(tercero.(map[string]interface{})["TerceroId"].(map[string]interface{})["NombreCompleto"].(string)),
				"Documento":     tercero.(map[string]interface{})["Numero"].(string),
				"Id":            tercero.(map[string]interface{})["TerceroId"].(map[string]interface{})["Id"],
				"Vinculaciones": docenteVinculacionesIds,
			}
			docentes = append(docentes, docente)
		}
	}

	// Si se encontraron docentes con vinculaciones, actualiza la respuesta
	if len(docentes) > 0 {
		response = requestmanager.APIResponseDTO(true, 200, docentes, "Registros de docentes encontrados")
	}

	return response
}

// BuscarDocentesPorNombreConVinculaciones busca docentes por nombre y sus vinculaciones de tipo docente
func BuscarDocentesPorNombreConVinculaciones2(nombre string) (response requestmanager.APIResponse) {
	// Inicialización de variables
	var docentes []map[string]interface{}
	idsVinculacionTipoDocente := []int64{293, 294, 296, 297, 298, 299}

	var wg sync.WaitGroup
	vinculacionesCh := make(chan map[string]interface{}, len(idsVinculacionTipoDocente))
	errCh := make(chan error, len(idsVinculacionTipoDocente))

	// Función auxiliar para buscar vinculaciones
	buscarVinculaciones := func(vinculacionID int64) {
		defer wg.Done()
		resVinculacion := []interface{}{}
		errResVinculacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"vinculacion?limit=0&query=TipoVinculacionId__in:"+fmt.Sprintf("%d", vinculacionID)+
			",Activo:true,TerceroPrincipalId.NombreCompleto__icontains:"+nombre+"&fields=TerceroPrincipalId,TipoVinculacionId", &resVinculacion)
		if errResVinculacion != nil {
			errCh <- errResVinculacion
			return
		}
		if fmt.Sprint(resVinculacion) == "[map[]]" {
			return
		}

		var tercerosIds string
		for _, vinculacion := range resVinculacion {
			tercerosIds += fmt.Sprintf("%v", vinculacion.(map[string]interface{})["TerceroPrincipalId"].(map[string]interface{})["Id"]) + "|"
		}
		if len(tercerosIds) > 0 {
			tercerosIds = tercerosIds[:len(tercerosIds)-1]
		}

		resDocumento := []interface{}{}
		errDocumento := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=Activo:true,TerceroId__in:"+tercerosIds+"&fields=Numero,TerceroId", &resDocumento)
		if errDocumento != nil {
			errCh <- errDocumento
			return
		}

		for _, vinculacion := range resVinculacion {
			for indexDocumento, documento := range resDocumento {
				if vinculacion.(map[string]interface{})["TerceroPrincipalId"].(map[string]interface{})["Id"] == documento.(map[string]interface{})["TerceroId"].(map[string]interface{})["Id"] {
					docente := map[string]interface{}{
						"Nombre":        utils.Capitalize(vinculacion.(map[string]interface{})["TerceroPrincipalId"].(map[string]interface{})["NombreCompleto"].(string)),
						"Documento":     documento.(map[string]interface{})["Numero"].(string),
						"Id":            vinculacion.(map[string]interface{})["TerceroPrincipalId"].(map[string]interface{})["Id"],
						"Vinculaciones": []int64{vinculacionID},
					}
					vinculacionesCh <- docente
					resDocumento = append(resDocumento[:indexDocumento], resDocumento[indexDocumento+1:]...)
					break
				}
			}
		}
	}

	for _, vinculacionID := range idsVinculacionTipoDocente {
		wg.Add(1)
		go buscarVinculaciones(vinculacionID)
	}

	go func() {
		wg.Wait()
		close(vinculacionesCh)
		close(errCh)
	}()

	docentesMap := make(map[string]map[string]interface{})

	for docente := range vinculacionesCh {
		id := fmt.Sprintf("%v", docente["Id"])
		if existingDocente, exists := docentesMap[id]; exists {
			existingDocente["Vinculaciones"] = append(existingDocente["Vinculaciones"].([]int64), docente["Vinculaciones"].([]int64)...)
		} else {
			docentesMap[id] = docente
		}
	}

	for _, docente := range docentesMap {
		docentes = append(docentes, docente)
	}

	if len(docentes) > 0 {
		response = requestmanager.APIResponseDTO(true, 200, docentes, "Registros de docentes encontrados")
	} else {
		select {
		case err := <-errCh:
			logs.Error(err)
			response = requestmanager.APIResponseDTO(false, 500, nil, "Error en la búsqueda de docentes")
		default:
			response = requestmanager.APIResponseDTO(false, 404, nil, "No se encontraron registros de docentes")
		}
	}

	return response
}
