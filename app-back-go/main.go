package main

import (
	"app-back-go/api/handlers"
	"app-back-go/internal/service"
	"fmt"
	"log"
	"net/http"
)

func main() {

	downloadService := service.NewDownloadService()
	downloadHandler := handlers.NewDownloadHandler(downloadService)
	http.HandleFunc("/download", downloadHandler.HandleDownload)

	port := 8080
	fmt.Printf("Servidor escuchando en http://localhost:%d\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
