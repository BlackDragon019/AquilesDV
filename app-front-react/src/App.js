import React from "react";
import './App.css';

function App() {
    const [videoUrl, setVideoUrl] = React.useState("");
    const [metadata, setMetadata] = React.useState(null);
    const [loading, setLoading] = React.useState(false);
    const [downloading, setDownloading] = React.useState(false);
    const [error, setError] = React.useState(null);
    const [videoBlob, setVideoBlob] = React.useState(null);
    const [debounceTimer, setDebounceTimer] = React.useState(null);

    const backendUrl = "/api";

    // Cargar metadatos automáticamente cuando el URL cambia (con debounce)
    React.useEffect(() => {
        if (debounceTimer) {
            clearTimeout(debounceTimer);
        }

        if (!videoUrl.trim()) {
            setMetadata(null);
            setVideoBlob(null);
            setError(null);
            return;
        }

        const timer = setTimeout(() => {
            handleFetchMetadata(videoUrl);
        }, 800); // Esperar 0.8 segundos después de que el usuario deje de escribir

        setDebounceTimer(timer);

        return () => {
            if (timer) clearTimeout(timer);
        };
    }, [videoUrl]);

    const handleFetchMetadata = async (url) => {
        setLoading(true);
        setError(null);
        setMetadata(null);
        setVideoBlob(null);

        try {
            const response = await fetch(`${backendUrl}/metadata?url=${encodeURIComponent(url)}`);

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`Error ${response.status}: ${errorText}`);
            }

            const data = await response.json();
            setMetadata(data);
            // No cargar el video aquí, solo cuando el usuario haga clic en descargar
        } catch (err) {
            console.error('Error al obtener metadatos:', err);
            setError(`No se pudieron obtener los metadatos: ${err.message}`);
        } finally {
            setLoading(false);
        }
    };

    const loadVideo = async (url) => {
        try {
            const response = await fetch(`${backendUrl}/download`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ url: url }),
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`Error ${response.status}: ${errorText}`);
            }

            const blob = await response.blob();
            setVideoBlob(blob);
        } catch (err) {
            console.error('Error al cargar el video:', err);
            setError(`No se pudo cargar el video: ${err.message}`);
        }
    };

    const handleDownloadVideo = async () => {
        setDownloading(true);
        setError(null);

        if (!videoUrl) {
            setError('Por favor, ingresa una URL de video para descargar.');
            setDownloading(false);
            return;
        }

        try {
            let blob = videoBlob;

            // Si no tenemos el blob cargado, descargarlo ahora
            if (!blob) {
                const response = await fetch(`${backendUrl}/download`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ url: videoUrl }),
                });

                if (!response.ok) {
                    const errorText = await response.text();
                    throw new Error(`Error ${response.status}: ${errorText}`);
                }

                blob = await response.blob();
            }

            // Obtener el nombre del archivo de los metadatos o usar un nombre por defecto
            const fileName = metadata?.Title ? `${metadata.Title.substring(0, 50)}.mp4` : 'video.mp4';

            // Crear un URL del blob
            const url = window.URL.createObjectURL(blob);

            // Crear un elemento <a> y simular un click para descargar
            const link = document.createElement('a');
            link.href = url;
            link.download = fileName;
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);

            // Liberar el URL del objeto
            window.URL.revokeObjectURL(url);

            console.log(`Archivo descargado: ${fileName}`);
        } catch (err) {
            console.error('Error al iniciar la descarga:', err);
            setError(`No se pudo iniciar la descarga: ${err.message}`);
        } finally {
            setDownloading(false);
        }
    };


    return (
        <div className="App">
            <header className="App-header">
                <h1>Descargador de Videos 👾</h1>
                <p>TikTok / Instagram</p>
            </header>

            <div className="container">
                <input
                    type="text"
                    placeholder="Pega la URL del video de TikTok o Instagram aquí"
                    value={videoUrl}
                    onChange={(e) => setVideoUrl(e.target.value)}
                    style={{ marginBottom: '20px' }}
                />

                {error && <p className="error-message">{error}</p>}

                {loading && !metadata && (
                    <div style={{ textAlign: 'center', padding: '20px', color: '#666' }}>
                        <p>Cargando vista previa...</p>
                    </div>
                )}

                {metadata && (
                    <div className="video-preview">
                        <h2>{metadata.Title}</h2>
                        
                        {/* Video Player - Solo si está cargado */}
                        {videoBlob && (
                            <div style={{ marginBottom: '20px', textAlign: 'center' }}>
                                <video 
                                    width="100%" 
                                    height="auto" 
                                    controls 
                                    style={{ maxWidth: '600px', backgroundColor: '#000' }}
                                >
                                    <source src={window.URL.createObjectURL(videoBlob)} type="video/mp4" />
                                    Tu navegador no soporta el video HTML5. Por favor, descarga el video.
                                </video>
                            </div>
                        )}

                        {/* Thumbnail - Previsualización rápida */}
                        {metadata.Thumbnail && !videoBlob && (
                            <div style={{ marginBottom: '20px', textAlign: 'center' }}>
                                <img
                                    src={metadata.Thumbnail}
                                    alt="Miniatura del video"
                                    className="thumbnail"
                                    style={{ maxWidth: '100%', maxHeight: '400px', borderRadius: '8px' }}
                                />
                            </div>
                        )}
                        
                        <p className="original-url" style={{ textAlign: 'center', marginBottom: '20px' }}>
                            <a href={metadata.OriginalURL} target="_blank" rel="noopener noreferrer">
                                Ver en original
                            </a>
                        </p>
                    </div>
                )}

                {/* Botón de descarga siempre visible si hay URL o metadatos */}
                {(videoUrl.trim() || metadata) && (
                    <div style={{ textAlign: 'center', marginTop: '20px' }}>
                        <button onClick={handleDownloadVideo} disabled={downloading}>
                            {downloading ? 'Descargando...' : 'Descargar Video'}
                        </button>
                    </div>
                )}
            </div>
        </div>
    );
}

export default App;
