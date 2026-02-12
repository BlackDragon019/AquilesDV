package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"app-back-go/api/models"
	"app-back-go/internal/service"
)

type DownloadHandler struct {
	downloadService service.DownloadService
}

func NewDownloadHandler(s service.DownloadService) *DownloadHandler {
	return &DownloadHandler{
		downloadService: s,
	}
}

func (h *DownloadHandler) HandleDownload(w http.ResponseWriter, r *http.Request) {
	// Asegurarse de que solo aceptamos solicitudes POST
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.DownloadRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Error al decodificar la solicitud: %v\n", err)
		http.Error(w, "Solicitud JSON inválida", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "El campo 'url' es requerido", http.StatusBadRequest)
		return
	}

	// Llamar al servicio para procesar la descarga
	// Ahora ProcessVideoDownload devuelve la RUTA AL ARCHIVO LOCAL
	localFilePath, err := h.downloadService.ProcessVideoDownload(req.URL)
	if err != nil {
		log.Printf("Error en el servicio de descarga: %v\n", err)
		http.Error(w, fmt.Sprintf("Error al procesar la descarga: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	defer func() {
		if err := os.Remove(localFilePath); err != nil {
			log.Printf("Error al eliminar el archivo temporal %s: %v", localFilePath, err)
		} else {
			fmt.Printf("Archivo temporal %s eliminado exitosamente.\n", localFilePath)
		}
	}()

	file, err := os.Open(localFilePath)
	if err != nil {
		log.Printf("Error al abrir el archivo local %s: %v", localFilePath, err)
		http.Error(w, "Error interno al acceder al video descargado", http.StatusInternalServerError)
		return
	}
	defer file.Close() // Cerrar el archivo después de servirlo

	// Obtener información del archivo para determinar el tamaño y el nombre
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("Error al obtener información del archivo %s: %v", localFilePath, err)
		http.Error(w, "Error interno al obtener información del video", http.StatusInternalServerError)
		return
	}

	// Enviar una respuesta de éxito al cliente con la ruta del archivo local
	w.Header().Set("Content-Type", "video/mp4")
	downloadFileName := filepath.Base(localFilePath)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", downloadFileName))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	fmt.Printf("Sirviendo archivo %s al cliente. Tamaño: %d bytes.\n", downloadFileName, fileInfo.Size())

	// Enviar el archivo al cliente
	http.ServeContent(w, r, downloadFileName, fileInfo.ModTime(), file)
}
