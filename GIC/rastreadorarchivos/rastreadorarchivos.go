package rastreadorarchivos

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

const n_directorios_ruta_montados = 4

var nombre string

type Archivo struct {
	Indice    int
	Nombre    string
	Ruta      string
	Extension string
}

var archivo Archivo
var Slice_archivos []Archivo
var slice_ruta []string
var slice_ruta_reducido []string
var N_archivos int = 0
var slice_formatos []string
var Opc_tipos int

func (archivo Archivo) Agregar() {

	Slice_archivos = append(Slice_archivos, archivo)
	N_archivos++
}

func printFile(path string, info os.FileInfo, opc_tipos int, err error) error {

	if err != nil {
		log.Print(err)
		return nil
	}

	slice_ruta = strings.Split(path, "/")
	nombre = slice_ruta[len(slice_ruta)-1]

	switch opc_tipos {

	case 1:

		slice_formatos = []string{".*"}

	case 2:

		slice_formatos = []string{".avi", ".mp4", ".mkv", ".mpg", ".mov", ".mpeg", ".wmv"}

	case 3:

		slice_formatos = []string{".mp3", ".ogg", ".flac", ".wav"}

	}

	var nombre_simp string

	for _, formato := range slice_formatos {

		if filepath.Ext(path) == formato && opc_tipos != 1 {

			for indice, caracter := range nombre {

				if caracter == '.' && indice < len(nombre)-4 {

					nombre_simp += " "

				} else if caracter == '(' || caracter == '[' {

					nombre_simp += filepath.Ext(path)

					break

				} else {

					nombre_simp += string(caracter)
				}

				archivo.Indice = indice
			}

			slice_ruta_reducido = slice_ruta[n_directorios_ruta_montados : len(slice_ruta)-1]

			for _, elemento := range slice_ruta_reducido {

				archivo.Ruta += string("/" + elemento)

			}

			archivo.Nombre = string(nombre_simp)
			archivo.Extension = filepath.Ext(path)
			archivo.Agregar()
			archivo.Ruta = ""

		} else if opc_tipos == 1 {

			archivo.Nombre = string(nombre)
			archivo.Extension = filepath.Ext(path)
			archivo.Agregar()
			archivo.Ruta = ""

		}

	}

	return nil

}

func RastrearDispB(ruta string) {

	err := filepath.Walk(ruta,

		func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			printFile(path, info, Opc_tipos, err)

			return nil

		})

	if err != nil {

		log.Println(err)

	}

}
