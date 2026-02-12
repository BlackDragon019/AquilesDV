package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"app-back-go/internal/service"
)

// MetadataHandler depende del servicio de descarga para obtener metadatos
type MetadataHandler struct {
	Service service.DownloadService
}

// NewMetadataHandler crea un nuevo handler de metadatos con el servicio inyectado
func NewMetadataHandler(s service.DownloadService) *MetadataHandler {
	return &MetadataHandler{Service: s}
}

// HandleGetMetadata maneja las solicitudes GET para obtener metadatos de videos
func (h *MetadataHandler) HandleGetMetadata(w http.ResponseWriter, r *http.Request) {
	// Asegurarse de que solo aceptamos solicitudes GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// La URL del video se pasará como parámetro de consulta
	videoURL := r.URL.Query().Get("url")
	if videoURL == "" {
		http.Error(w, "Parámetro 'url' no proporcionado", http.StatusBadRequest)
		return
	}

	// Llamar al servicio para obtener los metadatos
	metadata, err := h.Service.GetVideoMetadata(videoURL)
	if err != nil {
		log.Printf("Error al obtener metadatos para URL %s: %v\n", videoURL, err)
		http.Error(w, fmt.Sprintf("Error al obtener metadatos: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Enviar los metadatos como respuesta JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metadata)
}
