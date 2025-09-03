<script lang="ts">
	import type { Artist, Location } from '@/types';
	import { Map, TileLayer, Control, Marker, DivIcon } from 'leaflet';
	import { onMount, untrack } from 'svelte';

	// Import des polices
	const link1 = document.createElement('link');
	link1.rel = 'stylesheet';
	link1.href = 'https://fonts.googleapis.com/css2?family=Russo+One&display=swap';
	document.head.appendChild(link1);

	const link2 = document.createElement('link');
	link2.rel = 'stylesheet';
	link2.href = 'https://fonts.googleapis.com/css2?family=Jost:wght@400;500;600&display=swap';
	document.head.appendChild(link2);

	const link3 = document.createElement('link');
	link3.rel = 'stylesheet';
	link3.href =
		'https://fonts.googleapis.com/css2?family=Montserrat:wght@400;500;600;700&display=swap';
	document.head.appendChild(link3);

	// Import local de la police Hanson
	const hansonStyle = document.createElement('style');
	hansonStyle.textContent = `
		@font-face {
			font-family: 'Hanson';
			src: url('/font/Hanson-Bold.ttf') format('truetype');
			font-weight: bold;
			font-style: normal;
		}
	`;
	document.head.appendChild(hansonStyle);

	const addressTypes = {
		country: 'pays',
		state: 'état',
		county: 'comté',
		city: 'ville',
		town: 'village',
		municipality: 'municipalité'
	};

	let { artist }: { artist: Artist } = $props();

	let mapRef = $state<HTMLElement | string>('map');
	let map: Map = $derived(new Map(mapRef).setView([49.43814906784801, 1.1023802854095484], 13));

	let loading = $state(true);
	let abortController = $state<AbortController | null>(null);
	let currentArtistLocations = $state<(Location & { slug: string })[]>([]);

	function clearAllMarkers() {
		map.eachLayer((layer) => {
			if (layer instanceof Marker) {
				layer.remove();
			}
		});
	}

	onMount(async () => {
		const layer = new TileLayer('https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png', {
			subdomains: 'abcd',
			maxZoom: 20,
			minZoom: 0
		});

		map.addLayer(layer);
	});

	$effect(() => {
		if (artist) {
			untrack(() => {
				fetchLocations();
			});
		}
	});

	async function fetchLocations() {
		abortController?.abort('Artist changed');
		abortController = new AbortController();
		clearAllMarkers();

		loading = true;
		const response = await fetch(`http://localhost:8080/locations/${artist.id}`, {
			signal: abortController?.signal
		});
		const data: Record<string, Location> = await response.json();

		currentArtistLocations = Object.entries(data).map(([slug, location]) => ({
			...location,
			slug
		}));

		for (const location of currentArtistLocations) {
			new Marker([Number(location.lat), Number(location.lon)], {
				title: `(${location.dates.length} concerts)`,
				icon: new DivIcon({
					className: '',
					iconSize: [32, 32],
					iconAnchor: [16, 32],
					html: `<svg class="map-marker" width="32px" height="32px" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"
										fill="var(--dark-vibrant)">
										<path fill="inherit"
											d="M11.291 21.706 12 21l-.709.706zM12 21l.708.706a1 1 0 0 1-1.417 0l-.006-.007-.017-.017-.062-.063a47.708 47.708 0 0 1-1.04-1.106 49.562 49.562 0 0 1-2.456-2.908c-.892-1.15-1.804-2.45-2.497-3.734C4.535 12.612 4 11.248 4 10c0-4.539 3.592-8 8-8 4.408 0 8 3.461 8 8 0 1.248-.535 2.612-1.213 3.87-.693 1.286-1.604 2.585-2.497 3.735a49.583 49.583 0 0 1-3.496 4.014l-.062.063-.017.017-.006.006L12 21z" />
										<circle class="map-marker-circle" cx="12" cy="10" r="3" />
								</svg>`
				})
			}).addTo(map);
		}

		handleMoveToGeneralView();

		loading = false;
	}

	function handleMoveToGeneralView() {
		// Si aucune location, on ne fait rien
		if (currentArtistLocations.length === 0) return;

		// Si une seule location, on centre sur celle-ci
		if (currentArtistLocations.length === 1) {
			const location = currentArtistLocations[0];
			handleSelectLocation(location);
			return;
		}

		// Calcul des limites globales pour toutes les locations
		let minLat = Number.POSITIVE_INFINITY;
		let maxLat = Number.NEGATIVE_INFINITY;
		let minLon = Number.POSITIVE_INFINITY;
		let maxLon = Number.NEGATIVE_INFINITY;

		// Parcours de toutes les locations pour trouver les limites
		for (const location of currentArtistLocations) {
			const lat = Number(location.lat);
			const lon = Number(location.lon);

			// Mise à jour des limites avec les coordonnées du centre
			minLat = Math.min(minLat, lat);
			maxLat = Math.max(maxLat, lat);
			minLon = Math.min(minLon, lon);
			maxLon = Math.max(maxLon, lon);

			// Si la location a un boundingbox, on l'utilise pour des limites plus précises
			if (location.boundingbox && location.boundingbox.length >= 4) {
				const bbox = location.boundingbox;
				minLat = Math.min(minLat, Number(bbox[0]));
				maxLat = Math.max(maxLat, Number(bbox[1]));
				minLon = Math.min(minLon, Number(bbox[2]));
				maxLon = Math.max(maxLon, Number(bbox[3]));
			}
		}

		// Ajout d'une marge pour que les marqueurs ne soient pas collés aux bords
		const latMargin = (maxLat - minLat) * 0.1;
		const lonMargin = (maxLon - minLon) * 0.1;

		// Application de la vue globale avec les limites calculées
		map.fitBounds([
			[minLat - latMargin, minLon - lonMargin],
			[maxLat + latMargin, maxLon + lonMargin]
		]);
	}

	function handleSelectLocation(location: Location) {
		map.setView([Number(location.lat), Number(location.lon)]);
		map.fitBounds([
			[Number(location.boundingbox[0]), Number(location.boundingbox[2])],
			[Number(location.boundingbox[1]), Number(location.boundingbox[3])]
		]);
	}
</script>

<link rel="stylesheet" href="https://unpkg.com/leaflet@2.0.0-alpha/dist/leaflet.css" />

<div class="map-container">
	<div bind:this={mapRef} class="map" aria-label="Carte des concerts"></div>
	<button class="general-view-button" onclick={() => handleMoveToGeneralView()}>
		Vue générale
	</button>
	<div class="side-panel">
		{#if loading}
			<div class="loading">
				<h4>Chargement des concerts...</h4>
				<p>Cela peut prendre jusqu'à 1 minute</p>
			</div>
		{:else}
			<div class="concerts-list">
				<div class="concerts-list-header">
					<h2>Liste des concerts</h2>
					<div class="header-separator"></div>
				</div>
				<div class="concerts-list-content">
					{#each currentArtistLocations as location, locationIndex (location.slug)}
						{#each location.dates as date, dateIndex (date)}
							<div class="concert-row">
								<button onclick={() => handleSelectLocation(location)} class="concert-item">
									<div class="concert-info">
										<span class="concert-location">{location.name}</span>
										<span class="concert-type"
											>({location.type in addressTypes
												? addressTypes[location.type as keyof typeof addressTypes]
												: location.type})</span
										>
										<span class="concert-dates">{date}</span>
									</div>
								</button>
								{#if !(locationIndex === currentArtistLocations.length - 1 && dateIndex === location.dates.length - 1)}
									<div class="concert-separator"></div>
								{/if}
							</div>
						{/each}
					{/each}
				</div>
			</div>
		{/if}
	</div>
</div>

<style>
	.map-container {
		position: relative;
		width: 100%;
		height: calc(100vh - 8rem);
	}

	.side-panel {
		position: absolute;
		top: 2rem;
		left: 2rem;
		width: 350px;
		max-height: calc(100% - 4rem);
		z-index: 1000;
		display: flex;
		flex-direction: column;
		padding: 1.5rem;
		background-color: var(--dark-vibrant);
		backdrop-filter: blur(10px);
		border-radius: 1.5rem;
		box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
	}

	.map {
		width: 100%;
		height: 100%;
		border-radius: 1.5rem;
		overflow: hidden;
		background-color: var(--dark-muted);
	}

	.general-view-button {
		position: absolute;
		top: 2rem;
		right: 2rem;
		background-color: var(--dark-vibrant);
		color: var(--light-vibrant);
		border: none;
		padding: 0.5rem 1rem;
		border-radius: 50px;
		cursor: pointer;
		z-index: 1000;
		font-family: 'Jost', sans-serif;
		font-weight: 500;
		transition: background-color 0.3s ease;
	}

	.general-view-button:hover {
		background-color: var(--dark-muted);
	}

	@media (max-width: 1000px) {
		.side-panel {
			position: relative;
			top: 0;
			left: 0;
			width: 100%;
			max-height: 35vh;
			margin-bottom: 1rem;
		}

		.map {
			height: calc(65vh - 2rem);
		}
	}

	.concerts-list-header {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.5rem;
	}

	.header-separator {
		width: 50%;
		height: 1px;
		background: color-mix(in srgb, var(--light-vibrant), transparent 60%);
		margin: 0.5rem;
	}

	.concerts-list {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		height: 100%;
		overflow-y: hidden;
	}

	.concerts-list-content {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		overflow-y: auto;
		max-height: 100%;
		scrollbar-width: thin;
		scrollbar-color: var(--light-vibrant) var(--dark-muted);
	}

	.concerts-list-content::-webkit-scrollbar {
		width: 8px;
	}

	.concerts-list-content::-webkit-scrollbar-track {
		background: var(--dark-muted);
		border-radius: 4px;
	}

	.concerts-list-content::-webkit-scrollbar-thumb {
		background: var(--light-vibrant);
		border-radius: 4px;
		transition: background 0.3s ease;
	}

	.concerts-list-content::-webkit-scrollbar-thumb:hover {
		background: var(--muted);
	}

	.concert-row {
		display: flex;
		flex-direction: column;
	}

	.concert-item {
		border: none;
		outline: none;
		color: inherit;
		transition: background-color 0.2s ease;
		padding: 0.8rem;
		border-radius: 0.5rem;
		background-color: transparent;
		width: 100%;
		text-align: left;
	}

	.concert-item:hover {
		background-color: var(--dark-muted);
	}

	.concert-info {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		flex-wrap: wrap;
	}

	.concert-location {
		font-family: 'Hanson', sans-serif;
		font-weight: bold;
		font-size: 0.9rem;
		color: var(--light-vibrant);
	}

	.concert-type {
		font-size: 0.75rem;
		color: color-mix(in srgb, var(--light-vibrant), transparent 30%);
		font-style: italic;
	}

	.concert-dates {
		font-family: 'Hanson', sans-serif;
		font-weight: bold;
		font-size: 0.7rem;
		color: color-mix(in srgb, var(--light-vibrant), transparent 20%);
		margin-left: auto;
	}

	.concert-separator {
		height: 1px;
		background: color-mix(in srgb, var(--light-vibrant), transparent 80%);
		margin: 0.5rem 0;
	}

	.loading {
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		justify-content: center;
		gap: 1rem;
		height: 100%;
	}

	.loading h4 {
		font-size: 1.3rem;
		margin: 0;
		font-family: 'Russo One', sans-serif;
	}

	.loading p {
		margin: 0;
		font-size: 0.9rem;
		color: color-mix(in srgb, var(--light-vibrant), transparent 30%);
		font-family: 'Jost', sans-serif;
		font-weight: 400;
	}

	.concerts-list h2 {
		font-family: 'Bukhari Script', cursive;
		font-weight: 400;
	}

	.concerts-list button {
		font-family: 'Jost', sans-serif;
		font-weight: 500;
	}

	.concert-item {
		font-family: 'Montserrat', sans-serif;
		font-weight: 600;
	}

	.concert-item-type {
		font-family: 'Jost', sans-serif;
		font-weight: 400;
	}

	.concert-item-date {
		font-family: 'Jost', sans-serif;
		font-weight: 400;
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

	:global(.leaflet-control-attribution) {
		display: none;
	}

	:global(.map-marker) {
		filter: drop-shadow(0 0 3px rgba(0, 0, 0, 0.2));
	}

	:global(.map-marker-circle) {
		fill: color-mix(in srgb, var(--dark-vibrant), white 90%);
	}
</style>
