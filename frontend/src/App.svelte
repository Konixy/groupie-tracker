<script lang="ts">
	import { onMount } from 'svelte';
	import Logo from './lib/Logo.svelte';
	import Carousel from './lib/Carousel.svelte';
	import SearchBar from './lib/SearchBar.svelte';
	import ArtistDetail from './lib/ArtistDetail.svelte';
	import ConcertsMap from './lib/ConcertsMap.svelte';
	import type { Artist } from './types';
	import Footer from './lib/Footer.svelte';

	let artists = $state<Artist[]>([]);
	let selectedArtist = $state<Artist | null>(null);

	onMount(async () => {
		const res = await fetch('http://localhost:8080/artists');
		const data = await res.json();
		// Supporte à la fois l'API locale (data.artists) et distante (data)
		artists = Array.isArray(data) ? data : data.artists;
	});
</script>

<main>
	<Logo />
	<SearchBar {artists} bind:selectedArtist />
	<Carousel bind:selectedArtist {artists} />
	{#if selectedArtist}
		<ArtistDetail bind:artist={selectedArtist} onClose={() => (selectedArtist = null)} />
	{/if}
	<ConcertsMap />
	<Footer />
</main>

<style>
	main {
		min-height: 100vh;
		box-sizing: border-box;
		padding: 2rem 0 4rem 0;
		display: flex;
		flex-direction: column;
		align-items: center;
	}
</style>
