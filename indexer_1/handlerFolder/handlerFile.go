package handlerfolder

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AdrianaPerezC/indexer_1/models"
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

// GetAllFiles recorre el directorio y obtiene información de todos los archivos
func GetAllFiles(rootDir string) ([]FileInfo, error) {
	var files []FileInfo

	// Verificar si el directorio existe
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("el directorio '%s' no existe", rootDir)
	}

	// Función para recorrer el directorio
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error al acceder a '%s': %v", path, err)
		}

		// Agregar solo los archivos a la lista
		if !info.IsDir() {
			files = append(files, FileInfo{
				Path: path,
				Name: info.Name(),
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

// ProcessFiles procesa una lista de paths de archivos y genera un JSON con los datos
func ProcessFiles(paths []FileInfo) error {

	for _, path := range paths {

		go func(path string) {
			email, err := ParseEmail(path)
			if err != nil {
				fmt.Printf("Error procesando archivo %s: %v\n", path, err)
			}
			// Agregar el resultado a la lista de emails
			Emails = append(Emails, email)
		}(path.Path) // Pasar path por valor para evitar conflictos
	}

	// Guardar los resultados en un archivo JSON
	// jsonData, err := json.MarshalIndent(Emails, "", "  ")
	// if err != nil {
	// 	return fmt.Errorf("error al convertir datos a JSON: %v", err)
	// }

	// err = os.WriteFile(outputFile, jsonData, 0644)
	// if err != nil {
	// 	return fmt.Errorf("error al escribir archivo JSON: %v", err)
	// }

	fmt.Printf("Datos procesados...")
	return nil
}
