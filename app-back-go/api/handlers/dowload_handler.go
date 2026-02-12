package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

	// Enviar una respuesta de éxito al cliente con la ruta del archivo local
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":         "Video descargado exitosamente a la ubicación local",
		"original_url":    req.URL,
		"local_file_path": localFilePath, // Añadimos la ruta del archivo local aquí
	})
}
