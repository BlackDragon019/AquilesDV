package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type DownloadService interface {
	ProcessVideoDownload(videoURL string) (string, error)
}

type dowloaderService struct {
}

func NewDownloadService() DownloadService {
	return &dowloaderService{}
}

type YtDlpOutput struct {
	Ext   string `json:"ext"`   // Extensión del archivo (ej. mp4)
	Title string `json:"title"` // Título del video
}

func (s *dowloaderService) ProcessVideoDownload(videoURL string) (string, error) {
	fmt.Printf("Servicio: Iniciando procesamiento de descarga para URL: %s\n", videoURL)

	if videoURL == "" {
		return "", errors.New("URL de video no puede estar vacía")
	}

	tempDir := "downloads" // Directorio donde se guardarán los videos

	// Asegurarse de que el directorio de descargas existe
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("error al crear el directorio de descargas: %v", err)
	}

	// --- 1. Obtener el título y extensión del video usando yt-dlp (sin descargarlo aún) ---
	// Esto nos permite crear un nombre de archivo seguro ANTES de la descarga
	cmdInfo := exec.Command("yt-dlp", "--print-json", "--flat-playlist", "--skip-download", videoURL)

	var stdoutInfo, stderrInfo bytes.Buffer
	cmdInfo.Stdout = &stdoutInfo
	cmdInfo.Stderr = &stderrInfo

	err := cmdInfo.Run()
	if err != nil {
		logError := fmt.Errorf("error al obtener info con yt-dlp para URL %s: %v\nStderr: %s", videoURL, err, stderrInfo.String())
		fmt.Printf("ERROR: %s\n", logError)
		return "", fmt.Errorf("no se pudo obtener información del video. Error: %s", strings.TrimSpace(stderrInfo.String()))
	}

	var ytDlpInfo YtDlpOutput
	err = json.Unmarshal(stdoutInfo.Bytes(), &ytDlpInfo)
	if err != nil {
		logError := fmt.Errorf("error al decodificar la salida JSON de yt-dlp info para URL %s: %v\nStdout: %s", videoURL, err, stdoutInfo.String())
		fmt.Printf("ERROR: %s\n", logError)
		return "", errors.New("error al parsear la información de yt-dlp")
	}

	if ytDlpInfo.Ext == "" {
		ytDlpInfo.Ext = "mp4" // Fallback si no se obtiene la extensión
	}
	if ytDlpInfo.Title == "" {
		ytDlpInfo.Title = "video_descargado" // Fallback si no se obtiene el título
	}
	// Crear un nombre de archivo seguro
	fileName := fmt.Sprintf("%s.%s", sanitizeFilename(ytDlpInfo.Title), "mp4")
	filePath := filepath.Join(tempDir, fileName)

	fmt.Printf("Servicio: Título del video: %s, Extensión: %s\n", ytDlpInfo.Title, ytDlpInfo.Ext)
	fmt.Printf("Servicio: Intentando descargar directamente con yt-dlp a: %s\n", filePath)

	// --- 2. Descargar el video directamente con yt-dlp a la ruta especificada ---
	// -o: Especifica el nombre de archivo de salida
	// --restrict-filenames: Ayuda a evitar caracteres problemáticos en nombres de archivo
	cmdDownload := exec.Command("yt-dlp",
		"-f", "bv*[ext=mp4][vcodec!=hevc][vcodec!=h265]/bv*[ext=mp4]/b[ext=mp4]/best",
		"-o", filePath,
		"--restrict-filenames",
		videoURL,
	)

	var stdoutDownload, stderrDownload bytes.Buffer
	cmdDownload.Stdout = &stdoutDownload
	cmdDownload.Stderr = &stderrDownload

	err = cmdDownload.Run()
	if err != nil {
		logError := fmt.Errorf("error al descargar el video con yt-dlp para URL %s: %v\nStderr: %s", videoURL, err, stderrDownload.String())
		fmt.Printf("ERROR: %s\n", logError)
		// Intentar borrar el archivo incompleto si se creó
		os.Remove(filePath)
		return "", fmt.Errorf("no se pudo descargar el video. Error: %s", strings.TrimSpace(stderrDownload.String()))
	}

	fmt.Printf("Servicio: Video descargado exitosamente por yt-dlp a: %s\n", filePath)
	return filePath, nil // Retornar la ruta del archivo local
}

// sanitizeFilename limpia una cadena para usarla como nombre de archivo seguro
func sanitizeFilename(s string) string {
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "\\", "_")
	s = strings.ReplaceAll(s, ":", "_")
	s = strings.ReplaceAll(s, "*", "_")
	s = strings.ReplaceAll(s, "?", "_")
	s = strings.ReplaceAll(s, "\"", "_")
	s = strings.ReplaceAll(s, "<", "_")
	s = strings.ReplaceAll(s, ">", "_")
	s = strings.ReplaceAll(s, "|", "_")
	s = strings.ReplaceAll(s, " ", "_") // Reemplazar espacios por guiones bajos
	// Limitar longitud para evitar problemas con algunos sistemas de archivos
	if len(s) > 100 {
		s = s[:100]
	}
	return s
}
