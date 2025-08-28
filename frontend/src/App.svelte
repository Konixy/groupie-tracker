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

	let artists = $state<Artist[]>([]);
	let selectedArtist = $state<Artist | null>(null);
	let loading = $state(true);
	let currentIndex = $state(Math.floor(Math.random() * 52));

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
	}
</script>

<main>
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
			<SearchBar {artists} bind:current={currentIndex} />
			<Carousel bind:selectedArtist {artists} bind:loading bind:current={currentIndex} />
			{#if selectedArtist}
				<ArtistDetail bind:artist={selectedArtist} onClose={closeArtistDetail} />
			{/if}
		</div>
		<div class="side-navigation">
			<button class="side-nav-button left" onclick={() => scrollToSection('map-section')}>
				Concerts
			</button>
		</div>
	</section>

	<!-- Section 3: Map -->
	<section id="map-section" class="content-section">
		<div class="title">
			<h1>Concerts</h1>
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

	.title {
		width: 100%;
		text-align: center;
		margin-bottom: 1rem;
		border-radius: 1rem;
		font-size: 1.5rem;
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

	.side-navigation {
		position: absolute;
		bottom: 3rem;
		left: 0;
		right: 0;
		display: flex;
		justify-content: space-between;
		padding: 0 3rem;
		pointer-events: none;
	}

	.side-nav-button {
		background: var(--dark-muted);
		color: var(--light-vibrant);
		border: none;
		padding: 0.8rem 1.5rem;
		border-radius: 8px;
		cursor: pointer;
		font-size: 1rem;
		transition: all 0.3s ease;
		pointer-events: auto;
		opacity: 0.8;
	}

	.side-nav-button:hover {
		background: var(--muted);
		opacity: 1;
		transform: translateY(-2px);
	}

	.side-nav-button.left {
		margin-right: auto;
	}
</style>
