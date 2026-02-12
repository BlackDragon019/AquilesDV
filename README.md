# 🎬 Descargador de Videos - TikTok & Instagram

Una aplicación web moderna para descargar videos de **TikTok** e **Instagram** de forma rápida y sencilla.

## 📋 Características

✨ **Descarga fácil**
- Pega la URL del video y descárgalo automáticamente
- Soporte para TikTok e Instagram
- Descarga directa en formato MP4

⚡ **Rendimiento optimizado**
- Carga rápida de previsualización (metadatos e imagen miniatura)
- Interfaz responsiva e intuitiva
- Sin instalación de herramientas externas requerida

🎥 **Vista previa**
- Visualiza el título y miniatura del video
- Acceso directo a la URL original
- Previsualización instantánea

🔧 **Tecnología moderna**
- Backend en Go
- Frontend en React
- yt-dlp descargado automáticamente
- CORS habilitado para comunicación fluida

## 🚀 Instalación

### Requisitos previos
- **Node.js** (v14 o superior)
- **Go** (v1.25 o superior)
- **Git**

### Clonar el repositorio

```bash
git clone https://github.com/tu-usuario/AquilesDV.git
cd AquilesDV
```

### Instalar dependencias del frontend

```bash
cd app-front-react
npm install
```

## 🏃 Ejecutar la aplicación

### Terminal 1: Iniciar el backend (Go)

```bash
cd app-back-go
go build -o app.exe
.\app.exe
```

El servidor estará disponible en `http://localhost:8080`

El backend descargará automáticamente `yt-dlp` en la carpeta `tools/` la primera vez que se ejecute.

### Terminal 2: Iniciar el frontend (React)

```bash
cd app-front-react
npm start
```

La aplicación estará disponible en `http://localhost:3000`

## 📖 Cómo usar

1. Abre `http://localhost:3000` en tu navegador
2. Pega la URL del video de TikTok o Instagram en el input
3. **Espera ~1 segundo** → Se cargarán automáticamente los metadatos
4. Verás el título y miniatura del video (vista previa)
5. Haz clic en el botón **"Descargar Video"**
6. El navegador abrirá el diálogo para guardar el archivo en tu PC

## 📁 Estructura del proyecto

```
AquilesDV/
├── app-back-go/                 # Backend en Go
│   ├── main.go                  # Punto de entrada
│   ├── go.mod                   # Módulos de Go
│   ├── api/
│   │   ├── handlers/            # Manejadores HTTP
│   │   │   ├── metadata_handler.go
│   │   │   └── dowload_handler.go
│   │   └── models/              # Modelos de datos
│   │       └── dowload_request.go
│   ├── internal/
│   │   ├── service/             # Lógica de negocio
│   │   │   └── dowload_service.go
│   │   └── tools/               # Herramientas utilitarias
│   │       └── ytdlp.go         # Descarga automática de yt-dlp
│   ├── downloads/               # Carpeta de descargas temporales
│   └── tools/                   # yt-dlp ejecutable (generado)
│
├── app-front-react/             # Frontend en React
│   ├── src/
│   │   ├── App.js               # Componente principal
│   │   ├── App.css              # Estilos
│   │   ├── index.js             # Punto de entrada
│   │   ├── index.css            # Estilos globales
│   │   └── setupProxy.js        # Configuración de proxy
│   ├── public/
│   │   ├── index.html           # HTML principal
│   │   └── ...iconos
│   ├── package.json             # Dependencias
│   └── package-lock.json        # Lock file
│
└── README.md                    # Este archivo
```

## 🔌 API Endpoints

### GET `/metadata`
Obtiene los metadatos del video (título, miniatura, URL original)

**Parámetro:** `url` (URL encoded del video)

**Respuesta:**
```json
{
  "title": "Título del video",
  "thumbnail": "https://...",
  "original_url": "https://www.tiktok.com/..."
}
```

### POST `/download`
Descarga el video y lo devuelve como archivo binario

**Body:**
```json
{
  "url": "https://www.tiktok.com/..."
}
```

**Respuesta:** Stream de video MP4

## 🛠️ Configuración

### Puerto del backend
- Edita `app-back-go/main.go` para cambiar el puerto (por defecto: 8080)

### Puerto del frontend
- Edita `app-front-react/package.json` en scripts → "start" para cambiar el puerto

### Tiempo de debounce
- Edita `app-front-react/src/App.js` en el useEffect para cambiar el tiempo de espera antes de cargar (por defecto: 0.8s)

## 🐛 Troubleshooting

### "No se pudieron obtener los metadatos"
- Verifica que el backend esté corriendo en `http://localhost:8080`
- Verifica que la URL del video sea válida
- Comprueba la consola del navegador (F12) para más detalles

### El video no se descarga
- Asegúrate de que el backend esté ejecutándose
- Verifica que haya espacio disponible en tu disco duro
- Revisa la carpeta `app-back-go/downloads/` donde se guardan temporalmente

### yt-dlp no se descarga
- Verifica tu conexión a Internet
- La versión se obtiene automáticamente desde GitHub
- Si hay problemas, puedes descargar manualmente desde https://github.com/yt-dlp/yt-dlp

## 📄 Notas importantes

⚠️ **Respeto de derechos de autor**
- Este software es solo para uso personal
- Respeta los términos de servicio de TikTok e Instagram
- No descargues contenido protegido sin permiso

📝 **Archivos temporales**
- Los videos se descargan temporalmente en `app-back-go/downloads/`
- Se eliminan automáticamente después de servir al usuario
- El directorio se crea automáticamente

## 🤝 Contribuciones

Las contribuciones son bienvenidas. Para cambios importantes:
1. Haz fork del repositorio
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Haz commit de tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Haz push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📝 Licencia

Este proyecto está bajo licencia MIT. Ver el archivo LICENSE para más detalles.

## 🙋 Soporte

Si encuentras problemas o tienes preguntas, por favor abre un issue en el repositorio.

---

**Última actualización:** 12 de febrero de 2026

**Versión:** 1.0.0
