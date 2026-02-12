import React from "react";
import './App.css';

function App() {
    const [videoUrl, setVideoUrl] = React.useState("");
    const [metadata, setMetadata] = React.useState(null);
    const [loading, setLoading] = React.useState(false);
    const [error, setError] = React.useState(null);

    const backendUrl = "http://localhost:8080";

    const handleFetchMetadata = async () => {
        setLoading(true);
        setError(null);
        setMetadata(null); 

        if (!videoUrl) {
            setError("Please enter a video URL.");
            setLoading(false);
            return;
        }

        try {
            const response = await fetch(`${backendUrl}/metadata?url=${encodeURIComponent(videoUrl)}`);

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`Error ${response.status}: ${errorText}`);
            }

            const data = await response.json();
            setMetadata(data);
        } catch (err) {
            console.error('Error al obtener metadatos:', err);
            setError(`No se pudieron obtener los metadatos: ${err.message}`);
        } finally {
            setLoading(false);
        }
    };

    const handleDownloadVideo = async () => {
        setLoading(true);
        setError(null);

        if (!videoUrl) {
            setError('Por favor, ingresa una URL de video para descargar.');
            setLoading(false);
            return;
        }
        try {
            const response = await fetch(`${backendBaseUrl}/download`, { // Usamos la URL directa del backend
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
            console.log('Solicitud de descarga enviada. El navegador debería iniciar la descarga.');
            alert('¡Descarga iniciada! Tu navegador debería mostrar un diálogo para guardar el archivo.');
        } catch (err) {
            console.error('Error al iniciar la descarga:', err);
            setError(`No se pudo iniciar la descarga: ${err.message}`);
        } finally {
            setLoading(false);
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
onKeyDown={(e) => { // Permite previsualizar con Enter
if (e.key === 'Enter') {
handleFetchMetadata();
}
}}
/>
<button onClick={handleFetchMetadata} disabled={loading}>
{loading && metadata === null ? 'Cargando...' : 'Previsualizar'}
</button>
{metadata && (
<button onClick={handleDownloadVideo} disabled={loading}>
{loading ? 'Preparando descarga...' : 'Descargar Video'}
</button>
)}
</div>

{error && <p className="error-message">{error}</p>}

{metadata && (
<div className="video-preview">
<h2>{metadata.Title}</h2>
{metadata.Thumbnail && (
<img
src={metadata.Thumbnail}
alt="Miniatura del video"
className="thumbnail"
/>
)}
<p className="original-url">URL original: <a href={metadata.OriginalURL} target="_blank" rel="noopener noreferrer">{metadata.OriginalURL}</a></p>
</div>
)}
</div>
);
}

export default App;

// Nota: Asegúrate de que el backend esté corriendo en http://localhost:8080 y que el proxy esté configurado correctamente para evitar problemas de CORS.
