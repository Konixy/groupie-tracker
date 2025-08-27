<script lang="ts">
	import type { Artist, Location } from '@/types';
	import { Map, TileLayer, Control, LatLngBounds } from 'leaflet';
	import { onMount, untrack } from 'svelte';

	const addressTypes = {
		country: 'pays',
		state: 'état',
		county: 'comté',
		city: 'ville',
		venue: 'lieu'
	};

	let { artist }: { artist: Artist } = $props();

	let mapRef = $state<HTMLElement | string>('map');
	let map: Map = $derived(new Map(mapRef).setView([49.43814906784801, 1.1023802854095484], 13));

	let loading = $state(true);
	let abortController = $state<AbortController | null>(null);
	let currentArtistLocations = $state<(Location & { slug: string })[]>([]);

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

		map.addLayer(lightLayer);
	});

	$effect(() => {
		if (artist) {
			untrack(() => {
				abortController?.abort('Artist changed');
				abortController = new AbortController();
				fetchLocations();
			});
		}
	});

	async function fetchLocations() {
		loading = true;
		const response = await fetch(`http://localhost:8080/locations/${artist.id}`, {
			signal: abortController?.signal
		});
		const data: Record<string, Location> = await response.json();

		currentArtistLocations = Object.entries(data).map(([slug, location]) => ({
			...location,
			slug
		}));
		handleSelectLocation(currentArtistLocations[0]);
		loading = false;
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
	<div class="map-left">
		<div bind:this={mapRef} class="map" aria-label="Carte des concerts"></div>
	</div>
	<div class="map-right">
		{#if loading}
			<div class="loading">
				<h4>Chargement des concerts...</h4>
				<p>Cela peut prendre jusqu'à 1 minute</p>
			</div>
		{:else}
			<div class="concerts-list">
				<h2>Liste des concerts</h2>
				<div class="concerts-list-content">
					{#each currentArtistLocations as location (location.slug)}
						<div>
							<button onclick={() => handleSelectLocation(location)} class="concert-item">
								{location.name}
								<span class="concert-item-type"
									>({location.type in addressTypes
										? addressTypes[location.type as keyof typeof addressTypes]
										: location.type})</span
								>
							</button>
							<div class="concert-item-dates">
								{#each location.dates as date}
									<span class="concert-item-date">{date}</span>
								{/each}
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}
	</div>
</div>

<style>
	.map-container {
		display: flex;
		width: 100%;
		gap: 2rem;
		padding: 0 2rem;
		min-height: 100%;
	}

	.map-left {
		flex: 4;
		display: flex;
		flex-direction: column;
	}

	.map-right {
		flex: 2;
		display: flex;
		flex-direction: column;
		justify-content: center;
		padding: 2rem;
		background-color: var(--dark-muted);
		border-radius: 1.5rem;
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
	}

	.concert-item {
		border: 2px solid color-mix(in srgb, var(--light-vibrant), transparent 70%);
		outline: none;
		color: inherit;
		transition: background-color 0.1s ease;
		padding: 0.5rem;
		border-radius: 0.5rem;
		background-color: var(--dark-muted);
	}

	.concert-item:hover {
		background-color: var(--dark-vibrant);
	}

	.concert-item-type {
		font-size: 0.8rem;
		color: color-mix(in srgb, var(--light-vibrant), transparent 70%);
	}

	.concert-item-dates {
		display: flex;
		flex-direction: column;
		gap: 0.2rem;
	}

	.concert-item-date {
		font-size: 0.8rem;
		color: color-mix(in srgb, var(--light-vibrant), transparent 70%);
	}

	@media (max-width: 1000px) {
		.map-container {
			flex-direction: column;
		}

		.map-left {
			flex: 1;
		}

		.map-right {
			flex: 1;
		}
	}

	.map {
		width: 100%;
		height: 100%;
		border-radius: 1.5rem;
		overflow: hidden;
		background-color: var(--dark-muted);
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
	}

	.loading p {
		margin: 0;
		font-size: 0.9rem;
		color: color-mix(in srgb, var(--light-vibrant), transparent 30%);
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
