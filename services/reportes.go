package services

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/k0kubun/pp"
	"github.com/phpdave11/gofpdf"
	"github.com/udistrital/sga_trabajo_docente_mid/models"
	"github.com/udistrital/sga_trabajo_docente_mid/utils"
	requestmanager "github.com/udistrital/utils_oas/requestresponse"
	"github.com/udistrital/utils_oas/xlsx2pdf"
	"github.com/xuri/excelize/v2"
)

// ReporteCargaLectiva ...
func RepCargaLectiva(docente, vinculacion, periodo int64, carga string) requestmanager.APIResponse {
	info, errinfo := obtenerInformacionRequeridaRepCargaLectiva(docente, vinculacion, periodo)
	if errinfo != nil {
		logs.Error(errinfo)
		return requestmanager.APIResponseDTO(false, 404, nil, errinfo.Error())
	}
	return generarReporteCargaLectiva(info, carga)
}

type resumenJson struct {
	HorasLectivas    float64 `json:"horas_lectivas,omitempty"`
	HorasActividades float64 `json:"horas_actividades,omitempty"`
	Observacion      string  `json:"observacion,omitempty"`
}

type infoRequeridaRepCL struct {
	datoIdenfTercero models.DatosIdentificacion
	datoVinculacion  models.Parametro
	datoPeriodo      models.Periodo
	datoPlanDocente  models.PlanDocente
	datoResumen      resumenJson
	datosCargaPlan   []models.CargaPlan
}

func obtenerInformacionRequeridaRepCargaLectiva(docente, vinculacion, periodo int64) (infoRequeridaRepCL, error) {
	resp, err := requestmanager.Get("http://"+beego.AppConfig.String("TercerosService")+
		fmt.Sprintf("datos_identificacion?query=Activo:true,TerceroId__Id:%d&fields=TerceroId,Numero,TipoDocumentoId&sortby=FechaExpedicion,Id&order=desc&limit=1", docente), requestmanager.ParseResonseNoFormat)
	if err != nil {
		logs.Error(err)
		return infoRequeridaRepCL{}, fmt.Errorf("TercerosService (datos_identificacion): " + err.Error())
	}
	datoIdenfTercero := models.DatosIdentificacion{}
	utils.ParseData(resp.([]interface{})[0], &datoIdenfTercero)

	resp, err = requestmanager.Get("http://"+beego.AppConfig.String("ParametroService")+fmt.Sprintf("parametro/%d", vinculacion), requestmanager.ParseResponseFormato1)
	if err != nil {
		logs.Error(err)
		return infoRequeridaRepCL{}, fmt.Errorf("ParametroService (parametro): " + err.Error())
	}
	datoVinculacion := models.Parametro{}
	utils.ParseData(resp, &datoVinculacion)

	resp, err = requestmanager.Get("http://"+beego.AppConfig.String("ParametroService")+fmt.Sprintf("periodo/%d", periodo), requestmanager.ParseResponseFormato1)
	if err != nil {
		logs.Error(err)
		return infoRequeridaRepCL{}, fmt.Errorf("ParametroService (periodo): " + err.Error())
	}
	datoPeriodo := models.Periodo{}
	utils.ParseData(resp, &datoPeriodo)

	resp, err = requestmanager.Get("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+
		fmt.Sprintf("plan_docente?query=activo:true,docente_id:%d,tipo_vinculacion_id:%d,periodo_id:%d&limit=1", docente, vinculacion, periodo), requestmanager.ParseResponseFormato1)
	if err != nil {
		logs.Error(err)
		return infoRequeridaRepCL{}, fmt.Errorf("PlanTrabajoDocenteService (plan_docente): " + err.Error())
	}
	datoPlanDocente := models.PlanDocente{}
	utils.ParseData(resp.([]interface{})[0], &datoPlanDocente)

	datoResumen := resumenJson{}
	json.Unmarshal([]byte(datoPlanDocente.Resumen), &datoResumen)

	resp, err = requestmanager.Get("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+
		fmt.Sprintf("carga_plan?query=activo:true,plan_docente_id:%s,&limit=0", datoPlanDocente.Id), requestmanager.ParseResponseFormato1)
	if err != nil {
		logs.Error(err)
		return infoRequeridaRepCL{}, fmt.Errorf("PlanTrabajoDocenteService (carga_plan): " + err.Error())
	}
	datosCargaPlan := []models.CargaPlan{}
	utils.ParseData(resp, &datosCargaPlan)

	for i := 0; i < len(datosCargaPlan); i++ {
		resp, err := requestmanager.Get("https://"+beego.AppConfig.String("HorarioService")+
			fmt.Sprintf("colocacion-espacio-academico/%s", datosCargaPlan[i].Colocacion_espacio_academico_id), requestmanager.ParseResponseFormato2)

		if err != nil {
			logs.Error(err)
			return infoRequeridaRepCL{}, fmt.Errorf("HorarioService (colocacion_espacio_academico): " + err.Error())
		}

		resumenColocacion := models.ResumenColocacion{}
		json.Unmarshal([]byte(resp.(map[string]interface{})["ResumenColocacionEspacioFisico"].(string)), &resumenColocacion)

		datosCargaPlan[i].Horario = string(resumenColocacion.Colocacion)
		datosCargaPlan[i].Sede_id = fmt.Sprint(resumenColocacion.EspacioFisico.SedeId)
		datosCargaPlan[i].Edificio_id = fmt.Sprint(resumenColocacion.EspacioFisico.EdificioId)
		datosCargaPlan[i].Salon_id = fmt.Sprint(resumenColocacion.EspacioFisico.SalonId)

	}

	return infoRequeridaRepCL{datoIdenfTercero, datoVinculacion, datoPeriodo, datoPlanDocente, datoResumen, datosCargaPlan}, nil
}

func generarReporteCargaLectiva(infoRequerida infoRequeridaRepCL, cargaTipo string) requestmanager.APIResponse {
	pp.Println(infoRequerida)
	inBog, _ := time.LoadLocation("America/Bogota")
	horaes := time.Now().In(inBog).Format("02/01/2006 15:04:05")

	path := beego.AppConfig.String("StaticPath")
	template, errt := excelize.OpenFile(path + "/templates/PTD.xlsx")
	if errt != nil {
		logs.Error(errt)
		return requestmanager.APIResponseDTO(false, 404, nil, "ReporteCargaLectiva (reading_template): "+errt.Error())
	}
	defer func() {
		if err := template.Close(); err != nil {
			logs.Error(err)
		}
	}()

	sheet := "PTD"
	nombreFormateado := utils.FormatNameTercero(infoRequerida.datoIdenfTercero.TerceroId)

	vinculacionFormateado := strings.ToLower(strings.Replace(infoRequerida.datoVinculacion.Nombre, "DOCENTE DE ", "", 1))
	vinculacionFormateado = strings.ToUpper(vinculacionFormateado[0:1]) + vinculacionFormateado[1:]

	footerstr := fmt.Sprintf("&L%s&C&CPágina &P de &N&R%s", "Oficina Asesora de Tecnologías e Información", "Fecha de generación: "+horaes)
	template.SetHeaderFooter(sheet, &excelize.HeaderFooterOptions{
		AlignWithMargins: boolPtr(true),
		ScaleWithDoc:     boolPtr(true),
		OddFooter:        footerstr,
	})
	// ? información del docente
	template.SetCellValue(sheet, "B8", nombreFormateado)
	template.SetCellValue(sheet, "V8", infoRequerida.datoIdenfTercero.TipoDocumentoId.CodigoAbreviacion+": "+infoRequerida.datoIdenfTercero.Numero)
	template.SetCellValue(sheet, "B11", infoRequerida.datoPeriodo.Nombre)
	template.SetCellValue(sheet, "V11", vinculacionFormateado)

	type coord struct {
		X float64 `json:"x"` // ? día
		Y float64 `json:"y"` // ? hora
	}

	type horario struct {
		HoraFormato string `json:"horaFormato,omitempty"`
		TipoCarga   int    `json:"tipo,omitempty"`
		Posicion    coord  `json:"finalPosition,omitempty"`
	}

	horarioIs := horario{}

	const (
		CargaLectiva int     = 1
		Actividades          = 2
		WidthX       float64 = 110
		HeightY      float64 = 18.75 // ? Altura de hora es 75px 1/4 => 18.75
	)

	ActividadStyle, _ := template.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Size: 6.5},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"e0ffff"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	CargaStyle, _ := template.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Size: 6.5},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"fafad2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	Lunes, Madrugada, _ := excelize.CellNameToCoordinates("G16") // ? Donde inicia cuadrícula de horario
	horamax := int(0)

	for _, carga := range infoRequerida.datosCargaPlan {
		json.Unmarshal([]byte(carga.Horario), &horarioIs)
		if cargaTipo == "C" && horarioIs.TipoCarga == Actividades {
			continue // ? Saltar a siguiente carga porque no es Carga Lectiva
		}
		if cargaTipo == "A" && horarioIs.TipoCarga == CargaLectiva {
			continue // ? Saltar a siguiente carga porque no es Actividad
		}
		// ? Añadir carga o actividad o las dos segun CargaTipo
		dia := int(horarioIs.Posicion.X/WidthX) * 5 // ? 5 => Cantidad de columnas por día cuadrícula excel
		horaIni := int(horarioIs.Posicion.Y / HeightY)
		horaFin := horaIni + int(carga.Duracion*4) // ? duración * 4 es para contar en cuartos de hora
		if horaFin >= horamax {
			horamax = horaFin
		}
		ini, _ := excelize.CoordinatesToCellName(Lunes+dia, Madrugada+horaIni)
		fin, _ := excelize.CoordinatesToCellName(Lunes+dia+4, Madrugada+horaFin-1) // ? +4 y -1 ajuste limite celdas
		template.MergeCell(sheet, ini, fin)

		nombreCarga := ""
		if horarioIs.TipoCarga == CargaLectiva {
			resp, err := requestmanager.Get("http://"+beego.AppConfig.String("EspaciosAcademicosService")+
				fmt.Sprintf("espacio-academico/%s", carga.Espacio_academico_id), requestmanager.ParseResponseFormato1)
			if err != nil {
				logs.Error(err)
				return requestmanager.APIResponseDTO(false, 404, nil, "EspaciosAcademicosService (espacio-academico): "+err.Error())
			}
			nombreCarga = resp.(map[string]interface{})["nombre"].(string) + " - " + resp.(map[string]interface{})["grupo"].(string)
			template.SetCellStyle(sheet, ini, fin, CargaStyle)
		} else if horarioIs.TipoCarga == Actividades {
			resp, err := requestmanager.Get("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+
				fmt.Sprintf("actividad/%s", carga.Actividad_id), requestmanager.ParseResponseFormato1)
			if err != nil {
				logs.Error(err)
				return requestmanager.APIResponseDTO(false, 404, nil, "PlanTrabajoDocenteService (actividad): "+err.Error())
			}
			nombreCarga = resp.(map[string]interface{})["nombre"].(string)
			template.SetCellStyle(sheet, ini, fin, ActividadStyle)
		}
		pp.Println(carga)
		infoEspacio, err := consultarInfoEspacioFisico(carga.Sede_id, carga.Edificio_id, carga.Salon_id)
		if err != nil {
			logs.Error(err)
			return requestmanager.APIResponseDTO(false, 404, nil, "OikosService (espacio_fisico): "+err.Error())
		}

		labelTag := fmt.Sprintf("*%s\n*%s - %s\n*%s\n*%s",
			nombreCarga,
			infoEspacio.(map[string]interface{})["sede"].(map[string]interface{})["CodigoAbreviacion"].(string),
			infoEspacio.(map[string]interface{})["edificio"].(map[string]interface{})["Nombre"].(string),
			infoEspacio.(map[string]interface{})["salon"].(map[string]interface{})["Nombre"].(string),
			horarioIs.HoraFormato,
		)
		template.SetCellValue(sheet, ini, labelTag)
	}

	// ? resumen
	template.SetCellValue(sheet, "M94", vinculacionFormateado)
	template.SetCellValue(sheet, "AD94", infoRequerida.datoResumen.HorasLectivas)
	template.SetCellValue(sheet, "M95", vinculacionFormateado)
	template.SetCellValue(sheet, "AD95", infoRequerida.datoResumen.HorasActividades)
	template.SetCellValue(sheet, "AD96", infoRequerida.datoResumen.HorasLectivas+infoRequerida.datoResumen.HorasActividades)
	template.SetCellValue(sheet, "B99", infoRequerida.datoResumen.Observacion)

	if cargaTipo == "C" { // ? si carga se borra actividades y total
		template.RemoveRow(sheet, 95)
		template.RemoveRow(sheet, 95)
	} else if cargaTipo == "A" { // ? si actividades se borra carga y total
		template.RemoveRow(sheet, 94)
		template.RemoveRow(sheet, 95)
	}

	if (Madrugada + horamax) <= 64 { // ? celda donde empieza la noche
		template.DeletePicture(sheet, "A87")
		template.DeletePicture(sheet, "AF87")
		for i := 0; i <= 20; i++ {
			template.RemoveRow(sheet, 64) // ? remover el horario de la noche
		}
		for i := 16; i <= 63; i++ {
			template.SetRowHeight(sheet, i, 9.8458) // ? ajustar altura del horario día si se quita la parte de la noche
		}
		template.AddPicture(sheet, "A66", path+"/img/logoud.jpeg", &excelize.GraphicOptions{
			ScaleX: 0.4,
			ScaleY: 0.324,
		})
		template.AddPicture(sheet, "AF66", path+"/img/logosga.jpeg", &excelize.GraphicOptions{
			ScaleX:  0.627,
			ScaleY:  0.5,
			OffsetX: 3,
		})
	}

	// * ----------
	// * Construcción de excel a pdf
	//

	pdf := gofpdf.New("P", "mm", "Letter", "")

	ExcelPdf := xlsx2pdf.Excel2PDF{
		Excel:  template,
		Pdf:    pdf,
		Sheets: make(map[string]xlsx2pdf.SheetInfo),
		WFx:    2.02,
		HFx:    2.85,
		Header: func() {},
		Footer: func() {},
	}

	ExcelPdf.Header = func() {
		pdf.SetFontSize(9)
		pdf.SetFontStyle("")
		lm, _, rm, _ := pdf.GetMargins()
		pw, _ := pdf.GetPageSize()
		x, y := pdf.GetXY()
		pdf.SetXY(lm, 8)
		//pdf.CellFormat(pw-lm-rm, 9, pdf.UnicodeTranslatorFromDescriptor("")("Plan Trabajo Docente"), "", 0, "CT", false, 0, "")
		pdf.BeginLayer(ExcelPdf.Layers.Imgs)
		pdf.ImageOptions(path+"/img/logoud.jpeg", lm, 8, 0, 15, false, gofpdf.ImageOptions{ImageType: "JPEG", ReadDpi: true}, 0, "")
		pdf.ImageOptions(path+"/img/logosga.jpeg", pw-rm-46.3157, 8, 46.3157, 0, false, gofpdf.ImageOptions{ImageType: "JPEG", ReadDpi: true}, 0, "")
		pdf.EndLayer()
		pdf.SetXY(x, y)
	}

	maxpages := ExcelPdf.EstimateMaxPages()
	ExcelPdf.Footer = func() {
		pdf.SetFontSize(9)
		pdf.SetFontStyle("")
		pagenum := pdf.PageNo()
		lm, _, rm, bm := pdf.GetMargins()
		pw, ph := pdf.GetPageSize()
		x, y := pdf.GetXY()
		pdf.SetXY(lm, ph-bm)
		w := (pw - lm - rm) / 3
		pdf.BeginLayer(ExcelPdf.Layers.Txts)
		pdf.CellFormat(w, 9, pdf.UnicodeTranslatorFromDescriptor("")("Oficina Asesora de Tecnologías e Información"), "", 0, "LT", false, 0, "")
		pdf.CellFormat(w, 9, pdf.UnicodeTranslatorFromDescriptor("")(fmt.Sprintf("Página %d de %d", pagenum, maxpages)), "", 0, "CT", false, 0, "")
		pdf.CellFormat(w, 9, pdf.UnicodeTranslatorFromDescriptor("")("Fecha de generación: "+horaes), "", 0, "RT", false, 0, "")
		pdf.EndLayer()
		pdf.SetXY(x, y)
	}

	ExcelPdf.ConvertSheets()

	// ? una vaina ahi para redimensionar las filas.. no coinciden en excel con respecto a pdf :(
	dim, _ := template.GetSheetDimension(sheet)
	_, maxrow, _ := excelize.CellNameToCoordinates(strings.Split(dim, ":")[1])
	for r := 1; r <= maxrow; r++ {
		h, _ := template.GetRowHeight(sheet, r)
		template.SetRowHeight(sheet, r, h*1.046)
	}

	// * ----------
	// * Convertir a base64
	//
	// ? excel
	bufferExcel, err := template.WriteToBuffer()
	if err != nil {
		logs.Error(err)
		return requestmanager.APIResponseDTO(false, 404, nil, "ReporteCargaLectiva (writing_file): "+err.Error())
		/* badAns, code := requestmanager.MidResponseFormat("ReporteCargaLectiva (writing_file)", "POST", false, map[string]interface{}{
			"response": nil,
			"error":    err.Error(),
		})
		c.Ctx.Output.SetStatus(code)
		c.Data["json"] = badAns
		c.ServeJSON()
		return */
	}
	encodedFileExcel := base64.StdEncoding.EncodeToString(bufferExcel.Bytes())

	// ? pdf
	var bufferPdf bytes.Buffer
	writer := bufio.NewWriter(&bufferPdf)
	pdf.Output(writer)
	writer.Flush()
	encodedFilePdf := base64.StdEncoding.EncodeToString(bufferPdf.Bytes())

	return requestmanager.APIResponseDTO(true, 200, map[string]interface{}{
		"excel": encodedFileExcel,
		"pdf":   encodedFilePdf,
	}, "Report Generation successful")
}

// ReporteVerificacionCumplimientoPTD ...
func RepCumplimiento(vigencia int64, proyecto string) requestmanager.APIResponse {
	info, errinfo := obtenerInformacionRequeridaRepCumplimiento(vigencia, proyecto)
	if errinfo != nil {
		logs.Error(errinfo)
		return requestmanager.APIResponseDTO(false, 404, nil, errinfo.Error())
	}
	return generarReporteCumplimiento(info)
}

type asignaturaPadreGrupos struct {
	Nombre string
	Grupos []string
}

type formatoCumplimiento struct {
	Nombre            string
	Documento         string
	Vinculacion       string
	Proyecto          string
	Actividades       map[string]float64
	Asignaturas       map[string]asignaturaPadreGrupos
	HorasLectivas     float64
	NumAsignaturas    int
	NumGrupos         int
	AsignaturasGrupos string
	Observacion       string
}

type infoRequeridaCumplimiento struct {
	PlanesPlanta  map[string]map[string]formatoCumplimiento
	PlanesTCO     map[string]map[string]formatoCumplimiento
	PlanesMTO     map[string]map[string]formatoCumplimiento
	ListaIdPlanes []string
}

func obtenerInformacionRequeridaRepCumplimiento(vigencia int64, proyectoFilter string) (infoRequeridaCumplimiento, error) {
	PlanesPlanta := map[string]map[string]formatoCumplimiento{} // ? tercero.proyecto.formatoCumplimiento
	PlanesTCO := map[string]map[string]formatoCumplimiento{}    // ? tercero.proyecto.formatoCumplimiento
	PlanesMTO := map[string]map[string]formatoCumplimiento{}    // ? tercero.proyecto.formatoCumplimiento
	ListaIdPlanes := []string{}

	plan_aprobado := "646fcf784c0bc253c1c720d4"
	resp, err := requestmanager.Get("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+
		fmt.Sprintf("plan_docente?query=activo:true,estado_plan_id:%s,periodo_id:%d&limit=0", plan_aprobado, vigencia), requestmanager.ParseResponseFormato1)
	if err != nil {
		logs.Error(err)
		return infoRequeridaCumplimiento{}, fmt.Errorf("PlanTrabajoDocenteService (plan_docente): " + err.Error())
		/* badAns, code := requestmanager.MidResponseFormat("PlanTrabajoDocenteService (plan_docente)", "GET", false, map[string]interface{}{
			"response": resp,
			"error":    err.Error(),
		})
		c.Ctx.Output.SetStatus(code)
		c.Data["json"] = badAns
		c.ServeJSON()
		return */
	}
	lista_planes := []models.PlanDocente{}
	utils.ParseData(resp, &lista_planes)
	for _, plan_docente := range lista_planes {

		resp, err = requestmanager.Get("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+
			fmt.Sprintf("carga_plan?query=activo:true,plan_docente_id:%s&limit=0", plan_docente.Id), requestmanager.ParseResponseFormato1)
		if err != nil {
			logs.Error(err)
			return infoRequeridaCumplimiento{}, fmt.Errorf("PlanTrabajoDocenteService (carga_plan): " + err.Error())
			/* badAns, code := requestmanager.MidResponseFormat("PlanTrabajoDocenteService (carga_plan)", "GET", false, map[string]interface{}{
				"response": resp,
				"error":    err.Error(),
			})
			c.Ctx.Output.SetStatus(code)
			c.Data["json"] = badAns
			c.ServeJSON()
			return */
		}
		carga_plan := []models.CargaPlan{}
		utils.ParseData(resp, &carga_plan)

		agrupacionEspacios := map[string]float64{}
		agrupacionActividades := map[string]float64{}
		for _, carga := range carga_plan {
			if carga.Espacio_academico_id != "" {
				agrupacionEspacios[carga.Espacio_academico_id] = agrupacionEspacios[carga.Espacio_academico_id] + carga.Duracion
			} else if carga.Actividad_id != "" {
				agrupacionActividades[carga.Actividad_id] = agrupacionActividades[carga.Actividad_id] + carga.Duracion
			}
		}

		for idEspAcad := range agrupacionEspacios {
			resp, err = requestmanager.Get("http://"+beego.AppConfig.String("EspaciosAcademicosService")+
				fmt.Sprintf("espacio-academico?query=_id:%s", idEspAcad), requestmanager.ParseResponseFormato1)
			if err != nil {
				logs.Error(err)
				return infoRequeridaCumplimiento{}, fmt.Errorf("EspaciosAcademicosService (espacio-academico): " + err.Error())
				/* badAns, code := requestmanager.MidResponseFormat("EspaciosAcademicosService (espacio-academico)", "GET", false, map[string]interface{}{
					"response": resp,
					"error":    err.Error(),
				})
				c.Ctx.Output.SetStatus(code)
				c.Data["json"] = badAns
				c.ServeJSON()
				return */
			}
			espacio_academico := models.EspacioAcademico{}
			utils.ParseData(resp.([]interface{})[0], &espacio_academico)

			docenteId := plan_docente.Docente_id
			projectId := fmt.Sprintf("%d", espacio_academico.Proyecto_academico_id)

			if (proyectoFilter == "0") || (proyectoFilter == projectId) {
				resp, err = requestmanager.Get("http://"+beego.AppConfig.String("ProyectoAcademicoService")+
					fmt.Sprintf("proyecto_academico_institucion/%s", projectId), requestmanager.ParseResonseNoFormat)
				if err != nil {
					logs.Error(err)
					return infoRequeridaCumplimiento{}, fmt.Errorf("ProyectoAcademicoService (proyecto_academico_institucion): " + err.Error())
					/* badAns, code := requestmanager.MidResponseFormat("ProyectoAcademicoService (proyecto_academico_institucion)", "GET", false, map[string]interface{}{
						"response": resp,
						"error":    err.Error(),
					})
					c.Ctx.Output.SetStatus(code)
					c.Data["json"] = badAns
					c.ServeJSON()
					return */
				}

				if plan_docente.Tipo_vinculacion_id == "293" || plan_docente.Tipo_vinculacion_id == "294" { // ? Carrera T Comp || Carrera Med T
					if _, ok := PlanesPlanta[docenteId]; !ok {
						PlanesPlanta[docenteId] = map[string]formatoCumplimiento{}
					}
					if _, ok := PlanesPlanta[docenteId][projectId]; !ok {
						PlanesPlanta[docenteId][projectId] = formatoCumplimiento{
							Asignaturas: map[string]asignaturaPadreGrupos{},
						}
					}
					PlanesPlanta[docenteId][projectId] = formatoCumplimiento{
						Proyecto:      resp.(map[string]interface{})["Nombre"].(string),
						HorasLectivas: PlanesPlanta[docenteId][projectId].HorasLectivas + agrupacionEspacios[idEspAcad],
						Asignaturas:   PlanesPlanta[docenteId][projectId].Asignaturas,
					}
					PlanesPlanta[docenteId][projectId].Asignaturas[espacio_academico.Espacio_academico_padre] = asignaturaPadreGrupos{
						Nombre: espacio_academico.Nombre,
						Grupos: append(PlanesPlanta[docenteId][projectId].Asignaturas[espacio_academico.Espacio_academico_padre].Grupos, espacio_academico.Grupo),
					}
				} else if plan_docente.Tipo_vinculacion_id == "296" { // ? T Comp Ocacional
					if _, ok := PlanesTCO[docenteId]; !ok {
						PlanesTCO[docenteId] = map[string]formatoCumplimiento{}
					}
					if _, ok := PlanesTCO[docenteId][projectId]; !ok {
						PlanesTCO[docenteId][projectId] = formatoCumplimiento{
							Asignaturas: map[string]asignaturaPadreGrupos{},
						}
					}
					PlanesTCO[docenteId][projectId] = formatoCumplimiento{
						Proyecto:      resp.(map[string]interface{})["Nombre"].(string),
						HorasLectivas: PlanesTCO[docenteId][projectId].HorasLectivas + agrupacionEspacios[idEspAcad],
						Asignaturas:   PlanesTCO[docenteId][projectId].Asignaturas,
					}
					PlanesTCO[docenteId][projectId].Asignaturas[espacio_academico.Espacio_academico_padre] = asignaturaPadreGrupos{
						Nombre: espacio_academico.Nombre,
						Grupos: append(PlanesTCO[docenteId][projectId].Asignaturas[espacio_academico.Espacio_academico_padre].Grupos, espacio_academico.Grupo),
					}
				} else if plan_docente.Tipo_vinculacion_id == "298" { // ? Med T Ocacional
					if _, ok := PlanesMTO[docenteId]; !ok {
						PlanesMTO[docenteId] = map[string]formatoCumplimiento{}
					}
					if _, ok := PlanesMTO[docenteId][projectId]; !ok {
						PlanesMTO[docenteId][projectId] = formatoCumplimiento{
							Asignaturas: map[string]asignaturaPadreGrupos{},
						}
					}
					PlanesMTO[docenteId][projectId] = formatoCumplimiento{
						Proyecto:      resp.(map[string]interface{})["Nombre"].(string),
						HorasLectivas: PlanesMTO[docenteId][projectId].HorasLectivas + agrupacionEspacios[idEspAcad],
						Asignaturas:   PlanesMTO[docenteId][projectId].Asignaturas,
					}
					PlanesMTO[docenteId][projectId].Asignaturas[espacio_academico.Espacio_academico_padre] = asignaturaPadreGrupos{
						Nombre: espacio_academico.Nombre,
						Grupos: append(PlanesMTO[docenteId][projectId].Asignaturas[espacio_academico.Espacio_academico_padre].Grupos, espacio_academico.Grupo),
					}
				}
				ListaIdPlanes = append(ListaIdPlanes, plan_docente.Id)
			}

		}

		resp, err := requestmanager.Get("http://"+beego.AppConfig.String("TercerosService")+
			fmt.Sprintf("datos_identificacion?query=Activo:true,TerceroId__Id:%v&fields=TerceroId,Numero,TipoDocumentoId&sortby=FechaExpedicion,Id&order=desc&limit=1",
				plan_docente.Docente_id), requestmanager.ParseResonseNoFormat)
		if err != nil {
			logs.Error(err)
			return infoRequeridaCumplimiento{}, fmt.Errorf("TercerosService (datos_identificacion): " + err.Error())
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
			fmt.Sprintf("parametro/%s", plan_docente.Tipo_vinculacion_id), requestmanager.ParseResponseFormato1)
		if err != nil {
			logs.Error(err)
			return infoRequeridaCumplimiento{}, fmt.Errorf("ParametroService (parametro): " + err.Error())
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

		datoResumen := map[string]interface{}{}
		json.Unmarshal([]byte(plan_docente.Resumen), &datoResumen)

		if plan_docente.Tipo_vinculacion_id == "293" || plan_docente.Tipo_vinculacion_id == "294" { // ? Carrera T Comp || Carrera Med T
			if _, ok := PlanesPlanta[plan_docente.Docente_id]; ok {
				PlanesPlanta[plan_docente.Docente_id]["actividades"] = formatoCumplimiento{
					Nombre:      utils.FormatNameTercero(datos_identificacion.TerceroId),
					Documento:   datos_identificacion.Numero,
					Vinculacion: infoVinculacion.Nombre,
					Actividades: agrupacionActividades,
					Observacion: datoResumen["observacion"].(string),
				}
			}
		} else if plan_docente.Tipo_vinculacion_id == "296" { // ? T Comp Ocacional
			if _, ok := PlanesTCO[plan_docente.Docente_id]; ok {
				PlanesTCO[plan_docente.Docente_id]["actividades"] = formatoCumplimiento{
					Nombre:      utils.FormatNameTercero(datos_identificacion.TerceroId),
					Documento:   datos_identificacion.Numero,
					Vinculacion: infoVinculacion.Nombre,
					Actividades: agrupacionActividades,
					Observacion: datoResumen["observacion"].(string),
				}
			}
		} else if plan_docente.Tipo_vinculacion_id == "298" { // ? Med T Ocacional
			if _, ok := PlanesMTO[plan_docente.Docente_id]; ok {
				PlanesMTO[plan_docente.Docente_id]["actividades"] = formatoCumplimiento{
					Nombre:      utils.FormatNameTercero(datos_identificacion.TerceroId),
					Documento:   datos_identificacion.Numero,
					Vinculacion: infoVinculacion.Nombre,
					Actividades: agrupacionActividades,
					Observacion: datoResumen["observacion"].(string),
				}
			}
		}

	}

	return infoRequeridaCumplimiento{PlanesPlanta, PlanesTCO, PlanesMTO, ListaIdPlanes}, nil
}

func generarReporteCumplimiento(infoRequerida infoRequeridaCumplimiento) requestmanager.APIResponse {
	inBog, _ := time.LoadLocation("America/Bogota")
	horaes := time.Now().In(inBog).Format("02/01/2006 15:04:05")

	path := beego.AppConfig.String("StaticPath")
	template, errt := excelize.OpenFile(path + "/templates/Verif_Cump_PTD.xlsx")
	if errt != nil {
		logs.Error(errt)
		return requestmanager.APIResponseDTO(false, 404, nil, "ReporteVerifCumpPTD (reading_template): "+errt.Error())
		/* badAns, code := requestmanager.MidResponseFormat("ReporteVerifCumpPTD (reading_template)", "GET", false, map[string]interface{}{
			"response": template,
			"error":    errt.Error(),
		})
		c.Ctx.Output.SetStatus(code)
		c.Data["json"] = badAns
		c.ServeJSON()
		return */
	}
	defer func() {
		if err := template.Close(); err != nil {
			logs.Error(err)
		}
	}()

	posicionActividades := map[string]interface{}{
		"647609c548f8405cfda2783f": "E",
		"64c0a7b2d1e67f3cb057f20b": "F",
		"64760a1e48f8405cfda27843": "G",
		"64c0a7e9d1e67fa0bd57f20e": "H",
		"64c0a81cd1e67f59af57f211": "I",
		"64c0a862d1e67fda0f57f214": "J",
		"64c0a89ed1e67f6b2557f217": "K",
		"64c0a8c9d1e67f5c7757f21a": "L",
		"64c0a8f3d1e67fbc4d57f21d": "M",
		"64c0a927d1e67f874557f220": "N",
		"64760be348f8405cfda27853": "O",
		"64760bff48f8405cfda27855": "P",
		"64760bd448f8405cfda27851": "Q",
		"64760af748f8405cfda2784b": "R",
		"64760a4048f8405cfda27845": "S",
		"647609f548f8405cfda27841": "T",
		"6476094048f8405cfda2783d": "U",
		"64760b8b48f8405cfda2784d": "V",
		"64760bb248f8405cfda2784f": "W",
		"64c0a945d1e67f813a57f223": "X",
		"64760a7648f8405cfda27847": "Y",
		"64c0a988d1e67f3c0857f226": "Z",
		"64c0a9c9d1e67f909457f229": "AA",
		"64c0a9f8d1e67fa11057f22c": "AB",
		"64c0aa2fd1e67f552157f22f": "AC",
		"64760a9d48f8405cfda27849": "AE",
	}

	rowPosition := 5
	sheet := "Planta"

	for docenteId := range infoRequerida.PlanesPlanta {
		for proyecto := range infoRequerida.PlanesPlanta[docenteId] {
			if proyecto != "actividades" {
				template.DuplicateRow(sheet, rowPosition)
			}
		}
	}

	footerstr := fmt.Sprintf("&L%s&C&CPágina &P de &N&R%s", "Oficina Asesora de Tecnologías e Información", "Fecha de generación: "+horaes)
	template.SetHeaderFooter(sheet, &excelize.HeaderFooterOptions{
		AlignWithMargins: boolPtr(true),
		ScaleWithDoc:     boolPtr(true),
		OddFooter:        footerstr,
	})

	for docenteId := range infoRequerida.PlanesPlanta {
		incrow := len(infoRequerida.PlanesPlanta[docenteId]) - 1
		template.MergeCell(sheet, fmt.Sprintf("A%d", rowPosition), fmt.Sprintf("A%d", rowPosition+incrow-1))
		template.SetCellValue(sheet, fmt.Sprintf("A%d", rowPosition), infoRequerida.PlanesPlanta[docenteId]["actividades"].Nombre)
		template.MergeCell(sheet, fmt.Sprintf("B%d", rowPosition), fmt.Sprintf("B%d", rowPosition+incrow-1))
		template.SetCellValue(sheet, fmt.Sprintf("B%d", rowPosition), infoRequerida.PlanesPlanta[docenteId]["actividades"].Documento)
		template.MergeCell(sheet, fmt.Sprintf("C%d", rowPosition), fmt.Sprintf("C%d", rowPosition+incrow-1))
		template.SetCellValue(sheet, fmt.Sprintf("C%d", rowPosition), infoRequerida.PlanesPlanta[docenteId]["actividades"].Vinculacion)
		for k := range posicionActividades {
			col := posicionActividades[k].(string)
			template.MergeCell(sheet, fmt.Sprintf("%s%d", col, rowPosition), fmt.Sprintf("%s%d", col, rowPosition+incrow-1))
			//template.SetCellValue(sheet, fmt.Sprintf("%s%d", col, rowPosition), 0)
		}

		iterProyects := int(0)
		for proyecto := range infoRequerida.PlanesPlanta[docenteId] {
			if proyecto != "actividades" {
				nombreAsignaturasGrupos := ""
				numGrupos := int(0)
				for k := range infoRequerida.PlanesPlanta[docenteId][proyecto].Asignaturas {
					nombreAsignaturasGrupos += "* " + infoRequerida.PlanesPlanta[docenteId][proyecto].Asignaturas[k].Nombre + " ("
					numGrupos = len(infoRequerida.PlanesPlanta[docenteId][proyecto].Asignaturas[k].Grupos)
					for i, g := range infoRequerida.PlanesPlanta[docenteId][proyecto].Asignaturas[k].Grupos {
						if (i + 1) != numGrupos {
							nombreAsignaturasGrupos += g + ", "
						} else {
							nombreAsignaturasGrupos += g + ")\n"
						}
					}
				}
				template.SetCellValue(sheet, fmt.Sprintf("D%d", rowPosition+iterProyects), infoRequerida.PlanesPlanta[docenteId][proyecto].Proyecto)
				template.SetCellValue(sheet, fmt.Sprintf("AG%d", rowPosition+iterProyects), infoRequerida.PlanesPlanta[docenteId][proyecto].HorasLectivas)
				template.SetCellValue(sheet, fmt.Sprintf("AI%d", rowPosition+iterProyects), len(infoRequerida.PlanesPlanta[docenteId][proyecto].Asignaturas))
				template.SetCellValue(sheet, fmt.Sprintf("AK%d", rowPosition+iterProyects), numGrupos)
				template.SetCellValue(sheet, fmt.Sprintf("AM%d", rowPosition+iterProyects), strings.TrimRight(nombreAsignaturasGrupos, "\n"))
				iterProyects++
			}
		}
		sumaHorasActividades := float64(0)
		for idActividad := range infoRequerida.PlanesPlanta[docenteId]["actividades"].Actividades {
			sumaHorasActividades += infoRequerida.PlanesPlanta[docenteId]["actividades"].Actividades[idActividad]
			template.SetCellValue(sheet, fmt.Sprintf("%s%d", posicionActividades[idActividad], rowPosition), infoRequerida.PlanesPlanta[docenteId]["actividades"].Actividades[idActividad])
		}
		template.MergeCell(sheet, fmt.Sprintf("AF%d", rowPosition), fmt.Sprintf("AF%d", rowPosition+incrow-1))
		template.SetCellValue(sheet, fmt.Sprintf("AF%d", rowPosition), sumaHorasActividades)
		template.MergeCell(sheet, fmt.Sprintf("AH%d", rowPosition), fmt.Sprintf("AH%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AH%d", rowPosition),
			fmt.Sprintf(`=SUM(AG%d:AG%d)`, rowPosition, rowPosition+incrow-1))
		template.MergeCell(sheet, fmt.Sprintf("AJ%d", rowPosition), fmt.Sprintf("AJ%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AJ%d", rowPosition),
			fmt.Sprintf(`=SUM(AI%d:AI%d)`, rowPosition, rowPosition+incrow-1))
		template.MergeCell(sheet, fmt.Sprintf("AL%d", rowPosition), fmt.Sprintf("AL%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AL%d", rowPosition),
			fmt.Sprintf(`=SUM(AK%d:AK%d)`, rowPosition, rowPosition+incrow-1))
		template.MergeCell(sheet, fmt.Sprintf("AN%d", rowPosition), fmt.Sprintf("AN%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AN%d", rowPosition),
			fmt.Sprintf(`=AF%d+AH%d`, rowPosition, rowPosition))
		template.MergeCell(sheet, fmt.Sprintf("AO%d", rowPosition), fmt.Sprintf("AO%d", rowPosition+incrow-1))
		template.SetCellValue(sheet, fmt.Sprintf("AO%d", rowPosition), infoRequerida.PlanesPlanta[docenteId]["actividades"].Observacion)

		template.MergeCell(sheet, fmt.Sprintf("AD%d", rowPosition), fmt.Sprintf("AD%d", rowPosition+incrow-1))
		dv := excelize.NewDataValidation(true)
		dv.Sqref = fmt.Sprintf("AD%d:AD%d", rowPosition, rowPosition+incrow-1)
		dv.SetDropList([]string{"Investigador Principal", "Co-Investigador"})
		template.AddDataValidation(sheet, dv)
		template.SetCellValue(sheet, fmt.Sprintf("AD%d", rowPosition), "Seleccionar")

		template.MergeCell(sheet, fmt.Sprintf("AP%d", rowPosition), fmt.Sprintf("AP%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AP%d", rowPosition),
			fmt.Sprintf(`=IF(VALUE(AN%d)=40,"Cumple",IF(VALUE(AN%d)>40,"Mas Horas","Menos Horas"))`, rowPosition, rowPosition))
		template.MergeCell(sheet, fmt.Sprintf("AQ%d", rowPosition), fmt.Sprintf("AQ%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AQ%d", rowPosition),
			fmt.Sprintf(`=IF(VALUE(AJ%d)=1,IF(VALUE(AH%d)>18,"Exceso Carga Lectiva",IF(VALUE(AH%d)<16,"Carga Lectiva Insuficiente","Cumple")),IF(VALUE(AH%d)>14,"Exceso Carga Lectiva",IF(VALUE(AH%d)<12,"Carga Lectiva Insuficiente","Cumple")))`,
				rowPosition, rowPosition, rowPosition, rowPosition, rowPosition))
		template.MergeCell(sheet, fmt.Sprintf("AR%d", rowPosition), fmt.Sprintf("AR%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AR%d", rowPosition),
			fmt.Sprintf(`=IF(VALUE(U%d)>VALUE(AH%d)/2,"Exceso Horas","Cumple")`, rowPosition, rowPosition))
		template.MergeCell(sheet, fmt.Sprintf("AS%d", rowPosition), fmt.Sprintf("AS%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AS%d", rowPosition),
			fmt.Sprintf(`=IF(VALUE(E%d)>8,"Exceso Horas","Cumple")`, rowPosition))
		template.MergeCell(sheet, fmt.Sprintf("AT%d", rowPosition), fmt.Sprintf("AT%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AT%d", rowPosition),
			fmt.Sprintf(`=IF(VALUE(N%d)=0,"N/A",IF(VALUE(N%d)>12,"Exceso Horas",IF(VALUE(N%d)<8,"Faltan Horas","Cumple")))`, rowPosition, rowPosition, rowPosition))
		template.MergeCell(sheet, fmt.Sprintf("AU%d", rowPosition), fmt.Sprintf("AU%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AU%d", rowPosition),
			fmt.Sprintf(`=IF(VALUE(J%d)=0,"N/A",IF(VALUE(J%d)>20,"Exceso Horas",IF(VALUE(J%d)<12,"Faltan Horas","Cumple")))`, rowPosition, rowPosition, rowPosition))
		template.MergeCell(sheet, fmt.Sprintf("AV%d", rowPosition), fmt.Sprintf("AV%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AV%d", rowPosition),
			fmt.Sprintf(`=IF(VALUE(K%d)=0,"N/A",IF(VALUE(K%d)>20,"Exceso Horas",IF(VALUE(K%d)<12,"Faltan Horas","Cumple")))`, rowPosition, rowPosition, rowPosition))
		template.MergeCell(sheet, fmt.Sprintf("AW%d", rowPosition), fmt.Sprintf("AW%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AW%d", rowPosition),
			fmt.Sprintf(`=IF(VALUE(L%d)=0,"N/A",IF(VALUE(L%d)>20,"Exceso Horas",IF(VALUE(L%d)<12,"Faltan Horas","Cumple")))`, rowPosition, rowPosition, rowPosition))
		template.MergeCell(sheet, fmt.Sprintf("AX%d", rowPosition), fmt.Sprintf("AX%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AX%d", rowPosition),
			fmt.Sprintf(`=IF(AND(AD%d="Investigador Principal",(VALUE(Z%d)+VALUE(Y%d)+VALUE(AA%d)+VALUE(AB%d)+VALUE(AC%d))>10),"Exeso de horas",IF(AND(AD%d="Co-Investigador",(VALUE(Z%d)+VALUE(Y%d)+VALUE(AA%d)+VALUE(AB%d)+VALUE(AC%d))>8),"Exceso de horas","Cumple"))`,
				rowPosition, rowPosition, rowPosition, rowPosition, rowPosition, rowPosition, rowPosition, rowPosition, rowPosition, rowPosition, rowPosition, rowPosition))
		rowPosition += incrow
	}
	template.SetSheetDimension(sheet, fmt.Sprintf("A1:AX%d", rowPosition-1))
	template.RemoveRow(sheet, rowPosition)

	rowPosition = 5
	sheet = "TCO"

	for docenteId := range infoRequerida.PlanesTCO {
		for proyecto := range infoRequerida.PlanesTCO[docenteId] {
			if proyecto != "actividades" {
				template.DuplicateRow(sheet, rowPosition)
			}
		}
	}

	footerstr = fmt.Sprintf("&L%s&C&CPágina &P de &N&R%s", "Oficina Asesora de Tecnologías e Información", "Fecha de generación: "+horaes)
	template.SetHeaderFooter(sheet, &excelize.HeaderFooterOptions{
		AlignWithMargins: boolPtr(true),
		ScaleWithDoc:     boolPtr(true),
		OddFooter:        footerstr,
	})

	for docenteId := range infoRequerida.PlanesTCO {
		incrow := len(infoRequerida.PlanesTCO[docenteId]) - 1
		template.MergeCell(sheet, fmt.Sprintf("A%d", rowPosition), fmt.Sprintf("A%d", rowPosition+incrow-1))
		template.SetCellValue(sheet, fmt.Sprintf("A%d", rowPosition), infoRequerida.PlanesTCO[docenteId]["actividades"].Nombre)
		template.MergeCell(sheet, fmt.Sprintf("B%d", rowPosition), fmt.Sprintf("B%d", rowPosition+incrow-1))
		template.SetCellValue(sheet, fmt.Sprintf("B%d", rowPosition), infoRequerida.PlanesTCO[docenteId]["actividades"].Documento)
		template.MergeCell(sheet, fmt.Sprintf("C%d", rowPosition), fmt.Sprintf("C%d", rowPosition+incrow-1))
		template.SetCellValue(sheet, fmt.Sprintf("C%d", rowPosition), infoRequerida.PlanesTCO[docenteId]["actividades"].Vinculacion)
		for k := range posicionActividades {
			col := posicionActividades[k].(string)
			template.MergeCell(sheet, fmt.Sprintf("%s%d", col, rowPosition), fmt.Sprintf("%s%d", col, rowPosition+incrow-1))
			//template.SetCellValue(sheet, fmt.Sprintf("%s%d", col, rowPosition), 0)
		}

		iterProyects := int(0)
		for proyecto := range infoRequerida.PlanesTCO[docenteId] {
			if proyecto != "actividades" {
				nombreAsignaturasGrupos := ""
				numGrupos := int(0)
				for k := range infoRequerida.PlanesTCO[docenteId][proyecto].Asignaturas {
					nombreAsignaturasGrupos += "* " + infoRequerida.PlanesTCO[docenteId][proyecto].Asignaturas[k].Nombre + " ("
					numGrupos = len(infoRequerida.PlanesTCO[docenteId][proyecto].Asignaturas[k].Grupos)
					for i, g := range infoRequerida.PlanesTCO[docenteId][proyecto].Asignaturas[k].Grupos {
						if (i + 1) != numGrupos {
							nombreAsignaturasGrupos += g + ", "
						} else {
							nombreAsignaturasGrupos += g + ")\n"
						}
					}
				}
				template.SetCellValue(sheet, fmt.Sprintf("D%d", rowPosition+iterProyects), infoRequerida.PlanesTCO[docenteId][proyecto].Proyecto)
				template.SetCellValue(sheet, fmt.Sprintf("AG%d", rowPosition+iterProyects), infoRequerida.PlanesTCO[docenteId][proyecto].HorasLectivas)
				template.SetCellValue(sheet, fmt.Sprintf("AI%d", rowPosition+iterProyects), len(infoRequerida.PlanesTCO[docenteId][proyecto].Asignaturas))
				template.SetCellValue(sheet, fmt.Sprintf("AK%d", rowPosition+iterProyects), numGrupos)
				template.SetCellValue(sheet, fmt.Sprintf("AM%d", rowPosition+iterProyects), strings.TrimRight(nombreAsignaturasGrupos, " \n"))
				iterProyects++
			}
		}
		sumaHorasActividades := float64(0)
		for idActividad := range infoRequerida.PlanesTCO[docenteId]["actividades"].Actividades {
			sumaHorasActividades += infoRequerida.PlanesTCO[docenteId]["actividades"].Actividades[idActividad]
			template.SetCellValue(sheet, fmt.Sprintf("%s%d", posicionActividades[idActividad], rowPosition), infoRequerida.PlanesTCO[docenteId]["actividades"].Actividades[idActividad])
		}
		template.MergeCell(sheet, fmt.Sprintf("AF%d", rowPosition), fmt.Sprintf("AF%d", rowPosition+incrow-1))
		template.SetCellValue(sheet, fmt.Sprintf("AF%d", rowPosition), sumaHorasActividades)
		template.MergeCell(sheet, fmt.Sprintf("AH%d", rowPosition), fmt.Sprintf("AH%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AH%d", rowPosition),
			fmt.Sprintf(`=SUM(AG%d:AG%d)`, rowPosition, rowPosition+incrow-1))
		template.MergeCell(sheet, fmt.Sprintf("AJ%d", rowPosition), fmt.Sprintf("AJ%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AJ%d", rowPosition),
			fmt.Sprintf(`=SUM(AI%d:AI%d)`, rowPosition, rowPosition+incrow-1))
		template.MergeCell(sheet, fmt.Sprintf("AL%d", rowPosition), fmt.Sprintf("AL%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AL%d", rowPosition),
			fmt.Sprintf(`=SUM(AK%d:AK%d)`, rowPosition, rowPosition+incrow-1))
		template.MergeCell(sheet, fmt.Sprintf("AN%d", rowPosition), fmt.Sprintf("AN%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AN%d", rowPosition),
			fmt.Sprintf(`=AF%d+AH%d`, rowPosition, rowPosition))
		template.MergeCell(sheet, fmt.Sprintf("AO%d", rowPosition), fmt.Sprintf("AO%d", rowPosition+incrow-1))
		template.SetCellValue(sheet, fmt.Sprintf("AO%d", rowPosition), infoRequerida.PlanesTCO[docenteId]["actividades"].Observacion)

		template.MergeCell(sheet, fmt.Sprintf("AD%d", rowPosition), fmt.Sprintf("AD%d", rowPosition+incrow-1))
		dv := excelize.NewDataValidation(true)
		dv.Sqref = fmt.Sprintf("AD%d:AD%d", rowPosition, rowPosition+incrow-1)
		dv.SetDropList([]string{"Investigador Principal", "Co-Investigador"})
		template.AddDataValidation(sheet, dv)
		template.SetCellValue(sheet, fmt.Sprintf("AD%d", rowPosition), "Seleccionar")

		template.MergeCell(sheet, fmt.Sprintf("AP%d", rowPosition), fmt.Sprintf("AP%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AP%d", rowPosition),
			fmt.Sprintf(`=IF(VALUE(AN%d)=40,"Cumple",IF(VALUE(AN%d)>40,"Mas Horas","Menos Horas"))`, rowPosition, rowPosition))
		template.MergeCell(sheet, fmt.Sprintf("AQ%d", rowPosition), fmt.Sprintf("AQ%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AQ%d", rowPosition),
			fmt.Sprintf(`=IF(VALUE(AH%d)>24,"Exceso Carga Lectiva",IF(VALUE(AH%d)<20,"Carga Lectiva insuficiente","Cumple"))`, rowPosition, rowPosition))
		template.MergeCell(sheet, fmt.Sprintf("AR%d", rowPosition), fmt.Sprintf("AR%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AR%d", rowPosition),
			fmt.Sprintf(`=IF(VALUE(U%d)>VALUE(AH%d)/2,"Exceso Horas","Cumple")`, rowPosition, rowPosition))
		rowPosition += incrow
	}
	template.SetSheetDimension(sheet, fmt.Sprintf("A1:AR%d", rowPosition-1))
	template.RemoveRow(sheet, rowPosition)

	rowPosition = 5
	sheet = "MTO"

	for docenteId := range infoRequerida.PlanesMTO {
		for proyecto := range infoRequerida.PlanesMTO[docenteId] {
			if proyecto != "actividades" {
				template.DuplicateRow(sheet, rowPosition)
			}
		}
	}

	footerstr = fmt.Sprintf("&L%s&C&CPágina &P de &N&R%s", "Oficina Asesora de Tecnologías e Información", "Fecha de generación: "+horaes)
	template.SetHeaderFooter(sheet, &excelize.HeaderFooterOptions{
		AlignWithMargins: boolPtr(true),
		ScaleWithDoc:     boolPtr(true),
		OddFooter:        footerstr,
	})

	for docenteId := range infoRequerida.PlanesMTO {
		incrow := len(infoRequerida.PlanesMTO[docenteId]) - 1
		template.MergeCell(sheet, fmt.Sprintf("A%d", rowPosition), fmt.Sprintf("A%d", rowPosition+incrow-1))
		template.SetCellValue(sheet, fmt.Sprintf("A%d", rowPosition), infoRequerida.PlanesMTO[docenteId]["actividades"].Nombre)
		template.MergeCell(sheet, fmt.Sprintf("B%d", rowPosition), fmt.Sprintf("B%d", rowPosition+incrow-1))
		template.SetCellValue(sheet, fmt.Sprintf("B%d", rowPosition), infoRequerida.PlanesMTO[docenteId]["actividades"].Documento)
		template.MergeCell(sheet, fmt.Sprintf("C%d", rowPosition), fmt.Sprintf("C%d", rowPosition+incrow-1))
		template.SetCellValue(sheet, fmt.Sprintf("C%d", rowPosition), infoRequerida.PlanesMTO[docenteId]["actividades"].Vinculacion)
		for k := range posicionActividades {
			col := posicionActividades[k].(string)
			template.MergeCell(sheet, fmt.Sprintf("%s%d", col, rowPosition), fmt.Sprintf("%s%d", col, rowPosition+incrow-1))
			//template.SetCellValue(sheet, fmt.Sprintf("%s%d", col, rowPosition), 0)
		}

		iterProyects := int(0)
		for proyecto := range infoRequerida.PlanesMTO[docenteId] {
			if proyecto != "actividades" {
				nombreAsignaturasGrupos := ""
				numGrupos := int(0)
				for k := range infoRequerida.PlanesMTO[docenteId][proyecto].Asignaturas {
					nombreAsignaturasGrupos += "* " + infoRequerida.PlanesMTO[docenteId][proyecto].Asignaturas[k].Nombre + " ("
					numGrupos = len(infoRequerida.PlanesMTO[docenteId][proyecto].Asignaturas[k].Grupos)
					for i, g := range infoRequerida.PlanesMTO[docenteId][proyecto].Asignaturas[k].Grupos {
						if (i + 1) != numGrupos {
							nombreAsignaturasGrupos += g + ", "
						} else {
							nombreAsignaturasGrupos += g + ")\n"
						}
					}
				}
				template.SetCellValue(sheet, fmt.Sprintf("D%d", rowPosition+iterProyects), infoRequerida.PlanesMTO[docenteId][proyecto].Proyecto)
				template.SetCellValue(sheet, fmt.Sprintf("AG%d", rowPosition+iterProyects), infoRequerida.PlanesMTO[docenteId][proyecto].HorasLectivas)
				template.SetCellValue(sheet, fmt.Sprintf("AI%d", rowPosition+iterProyects), len(infoRequerida.PlanesMTO[docenteId][proyecto].Asignaturas))
				template.SetCellValue(sheet, fmt.Sprintf("AK%d", rowPosition+iterProyects), numGrupos)
				template.SetCellValue(sheet, fmt.Sprintf("AM%d", rowPosition+iterProyects), strings.TrimRight(nombreAsignaturasGrupos, " \n"))
				iterProyects++
			}
		}
		sumaHorasActividades := float64(0)
		for idActividad := range infoRequerida.PlanesMTO[docenteId]["actividades"].Actividades {
			sumaHorasActividades += infoRequerida.PlanesMTO[docenteId]["actividades"].Actividades[idActividad]
			template.SetCellValue(sheet, fmt.Sprintf("%s%d", posicionActividades[idActividad], rowPosition), infoRequerida.PlanesMTO[docenteId]["actividades"].Actividades[idActividad])
		}
		template.MergeCell(sheet, fmt.Sprintf("AF%d", rowPosition), fmt.Sprintf("AF%d", rowPosition+incrow-1))
		template.SetCellValue(sheet, fmt.Sprintf("AF%d", rowPosition), sumaHorasActividades)
		template.MergeCell(sheet, fmt.Sprintf("AH%d", rowPosition), fmt.Sprintf("AH%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AH%d", rowPosition),
			fmt.Sprintf(`=SUM(AG%d:AG%d)`, rowPosition, rowPosition+incrow-1))
		template.MergeCell(sheet, fmt.Sprintf("AJ%d", rowPosition), fmt.Sprintf("AJ%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AJ%d", rowPosition),
			fmt.Sprintf(`=SUM(AI%d:AI%d)`, rowPosition, rowPosition+incrow-1))
		template.MergeCell(sheet, fmt.Sprintf("AL%d", rowPosition), fmt.Sprintf("AL%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AL%d", rowPosition),
			fmt.Sprintf(`=SUM(AK%d:AK%d)`, rowPosition, rowPosition+incrow-1))
		template.MergeCell(sheet, fmt.Sprintf("AN%d", rowPosition), fmt.Sprintf("AN%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AN%d", rowPosition),
			fmt.Sprintf(`=AF%d+AH%d`, rowPosition, rowPosition))
		template.MergeCell(sheet, fmt.Sprintf("AO%d", rowPosition), fmt.Sprintf("AO%d", rowPosition+incrow-1))
		template.SetCellValue(sheet, fmt.Sprintf("AO%d", rowPosition), infoRequerida.PlanesMTO[docenteId]["actividades"].Observacion)

		template.MergeCell(sheet, fmt.Sprintf("AD%d", rowPosition), fmt.Sprintf("AD%d", rowPosition+incrow-1))
		dv := excelize.NewDataValidation(true)
		dv.Sqref = fmt.Sprintf("AD%d:AD%d", rowPosition, rowPosition+incrow-1)
		dv.SetDropList([]string{"Investigador Principal", "Co-Investigador"})
		template.AddDataValidation(sheet, dv)
		template.SetCellValue(sheet, fmt.Sprintf("AD%d", rowPosition), "Seleccionar")

		template.MergeCell(sheet, fmt.Sprintf("AP%d", rowPosition), fmt.Sprintf("AP%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AP%d", rowPosition),
			fmt.Sprintf(`=IF(VALUE(AN%d)=20,"Cumple",IF(VALUE(AN%d)>20,"Mas Horas","Menos Horas"))`, rowPosition, rowPosition))
		template.MergeCell(sheet, fmt.Sprintf("AQ%d", rowPosition), fmt.Sprintf("AQ%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AQ%d", rowPosition),
			fmt.Sprintf(`=IF(VALUE(AH%d)>16,"Exceso Carga Lectiva",IF(VALUE(AH%d)<12,"Carga Lectiva insuficiente","Cumple"))`, rowPosition, rowPosition))
		template.MergeCell(sheet, fmt.Sprintf("AR%d", rowPosition), fmt.Sprintf("AR%d", rowPosition+incrow-1))
		template.SetCellFormula(sheet, fmt.Sprintf("AR%d", rowPosition),
			fmt.Sprintf(`=IF(VALUE(U%d)>VALUE(AH%d)/2,"Exceso Horas","Cumple")`, rowPosition, rowPosition))
		rowPosition += incrow
	}
	template.SetSheetDimension(sheet, fmt.Sprintf("A1:AR%d", rowPosition-1))
	template.RemoveRow(sheet, rowPosition)

	// * ----------
	// * Construcción de excel a pdf
	//

	pdf := gofpdf.New("L", "mm", "", "")

	ExcelPdf := xlsx2pdf.Excel2PDF{
		Excel:    template,
		Pdf:      pdf,
		Sheets:   make(map[string]xlsx2pdf.SheetInfo),
		WFx:      2.02,
		HFx:      2.925,
		FontDims: xlsx2pdf.FontDims{Size: 0.85},
		Header:   func() {},
		Footer:   func() {},
		CustomSize: xlsx2pdf.PageFormat{
			Orientation: "L",
			Wd:          215.9,
			Ht:          1778,
		},
	}

	ExcelPdf.Header = func() {
		if ExcelPdf.PageCount == 1 {
			pdf.SetFontSize(9)
			pdf.SetFontStyle("")
			lm, _, rm, _ := pdf.GetMargins()
			pw, _ := pdf.GetPageSize()
			x, y := pdf.GetXY()
			pdf.SetXY(lm, 8)
			pdf.BeginLayer(ExcelPdf.Layers.Imgs)
			pdf.ImageOptions(path+"/img/logoud.jpeg", lm+1.5, 12, 50, 0, false, gofpdf.ImageOptions{ImageType: "JPEG", ReadDpi: true}, 0, "")
			pdf.ImageOptions(path+"/img/logosigud.jpeg", (pw/5)-rm-60, 12, 60, 0, false, gofpdf.ImageOptions{ImageType: "JPEG", ReadDpi: true}, 0, "")
			pdf.EndLayer()
			pdf.SetXY(x, y)
		}
	}

	maxpages := ExcelPdf.EstimateMaxPages()
	ExcelPdf.Footer = func() {
		pdf.SetFontSize(9)
		pdf.SetFontStyle("")
		pagenum := pdf.PageNo()
		lm, _, rm, bm := pdf.GetMargins()
		pw, ph := pdf.GetPageSize()
		x, y := pdf.GetXY()
		pdf.SetXY(lm, ph-bm)
		w := (pw - lm - rm) / 3
		pdf.BeginLayer(ExcelPdf.Layers.Txts)
		pdf.CellFormat(w, 9, pdf.UnicodeTranslatorFromDescriptor("")("Oficina Asesora de Tecnologías e Información"), "", 0, "LT", false, 0, "")
		pdf.CellFormat(w, 9, pdf.UnicodeTranslatorFromDescriptor("")(fmt.Sprintf("Página %d de %d", pagenum, maxpages)), "", 0, "CT", false, 0, "")
		pdf.CellFormat(w, 9, pdf.UnicodeTranslatorFromDescriptor("")("Fecha de generación: "+horaes), "", 0, "RT", false, 0, "")
		pdf.EndLayer()
		pdf.SetXY(x, y)
	}

	ExcelPdf.ConvertSheets()

	// * ----------
	// * Convertir a base64
	//
	// ? excel
	buffer, err := template.WriteToBuffer()
	if err != nil {
		logs.Error(err)
		return requestmanager.APIResponseDTO(false, 404, nil, "ReporteVerifCumpPTD (writing_file): "+err.Error())
		/* badAns, code := requestmanager.MidResponseFormat("ReporteVerifCumpPTD (writing_file)", "POST", false, map[string]interface{}{
			"response": nil,
			"error":    err.Error(),
		})
		c.Ctx.Output.SetStatus(code)
		c.Data["json"] = badAns
		c.ServeJSON()
		return */
	}
	encodedFileExcel := base64.StdEncoding.EncodeToString(buffer.Bytes())

	// ? pdf
	var bufferPdf bytes.Buffer
	writer := bufio.NewWriter(&bufferPdf)
	pdf.Output(writer)
	writer.Flush()
	encodedFilePdf := base64.StdEncoding.EncodeToString(bufferPdf.Bytes())
	//
	// * ----------

	listaIdPlanesSimple := utils.RemoveDuplicated(infoRequerida.ListaIdPlanes, func(item interface{}) interface{} {
		return item.(string)
	})

	return requestmanager.APIResponseDTO(true, 200, map[string]interface{}{
		"excel":         encodedFileExcel,
		"pdf":           encodedFilePdf,
		"listaIdPlanes": listaIdPlanesSimple,
	}, "Report Generation successful")
}

// funciones transversales
func consultarInfoEspacioFisico(sede_id, edificio_id, salon_id string) (interface{}, error) {
	sede, err := requestmanager.Get("http://"+beego.AppConfig.String("OikosService")+fmt.Sprintf("espacio_fisico?query=Id:%s&fields=Id,Nombre,CodigoAbreviacion&limit=1", sede_id),
		requestmanager.ParseResonseNoFormat)

	if err != nil {
		return nil, err
	}
	edificio, err := requestmanager.Get("http://"+beego.AppConfig.String("OikosService")+fmt.Sprintf("espacio_fisico?query=Id:%s&fields=Id,Nombre,CodigoAbreviacion&limit=1", edificio_id),
		requestmanager.ParseResonseNoFormat)

	if err != nil {
		return nil, err
	}
	salon, err := requestmanager.Get("http://"+beego.AppConfig.String("OikosService")+fmt.Sprintf("espacio_fisico?query=Id:%s&fields=Id,Nombre,CodigoAbreviacion&limit=1", salon_id),
		requestmanager.ParseResonseNoFormat)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"sede":     sede.([]interface{})[0],
		"edificio": edificio.([]interface{})[0],
		"salon":    salon.([]interface{})[0],
	}, nil
}

func boolPtr(b bool) *bool {
	return &b
}
