<script lang="ts">
  import { onMount } from 'svelte';

  let { artists = [] } = $props();
  let mapContainer = $state<HTMLElement | null>(null);
  let map = $state<any>(null);
  let concerts = $state<Array<{name: string, date: string, location: string, coordinates: [number, number]}>>([]);
  let loading = $state(true);

  onMount(async () => {
    await loadMap();
    await loadConcertsData();
  });

  async function loadMap() {
    // Charger Leaflet CSS
    const link = document.createElement('link');
    link.rel = 'stylesheet';
    link.href = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.css';
    document.head.appendChild(link);

    // Charger Leaflet JS
    const script = document.createElement('script');
    script.src = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.js';
    script.onload = () => {
      initMap();
    };
    document.head.appendChild(script);
  }

  function initMap() {
    if (!mapContainer || typeof (window as any).L === 'undefined') return;

    const L = (window as any).L;
    // Initialiser la carte centrée sur l'Europe
    map = L.map(mapContainer).setView([48.8566, 2.3522], 4);

    // Ajouter la couche de tuiles OpenStreetMap
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '© OpenStreetMap contributors'
    }).addTo(map);
  }

  async function loadConcertsData() {
    try {
      loading = true;
      const allConcerts: Array<{name: string, date: string, location: string, coordinates: [number, number]}> = [];

      // Charger les concerts pour chaque artiste
      for (const artist of artists) { // Charger tous les artistes
        try {
          const response = await fetch(`http://localhost:8080/artists/${artist.id}`);
          if (response.ok) {
            const data = await response.json();
            const artistConcerts = data.concerts || {};
            
            // Traiter chaque concert de l'artiste
            for (const [date, locations] of Object.entries(artistConcerts)) {
              for (const location of locations as string[]) {
                const coordinates = await geocodeLocation(location);
                if (coordinates) {
                  allConcerts.push({
                    name: artist.name,
                    date: date,
                    location: location,
                    coordinates: coordinates
                  });
                }
              }
            }
          }
        } catch (err) {
          console.error(`Erreur lors du chargement des concerts pour ${artist.name}:`, err);
        }
      }

      concerts = allConcerts;
      addMarkersToMap();
    } catch (err) {
      console.error('Erreur lors du chargement des données de concerts:', err);
    } finally {
      loading = false;
    }
  }

  async function geocodeLocation(location: string): Promise<[number, number] | null> {
    try {
      // Utiliser l'API de géocodage gratuite Nominatim
      const response = await fetch(
        `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(location)}&limit=1`
      );
      const data = await response.json();
      
      if (data && data.length > 0) {
        return [parseFloat(data[0].lat), parseFloat(data[0].lon)];
      }
    } catch (err) {
      console.error(`Erreur de géocodage pour ${location}:`, err);
    }
    return null;
  }

  function addMarkersToMap() {
    if (!map || !concerts.length) return;

    const L = (window as any).L;
    if (!L) return;

    // Grouper les concerts par lieu pour éviter les marqueurs multiples
    const locationGroups = new Map<string, Array<{name: string, date: string, location: string, coordinates: [number, number]}>>();
    
    concerts.forEach(concert => {
      const key = `${concert.coordinates[0]},${concert.coordinates[1]}`;
      if (!locationGroups.has(key)) {
        locationGroups.set(key, []);
      }
      locationGroups.get(key)!.push(concert);
    });

    // Ajouter les marqueurs groupés
    locationGroups.forEach((concertsAtLocation, coordsKey) => {
      const [lat, lng] = coordsKey.split(',').map(Number);
      const location = concertsAtLocation[0].location;
      
      // Créer le contenu du popup
      const popupContent = `
        <div style="max-width: 300px;">
          <h3 style="margin: 0 0 10px 0; color: #333;">${location}</h3>
          <div style="font-size: 14px; color: #666;">
            ${concertsAtLocation.map(concert => 
              `<div style="margin: 5px 0; padding: 5px; background: #f5f5f5; border-radius: 4px;">
                <strong>${concert.name}</strong><br>
                ${formatDate(concert.date)}
              </div>`
            ).join('')}
          </div>
        </div>
      `;

      // Ajouter le marqueur
      const marker = L.marker([lat, lng])
        .addTo(map)
        .bindPopup(popupContent);

      // Style spécial pour les marqueurs avec plusieurs concerts
      if (concertsAtLocation.length > 1) {
        marker.setIcon(L.divIcon({
          className: 'custom-div-icon',
          html: `<div style="background-color: #e74c3c; width: 20px; height: 20px; display: block; border: 2px solid #fff; border-radius: 50%; text-align: center; line-height: 16px; color: white; font-size: 12px; font-weight: bold;">${concertsAtLocation.length}</div>`,
          iconSize: [20, 20],
          iconAnchor: [10, 10]
        }));
      }
    });
  }

  function formatDate(dateStr: string): string {
    try {
      const date = new Date(dateStr);
      return date.toLocaleDateString('fr-FR', {
        weekday: 'long',
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      });
    } catch {
      return dateStr;
    }
  }
</script>

<style>
  .map-section {
    margin: 3rem 0;
    padding: 2rem;
    background: white;
    border-radius: 20px;
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
    max-width: 1200px;
    width: 100%;
  }

  .map-header {
    text-align: center;
    margin-bottom: 2rem;
  }

  .map-header h2 {
    color: #333;
    font-size: 2rem;
    margin: 0 0 1rem 0;
  }

  .map-header p {
    color: #666;
    font-size: 1.1rem;
    margin: 0;
  }

  .map-container {
    height: 500px;
    border-radius: 12px;
    overflow: hidden;
    position: relative;
  }

  .loading-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(255, 255, 255, 0.9);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .loading-spinner {
    width: 40px;
    height: 40px;
    border: 4px solid #f3f3f3;
    border-top: 4px solid #667eea;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }

  .map-stats {
    display: flex;
    justify-content: center;
    gap: 2rem;
    margin-top: 1rem;
    color: #666;
    font-size: 0.9rem;
  }

  .stat-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .stat-number {
    font-weight: bold;
    color: #667eea;
  }
</style>

<div class="map-section">
  <div class="map-header">
    <h2>🌍 Carte des Concerts</h2>
    <p>Découvrez où se sont produits vos artistes préférés</p>
  </div>

  <div class="map-container" bind:this={mapContainer}>
    {#if loading}
      <div class="loading-overlay">
        <div class="loading-spinner"></div>
      </div>
    {/if}
  </div>

  <div class="map-stats">
    <div class="stat-item">
      <span>📍</span>
      <span><span class="stat-number">{concerts.length}</span> concerts affichés</span>
    </div>
    <div class="stat-item">
      <span>🎵</span>
      <span><span class="stat-number">{new Set(concerts.map(c => c.name)).size}</span> artistes</span>
    </div>
    <div class="stat-item">
      <span>🌍</span>
      <span><span class="stat-number">{new Set(concerts.map(c => c.location)).size}</span> lieux</span>
    </div>
  </div>
</div> 