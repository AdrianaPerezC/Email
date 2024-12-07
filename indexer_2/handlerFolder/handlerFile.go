package handlerfolder

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/AdrianaPerezC/indexer_2/models"
)

var Emails []models.Email

type FileInfo struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

func AllFiles(path string) ([]string, error) {
	rootDir := path
	folders, err := GetFolderNames(rootDir)
	if err != nil {
		return nil, fmt.Errorf("error al obtener las carpetas: %v\n", err)
	}
	return folders, nil
}

func GetAllFiles(rootDir string) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error al acceder a '%s': %v", path, err)
		}

		if !d.IsDir() {
			files = append(files, FileInfo{
				Path: path,
				Name: d.Name(),
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// GetFolderNames recorre un directorio y devuelve una lista con los nombres de las carpetas
func GetFolderNames(rootDir string) ([]string, error) {
	var folders []string

	// Función para recorrer el directorio
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Verificar si es un directorio (y no el rootDir)
		if info.IsDir() && path != rootDir {
			folders = append(folders, info.Name())
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return folders, nil
}

func ParseEmail(filePath string) (models.Email, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return models.Email{}, fmt.Errorf("error al abrir archivo %s: %v", filePath, err)
	}
	defer file.Close()

	email := models.Email{}
	var contentBuilder strings.Builder
	scanner := bufio.NewScanner(file)
	const maxBufferSize = 10 * 1024 * 1024 // 10 MB
	buf := make([]byte, maxBufferSize)
	scanner.Buffer(buf, maxBufferSize)
	inHeader := true

	for scanner.Scan() {
		line := scanner.Text()
		if inHeader {
			if line == "" { // Fin de la cabecera
				inHeader = false
				continue
			}
			// Procesar las líneas de la cabecera
			parts := strings.SplitN(line, ":", 2)
			if len(parts) < 2 {
				continue
			}
			key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
			switch key {
			case "Message-ID":
				email.MessageID = value
			case "Date":
				email.Date = value
			case "From":
				email.From = value
			case "To":
				email.To = value
			case "Subject":
				email.Subject = value
			case "Mime-Version":
				email.Mime_Version = value
			case "Content-Type":
				email.Content_Type = value
			case "Content-Transfer-Encoding":
				email.Content_Transfer_Encoding = value
			case "X-From":
				email.X_From = value
			case "X-To":
				email.X_To = value
			case "X-cc":
				email.X_cc = value
			case "X-bcc":
				email.X_bcc = value
			case "X-Folder":
				email.X_Folder = value
			case "X-Origin":
				email.X_Origin = value
			case "X-FileName":
				email.X_FileName = value
			}
		} else {

			contentBuilder.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return models.Email{}, fmt.Errorf("error al leer archivo %s: %v", filePath, err)
	}

	email.Body = contentBuilder.String()
	return email, nil
}

const batchSize = 500 // Tamaño del lote

func ProcessFiles(paths []FileInfo, maxWorkers int, indexName string) error {
	// Canal para limitar la concurrencia
	sem := make(chan struct{}, maxWorkers)
	var (
		wg           sync.WaitGroup
		currentBatch []models.Email
		mu           sync.Mutex
	)

	// Recorrer los archivos
	for _, path := range paths {
		wg.Add(1)
		sem <- struct{}{} // Adquirir un "slot" en el semáforo (límite de trabajadores concurrentes)

		go func(path string) {
			defer wg.Done()
			defer func() { <-sem }() // Liberar el "slot" en el semáforo

			// Parsear el archivo y obtener el email
			email, err := ParseEmail(path)
			if err != nil {
				fmt.Printf("Error procesando archivo %s: %v\n", path, err)
				return
			}

			// Proteger el acceso concurrente a currentBatch
			mu.Lock()
			currentBatch = append(currentBatch, email)

			// Si el lote alcanza el tamaño definido, enviarlo a ZincSearch
			if len(currentBatch) >= batchSize {
				// Enviar el lote
				sendBatchToZincSearch(indexName, currentBatch)
				// Limpiar el lote actual
				currentBatch = nil
			}
			mu.Unlock()
		}(path.Path)
	}

	// Esperar a que todas las goroutines terminen
	wg.Wait()

	// Si quedan correos electrónicos en el lote después de procesar todos los archivos, enviarlos
	if len(currentBatch) > 0 {
		sendBatchToZincSearch(indexName, currentBatch)
	}

	return nil
}

func sendBatchToZincSearch(indexName string, batch []models.Email) {
	// Usamos jsonToNDJSON para crear el archivo NDJSON
	filePath, err := jsonToNDJSON(indexName, batch)
	if err != nil {
		fmt.Printf("Error al crear el archivo NDJSON: %v\n", err)
		return
	}

	// Leer los datos desde el archivo NDJSON
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error al leer el archivo NDJSON: %v\n", err)
		return
	}
	url := os.Getenv("ZINC_URL") + "/api/_bulkv2"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		fmt.Printf("Error al crear la solicitud: %v\n", err)
		return
	}

	// Establecer autenticación
	req.SetBasicAuth(os.Getenv("ZINC_USER"), os.Getenv("ZINC_PASSWORD"))
	req.Header.Set("Content-Type", "application/octet-stream")

	// Realizar la petición HTTP
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error al enviar la solicitud: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Leer la respuesta de ZincSearch
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Respuesta de ZincSearch:", string(body))

	// Verificar si hubo error en la respuesta
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error en la respuesta de ZincSearch: %s\n", resp.Status)
	}
}

func jsonToNDJSON(index string, batch []models.Email) (string, error) {
	var r models.Request

	r.Index = index
	r.Records = batch

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

// func sendBatchToZincSearch(batch []models.Email) {
// 	// Serializar el lote y enviarlo a ZincSearch
// 	data, err := json.Marshal(batch)
// 	if err != nil {
// 		fmt.Printf("Error al serializar el lote: %v\n", err)
// 		return
// 	}

// 	url := os.Getenv("ZINC_URL") + "/api/_bulkv2"
// 	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
// 	if err != nil {
// 		fmt.Printf("Error al crear la solicitud: %v\n", err)
// 		return
// 	}

// 	req.SetBasicAuth(os.Getenv("ZINC_USER"), os.Getenv("ZINC_PASSWORD"))
// 	req.Header.Set("Content-Type", "application/json")

// 	client := http.DefaultClient
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Printf("Error al enviar el lote a ZincSearch: %v\n", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		fmt.Printf("Error en la respuesta de ZincSearch: %s\n", resp.Status)
// 	}
// }
