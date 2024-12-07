package main

import (
	"fmt"
	"os"
	"runtime"

	handlerfolder "github.com/AdrianaPerezC/indexer_2/handlerFolder"
)

// Documento a indexar
type Document struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func main() {
	// Inicializar y obtener el directorio de los argumentos

	directory, err := argumentInit()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Leer todos los archivos dentro del directorio
	files, err := handlerfolder.GetAllFiles(directory)
	if err != nil {
		fmt.Printf("Error al leer archivos: %v\n", err)
		os.Exit(1)
	}
	//Número de nucleos disponibles en la PC
	workers := runtime.NumCPU()
	handlerfolder.ProcessFiles(files, workers, "enron_mail")
}

// argumentInit valida los argumentos de línea de comandos y devuelve el directorio
func argumentInit() (string, error) {
	if len(os.Args) != 2 {
		return "", fmt.Errorf("ingrese el formato: ./indexer <directorio>")
	}
	return os.Args[1], nil
}
func ProfilingTest(directory string, indexName string) int {
	files, err := handlerfolder.GetAllFiles(directory)
	if err != nil {
		fmt.Printf("Error al leer archivos: %v\n", err)
		os.Exit(1)
	}
	//Número de nucleos disponibles en la PC
	workers := runtime.NumCPU()

	handlerfolder.ProcessFiles(files, workers, "enron_mail")
	return 1
}
