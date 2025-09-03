<script lang="ts">
	import { onMount } from 'svelte';
	import Logo from './lib/components/logo/Logo.svelte';
	import Carousel from './lib/components/Carousel.svelte';
	import SearchBar from './lib/components/SearchBar.svelte';
	import ArtistDetail from './lib/components/ArtistDetail.svelte';
	import Map from './lib/components/map/Map.svelte';

	import type { Artist } from './types';
	import Footer from './lib/components/Footer.svelte';
	import { fade } from 'svelte/transition';

	// Import des polices
	const link1 = document.createElement('link');
	link1.rel = 'stylesheet';
	link1.href = 'https://fonts.googleapis.com/css2?family=Bungee+Outline&display=swap';
	document.head.appendChild(link1);

	const link2 = document.createElement('link');
	link2.rel = 'stylesheet';
	link2.href =
		'https://fonts.googleapis.com/css2?family=Montserrat:wght@400;500;600;700&display=swap';
	document.head.appendChild(link2);

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

	let artists = $state<Artist[]>([]);
	let selectedArtist = $state<Artist | null>(null);
	let loading = $state(true);
	let currentIndex = $state(Math.floor(Math.random() * 52));
	let highlightedMember = $state<string | null>(null);
	let highlightedDate = $state<string | null>(null);
	let highlightedAlbum = $state<string | null>(null);

	onMount(async () => {
		await loadArtists();
	});

	async function loadArtists() {
		const res = await fetch('http://localhost:8080/artists');
		const data = await res.json();
		artists = Array.isArray(data) ? data : data.artists;
	}

	function scrollToSection(sectionId: string) {
		const section = document.getElementById(sectionId);
		if (section) {
			// Calculer la position exacte pour que la section soit entièrement visible
			const sectionTop = section.offsetTop;
			const windowHeight = window.innerHeight;
			const sectionHeight = section.offsetHeight;

			// Si la section est plus haute que la fenêtre, on scroll au début
			// Sinon, on centre la section dans la fenêtre
			const scrollPosition =
				sectionHeight > windowHeight ? sectionTop : sectionTop - (windowHeight - sectionHeight) / 2;

			window.scrollTo({
				top: Math.max(0, scrollPosition),
				behavior: 'smooth'
			});
		}
	}

	// Gestion de la sélection d'artiste
	onMount(() => {
		scrollToSection('hero-section');
	});

	function closeArtistDetail() {
		selectedArtist = null;
		highlightedMember = null;
		highlightedDate = null;
		highlightedAlbum = null;
	}

	function handleMemberSearch(memberName: string, artist: Artist) {
		highlightedMember = memberName;
		// Arrêter le clignotement après 3 secondes
		setTimeout(() => {
			highlightedMember = null;
		}, 3000);
	}

	function handleDateSearch(date: string, artist: Artist) {
		highlightedDate = date;
		// Arrêter le clignotement après 3 secondes
		setTimeout(() => {
			highlightedDate = null;
		}, 3000);
	}

	function handleAlbumSearch(album: string, artist: Artist) {
		highlightedAlbum = album;
		// Arrêter le clignotement après 3 secondes
		setTimeout(() => {
			highlightedAlbum = null;
		}, 3000);
	}
</script>

<main>
	<!-- Élément décoratif SVG -->
	<div class="decorative-svg">
		<img src="/fondaccueil.svg" alt="Décoration" />
	</div>

	<!-- Section 1: Titre Groupie Tracker -->
	<section id="hero-section" class="hero-section">
		<div class="hero-content">
			<Logo {loading} />
			{#if !loading}
				<button
					class="scroll-button"
					onclick={() => scrollToSection('carousel-section')}
					aria-label="Défiler vers le carousel"
					in:fade={{ duration: 1000, delay: 1000 }}
				>
					<svg
						width="24"
						height="24"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="3"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<path d="m7 6 5 5 5-5" /><path d="m7 13 5 5 5-5" />
					</svg>
				</button>
			{/if}
		</div>
	</section>

	<!-- Section 2: Carousel, recherche et boutons -->
	<section id="carousel-section" class="content-section carousel-section">
		<div class="carousel-content">
			<SearchBar
				{artists}
				bind:current={currentIndex}
				onArtistSelect={(artist) => (selectedArtist = artist)}
				onMemberSearch={handleMemberSearch}
				onDateSearch={handleDateSearch}
				onAlbumSearch={handleAlbumSearch}
			/>
			<Carousel
				bind:selectedArtist
				{artists}
				bind:loading
				bind:current={currentIndex}
				{highlightedMember}
				{highlightedDate}
				{highlightedAlbum}
			/>
			{#if selectedArtist}
				<ArtistDetail bind:artist={selectedArtist} onClose={closeArtistDetail} />
			{/if}
		</div>
		<div class="center-navigation">
			<button class="center-nav-button" onclick={() => scrollToSection('map-section')}>
				Concerts
				<svg
					width="24"
					height="24"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2.5"
					stroke-linecap="round"
					stroke-linejoin="round"
					style="display: block; margin: 0 auto;"
				>
					<path d="M6 9l6 6 6-6" />
				</svg>
			</button>
		</div>
	</section>

	<!-- Section 3: Map -->
	<section id="map-section" class="map-section">
		<div class="title">
			<h1 class="concerts-title">Concerts</h1>
			<h2 class="artist-subtitle">DE {artists?.[currentIndex]?.name || '...'}</h2>
		</div>
		<Map artist={artists?.[currentIndex]} />
	</section>

	<Footer />
</main>

<style>
	main {
		display: flex;
		flex-direction: column;
		scroll-snap-type: y mandatory;
		position: relative;
	}

	.decorative-svg {
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		pointer-events: none;
		z-index: -1;
		overflow: hidden;
	}

	.decorative-svg img {
		width: 190%;
		height: auto;
		position: absolute;
		top: 40%;
		left: 55%;
		transform: translate(-50%, -50%);
		opacity: 0.15;
	}

	.hero-section {
		scroll-snap-align: start;
		height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		box-sizing: border-box;
	}

	.hero-content {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 4rem;
	}

	.scroll-button {
		background: none;
		border: 3px solid currentColor;
		border-radius: 50%;
		width: 50px;
		height: 50px;
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		transition:
			all 0.3s ease,
			color 1s ease;
		color: inherit;
		opacity: 0.6;
	}

	.scroll-button:hover {
		transform: translateY(5px);
		background: rgba(255, 255, 255, 0.1);
		opacity: 1;
	}

	.scroll-button:not(:hover) svg {
		animation: bounce 2s infinite;
	}

	@keyframes bounce {
		0%,
		20%,
		50%,
		80%,
		100% {
			transform: translateY(0);
		}
		40% {
			transform: translateY(-10px);
		}
		60% {
			transform: translateY(-5px);
		}
	}

	.content-section {
		scroll-snap-align: start;
		min-height: 100vh;
		box-sizing: border-box;
		padding: 2rem 0;
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	.map-section {
		scroll-snap-align: start;
		height: 100vh;
		box-sizing: border-box;
		padding: 0;
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	.title {
		width: 100%;
		text-align: center;
		margin-top: 3rem;
		margin-bottom: 3rem;
		border-radius: 1rem;
		font-size: 1.5rem;
	}

	.concerts-title {
		font-family: 'Bukhari Script', cursive;
		font-weight: normal;
		margin: 0;
		font-size: 3.5rem;
	}

	.artist-subtitle {
		font-family: 'Hanson', sans-serif;
		font-weight: bold;
		margin: 0;
		font-size: 1.2rem;
		color: color-mix(in srgb, var(--light-vibrant), transparent 20%);
		text-transform: uppercase;
		letter-spacing: 0.1em;
	}

	.carousel-section {
		position: relative;
		justify-content: flex-start;
		padding-top: 1rem;
		height: 100vh;
		overflow-y: auto;
	}

	.carousel-content {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
		width: 100%;
		max-width: 800px;
		margin: 0 auto;
		margin-top: -1rem;
	}

	.center-navigation {
		position: absolute;
		bottom: 0;
		left: 50%;
		transform: translateX(-50%);
		display: flex;
		justify-content: space-between;
		padding: 0 3rem;
		pointer-events: none;
	}

	.center-nav-button {
		background: none;
		color: var(--light-vibrant);
		border: none;
		cursor: pointer;
		font-size: 1rem;
		transition: all 0.3s ease;
		pointer-events: auto;
		opacity: 0.8;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 0.2rem;
	}

	.center-nav-button:hover {
		opacity: 1;
		transform: translateY(-2px);
	}
</style>
