package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	handlerfolder "github.com/AdrianaPerezC/indexer_1/handlerFolder"
	"github.com/AdrianaPerezC/indexer_1/models"
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

	handlerfolder.ProcessFiles(files)
	sendToZincSearch("enron_mails")
}

// argumentInit valida los argumentos de línea de comandos y devuelve el directorio
func argumentInit() (string, error) {
	if len(os.Args) != 2 {
		return "", fmt.Errorf("ingrese el formato: ./indexer <directorio>")
	}
	return os.Args[1], nil
}

func sendToZincSearch(indexName string) {
	filePath, err := jsonToNDJSON(indexName)
	if err != nil {
		fmt.Println("Error al crear el archivo ndjson:", err)
		return
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return
	}

	url := "http://localhost:4080/api/_bulkv2"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		fmt.Println("Error al crear la solicitud:", err)
		return
	}

	// Establecer la autenticación
	req.SetBasicAuth("admin", "123")

	// Establecer encabezados
	req.Header.Set("Content-Type", "application/octet-stream")

	// Realizar la petición
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error al hacer la petición:", err)
		return
	}
	defer resp.Body.Close()

	// Leer la respuesta
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error al leer la respuesta:", err)
		return
	}

	// Imprimir la respuesta
	fmt.Println("Respuesta:", string(respBody))

}
func jsonToNDJSON(index string) (string, error) {
	var r models.Request

	r.Index = index
	r.Records = handlerfolder.Emails

	fileName := r.Index + ".ndjson" // Nombre del archivo NDJSON
	file, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("error al crear el archivo ndjson: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	jsonData, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		return "", fmt.Errorf("error al convertir a json: %v", err)
	}

	_, err = writer.Write(jsonData)
	if err != nil {
		return "", fmt.Errorf("error al escribir en el archivo ndjson: %v", err)
	}

	err = writer.Flush()
	if err != nil {
		return "", fmt.Errorf("error al vaciar el búfer en el archivo ndjson: %v", err)
	}

	return fileName, nil
}

func ProfilingTest(directory string, indexName string) int {
	files, err := handlerfolder.GetAllFiles(directory)
	if err != nil {
		fmt.Printf("Error al leer archivos: %v\n", err)
		os.Exit(1)
	}

	handlerfolder.ProcessFiles(files)
	sendToZincSearch(indexName)
	return 1
}
