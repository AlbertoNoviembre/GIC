package generadorexcel

import (
	"fmt"
	"strconv"
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
var valor_Total float64
var valor_Usado float64
var valor_Libre float64
var n_filas int
var nombre_hoja string

func EstablecerValoresUsoDisco(total, usado, libre float64) {

	valor_Total = total
	valor_Usado = usado
	valor_Libre = libre
}

func CrearArchivo(nombre string, slice_archivos *[]rastreadorarchivos.Archivo, usuario string) {

	archivo_excel, err = excelize.OpenFile("/home/" + usuario + "/Escritorio" + "Contenido_Discos.xlsx")

	if err != nil {

		archivo_excel = excelize.NewFile()
		archivo_excel.Path = "/home/" + usuario + "/Escritorio" + "Contenido_Discos.xlsx"
	}

	crearHoja(nombre)

	insertarInfoExcel(nombre, slice_archivos)
	time.Sleep(time.Second * 1)
	setEstilos(nombre, slice_archivos)
	time.Sleep(time.Second * 1)
	archivo_excel.DeleteSheet("Sheet1")
	time.Sleep(time.Second * 1)
	archivo_excel.DeleteSheet("TODOS LOS DISPOSITIVOS")
	archivo_excel.Save()

	for _, hoja := range archivo_excel.GetSheetList() {

		valor_celda_C2, err := archivo_excel.GetCellValue(hoja, "C2")

		if err != nil {

			fmt.Println(err)
			return

		}

		valor_int, err := strconv.Atoi(valor_celda_C2)

		if err != nil {

			fmt.Print("El error es: ")
			fmt.Println(err)
			return

		}

		for j := 0; j < valor_int; j++ {

			valor_celda_nombre_archivo, err := archivo_excel.GetCellValue(hoja, fmt.Sprintf("A%d", j+4))
			if err != nil {

				fmt.Print("¿Está aquí el error?")
				fmt.Println(err)
				return

			}

			valor_celda_ruta_archivo, err := archivo_excel.GetCellValue(hoja, fmt.Sprintf("B%d", j+4))

			if err != nil {

				fmt.Println(err)
				return

			}
			archivo := rastreadorarchivos.Archivo{Nombre: valor_celda_nombre_archivo, Ruta: valor_celda_ruta_archivo}
			slice_todos_disp = append(slice_todos_disp, archivo)

		}

	}

	time.Sleep(time.Second * 1)
	crearHoja("TODOS LOS DISPOSITIVOS")
	time.Sleep(time.Second * 1)
	setEstilos("TODOS LOS DISPOSITIVOS", &slice_todos_disp)
	insertarInfoExcel("TODOS LOS DISPOSITIVOS", &slice_todos_disp)
	time.Sleep(time.Second * 1)
	slice_todos_disp = nil
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
	archivo_excel.SetPanes(nombre, `
    {
        "freeze":true,
        "y_split":3,
        "top_left_cell":"A4",
        "active_pane":"bottomRight",
        "panes":[
            {"pane":"topLeft"},
            {"pane":"topRight"},
            {"pane":"bottomLeft"},
            {"active_cell":"A4", "sqref":"A4", "pane":"bottomRight"}
            ]
    }
`)

	borde_celda_simple = []excelize.Border{{Type: "left", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1}}

	borde_celda_grueso = []excelize.Border{{Type: "left", Color: "0000AA", Style: 5},
		{Type: "top", Color: "0000BB", Style: 5},
		{Type: "bottom", Color: "000099", Style: 5},
		{Type: "right", Color: "0000CC", Style: 5}}

	encabezado, err := archivo_excel.NewStyle(&excelize.Style{

		Border:    borde_celda_grueso,
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#081c69"}, Pattern: 1},
		Font:      &excelize.Font{Size: 22, Color: "#FFFFFF"},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		NumFmt:    0,
	})

	encabezado_uso_disco, err := archivo_excel.NewStyle(&excelize.Style{

		Border:    borde_celda_grueso,
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#2c3b78"}, Pattern: 1},
		Font:      &excelize.Font{Size: 22, Color: "#FFFFFF"},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		NumFmt:    0,
	})

	encabezado_titl_ruta, err := archivo_excel.NewStyle(&excelize.Style{

		Border:    borde_celda_grueso,
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#2c365c"}, Pattern: 1},
		Font:      &excelize.Font{Size: 22, Color: "#FFFFFF"},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		NumFmt:    0,
	})

	columna_nombres_archivos, err := archivo_excel.NewStyle(&excelize.Style{

		Border:    borde_celda_simple,
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#fcf4e3"}, Pattern: 1},
		Font:      &excelize.Font{Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "left"},
	})

	columna_rutas_archivos, err := archivo_excel.NewStyle(&excelize.Style{

		Border:    borde_celda_simple,
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#d8f2db"}, Pattern: 1},
		Font:      &excelize.Font{Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "left"},
	})

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

	rango := fmt.Sprintf("A%d:A%d", 4, len(*slice_archivos)+4)
	if err := archivo_excel.SetConditionalFormat(nombre, rango, duplicCond); err != nil {

		fmt.Println(err)
		return

	}

	if archivo_excel.SetCellStyle(nombre, "A1", "B1", encabezado); err != nil {

		fmt.Println(err)

		return

	}

	if archivo_excel.SetCellStyle(nombre, "C1", "C1", encabezado); err != nil {

		fmt.Println(err)

		return

	}

	if archivo_excel.SetCellStyle(nombre, "C2", "C2", encabezado); err != nil {

		fmt.Println(err)

		return

	}

	if archivo_excel.SetCellStyle(nombre, "A2", "B2", encabezado_uso_disco); err != nil {

		fmt.Println(err)

		return

	}

	if archivo_excel.SetCellStyle(nombre, "A3", "B3", encabezado_titl_ruta); err != nil {

		fmt.Println(err)

		return

	}

	if archivo_excel.SetCellStyle(nombre, "A4", "A"+fmt.Sprint(len(*slice_archivos)+3), columna_nombres_archivos); err != nil {

		fmt.Println(err)

		return

	}

	if archivo_excel.SetCellStyle(nombre, "A4", "A"+fmt.Sprint(len(*slice_archivos)+3), columna_nombres_archivos); err != nil {

		fmt.Println(err)

		return

	}

	if archivo_excel.SetCellStyle(nombre, "B4", "B"+fmt.Sprint(len(*slice_archivos)+3), columna_rutas_archivos); err != nil {

		fmt.Println(err)

		return

	}

}

func insertarInfoExcel(nombre string, slice_archivos *[]rastreadorarchivos.Archivo) {

	var col_a string
	var col_b string

	archivo_excel.SetColWidth(nombre, "A", "A", 100)
	archivo_excel.SetColWidth(nombre, "B", "B", 120)
	archivo_excel.SetColWidth(nombre, "C", "C", 40)
	archivo_excel.SetCellValue(nombre, "A3", "TÍTULO")
	archivo_excel.SetCellValue(nombre, "B3", "RUTA (CARPETA)")
	archivo_excel.MergeCell(nombre, "A2", "B2")
	archivo_excel.SetCellValue(nombre, "A2", fmt.Sprintf("Total: %.2f GB    Usado: %.2f GB    Disponible: %.2f GB", valor_Total, valor_Usado, valor_Libre))
	archivo_excel.SetCellInt(nombre, "C2", len(*slice_archivos))
	archivo_excel.SetCellValue("TODOS LOS DISPOSITIVOS", "A2", "LISTADO DE ARCHIVOS DE TODOS LOS MEDIOS DE ALMACENAMIENTO.")
	for indice, archivo := range *slice_archivos {
		col_a = fmt.Sprintf("A%d", indice+4)
		col_b = fmt.Sprintf("B%d", indice+4)
		archivo_excel.SetCellValue(nombre, col_a, strings.ToUpper(archivo.Nombre))
		archivo_excel.SetCellValue(nombre, col_b, archivo.Ruta)
		archivo_excel.SetRowHeight(nombre, indice+4, 20)
	}

	archivo_excel.SetCellValue(nombre, "C1", "TOTAL ARCHIVOS")
	archivo_excel.SetRowHeight(nombre, 1, 25)
	archivo_excel.SetRowHeight(nombre, 2, 25)
	archivo_excel.SetRowHeight(nombre, 3, 25)
	archivo_excel.AutoFilter(nombre, "A3", "B3", `{"sort":"ascending"}`)

}
