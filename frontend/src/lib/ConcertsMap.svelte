<script lang="ts">
	import L from 'leaflet';
  import { onMount } from 'svelte';

  const currentArtist = "victoria-australia"

  let mapRef = $state<HTMLElement | string>("map")
  let map: L.Map = $derived(L.map(mapRef).setView([49.43814906784801, 1.1023802854095484], 13));

// TODO: add state management for the artists api and share the currently selected artist in it.


    onMount(() => {
      // On peut mettre les logos des artistes directement sur la map comme icone de marker
        L.marker([49.43814906784801, 1.1023802854095484]).addTo(map).setIcon(L.icon({
          iconUrl: '/z01.png',
          iconSize: [38, 38],
          iconAnchor: [19, 19],
          popupAnchor: [0, -19],
          className: 'z01-icon',
        }))
            .bindPopup('📍 Zone01')

            L.tileLayer(
	    'https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png',
	    {
	      attribution: `&copy;<a href="https://www.openstreetmap.org/copyright" target="_blank">OpenStreetMap</a>,
	        &copy;<a href="https://carto.com/attributions" target="_blank">CARTO</a>`,
	      subdomains: 'abcd',
	      maxZoom: 20,
	    }
	  ).addTo(map);
    });
</script>

<div class="title">
  <h1>Concerts</h1>
</div>

<div bind:this={mapRef} class="map"></div>

<style>
  .map {
    width: 90vw;
    height: 90vh;
    margin: 1rem;
    border-radius: 1.5rem;
    overflow: hidden;
    margin-top: 2rem;
  }

  .title {
    width: 90vw;
    text-align: center;
    margin-top: 2rem;
    border-radius: 1rem;
    font-size: 1.5rem;
  }

  :global(.z01-icon) {
    border-radius: 100%;
  }

  :global(.leaflet-control-zoom) {
    border: none !important;
    display: flex !important;
    flex-direction: column !important;
    gap: 0.3rem !important;
  }

  :global(.leaflet-control-zoom-in), :global(.leaflet-control-zoom-out) {
    border-radius: 100% !important;
    box-shadow: 0 0 0 0.1rem rgba(0, 0, 0, 0.1) !important;
  }

</style>  
<link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" />
