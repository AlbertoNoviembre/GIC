package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/AlbertoNoviembre/GIC/generadorexcel"
	"github.com/AlbertoNoviembre/GIC/rastreadorarchivos"
	"github.com/AlbertoNoviembre/GIC/usodisco"
	"github.com/tawesoft/golib/v2/dialog"
	"github.com/zcalusic/sysinfo"
)

var num_dispvs int

var disps []string
var data = binding.BindStringList(&disps)
var datosvacios = binding.BindStringList(&[]string{})
var ruta string
var disp_selec string
var ruta_slice []string
var progreso float64
var infoDisco usodisco.EstadoDispAlmac
var canal_nombres_archivo chan string

func main() {

	fmt.Printf("Estás usando el sistema operativo %s\n-----------------------------------------\n", getNombreSO())

	canal := make(chan []string)

	go getListaDispExter(canal)

	fmt.Println("Iniciando...")

	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("GIC - Alberto Álvarez Portero (2022)")

	var barra_de_progreso = widget.NewProgressBar()
	barra_de_progreso.Hide()

	lbl_lista_disp := widget.NewLabel("DISPOSITIVOS MONTADOS:")
	lbl_lista_disp.Alignment = fyne.TextAlign(1)

	lista := widget.NewListWithData(data,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))

		})

	lbl_nombres_archivo := widget.NewLabel("")
	lbl_nombres_archivo.Move(fyne.NewPos(10, 395))
	lbl_nombres_archivo.Resize(fyne.NewSize(460, 35))
	cuadro_lbl_nombres_archivo := canvas.NewRectangle(color.White)
	cuadro_lbl_nombres_archivo.Move(fyne.NewPos(10, 395))
	cuadro_lbl_nombres_archivo.Resize(fyne.NewSize(460, 35))
	cuadro_lbl_nombres_archivo.FillColor = color.Transparent
	cuadro_lbl_nombres_archivo.StrokeColor = color.RGBA{200, 100, 62, 1}
	cuadro_lbl_nombres_archivo.StrokeWidth = 0.5

	btn_gSQLite := widget.NewButton("Generar BD SQLite", func() {

		dialog.Info("¡CUÁN TENTADOR ES PULSAR UN BOTÓN\nPARA VER SI 'SUENA LA FLAUTA'! ¿EH?\nPERO AÚN NO TIENE FUNCIONALIDAD.\nCOMO UN COCHE BONITO, PERO SIN MOTOR.\nJAJAJA.\nXD - ALBERTO -")
		dialog.Info("DE MOMENTO NO PUEDES CREAR BASE DE DATOS SQLITE3,\nPERO SÍ UN LISTADO EXCEL.\nDISCULPA LAS MOLESTIAS. -ALBERTO-")
	})

	btn_gExcel := widget.NewButton("Generar archivo EXCEL", func() {

		barra_de_progreso.Show()

		go func() {

			rastreadorarchivos.RastrearDispB(ruta + "/" + disp_selec)

			for indice, archivo := range rastreadorarchivos.Slice_archivos {

				fmt.Println(archivo.Ruta)
				fmt.Println(archivo.Nombre)
				time.Sleep(time.Millisecond * 10)
				progreso = float64(indice) / float64(len(rastreadorarchivos.Slice_archivos)-1)
				barra_de_progreso.SetValue(progreso)

				if len(archivo.Nombre) > 52 {

					lbl_nombres_archivo.SetText("(Este nombre no se muestra por su longitud).")

				} else {

					lbl_nombres_archivo.SetText(archivo.Nombre)

				}

				if progreso > 0.98 {

					progreso = 1.0
				}
			}

			lbl_nombres_archivo.SetText("Nombres de archivo...")

			generadorexcel.EstablecerValoresUsoDisco(float64(infoDisco.Total)/float64(usodisco.GB), float64(infoDisco.Usado)/float64(usodisco.GB), float64(infoDisco.Libre)/float64(usodisco.GB))
			generadorexcel.CrearArchivo(disp_selec, &rastreadorarchivos.Slice_archivos, getNombreUsuario())
			rastreadorarchivos.Slice_archivos = make([]rastreadorarchivos.Archivo, 0)
			barra_de_progreso.Hide()
			barra_de_progreso.SetValue(0.0)
			dialog.Info("¡Proceso terminado! El archivo ha sido generado correctamente. \n - ALBERTO -")

		}()

	})

	go func() {

		for {

			fmt.Println(<-canal)
			time.Sleep(time.Millisecond * 500)

			disps = <-canal
			fmt.Printf("Hay %d dispositivos montados.\n", len(disps))
			if barra_de_progreso.Hidden {

				btn_gExcel.Enable()
				btn_gExcel.Text = "Generar Archivo EXCEL"
				btn_gSQLite.Enable()
			} else {

				btn_gExcel.Text = "EXPLORANDO..."
				btn_gExcel.Disable()
				btn_gSQLite.Disable()
			}
			data.Reload()
			btn_gExcel.Refresh()
			btn_gSQLite.Refresh()
		}

	}()

	btn_limpiar_narch := widget.NewButton("Limpiar nombres de archivos", func() {

	})

	radbox_tipos_archivo := widget.NewRadioGroup([]string{"Todos los archivos", "Archivos de Vídeo", "Archivos de Audio", "Audio y Vídeo"}, func(seleccionado string) {

		//dialog.Info("Has seleccionado: " + seleccionado)

		switch seleccionado {

		case "Todos los archivos":

			rastreadorarchivos.Opc_tipos = 1

		case "Archivos de Vídeo":

			rastreadorarchivos.Opc_tipos = 2

		case "Archivos de Audio":

			rastreadorarchivos.Opc_tipos = 3

		case "Audio y Vídeo":

			rastreadorarchivos.Opc_tipos = 4

		}

	})

	radbox_tipos_archivo.SetSelected("Todos los archivos")

	btn_salir := widget.NewButton("Salir", func() { os.Exit(0) })

	btn_gExcel.Disable()

	lista.OnSelected = func(id int) {
		fmt.Println(disps[id])
		disp_selec = disps[id]

		ruta += "/" + disp_selec
		infoDisco = usodisco.UsoDispAlmac(&ruta)

		fmt.Printf("Total: %.2f GB\nUsado: %.2f GB\nLibre: %.2f GB\n\n", float64(infoDisco.Total)/float64(usodisco.GB),
			float64(infoDisco.Usado)/float64(usodisco.GB),
			float64(infoDisco.Libre)/float64(usodisco.GB))
		btn_gExcel.Enable()
	}

	lbl_controles := widget.NewLabel("CONTROLES:")
	lbl_controles.Alignment = fyne.TextAlign(1)
	btn_gSQLite.Enable()
	btn_limpiar_narch.Disable()

	cuadro_lista := canvas.NewRectangle(color.White)
	cuadro_lista.Move(fyne.NewPos(10, 40))
	cuadro_lista.Resize(fyne.NewSize(241, 300))
	cuadro_lista.FillColor = color.Transparent
	cuadro_lista.StrokeColor = color.White
	cuadro_lista.StrokeWidth = 0.5

	cuadro_controles := canvas.NewRectangle(color.White)
	cuadro_controles.Move(fyne.NewPos(250, 40))
	cuadro_controles.Resize(fyne.NewSize(220, 300))
	cuadro_controles.FillColor = color.Transparent
	cuadro_controles.StrokeColor = color.White
	cuadro_controles.StrokeWidth = 0.5

	barra_de_progreso.Move(fyne.NewPos(10, 350))
	barra_de_progreso.Resize(fyne.NewSize(460, 35))
	cuadro_progreso := canvas.NewRectangle(color.White)
	cuadro_progreso.Move(fyne.NewPos(10, 350))
	cuadro_progreso.Resize(fyne.NewSize(460, 35))
	cuadro_progreso.FillColor = color.Transparent
	cuadro_progreso.StrokeColor = color.RGBA{0, 200, 0, 1}
	cuadro_progreso.StrokeWidth = 0.5

	radbox_tipos_archivo.Move(fyne.NewPos(265, 150))
	radbox_tipos_archivo.Resize(fyne.NewSize(200, 100))

	contenido := container.NewWithoutLayout(cuadro_lista, cuadro_controles, cuadro_progreso, lbl_lista_disp,
		lbl_controles, radbox_tipos_archivo, barra_de_progreso, btn_gExcel, btn_gSQLite, lista,
		btn_salir, lbl_nombres_archivo, cuadro_lbl_nombres_archivo)

	lbl_lista_disp.Move(fyne.NewPos(129, 5))
	lbl_controles.Move(fyne.NewPos(lbl_controles.MinSize().Width+250, 5))
	lista.Move(fyne.NewPos(10, lbl_lista_disp.MinSize().Height+5))
	lista.Resize(fyne.NewSize(241, 300))

	btn_gExcel.Move(fyne.NewPos(260, lbl_controles.MinSize().Height+10))
	btn_gExcel.Resize(fyne.NewSize(200, 50))
	cuadro_btn_gExcel := canvas.NewRectangle(color.White)
	cuadro_btn_gExcel.Move(fyne.NewPos(260, lbl_controles.MinSize().Height+100))
	cuadro_btn_gExcel.Resize(fyne.NewSize(200, 50))
	cuadro_btn_gExcel.FillColor = color.White

	btn_gSQLite.Move(fyne.NewPos(260, btn_gExcel.Position().Y+btn_gExcel.MinSize().Height+20))
	btn_gSQLite.Resize(fyne.NewSize(200, 50))

	btn_salir.Move(fyne.NewPos(260, btn_gExcel.MinSize().Height+260))
	btn_salir.Resize(fyne.NewSize(200, 35))

	w.SetContent(contenido)

	w.Resize(fyne.NewSize(489, 450))
	w.SetFixedSize(true)
	w.ShowAndRun()
	w.Content().Refresh()

}

func getListaDispExter(canal chan []string) {

	var slice_disp_extern []string

	var si sysinfo.SysInfo

	si.GetSysInfo()

	data := si.OS.Name

	fmt.Println(string(data))

	for {

		switch os := runtime.GOOS; os {

		case "darwin":

			ruta = "/Volumes"

		case "linux":

			usuario, err := user.Current()
			if err != nil {

				log.Fatalf(err.Error())

			}

			nombre_usuario := usuario.Name

			if data == "Manjaro Linux" {

				ruta = "/run/media/" + strings.ToLower(nombre_usuario)

			} else {

				ruta = "/media/" + strings.ToLower(nombre_usuario)

			}
			//ruta = "/media/"

		}

		slice_disp_extern = []string{}

		files, err := ioutil.ReadDir(ruta)

		if err != nil {
			fmt.Println("No hay ningún volumen disponible.")
			//log.Fatal(err)
			slice_disp_extern = []string{"Sin dispositivos"}
		} else {

			for _, file := range files {

				if file.IsDir() {

					slice_disp_extern = append(slice_disp_extern, file.Name())

				}

			}

		}

		time.Sleep(time.Millisecond * 250)
		canal <- slice_disp_extern

	}

}

func getNombreSO() string {

	var nombre_sistema string

	if runtime.GOOS == "linux" {

		nombre_sistema = "LINUX."

	} else if runtime.GOOS == "darwin" {

		nombre_sistema = "MACOS."

	} else if runtime.GOOS == "Windows" {

		nombre_sistema = "WINDOWS"

	}

	return nombre_sistema

}

func getNombreUsuario() string {

	usuario, err := user.Current()

	if err != nil {

		log.Fatalf(err.Error())

	}

	nombre_usuario := strings.ToLower(usuario.Name)

	return nombre_usuario

}
