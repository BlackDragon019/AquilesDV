package main

import (
	"app-back-go/api/handlers"
	"app-back-go/internal/service"
	"app-back-go/internal/tools"
	"fmt"
	"log"
	"net/http"
)

// corsMiddleware agrega headers CORS a todas las respuestas
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Manejar preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	// Asegurarse de que yt-dlp está disponible (descargarlo si no existe)
	_, err := tools.EnsureYtDlp()
	if err != nil {
		log.Fatalf("Error al obtener yt-dlp: %v", err)
	}

	// Inicializar el servicio de descarga
	downloadService := service.NewDownloadService()

	// Inicializar el handler de descarga con el servicio inyectado
	downloadHandler := handlers.NewDownloadHandler(downloadService)
	// Inicializar el handler de metadatos con el servicio inyectado
	metadataHandler := handlers.NewMetadataHandler(downloadService)

	// Crear un mux para las rutas
	mux := http.NewServeMux()

	// Definir las rutas y sus manejadores
	mux.HandleFunc("/download", downloadHandler.HandleDownload)
	mux.HandleFunc("/metadata", metadataHandler.HandleGetMetadata) // Nueva ruta para metadatos

	// Aplicar middleware CORS
	handler := corsMiddleware(mux)

	// Configurar el puerto del servidor
	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	// Iniciar el servidor
	log.Fatal(http.ListenAndServe(port, handler))
}
