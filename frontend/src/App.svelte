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
	let firstRender = $state(true);

	onMount(async () => {
		const res = await fetch('http://localhost:8080/artists');
		const data = await res.json();
		// Supporte à la fois l'API locale (data.artists) et distante (data)
		artists = Array.isArray(data) ? data : data.artists;
	});

	function scrollToContent() {
		const contentSection = document.querySelector('.content-section');
		if (contentSection) {
			contentSection.scrollIntoView({ behavior: 'smooth' });
		}
	}

	function scrollToMap() {
		const mapSection = document.querySelector('.map-section');
		if (mapSection) {
			mapSection.scrollIntoView({ behavior: 'smooth' });
		}
	}
</script>

<main>
	<!-- D'abord le titre qu'on va animer après -->
	<section class="hero-section">
		<div class="hero-content">
			<Logo {firstRender} />
			{#if !firstRender}
				<button
					class="scroll-button"
					onclick={scrollToContent}
					aria-label="Défiler vers le contenu"
				>
					<svg
						width="24"
						height="24"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
					>
						<path d="M7 13l5 5 5-5" />
						<path d="M7 6l5 5 5-5" />
					</svg>
				</button>
			{/if}
		</div>
	</section>

	<!-- Ensuite le contenu de recherche -->
	<section class="content-section">
		<SearchBar {artists} bind:selectedArtist />
		<Carousel bind:selectedArtist {artists} bind:firstRender />
		{#if selectedArtist}
			<ArtistDetail bind:artist={selectedArtist} onClose={() => (selectedArtist = null)} />
		{/if}
		<button class="scroll-button" onclick={scrollToMap} aria-label="Défiler vers la carte">
			<svg
				width="24"
				height="24"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
			>
				<path d="M7 13l5 5 5-5" />
				<path d="M7 6l5 5 5-5" />
			</svg>
		</button>
	</section>
	<section class="content-section map-section">
		<div class="map-container">
			<div class="map-left">
				<ConcertsMap />
			</div>
			<div class="map-right">
				<!-- Espace pour contenu futur -->
			</div>
		</div>
		<Footer />
	</section>
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
		gap: 2rem;
	}

	.scroll-button {
		background: none;
		border: 2px solid currentColor;
		border-radius: 50%;
		width: 50px;
		height: 50px;
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		transition: all 0.3s ease;
		color: inherit;
	}

	.scroll-button:hover {
		transform: translateY(5px);
		background: rgba(255, 255, 255, 0.1);
	}

	.scroll-button svg {
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

	.map-container {
		display: flex;
		width: 100%;
		height: 100vh;
		gap: 2rem;
		padding: 0 2rem;
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
