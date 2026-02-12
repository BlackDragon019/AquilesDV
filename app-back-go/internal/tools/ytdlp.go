package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

// GitHubRelease contiene información sobre la última release
type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

// GetLatestYtDlpVersion obtiene el tag de la última versión de yt-dlp desde GitHub
func GetLatestYtDlpVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/yt-dlp/yt-dlp/releases/latest")
	if err != nil {
		return "", fmt.Errorf("error al obtener información de releases: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error al obtener releases: código de estado %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("error al parsear releases: %v", err)
	}

	return release.TagName, nil
}

// EnsureYtDlp descarga yt-dlp si no existe en el directorio tools
func EnsureYtDlp() (string, error) {
	// Crear directorio tools si no existe
	toolsDir := filepath.Join(".", "tools")
	if err := os.MkdirAll(toolsDir, 0755); err != nil {
		return "", fmt.Errorf("error al crear directorio tools: %v", err)
	}

	// Determinar el nombre del ejecutable según el SO
	var ytdlpName string
	var downloadURL string

	// Obtener la última versión
	version, err := GetLatestYtDlpVersion()
	if err != nil {
		fmt.Printf("Advertencia: no se pudo obtener la última versión, usando fallback: %v\n", err)
		version = "2024.12.06" // Fallback a una versión conocida
	}

	if runtime.GOOS == "windows" {
		ytdlpName = "yt-dlp.exe"
		downloadURL = fmt.Sprintf("https://github.com/yt-dlp/yt-dlp/releases/download/%s/yt-dlp.exe", version)
	} else {
		ytdlpName = "yt-dlp"
		downloadURL = fmt.Sprintf("https://github.com/yt-dlp/yt-dlp/releases/download/%s/yt-dlp", version)
	}

	ytdlpPath := filepath.Join(toolsDir, ytdlpName)

	// Si ya existe, devolverlo
	if _, err := os.Stat(ytdlpPath); err == nil {
		fmt.Printf("yt-dlp encontrado en: %s\n", ytdlpPath)
		return ytdlpPath, nil
	}

	fmt.Printf("Descargando yt-dlp (%s) desde: %s\n", version, downloadURL)

	// Descargar yt-dlp
	resp, err := http.Get(downloadURL)
	if err != nil {
		return "", fmt.Errorf("error al descargar yt-dlp: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error al descargar yt-dlp: código de estado %d (URL: %s)", resp.StatusCode, downloadURL)
	}

	// Crear archivo
	out, err := os.Create(ytdlpPath)
	if err != nil {
		return "", fmt.Errorf("error al crear archivo yt-dlp: %v", err)
	}
	defer out.Close()

	// Copiar contenido
	if _, err := io.Copy(out, resp.Body); err != nil {
		os.Remove(ytdlpPath)
		return "", fmt.Errorf("error al escribir archivo yt-dlp: %v", err)
	}

	// Hacer ejecutable (en Windows no es necesario, pero en Linux sí)
	if runtime.GOOS != "windows" {
		if err := os.Chmod(ytdlpPath, 0755); err != nil {
			return "", fmt.Errorf("error al hacer yt-dlp ejecutable: %v", err)
		}
	}

	fmt.Printf("yt-dlp descargado exitosamente en: %s\n", ytdlpPath)
	return ytdlpPath, nil
}
