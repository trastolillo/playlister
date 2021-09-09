package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var carpeta string
var nombresCompletos []string

func main() {

	// Proceso las opciones del comando
	procesarFlags()

	// Obtengo la lista de los vídeos
	listaArchivos := fileList(carpeta)

	// Ordeno la lista según el número de vídeo
	sort.Slice(listaArchivos, func(p, q int) bool {
		return listaArchivos[p].numero < listaArchivos[q].numero
	})

	// Obtengo la lista de los nombres completos
	for _, archivo := range listaArchivos {
		nombresCompletos = append(nombresCompletos, archivo.nombreCompleto)
	}
	// fmt.Println(carpeta + "/playlist.m3u")

	// Escribo los nombres de los vídeos a un archivo
	playlist, err := os.Create(carpeta + "/playlist.m3u")
	if err != nil {
		fmt.Println(err)
		playlist.Close()
		os.Exit(2)
	}
	for _, ruta := range nombresCompletos {
		_, err = fmt.Fprintln(playlist, ruta)
		if err != nil {
			log.Fatal(err)
			playlist.Close()
		}
		fmt.Println("Añadido:", ruta)
	}
	playlist.Close()
	fmt.Println("Lista de reproducción creada con éxito")

}

type Archivo struct {
	numero         byte
	nombre         string
	nombreCompleto string
}

func (this *Archivo) extraerNumero() byte {
	re := regexp.MustCompile(`^[0-9]+`)
	resultadoString := re.Find([]byte(this.nombre))
	resultadoInt, err := strconv.ParseUint(string(resultadoString), 10, 8)
	if err != nil {
		log.Fatal(err)
	}
	this.numero = byte(resultadoInt)
	return byte(resultadoInt)
}

func procesarFlags() {
	// Declaraciones de Flag
	flag.StringVar(&carpeta, "d", "./", "Ubicación de los vídeos a incluir en la playlist (por defecto './')")
	// Validación
	if carpeta == "" {
		flag.Usage()
		os.Exit(2)
	}
	// Análisis de opciones
	flag.Parse()
}

func fileList(carpeta string) []Archivo {
	var sliceArchivos []Archivo
	files, err := ioutil.ReadDir(carpeta)
	// Manejo del error
	if err != nil {
		log.Fatal(err)
	}
	// Lista de archivos
	for _, file := range files {
		if !file.IsDir() {
			extension := path.Ext(file.Name())
			if extension == ".mp4" {
				rutaAbsoluta, err := filepath.Abs(carpeta)
				if err != nil {
					log.Fatal(err)
				}
				nombreCompleto := filepath.Join(rutaAbsoluta, file.Name())
				nombre := strings.TrimSuffix(file.Name(), extension)
				nuevoArchivo := Archivo{nombre: nombre, nombreCompleto: nombreCompleto}
				nuevoArchivo.extraerNumero()
				sliceArchivos = append(sliceArchivos, nuevoArchivo)
			}
		}
	}
	return sliceArchivos
}
