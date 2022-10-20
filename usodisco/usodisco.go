package usodisco

import (
	"fmt"
	"syscall"
)

type EstadoDispAlmac struct {
	Total uint64 `json:"total"`
	Usado uint64 `json:"usado"`
	Libre uint64 `json:"libre"`
}

func UsoDispAlmac(ruta *string) (Dispositivo EstadoDispAlmac) {

	sis_archivos := syscall.Statfs_t{}
	err := syscall.Statfs(*ruta, &sis_archivos)
	if err != nil {

		fmt.Println(err)
		return
	}

	Dispositivo.Total = sis_archivos.Blocks * uint64(sis_archivos.Bsize)
	Dispositivo.Libre = sis_archivos.Bfree * uint64(sis_archivos.Bsize)
	Dispositivo.Usado = Dispositivo.Total - Dispositivo.Usado

	return
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
	TB = 1024 * GB
)
