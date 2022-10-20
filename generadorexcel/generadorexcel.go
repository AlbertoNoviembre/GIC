package generadorexcel

import (
	"fmt"
	"strings"
	"time"

	"github.com/AlbertoNoviembre/GIC/rastreadorarchivos"
	"github.com/xuri/excelize/v2"
)

var archivo_excel *excelize.File
var err error
var valor_celda string
var borde_celda_simple, borde_celda_grueso []excelize.Border
var estilo1, estilo2, estilo3, estilo4 excelize.Style
var slice_todos_disp []rastreadorarchivos.Archivo //Aquí he declarado la variable como tipo SLICE del STRUCT 'Archivo'. Este STRUCT está definido en el paquete
// rastreadorarchivos
var n_filas int

func CrearArchivo(nombre string, slice_archivos *[]rastreadorarchivos.Archivo, usuario string) {

	archivo_excel, err = excelize.OpenFile("/home/" + usuario + "/" + "Contenido_Discos.xlsx")

	if err != nil {

		archivo_excel = excelize.NewFile()
		archivo_excel.Path = "/home/" + usuario + "/" + "Contenido_Discos.xlsx"
	}

	crearHoja(nombre)
	crearHoja("TODOS LOS DISPOSITIVOS")

	insertarInfoExcel(nombre, slice_archivos)

	for _, archivo := range *slice_archivos {

		slice_todos_disp = append(slice_todos_disp, archivo)

	}

	for i := 0; i < archivo_excel.SheetCount-1; i++ {

		n_filas += len(*slice_archivos)

		archivo_excel.SetActiveSheet(i)

		insertarInfoExcel("TODOS LOS DISPOSITIVOS", &slice_todos_disp)

	}

	time.Sleep(time.Second * 1)
	setEstilos(nombre, slice_archivos)
	time.Sleep(time.Second * 1)
	setEstilos("TODOS LOS DISPOSITIVOS", &slice_todos_disp)

	archivo_excel.Save()
	fmt.Print("\n\n¡EL ARCHIVO ESTÁ CREADO!\n\n")
}

func crearHoja(nombre string) {

	archivo_excel.DeleteSheet(nombre)
	time.Sleep(time.Second * 1)
	archivo_excel.NewSheet(nombre)
	archivo_excel.MergeCell(nombre, "A1", "B1")
	archivo_excel.SetCellValue(nombre, "A1", nombre)

	if err != nil {

		fmt.Println(err)

		return

	}

}

func setEstilos(nombre string, slice_archivos *[]rastreadorarchivos.Archivo) {

	borde_celda_simple = []excelize.Border{{Type: "left", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1}}

	borde_celda_grueso = []excelize.Border{{Type: "left", Color: "0000AA", Style: 5},
		{Type: "top", Color: "0000BB", Style: 5},
		{Type: "bottom", Color: "000099", Style: 5},
		{Type: "right", Color: "0000CC", Style: 5}}

	estilo1, err := archivo_excel.NewStyle(&excelize.Style{

		Border:    borde_celda_simple,
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#fcf4e3"}, Pattern: 1},
		Font:      &excelize.Font{Size: 12},
		Alignment: &excelize.Alignment{Horizontal: "left"},
	})

	estilo2, err := archivo_excel.NewStyle(&excelize.Style{

		Border:    borde_celda_grueso,
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#DFE8F7", "#ADA6C2", "#8A8290"}, Pattern: 1},
		Font:      &excelize.Font{Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	estilo3, err := archivo_excel.NewStyle(&excelize.Style{

		Border:    borde_celda_grueso,
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#DFE8F7", "#ADA6C2", "#8A8290"}, Pattern: 1},
		Font:      &excelize.Font{Size: 18},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	estilo4, err := archivo_excel.NewStyle(&excelize.Style{

		Border:    borde_celda_simple,
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#d8f2db"}, Pattern: 1},
		Font:      &excelize.Font{Size: 11},
		Alignment: &excelize.Alignment{Horizontal: "left"},
	})

	estilo5, err := archivo_excel.NewStyle(&excelize.Style{

		Border:    borde_celda_grueso,
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#DFE8F7", "#ADA6C2", "#8A8290"}, Pattern: 1},
		Font:      &excelize.Font{Size: 18},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		NumFmt:    0,
	})

	archivo_excel.SetCellStyle(nombre, "C2", "C2", estilo5)

	archivo_excel.SetCellStyle(nombre, "A2", "B2", estilo2)

	red, err := archivo_excel.NewConditionalStyle(`{
		
			"font":{

				"color":"#9A0511"

			},

			"fill":{

				"type":"pattern",
				"color":["#FEC7CE"],
				"pattern":1

			}
			
		}`)

	duplicCond := fmt.Sprintf(`[
		
		{

			"type":"duplicate",
			"criteria":"=",
			"format":%d



		}
	
	
	
	
	]`, red)

	rango := fmt.Sprintf("A%d:A%d", 3, len(slice_todos_disp)+3)
	if err := archivo_excel.SetConditionalFormat(nombre, rango, duplicCond); err != nil {

		fmt.Println(err)
		return

	}

	if archivo_excel.SetCellStyle(nombre, "A1", "B1", estilo3); err != nil {

		fmt.Println(err)

		return

	}

	if archivo_excel.SetCellStyle(nombre, "A3", "A"+fmt.Sprint(len(*slice_archivos)+2), estilo1); err != nil {

		fmt.Println(err)

		return

	}

	if archivo_excel.SetCellStyle(nombre, "A3", "A"+fmt.Sprint(len(*slice_archivos)+2), estilo1); err != nil {

		fmt.Println(err)

		return

	}

	if archivo_excel.SetCellStyle(nombre, "B3", "B"+fmt.Sprint(len(*slice_archivos)+2), estilo4); err != nil {

		fmt.Println(err)

		return

	}

}

func insertarInfoExcel(nombre string, slice_archivos *[]rastreadorarchivos.Archivo) {

	var col_a string
	var col_b string

	archivo_excel.SetColWidth(nombre, "A", "A", 100)
	archivo_excel.SetColWidth(nombre, "B", "B", 120)
	archivo_excel.SetCellValue(nombre, "A2", "TÍTULO")
	archivo_excel.SetCellValue(nombre, "B2", "RUTA (CARPETA)")
	archivo_excel.SetCellValue(nombre, "C2", len(*slice_archivos))

	for indice, archivo := range *slice_archivos {
		col_a = fmt.Sprintf("A%d", indice+3)
		col_b = fmt.Sprintf("B%d", indice+3)
		archivo_excel.SetCellValue(nombre, col_a, strings.ToUpper(archivo.Nombre))
		archivo_excel.SetCellValue(nombre, col_b, archivo.Ruta)

	}
}
