<script lang="ts">
	import { onMount } from 'svelte';
	import type { Artist } from '../../types';

	// Props reçues du composant parent
	let { artists } = $props<{ artists: Artist[] }>();

	// Types pour les statistiques
	type LocationInfo = {
		city: string;
		country: string;
		continent: string;
		date: string;
	};

	type ArtistStats = {
		id: number;
		name: string;
		image: string;
		concertCount: number;
		years: number[];
		locations: LocationInfo[];
		continents: Record<string, number>;
	};

	// État des statistiques
	let artistsStats = $state<ArtistStats[]>([]);
	let availableYears = $state<number[]>([]);
	let statsLoading = $state(false);
	let statsError = $state<string | null>(null);

	// Filtres pour les statistiques
	let selectedYear = $state<number | null>(null);
	let yearRange = $state({ start: 1990, end: 2030 });
	let useRange = $state(false);
	let selectedContinent = $state<string | null>(null);
	let selectedCountry = $state<string | null>(null);

	// Chargement des statistiques au montage
	onMount(async () => {
		await loadStats();
	});

	async function loadStats() {
		statsLoading = true;
		statsError = null;

		try {
			const response = await fetch('http://localhost:8080/stats');
			if (!response.ok) {
				throw new Error('Erreur lors du chargement des statistiques');
			}
			const data = await response.json();
			artistsStats = data.artists || [];
			availableYears = data.availableYears || [];
		} catch (err) {
			statsError = err instanceof Error ? err.message : 'Erreur inconnue';
		} finally {
			statsLoading = false;
		}
	}

	// Métriques calculées pour les statistiques
	let filteredArtists = $derived(() => {
		let filtered = artistsStats;

		if (useRange) {
			filtered = filtered.filter(
				(artist) =>
					artist.years?.some((year) => year >= yearRange.start && year <= yearRange.end) || false
			);
		} else if (selectedYear !== null) {
			filtered = filtered.filter((artist) => artist.years?.includes(selectedYear!) || false);
		}

		if (selectedContinent !== null) {
			filtered = filtered.filter(
				(artist) => artist.continents && artist.continents[selectedContinent!] > 0
			);
		}

		if (selectedCountry !== null) {
			filtered = filtered.filter((artist) =>
				artist.locations.some((location) => location.country === selectedCountry)
			);
		}

		return filtered;
	});

	let totalConcerts = $derived(
		filteredArtists().reduce((sum, artist) => sum + artist.concertCount, 0)
	);
	let averageConcerts = $derived(
		filteredArtists().length > 0 ? totalConcerts / filteredArtists().length : 0
	);
	let topArtist = $derived(filteredArtists().length > 0 ? filteredArtists()[0] : null);
	let categories = $derived({
		high: filteredArtists().filter((a) => a.concertCount >= 15).length,
		medium: filteredArtists().filter((a) => a.concertCount >= 8 && a.concertCount < 15).length,
		low: filteredArtists().filter((a) => a.concertCount < 8).length
	});

	function handleArtistClick(artistId: number) {
		const artist = artists.find((a: Artist) => a.id === artistId);
		if (artist) {
			// Dispatch d'un événement personnalisé pour informer le parent
			const event = new CustomEvent('artistSelected', {
				detail: { artist }
			});
			document.dispatchEvent(event);
		}
	}

	function clearStatsFilters() {
		selectedYear = null;
		yearRange = { start: 1990, end: 2030 };
		useRange = false;
		selectedContinent = null;
		selectedCountry = null;
	}

	function getContinentName(continent: string): string {
		const continentNames: Record<string, string> = {
			'North America': 'Amérique du Nord',
			'South America': 'Amérique du Sud',
			Europe: 'Europe',
			Asia: 'Asie',
			Africa: 'Afrique',
			Oceania: 'Océanie'
		};
		return continentNames[continent] || continent;
	}
</script>

<!-- Section 4: Statistiques -->

<div class="stats-container">
	<div class="stats-header">
		<h1>Statistiques des Artistes</h1>
		<p>Explorez les données de concerts et performances</p>
	</div>

	{#if statsLoading}
		<div class="loading" role="status" aria-live="polite">Chargement des statistiques...</div>
	{:else if statsError}
		<div class="error" role="alert">
			Erreur: {statsError}
		</div>
	{:else}
		<!-- Filtres -->
		<div class="filters-section">
			<h2>Filtres</h2>
			<div class="filters-grid">
				<div class="filter-group">
					<label for="year-select">Année</label>
					<select id="year-select" bind:value={selectedYear} disabled={useRange}>
						<option value={null}>Toutes les années</option>
						{#each availableYears as year}
							<option value={year}>{year}</option>
						{/each}
					</select>
				</div>

				<div class="filter-group">
					<label>
						<input type="checkbox" bind:checked={useRange} onchange={() => (selectedYear = null)} />
						Plage d'années
					</label>
					{#if useRange}
						<div class="range-inputs">
							<input
								type="number"
								bind:value={yearRange.start}
								min="1900"
								max="2030"
								placeholder="Début"
							/>
							<span>-</span>
							<input
								type="number"
								bind:value={yearRange.end}
								min="1900"
								max="2030"
								placeholder="Fin"
							/>
						</div>
					{/if}
				</div>

				<div class="filter-group">
					<label for="continent-select">Continent</label>
					<select id="continent-select" bind:value={selectedContinent}>
						<option value={null}>Tous les continents</option>
						<option value="Europe">Europe</option>
						<option value="North America">Amérique du Nord</option>
						<option value="South America">Amérique du Sud</option>
						<option value="Asia">Asie</option>
						<option value="Africa">Afrique</option>
						<option value="Oceania">Océanie</option>
					</select>
				</div>
			</div>

			<button class="clear-filters" onclick={clearStatsFilters}> Effacer les filtres </button>
		</div>

		<!-- Métriques principales -->
		<div class="metrics-section">
			<h2>Métriques</h2>
			<div class="metrics-grid">
				<div class="metric-card">
					<h3>Artistes</h3>
					<p class="metric-value">{filteredArtists().length}</p>
				</div>
				<div class="metric-card">
					<h3>Total Concerts</h3>
					<p class="metric-value">{totalConcerts}</p>
				</div>
				<div class="metric-card">
					<h3>Moyenne par Artiste</h3>
					<p class="metric-value">{averageConcerts.toFixed(1)}</p>
				</div>
				<div class="metric-card">
					<h3>Top Artiste</h3>
					<p class="metric-value">{topArtist?.name || 'N/A'}</p>
				</div>
			</div>
		</div>

		<!-- Répartition par catégorie -->
		<div class="categories-section">
			<h2>Répartition par nombre de concerts</h2>
			<div class="categories-grid">
				<div class="category-card high">
					<h3>Élevé (15+ concerts)</h3>
					<p class="category-count">{categories.high}</p>
				</div>
				<div class="category-card medium">
					<h3>Moyen (8-14 concerts)</h3>
					<p class="category-count">{categories.medium}</p>
				</div>
				<div class="category-card low">
					<h3>Faible (&lt; 8 concerts)</h3>
					<p class="category-count">{categories.low}</p>
				</div>
			</div>
		</div>

		<!-- Liste des artistes -->
		<div class="artists-section">
			<h2>Artistes ({filteredArtists().length})</h2>
			<div class="artists-grid">
				{#each filteredArtists() as artist}
					<button
						class="artist-card"
						onclick={() => handleArtistClick(artist.id)}
						aria-label="Voir les détails de {artist.name}"
					>
						<img src={artist.image} alt="Photo de {artist.name}" class="artist-image" />
						<div class="artist-info">
							<h3>{artist.name}</h3>
							<p>{artist.concertCount} concerts</p>
							<div class="artist-years">
								{#each artist.years?.slice(0, 3) as year}
									<span class="year-tag">{year}</span>
								{/each}
								{#if artist.years && artist.years.length > 3}
									<span class="year-tag">+{artist.years.length - 3}</span>
								{/if}
							</div>
						</div>
					</button>
				{/each}
			</div>
		</div>
	{/if}
</div>

<style>
	.stats-container {
		width: 100%;
		max-width: 1200px;
		margin: 0 auto;
		height: 100%;
	}

	.stats-header {
		text-align: center;
		margin-bottom: 3rem;
	}

	.stats-header h1 {
		color: var(--light-vibrant);
		font-size: 3rem;
		margin: 0 0 1rem 0;
	}

	.stats-header p {
		color: var(--light-muted);
		font-size: 1.2rem;
		margin: 0;
	}

	.loading,
	.error {
		text-align: center;
		padding: 3rem;
		color: var(--light-vibrant);
		font-size: 1.2rem;
	}

	.error {
		color: #ff6b6b;
		background: rgba(255, 107, 107, 0.1);
		border-radius: 8px;
		margin: 2rem 0;
	}

	.filters-section {
		background: var(--light-vibrant);
		border-radius: 16px;
		padding: 2rem;
		margin-bottom: 2rem;
		box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
	}

	.filters-section h2 {
		color: var(--dark-vibrant);
		margin: 0 0 1.5rem 0;
		font-size: 1.5rem;
	}

	.filters-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
		gap: 1.5rem;
		margin-bottom: 1.5rem;
	}

	.filter-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.filter-group label {
		color: var(--dark-vibrant);
		font-weight: 500;
	}

	.filter-group select,
	.filter-group input[type='number'] {
		padding: 0.8rem;
		border: 2px solid var(--light-muted);
		border-radius: 8px;
		font-size: 1rem;
		background: var(--light-vibrant);
		color: var(--dark-vibrant);
		transition: border-color 0.3s ease;
	}

	.filter-group select:focus,
	.filter-group input[type='number']:focus {
		outline: none;
		border-color: var(--vibrant);
	}

	.range-inputs {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.range-inputs input {
		width: 100px;
	}

	.clear-filters {
		background: var(--dark-muted);
		color: var(--light-vibrant);
		border: none;
		padding: 0.8rem 1.5rem;
		border-radius: 8px;
		cursor: pointer;
		font-size: 1rem;
		transition: background 0.3s ease;
	}

	.clear-filters:hover {
		background: var(--muted);
	}

	.metrics-section {
		background: var(--light-vibrant);
		border-radius: 16px;
		padding: 2rem;
		margin-bottom: 2rem;
		box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
	}

	.metrics-section h2 {
		color: var(--dark-vibrant);
		margin: 0 0 1.5rem 0;
		font-size: 1.5rem;
	}

	.metrics-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: 1.5rem;
	}

	.metric-card {
		background: var(--light-muted);
		padding: 1.5rem;
		border-radius: 12px;
		text-align: center;
		transition: transform 0.2s ease;
	}

	.metric-card:hover {
		transform: translateY(-2px);
	}

	.metric-card h3 {
		color: var(--dark-vibrant);
		margin: 0 0 0.5rem 0;
		font-size: 1.1rem;
	}

	.metric-value {
		color: var(--dark-muted);
		font-size: 2rem;
		font-weight: bold;
		margin: 0;
	}

	.categories-section {
		background: var(--light-vibrant);
		border-radius: 16px;
		padding: 2rem;
		margin-bottom: 2rem;
		box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
	}

	.categories-section h2 {
		color: var(--dark-vibrant);
		margin: 0 0 1.5rem 0;
		font-size: 1.5rem;
	}

	.categories-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: 1.5rem;
	}

	.category-card {
		padding: 1.5rem;
		border-radius: 12px;
		text-align: center;
		transition: transform 0.2s ease;
	}

	.category-card.high {
		background: #e8f5e8;
		border-left: 4px solid #4caf50;
	}

	.category-card.medium {
		background: #fff3e0;
		border-left: 4px solid #ff9800;
	}

	.category-card.low {
		background: #ffebee;
		border-left: 4px solid #f44336;
	}

	.category-card:hover {
		transform: translateY(-2px);
	}

	.category-card h3 {
		color: var(--dark-vibrant);
		margin: 0 0 0.5rem 0;
		font-size: 1.1rem;
	}

	.category-count {
		color: var(--dark-muted);
		font-size: 2rem;
		font-weight: bold;
		margin: 0;
	}

	.artists-section {
		background: var(--light-vibrant);
		border-radius: 16px;
		padding: 2rem;
		box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
	}

	.artists-section h2 {
		color: var(--dark-vibrant);
		margin: 0 0 1.5rem 0;
		font-size: 1.5rem;
	}

	.artists-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: 1.5rem;
	}

	.artist-card {
		background: var(--light-muted);
		border-radius: 12px;
		padding: 1.5rem;
		cursor: pointer;
		transition:
			transform 0.2s ease,
			box-shadow 0.2s ease;
		display: flex;
		align-items: center;
		gap: 1rem;
		border: none;
		width: 100%;
		text-align: left;
	}

	.artist-card:hover {
		transform: translateY(-2px);
		box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
	}

	.artist-image {
		width: 80px;
		height: 80px;
		object-fit: cover;
		border-radius: 8px;
	}

	.artist-info {
		flex: 1;
	}

	.artist-info h3 {
		color: var(--dark-vibrant);
		margin: 0 0 0.5rem 0;
		font-size: 1.2rem;
	}

	.artist-info p {
		color: var(--dark-muted);
		margin: 0 0 0.5rem 0;
		font-size: 1rem;
	}

	.artist-years {
		display: flex;
		flex-wrap: wrap;
		gap: 0.3rem;
	}

	.year-tag {
		background: var(--vibrant);
		color: var(--light-vibrant);
		padding: 0.2rem 0.5rem;
		border-radius: 12px;
		font-size: 0.8rem;
		font-weight: 500;
	}

	@media (max-width: 768px) {
		.stats-container {
			padding: 1rem;
		}

		.stats-header h1 {
			font-size: 2rem;
		}

		.filters-grid,
		.metrics-grid,
		.categories-grid {
			grid-template-columns: 1fr;
		}

		.artists-grid {
			grid-template-columns: 1fr;
		}
	}

	.content-section {
		min-height: 100vh;
		height: 100vh;
		box-sizing: border-box;
		padding: 2rem 0;
		display: flex;
		flex-direction: column;
		align-items: center;
		overflow-y: auto;
	}
</style>
