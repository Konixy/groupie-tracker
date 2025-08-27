<script lang="ts">
	import { fade, scale } from 'svelte/transition';
	import type { Artist } from '@/types';

	let {
		artists = [],
		current = $bindable()
	}: {
		artists: Artist[];
		current: number;
	} = $props();

	let query = $state('');
	let results = $state<Artist[]>([]);

	$effect(() => {
		if (query !== null) {
			handleSearch();
		}
	});

	function handleSearch() {
		if (!query.trim()) {
			results = [];
			return;
		}

		const searchTerm = query.toLowerCase();

		// Créer un tableau avec les artistes et leur score de pertinence
		const scoredResults = artists.map((artist) => {
			const artistName = artist.name.toLowerCase();
			let score = 0;

			// Score le plus élevé : commence par la recherche
			if (artistName.startsWith(searchTerm)) {
				score = 1000;
			}
			// Score élevé : contient la recherche quelque part dans le nom
			else if (artistName.includes(searchTerm)) {
				const index = artistName.indexOf(searchTerm);
				score = 1000 - index; // Plus tôt dans le nom = score plus élevé
			}

			return { artist, score };
		});

		// Filtrer les artistes qui correspondent et les trier par score
		const finalResults = scoredResults
			.filter((item) => item.score > 0)
			.sort((a, b) => {
				// Tri par score décroissant
				if (b.score !== a.score) {
					return b.score - a.score;
				}
				// Si même score, trier par ordre alphabétique
				return a.artist.name.localeCompare(b.artist.name);
			})
			.map((item) => item.artist)
			// Éliminer les doublons en gardant seulement la première occurrence
			.filter((artist, index, self) => index === self.findIndex((a) => a.id === artist.id));

		results = finalResults;
	}

	function handleSelect(artist: Artist) {
		current = artist.id - 1;
		query = '';
	}

	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Enter' && results.length > 0) {
			handleSelect(results[0]);
		}
	}
</script>

<div class="search-bar">
	<svg
		xmlns="http://www.w3.org/2000/svg"
		width="20"
		height="20"
		viewBox="0 0 24 24"
		fill="none"
		stroke="currentColor"
		stroke-width="3"
		stroke-linecap="round"
		stroke-linejoin="round"
		class="search-icon"
		aria-hidden="true"
	>
		<path d="m21 21-4.34-4.34" />
		<circle cx="11" cy="11" r="8" />
	</svg>

	<input
		type="text"
		placeholder="Chercher un artiste ou un groupe…"
		bind:value={query}
		onkeydown={handleKeyDown}
		aria-label="Rechercher un artiste"
	/>

	{#if results.length > 0}
		<div class="search-results" transition:scale={{ duration: 100, start: 0.5 }}>
			{#each results.slice(0, 10) as result}
				<button
					class="search-result"
					onclick={() => handleSelect(result)}
					aria-label="Sélectionner {result.name}"
				>
					{result.name}
				</button>
			{/each}
		</div>
	{/if}
</div>

<style>
	.search-bar {
		position: relative;
		display: flex;
		justify-content: center;
		align-items: center;
		margin-top: 4rem;
	}

	input {
		border-radius: 20px;
		border: none;
		padding: 1rem;
		width: 350px;
		font-size: 1rem;
		font-family: inherit;
		outline: none;
		background: var(--light-vibrant);
		color: var(--dark-vibrant);
		padding-left: 2.8rem;
		margin-right: -1rem;
		transition:
			background-color 0.3s ease,
			color 0.3s ease;
	}

	input::placeholder {
		color: currentColor;
		opacity: 0.7;
	}

	.search-icon {
		margin-right: -2.2rem;
		color: var(--dark-vibrant);
		z-index: 100;
		transition: color 0.3s ease;
	}

	.search-results {
		display: flex;
		position: absolute;
		top: 3rem;
		padding: 0.5rem;
		border-radius: 20px;
		background: var(--light-vibrant);
		width: 350px;
		z-index: 100;
		flex-direction: column;
		align-items: center;
		margin-top: 1rem;
		transition: all 0.3s ease;
		box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
	}

	.search-result {
		width: 100%;
		text-align: left;
		padding: 0.75rem;
		border-radius: 15px;
		cursor: pointer;
		transition:
			background-color 0.3s ease,
			color 0.3s ease;
		border: none;
		background: none;
		outline: none;
		font-family: inherit;
		font-size: 1rem;
		color: var(--dark-vibrant);
	}

	.search-result:hover {
		background: var(--dark-muted);
		color: var(--light-vibrant);
	}

	.search-result:focus {
		background: var(--dark-muted);
		color: var(--light-vibrant);
	}
</style>
