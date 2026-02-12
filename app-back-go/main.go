package main

import (
	"app-back-go/api/handlers"
	"app-back-go/internal/service"
	"fmt"
	"log"
	"net/http"
)

func main() {

	// Inicializar el servicio de descarga
	downloadService := service.NewDownloadService()

	// Inicializar el handler de descarga con el servicio inyectado
	downloadHandler := handlers.NewDownloadHandler(downloadService)
	// Inicializar el handler de metadatos con el servicio inyectado
	metadataHandler := handlers.NewMetadataHandler(downloadService)

	// Definir las rutas y sus manejadores
	http.HandleFunc("/download", downloadHandler.HandleDownload)
	http.HandleFunc("/metadata", metadataHandler.HandleGetMetadata) // Nueva ruta para metadatos

	// Configurar el puerto del servidor
	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	// Iniciar el servidor
	log.Fatal(http.ListenAndServe(port, nil))
}
