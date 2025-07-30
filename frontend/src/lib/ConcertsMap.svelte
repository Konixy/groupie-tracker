<script lang="ts">
	import { Marker, Map, TileLayer, Icon, Control, Popup } from 'leaflet';
	import { onMount } from 'svelte';

	const currentArtist = 'victoria-australia';

	let mapRef = $state<HTMLElement | string>('map');
	let map: Map = $derived(new Map(mapRef).setView([49.43814906784801, 1.1023802854095484], 13));

	// TODO: add state management for the artists api and share the currently selected artist in it.

	function createMarker(name: string) {
		const el = document.createElement('div');
		el.innerHTML = `
			<div class="z01-popup-content">
				<h3>${name}</h3>
			</div>
		`;
		return el;
	}

	onMount(async () => {
		// On peut mettre les logos des artistes directement sur la map comme icone de marker
		// new Marker([49.43814906784801, 1.1023802854095484])
		// 	.addTo(map)
		// 	.setIcon(
		// 		new Icon({
		// 			iconUrl: '/z01.png',
		// 			iconSize: [38, 38],
		// 			iconAnchor: [19, 19],
		// 			popupAnchor: [0, -19],
		// 			className: 'z01-icon'
		// 		})
		// 	)
		// 	.bindPopup(new Popup({ className: 'z01-popup', content: createMarker('Zone 01') }));

		const darkLayer = new TileLayer(
			'https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png',
			{
				subdomains: 'abcd',
				maxZoom: 20,
				minZoom: 0
			}
		);

		const lightLayer = new TileLayer(
			'https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png',
			{
				subdomains: 'abcd',
				maxZoom: 20,
				minZoom: 0
			}
		);

		new Control.Layers({
			light: lightLayer,
			dark: darkLayer
		}).addTo(map);

		map.addLayer(darkLayer);

		const locations = await fetchLocations();
	});

	async function fetchLocations() {
		const response = await fetch('http://localhost:8080/locations');
		const data = await response.json();
		return data;
	}
</script>

<div bind:this={mapRef} class="map"></div>
<link rel="stylesheet" href="https://unpkg.com/leaflet@2.0.0-alpha/dist/leaflet.css" />

<style>
	.map {
		width: 100%;
		height: 100%;
		border-radius: 1.5rem;
		overflow: hidden;
		background-color: var(--dark-muted);
	}

	:global(.z01-icon) {
		border-radius: 100%;
	}

	:global(.leaflet-top),
	:global(.leaflet-bottom) {
		z-index: 800 !important;
	}

	:global(.leaflet-control-zoom) {
		display: none;
	}

	:global(.leaflet-control-layers) {
		border-radius: 14px !important;
		background-color: var(--dark-muted) !important;
		color: var(--light-vibrant) !important;
		box-shadow: none !important;
		border: none !important;
		transition: all 0.1s ease !important;
	}

	:global(.leaflet-control-layers-list) {
		border: none !important;
		padding: 0.2rem !important;
	}

	:global(.leaflet-control-layers-base) {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.2rem;
	}

	:global(.leaflet-control-layers-base span) {
		display: flex;
		align-items: center;
		gap: 0.4rem;
	}

	:global(.leaflet-control-layers-selector) {
		-webkit-appearance: none;
		appearance: none;
		background-color: var(--dark-muted);
		margin: 0;
		font: inherit;
		color: currentColor;
		width: 1.15em;
		height: 1.15em;
		border: 0.15em solid var(--light-vibrant);
		border-radius: 50%;
		transform: translateY(-0.075em);
		display: grid;
		place-content: center;
	}

	:global(.leaflet-control-layers-selector::before) {
		content: '';
		width: 0.65em;
		height: 0.65em;
		border-radius: 50%;
		transform: scale(0);
		transition: 120ms transform ease-in-out;
		box-shadow: inset 1em 1em var(--light-vibrant);
	}

	:global(.leaflet-control-layers-selector:checked::before) {
		transform: scale(2);
	}

	:global(.leaflet-control-attribution) {
		display: none;
	}
</style>
