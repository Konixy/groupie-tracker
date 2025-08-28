<script lang="ts">
	import { scale } from 'svelte/transition';
	import type { Artist } from '@/types';

	let {
		artists = [],
		current = $bindable(),
		onArtistSelect = () => {},
		onMemberSearch = () => {},
		onDateSearch = () => {},
		onAlbumSearch = () => {}
	}: {
		artists: Artist[];
		current: number;
		onArtistSelect?: (artist: Artist) => void;
		onMemberSearch?: (memberName: string, artist: Artist) => void;
		onDateSearch?: (date: string, artist: Artist) => void;
		onAlbumSearch?: (album: string, artist: Artist) => void;
	} = $props();

	let query = $state('');
	let results = $state<SearchResult[]>([]);
	let showResults = $state(false);
	let allLocations = $state<Record<string, string[]>>({});

	// Types de résultats de recherche
	type SearchResultType = 'artist' | 'member' | 'location' | 'creation_date' | 'first_album';

	interface SearchResult {
		type: SearchResultType;
		artist: Artist;
		displayText: string;
		subText: string;
		score: number;
	}

	// Charger tous les lieux disponibles au montage du composant
	$effect(() => {
		loadAllLocations();
	});

	async function loadAllLocations() {
		try {
			const response = await fetch('http://localhost:8080/all-locations');
			if (response.ok) {
				allLocations = await response.json();
			}
		} catch (error) {
			console.error('Erreur lors du chargement des lieux:', error);
		}
	}

	$effect(() => {
		if (query !== null) {
			handleSearch();
		}
	});

	function handleSearch() {
		if (!query.trim()) {
			results = [];
			showResults = false;
			return;
		}

		const searchTerm = query.toLowerCase();
		let allResults: SearchResult[] = [];

		// 1. Recherche par nom d'artiste/groupe
		artists.forEach((artist) => {
			const artistName = artist.name.toLowerCase();
			let score = 0;

			if (artistName.startsWith(searchTerm)) {
				score = 1000;
			} else if (artistName.includes(searchTerm)) {
				const index = artistName.indexOf(searchTerm);
				score = 1000 - index;
			}

			if (score > 0) {
				allResults.push({
					type: 'artist',
					artist,
					displayText: artist.name,
					subText: 'Artiste/Groupe',
					score
				});
			}
		});

		// 2. Recherche par membres
		artists.forEach((artist) => {
			artist.members.forEach((member) => {
				const memberName = member.toLowerCase();
				let score = 0;

				if (memberName.startsWith(searchTerm)) {
					score = 900;
				} else if (memberName.includes(searchTerm)) {
					const index = memberName.indexOf(searchTerm);
					score = 900 - index;
				}

				if (score > 0) {
					allResults.push({
						type: 'member',
						artist,
						displayText: member,
						subText: `Membre de ${artist.name}`,
						score
					});
				}
			});
		});

		// 3. Recherche par emplacements de concerts
		Object.entries(allLocations).forEach(([location, artistNames]) => {
			const locationName = location.toLowerCase();
			if (locationName.includes(searchTerm)) {
				const artist = artists.find((a) => artistNames.includes(a.name));
				if (artist) {
					allResults.push({
						type: 'location',
						artist,
						displayText: `Concert de ${artist.name}`,
						subText: `Lieu: ${location}`,
						score: 800
					});
				}
			}
		});

		// 4. Recherche par date de création
		artists.forEach((artist) => {
			const creationYear = artist.creationDate.toString();
			if (creationYear.includes(searchTerm)) {
				allResults.push({
					type: 'creation_date',
					artist,
					displayText: artist.name,
					subText: `Créé en ${artist.creationDate}`,
					score: 700
				});
			}
		});

		// 5. Recherche par date du premier album
		artists.forEach((artist) => {
			if (artist.firstAlbum) {
				const albumYear = artist.firstAlbum.toLowerCase();
				if (albumYear.includes(searchTerm)) {
					allResults.push({
						type: 'first_album',
						artist,
						displayText: artist.name,
						subText: `Premier album: ${artist.firstAlbum}`,
						score: 600
					});
				}
			}
		});

		// Trier par score et type, éliminer les doublons
		const finalResults = allResults
			.sort((a, b) => {
				// D'abord par score décroissant
				if (b.score !== a.score) {
					return b.score - a.score;
				}
				// Puis par ordre de priorité des types
				const typePriority = {
					artist: 1,
					member: 2,
					location: 3,
					creation_date: 4,
					first_album: 5
				};
				return typePriority[a.type] - typePriority[b.type];
			})
			.filter(
				(result, index, self) =>
					index ===
					self.findIndex((r) => r.artist.id === result.artist.id && r.type === result.type)
			)
			.slice(0, 10);

		results = finalResults;
		showResults = finalResults.length > 0;
	}

	function handleSelect(result: SearchResult) {
		// Rediriger selon le type de résultat
		switch (result.type) {
			case 'location':
				// Aller à la section map avec l'artiste sélectionné
				current = result.artist.id - 1;
				scrollToMap();
				break;
			case 'member':
				// Pour les membres, aller à l'artiste ET déclencher l'effet de clignotement
				current = result.artist.id - 1;
				onMemberSearch(result.displayText, result.artist);
				onArtistSelect(result.artist);
				break;
			case 'creation_date':
				// Pour les dates de création, aller à l'artiste ET déclencher l'effet de clignotement
				current = result.artist.id - 1;
				onDateSearch(result.artist.creationDate.toString(), result.artist);
				onArtistSelect(result.artist);
				break;
			case 'first_album':
				// Pour les premiers albums, aller à l'artiste ET déclencher l'effet de clignotement
				current = result.artist.id - 1;
				onAlbumSearch(result.artist.firstAlbum, result.artist);
				onArtistSelect(result.artist);
				break;
			default:
				// Pour tous les autres types, aller à l'artiste dans le carousel ET ouvrir la card
				current = result.artist.id - 1;
				onArtistSelect(result.artist);
				break;
		}

		query = '';
		showResults = false;
	}

	function scrollToMap() {
		const mapSection = document.getElementById('map-section');
		if (mapSection) {
			mapSection.scrollIntoView({ behavior: 'smooth' });
		}
	}

	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Enter' && results.length > 0) {
			handleSelect(results[0]);
		} else if (event.key === 'Escape') {
			showResults = false;
			query = '';
		}
	}

	function handleInputFocus() {
		if (results.length > 0) {
			showResults = true;
		}
	}

	function handleInputBlur() {
		// Délai pour permettre le clic sur les résultats
		setTimeout(() => {
			showResults = false;
		}, 200);
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
		placeholder="Chercher un artiste, membre, lieu ou date…"
		bind:value={query}
		onkeydown={handleKeyDown}
		onfocus={handleInputFocus}
		onblur={handleInputBlur}
		aria-label="Recherche universelle"
	/>

	{#if showResults && results.length > 0}
		<div class="search-results" transition:scale={{ duration: 100, start: 0.5 }}>
			{#each results as result}
				<button
					class="search-result"
					onclick={() => handleSelect(result)}
					aria-label="Sélectionner {result.displayText}"
				>
					<div class="result-content">
						<div class="result-main">{result.displayText}</div>
						<div class="result-sub">{result.subText}</div>
					</div>
					<div class="result-type-badge {result.type}">
						{#if result.type === 'artist'}Artiste{/if}
						{#if result.type === 'member'}Membre{/if}
						{#if result.type === 'location'}Lieu{/if}
						{#if result.type === 'creation_date'}Création{/if}
						{#if result.type === 'first_album'}Album{/if}
					</div>
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
		width: 400px;
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
		width: 400px;
		z-index: 100;
		flex-direction: column;
		align-items: center;
		margin-top: 1rem;
		transition: all 0.3s ease;
		box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
		max-height: 400px;
		overflow-y: auto;
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
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.search-result:hover {
		background: var(--dark-muted);
		color: var(--light-vibrant);
	}

	.search-result:focus {
		background: var(--dark-muted);
		color: var(--light-vibrant);
	}

	.result-content {
		flex: 1;
		text-align: left;
	}

	.result-main {
		font-weight: 600;
		margin-bottom: 0.25rem;
	}

	.result-sub {
		font-size: 0.85rem;
		opacity: 0.8;
	}

	.result-type-badge {
		font-size: 0.7rem;
		font-weight: 600;
		margin-left: 0.5rem;
		padding: 0.3rem 0.6rem;
		border-radius: 6px;
		background: var(--dark-muted);
		color: var(--light-vibrant);
		min-width: 60px;
		text-align: center;
		white-space: nowrap;
	}

	/* Responsive */
	@media (max-width: 768px) {
		input {
			width: 300px;
		}

		.search-results {
			width: 300px;
		}
	}
</style>
