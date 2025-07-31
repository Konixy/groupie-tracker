<script lang="ts">
	import { onMount } from 'svelte';
	import Logo from './lib/Logo.svelte';
	import Carousel from './lib/Carousel.svelte';
	import SearchBar from './lib/SearchBar.svelte';
	import ArtistDetail from './lib/ArtistDetail.svelte';
	import ConcertsMap from './lib/ConcertsMap.svelte';
	import type { Artist } from './types';
	import Footer from './lib/Footer.svelte';
	import { fade } from 'svelte/transition';
	import ConcertsView from './lib/ConcertsView.svelte';
	import StatsView from './lib/StatsView.svelte';

	let artists = $state<Artist[]>([]);
	let selectedArtist = $state<Artist | null>(null);
	let firstRender = $state(true);
	let currentPage = $state(window.location.pathname);

	onMount(async () => {
		const res = await fetch('http://localhost:8080/artists');
		const data = await res.json();
		// Supporte à la fois l'API locale (data.artists) et distante (data)
		artists = Array.isArray(data) ? data : data.artists;

		window.scrollTo(0, 0);
	});

	function scrollToCarousel() {
		const carouselSection = document.querySelector('.carousel-section');
		if (carouselSection) {
			carouselSection.scrollIntoView({ behavior: 'smooth' });
		}
	}

	function scrollToMap() {
		const mapSection = document.querySelector('.map-section');
		if (mapSection) {
			mapSection.scrollIntoView({ behavior: 'smooth' });
		}
	}

	function navigate(path: string) {
		window.history.pushState(null, '', path);
		currentPage = path;
		window.scrollTo(0, 0);
	}
</script>

<main>
	{#if currentPage === '/' || currentPage === ''}
		<!-- D'abord le titre qu'on va animer après -->
		<section class="hero-section">
			<div class="hero-content">
				<Logo {firstRender} />
				{#if !firstRender}
					<button
						class="scroll-button"
						onclick={scrollToCarousel}
						aria-label="Défiler vers le contenu"
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

		<!-- Ensuite le contenu de recherche -->
		<section class="content-section carousel-section">
			<SearchBar {artists} bind:selectedArtist />
			<Carousel bind:selectedArtist {artists} bind:firstRender />
			{#if selectedArtist}
				<ArtistDetail bind:artist={selectedArtist} onClose={() => (selectedArtist = null)} />
			{/if}
			<div class="navigation-buttons">
				<button class="nav-button" onclick={() => navigate('/concerts')}> Concerts </button>
				<button class="nav-button" onclick={() => navigate('/stats')}> Statistiques </button>
			</div>
			<button class="scroll-button" onclick={scrollToMap} aria-label="Défiler vers la carte">
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
					<path d="M7 13l5 5 5-5" />
					<path d="M7 6l5 5 5-5" />
				</svg>
			</button>
		</section>
		<section class="content-section map-section">
			<div class="title">
				<h1>Concerts</h1>
			</div>
			<div class="map-container">
				<div class="map-left">
					<ConcertsMap />
				</div>
				<div class="map-right">
					<!-- Espace pour contenu futur -->
				</div>
			</div>
		</section>
	{:else if currentPage === '/concerts'}
		<section class="content-section">
			<h1>Concerts</h1>
			<ConcertsView artist={selectedArtist} onBack={() => navigate('/')} />
		</section>
	{:else if currentPage === '/stats'}
		<section class="content-section">
			<StatsView onBack={() => navigate('/')} />
		</section>
	{:else}
		<section class="content-section">
			<h1>404</h1>
		</section>
	{/if}
	<Footer />
</main>

<style>
	main {
		display: flex;
		flex-direction: column;
	}

	.hero-section {
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
		min-height: 100vh;
		box-sizing: border-box;
		padding: 2rem 0 4rem 0;
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

	.map-container {
		display: flex;
		width: 100%;
		height: 100vh;
		gap: 2rem;
		padding: 0 2rem;
	}

	@media (max-width: 1000px) {
		.map-container {
			flex-direction: column;
		}
	}

	.map-left {
		flex: 1;
		display: flex;
		flex-direction: column;
	}

	.map-right {
		flex: 1;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		padding: 2rem;
		background-color: var(--dark-muted);
		border-radius: 1.5rem;
	}
</style>
