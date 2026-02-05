package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/trabajo_docente_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	requestmanager "github.com/udistrital/utils_oas/requestresponse"
)

// DocenteController operations for Asignacion
type DocenteController struct {
	beego.Controller
}

// URLMapping ...
func (c *DocenteController) URLMapping() {
	c.Mapping("DocumentoDocenteVinculacion", c.DocumentoDocenteVinculacion)
	c.Mapping("NombreDocenteVinculacion", c.NombreDocenteVinculacion)
}

// DocumentoDocenteVinculacion ...
// @Title DocumentoDocenteVinculacion
// @Description Listar los docentes de acuerdo a la vinculacion y su documento
// @Param	documento		query 	string	true		"Documento docente"
// @Param	vinculacion		query 	int	true			"Id tipo de vinculación"
// @Success 200 {}
// @Failure 404 he request contains an incorrect parameter or no record exist
// @router /documento [get]
func (c *DocenteController) DocumentoDocenteVinculacion() {
	defer errorhandler.HandlePanic(&c.Controller)

	documento := c.GetString("documento")
	vinculacion, errvin := c.GetInt64("vinculacion")

	// Verificar si se proporciona el parametro de vinculacion, si no se proporciona se busca el docente con todas sus vinculaciones
	if errvin == nil { // Si no hay error en la conversion del parametro de vinculacion entonces la vinculacion fue proporcionada
		if documento == "" || vinculacion <= 0 { // Verifica si el documento o la vinculacion no son validos
			logs.Error(documento, vinculacion)                                                                            // Imprime los valores de los parametros
			c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) con valores no válidos") // Responde con un error 400
			c.Ctx.Output.SetStatus(400)                                                                                   // Establece el status de la respuesta
		} else { // Si los parametros son validos
			resultado := services.ListaDocentesxDocumentoVinculacion(documento, vinculacion) // Busca los docentes con el documento y la vinculacion proporcionados
			c.Data["json"] = resultado                                                       // Establece la respuesta
			c.Ctx.Output.SetStatus(resultado.Status)                                         // Establece el status de la respuesta
		}
	} else { // Si no se proporciona la vinculacion
		if documento == "" { // Verifica si el documento no es valido
			logs.Error(documento)                                                                                             // Imprime el valor del documento
			c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) no válido(s) o faltante(s)") // Responde con un error 400
			c.Ctx.Output.SetStatus(400)                                                                                       // Establece el status de la respuesta
		} else { // Si el documento es valido
			resultado := services.BuscarDocentesPorDocumentoConVinculaciones(documento) // Busca el docente con todas sus vinculaciones
			c.Data["json"] = resultado                                                  // Establece la respuesta
			c.Ctx.Output.SetStatus(resultado.Status)                                    // Establece el status de la respuesta
		}
	}

	c.ServeJSON()
}

// NombreDocenteVinculacion ...
// @Title NombreDocenteVinculacion
// @Description Listar los docentes de acuerdo a la vinculacion y su nombre
// @Param	nombre			query 	string	true		"Nombre docente"
// @Param	vinculacion		query	int	true			"Id tipo de vinculación"
// @Success 200 {}
// @Failure 404 he request contains an incorrect parameter or no record exist
// @router /nombre [get]
func (c *DocenteController) NombreDocenteVinculacion() {
	defer errorhandler.HandlePanic(&c.Controller)

	nombre := c.GetString("nombre")
	vinculacion, errvin := c.GetInt64("vinculacion")

	if errvin == nil {
		if nombre == "" || vinculacion <= 0 {
			logs.Error(nombre, vinculacion)
			c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) con valores no válidos")
			c.Ctx.Output.SetStatus(400)
		} else {
			resultado := services.ListaDocentesxNombreVinculacion(nombre, vinculacion)
			c.Data["json"] = resultado
			c.Ctx.Output.SetStatus(resultado.Status)
		}
	} else {
		if nombre == "" {
			logs.Error(nombre)
			c.Data["json"] = requestmanager.APIResponseDTO(false, 400, nil, "Error: Parámetro(s) no válido(s) o faltante(s)")
			c.Ctx.Output.SetStatus(400)
		} else {
			resultado := services.BuscarDocentesPorNombreConVinculaciones2(nombre)
			c.Data["json"] = resultado
			c.Ctx.Output.SetStatus(resultado.Status)
		}
	}
	c.ServeJSON()
}
