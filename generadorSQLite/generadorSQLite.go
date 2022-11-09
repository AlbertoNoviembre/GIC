package generadorSQLite

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var fecha string
var hora string

type Dispositivo struct {
	id int
	//mu         sync.Mutex
	Nombre           string
	Total            float64
	Utilizado        float64
	Disponible       float64
	Fecha_Actualizac string
	Hora_Actualizac  string
}
type Archivo struct {
	mu     sync.Mutex
	id     int
	Nombre string
	Extens string
	Ruta   string
}

const tabla_dispositivos string = `

		CREATE TABLE IF NOT EXISTS DispositivosAlmac (

			id INTEGER PRIMARY KEY AUTOINCREMENT,
			Nombre TEXT,
			Espacio_Total REAL,
			Espacio_Utilizado REAL,
			Espacio_Disponible REAL,
			Fecha_Actualizac TEXT,
			Hora_Actualizac TEXT
			
			

);`

const tabla_archivos string = `

		CREATE TABLE IF NOT EXISTS Archivos (

			id INTEGER NOT NULL PRIMARY KEY,
			Nombre TEXT,
			Extens TEXT,	
			Ruta TEXT
			

);`

func CrearDB() {

	db, _ = sql.Open("sqlite3", "gic.db")

	statement, err := db.Prepare(tabla_dispositivos)

	if err != nil {

		log.Fatal(err)

	}

	statement.Exec()

	statement, err = db.Prepare(tabla_archivos)

	if err != nil {

		log.Fatal(err)

	}

	statement.Exec()

	fmt.Println("¡Bien! ¡Has creado el archivo de base de datos!")

}

func GetDispositivosRegistrados() []Dispositivo {

	var disps []Dispositivo = nil
	var disp Dispositivo
	var id int
	var Nombre string
	var Espacio_Total, Espacio_Utilizado, Espacio_Disponible float64
	var Fecha_Actualizac, Hora_Actualizac string

	fila, err := db.Query("SELECT * FROM DispositivosAlmac ORDER BY Nombre;")

	if err != nil {

		log.Fatal(err)

	}

	defer fila.Close()

	for fila.Next() {

		fila.Scan(&id, &Nombre, &Espacio_Total, &Espacio_Utilizado, &Espacio_Disponible, &Fecha_Actualizac, &Hora_Actualizac)
		disp.id = id
		disp.Nombre = Nombre
		disp.Total = Espacio_Total
		disp.Utilizado = Espacio_Utilizado
		disp.Disponible = Espacio_Disponible
		disp.Fecha_Actualizac = Fecha_Actualizac
		disp.Hora_Actualizac = Hora_Actualizac
		disps = append(disps, disp)

	}

	return disps
}

func Actualizar(nombre_tabla string, etiqueta_disco string, espacio_total float64, espacio_utilizado float64, espacio_disponible float64) {

	fecha = fmt.Sprintf("%d-%d-%d", time.Now().Day(), time.Now().Month(), time.Now().Year())
	hora = fmt.Sprintf("%d:%d:%d", time.Now().Hour(), time.Now().Minute(), time.Now().Second())

	statement, err := db.Prepare(`UPDATE DispositivosAlmac SET Espacio_Total=?,
		Espacio_Utilizado=?,
		Espacio_Disponible=?,
		Fecha_Actualizac=?,
		Hora_Actualizac=?
		WHERE Nombre=? 
	`)

	if err != nil {

		fmt.Println("Se produjo un error al intentar actualizar registros en la tabla " + nombre_tabla)
		fmt.Println(err)

	}

	statement.Exec(espacio_total, espacio_utilizado, espacio_disponible, fecha, hora, etiqueta_disco)
}

func Insertar(nombre_tabla string, etiqueta_disco string, espacio_total float64, espacio_utilizado float64, espacio_disponible float64) {

	fecha = fmt.Sprintf("%d-%d-%d", time.Now().Day(), time.Now().Month(), time.Now().Year())
	hora = fmt.Sprintf("%d:%d:%d", time.Now().Hour(), time.Now().Minute(), time.Now().Second())

	//var valores string = fmt.Sprintf("%s,%.2f,%.2f,%.2f", etiqueta_disco, espacio_total, espacio_utilizado, espacio_disponible)
	//	statement, _ := db.Prepare("INSERT INTO DispositivosAlmac VALUES ('1','Disco1','300','200','100');")
	statement, err := db.Prepare("INSERT INTO DispositivosAlmac (Nombre,Espacio_Total,Espacio_Utilizado,Espacio_Disponible,Fecha_Actualizac,Hora_Actualizac) VALUES (?,?,?,?,?,?);")

	if err != nil {

		fmt.Println("Se produjo un error al intentar insertar registros en la tabla " + nombre_tabla)
		fmt.Println(err)

	}

	statement.Exec(etiqueta_disco, espacio_total, espacio_utilizado, espacio_disponible, fecha, hora)

}
