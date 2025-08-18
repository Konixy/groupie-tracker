<script lang="ts">
	import type { Artist } from '../types';
	import { Map, TileLayer, Control } from 'leaflet';
	import { onMount } from 'svelte';

	let { artist }: { artist: Artist } = $props();

	let mapRef = $state<HTMLElement | string>('map');
	let map: Map = $derived(new Map(mapRef).setView([49.43814906784801, 1.1023802854095484], 13));

	function createMarker(name: string) {
		const el = document.createElement('div');
		el.innerHTML = `
			<div class="popup-content">
				<h3>${name}</h3>
			</div>
		`;
		return el;
	}

	onMount(async () => {
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
	});

	// $effect(() => {
	// 	if (artist) {
	// 		fetchLocations();
	// 	}
	// });

	async function fetchLocations() {
		const response = await fetch(`http://localhost:8080/locations/${artist.id}`);
		const data = await response.json();
		processLocations(data.datesLocations);
	}

	async function processLocations(datesLocations: Record<string, string[]>) {
		for (const location in datesLocations) {
			const dates = datesLocations[location];
			console.log(location, dates);
		}
		try {
			const response = await fetch('http://localhost:8080/locations');
			const data = await response.json();
			return data;
		} catch (error) {
			console.error('Erreur lors de la récupération des lieux:', error);
			return null;
		}
	}
</script>

<div bind:this={mapRef} class="map" aria-label="Carte des concerts"></div>
<link rel="stylesheet" href="https://unpkg.com/leaflet@2.0.0-alpha/dist/leaflet.css" />

<style>
	.map {
		width: 100%;
		height: 100%;
		border-radius: 1.5rem;
		overflow: hidden;
		background-color: var(--dark-muted);
	}

	:global(.popup-content) {
		text-align: center;
		padding: 0.5rem;
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
