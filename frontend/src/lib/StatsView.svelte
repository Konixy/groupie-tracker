<script lang="ts">
	import { onMount } from 'svelte';

	let { onBack, onArtistClick }: { onBack: () => void; onArtistClick: (artistId: number) => void } =
		$props();

	// Types
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

	// État des données
	let artistsStats = $state<ArtistStats[]>([]);
	let availableYears = $state<number[]>([]);
	let loading = $state(false);
	let error = $state<string | null>(null);

	// Filtres
	let selectedYear = $state<number | null>(null);
	let yearRange = $state({ start: 1990, end: 2030 });
	let useRange = $state(false);
	let selectedContinent = $state<string | null>(null);
	let selectedCountry = $state<string | null>(null);

	// Métriques calculées basées sur les filtres
	let filteredArtists = $derived(() => {
		let filtered = artistsStats;

		// Filtre par année
		if (useRange) {
			filtered = filtered.filter(
				(artist) =>
					artist.years?.some((year) => year >= yearRange.start && year <= yearRange.end) || false
			);
		} else if (selectedYear !== null) {
			filtered = filtered.filter((artist) => artist.years?.includes(selectedYear!) || false);
		}

		// Filtre par continent
		if (selectedContinent !== null) {
			filtered = filtered.filter(
				(artist) => artist.continents && artist.continents[selectedContinent!] > 0
			);
		}

		// Filtre par pays
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

	onMount(async () => {
		await loadStats();
	});

	async function loadStats() {
		loading = true;
		error = null;

		try {
			const response = await fetch('http://localhost:8080/stats');

			if (!response.ok) {
				throw new Error('Erreur lors du chargement des statistiques');
			}

			const data = await response.json();
			artistsStats = data.artists || [];
			availableYears = data.availableYears || [];

			// Trier par nombre de concerts décroissant
			artistsStats.sort((a, b) => b.concertCount - a.concertCount);

			// Initialiser les filtres
			if (availableYears.length > 0) {
				yearRange.start = Math.min(...availableYears);
				yearRange.end = Math.max(...availableYears);
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Erreur inconnue';
			console.error('Erreur:', err);
		} finally {
			loading = false;
		}
	}

	function clearFilters() {
		selectedYear = null;
		useRange = false;
		selectedContinent = null;
		selectedCountry = null;
		yearRange = { start: Math.min(...availableYears), end: Math.max(...availableYears) };
	}

	function applyYearFilter(year: number | null) {
		selectedYear = year;
		useRange = false;
	}

	function applyRangeFilter() {
		selectedYear = null;
		useRange = true;
	}

	function applyContinentFilter(continent: string | null) {
		selectedContinent = continent;
		selectedCountry = null; // Réinitialiser le pays quand on change de continent
	}

	function applyCountryFilter(country: string | null) {
		selectedCountry = country;
	}

	function handleArtistClick(artistId: number) {
		onArtistClick(artistId);
	}

	// Obtenir les pays disponibles selon le continent sélectionné
	let availableCountries = $derived(() => {
		if (!selectedContinent) return [];

		const countries = new Set<string>();
		artistsStats.forEach((artist) => {
			if (artist.continents && artist.continents[selectedContinent!] > 0) {
				artist.locations.forEach((location) => {
					if (location.continent === selectedContinent) {
						countries.add(location.country);
					}
				});
			}
		});

		return Array.from(countries).sort();
	});
</script>

<div class="stats-container">
	<button class="back-button" onclick={onBack}> ← Retour </button>

	<div class="stats-content">
		<h1>📊 Statistiques des Concerts</h1>
		<p class="subtitle">Analyse complète des concerts répertoriés dans l'API</p>

		{#if loading}
			<div class="loading">Chargement des statistiques...</div>
		{:else if error}
			<div class="error">Erreur: {error}</div>
		{:else if artistsStats.length > 0}
			<!-- Filtres -->
			<div class="filters-section">
				<h2>🔍 Filtres par Année</h2>
				<div class="filters-container">
					<div class="filter-group">
						<label>Année spécifique:</label>
						<select
							bind:value={selectedYear}
							onchange={() => applyYearFilter(selectedYear)}
							disabled={useRange}
						>
							<option value={null}>Toutes les années</option>
							{#each availableYears as year}
								<option value={year}>{year}</option>
							{/each}
						</select>
					</div>

					<div class="filter-group">
						<label>
							<input type="checkbox" bind:checked={useRange} onchange={applyRangeFilter} />
							Tranche d'années
						</label>
						<div class="range-inputs" class:disabled={!useRange}>
							<input
								type="number"
								bind:value={yearRange.start}
								min={Math.min(...availableYears)}
								max={Math.max(...availableYears)}
								disabled={!useRange}
							/>
							<span>à</span>
							<input
								type="number"
								bind:value={yearRange.end}
								min={Math.min(...availableYears)}
								max={Math.max(...availableYears)}
								disabled={!useRange}
							/>
						</div>
					</div>

					<div class="filter-group">
						<label>🌍 Continent</label>
						<select
							value={selectedContinent}
							onchange={(e) => applyContinentFilter((e.target as HTMLSelectElement).value || null)}
						>
							<option value={null}>Tous les continents</option>
							<option value="Europe">Europe</option>
							<option value="Amérique du Nord">Amérique du Nord</option>
							<option value="Amérique du Sud">Amérique du Sud</option>
							<option value="Asie">Asie</option>
							<option value="Afrique">Afrique</option>
							<option value="Océanie">Océanie</option>
						</select>
					</div>

					{#if selectedContinent && availableCountries().length > 0}
						<div class="filter-group">
							<label>🏳️ Pays</label>
							<select
								value={selectedCountry}
								onchange={(e) => applyCountryFilter((e.target as HTMLSelectElement).value || null)}
							>
								<option value={null}>Tous les pays</option>
								{#each availableCountries() as country}
									<option value={country}>{country}</option>
								{/each}
							</select>
						</div>
					{/if}

					<button class="clear-filters" onclick={clearFilters}> 🗑️ Effacer les filtres </button>
				</div>

				{#if selectedYear !== null || useRange || selectedContinent !== null || selectedCountry !== null}
					<div class="active-filters">
						{#if selectedYear !== null}
							<span class="filter-badge">
								Année: {selectedYear}
							</span>
						{/if}
						{#if useRange}
							<span class="filter-badge">
								Tranche: {yearRange.start}-{yearRange.end}
							</span>
						{/if}
						{#if selectedContinent !== null}
							<span class="filter-badge">
								Continent: {selectedContinent}
							</span>
						{/if}
						{#if selectedCountry !== null}
							<span class="filter-badge">
								Pays: {selectedCountry}
							</span>
						{/if}
						<span class="results-count">
							{filteredArtists.length} artiste{filteredArtists.length > 1 ? 's' : ''} trouvé{filteredArtists.length >
							1
								? 's'
								: ''}
						</span>
					</div>
				{/if}
			</div>

			<!-- Métriques principales -->
			<div class="metric-card">
				<div class="metric-icon">🎸</div>
				<div class="metric-value">{totalConcerts}</div>
				<div class="metric-label">Concerts Total</div>
			</div>
			<div class="metric-card">
				<div class="metric-icon">📈</div>
				<div class="metric-value">{averageConcerts.toFixed(1)}</div>
				<div class="metric-label">Moyenne/Artiste</div>
			</div>
			<div class="metric-card">
				<div class="metric-icon">🏆</div>
				<div class="metric-value">{topArtist?.name || 'N/A'}</div>
				<div class="metric-label">Top Artiste</div>
			</div>

			<!-- Graphique des catégories -->
			<div class="chart-section">
				<h2>📊 Répartition par Catégorie</h2>
				<div class="chart-container">
					<div class="chart-bar">
						<div class="bar-label">Étoiles (15+ concerts)</div>
						<div class="bar-track">
							<div
								class="bar-fill high"
								style="width: {(categories.high / filteredArtists.length) * 100}%"
							></div>
						</div>
						<div class="bar-value">{categories.high}</div>
					</div>
					<div class="chart-bar">
						<div class="bar-label">Populaires (8-14 concerts)</div>
						<div class="bar-track">
							<div
								class="bar-fill medium"
								style="width: {(categories.medium / filteredArtists.length) * 100}%"
							></div>
						</div>
						<div class="bar-value">{categories.medium}</div>
					</div>
					<div class="chart-bar">
						<div class="bar-label">Émergents (&lt; 8 concerts)</div>
						<div class="bar-track">
							<div
								class="bar-fill low"
								style="width: {(categories.low / filteredArtists.length) * 100}%"
							></div>
						</div>
						<div class="bar-value">{categories.low}</div>
					</div>
				</div>
			</div>

			<!-- Top 10 des artistes -->
			<div class="top-section">
				<h2>🏆 Top 10 des Artistes</h2>
				<div class="top-grid">
					{#each filteredArtists().slice(0, 10) as artist, index}
						<div
							class="top-card clickable"
							style="--rank: {index + 1}"
							onclick={() => handleArtistClick(artist.id)}
						>
							<div class="rank-number">#{index + 1}</div>
							<div class="artist-image">
								<img src={artist.image} alt={artist.name} />
							</div>
							<div class="artist-info">
								<h3>{artist.name}</h3>
								<div class="concert-count">
									<span class="count">{artist.concertCount}</span>
									<span class="label">concerts</span>
								</div>
							</div>
							<div class="percentage">
								{((artist.concertCount / totalConcerts) * 100).toFixed(1)}%
							</div>
						</div>
					{/each}
				</div>
			</div>

			<!-- Tous les artistes -->
			<div class="all-section">
				<h2>🎵 Tous les Artistes</h2>
				<div class="stats-grid">
					{#each filteredArtists() as artist, index}
						<div
							class="artist-card clickable"
							style="--rank: {index + 1}"
							onclick={() => handleArtistClick(artist.id)}
						>
							<div class="artist-image">
								<img src={artist.image} alt={artist.name} />
							</div>
							<div class="artist-info">
								<h3>{artist.name}</h3>
								<div class="concert-count">
									<span class="count">{artist.concertCount}</span>
									<span class="label">concerts</span>
								</div>
								<div class="continents-info">
									{Object.entries(artist.continents || {})
										.map(([continent, count]) => `${continent}: ${count}`)
										.join(', ') || 'Aucun continent'}
								</div>
								<div class="percentage-small">
									{((artist.concertCount / totalConcerts) * 100).toFixed(1)}% du total
								</div>
							</div>
							<div class="rank-badge">#{index + 1}</div>
						</div>
					{/each}
				</div>
			</div>
		{:else}
			<div class="no-data">Aucune donnée disponible</div>
		{/if}
	</div>
</div>

<style>
	/* Forcer tout le texte en blanc */
	.stats-container,
	.stats-container * {
		color: #ffffff !important;
	}

	.stats-container {
		width: 100%;
		max-width: 1400px;
		margin: 0 auto;
		padding: 2rem;
	}

	.back-button {
		background: none;
		border: 2px solid currentColor;
		border-radius: 8px;
		padding: 0.5rem 1rem;
		color: inherit;
		cursor: pointer;
		transition: all 0.3s ease;
		margin-bottom: 2rem;
		font-size: 1rem;
	}

	.back-button:hover {
		background: rgba(255, 255, 255, 0.1);
		transform: translateY(-2px);
	}

	.stats-content {
		text-align: center;
		padding: 2rem;
		background-color: var(--dark-muted);
		border-radius: 1.5rem;
	}

	.stats-content h1 {
		font-size: 2.5rem;
		margin-bottom: 0.5rem;
		color: var(--light-muted);
	}

	.subtitle {
		font-size: 1.2rem;
		color: var(--light-muted);
		opacity: 0.8;
		margin-bottom: 2rem;
	}

	.loading,
	.error,
	.no-data {
		font-size: 1.2rem;
		padding: 2rem;
		color: var(--light-muted);
	}

	.error {
		color: #ff6b6b;
	}

	/* Filtres */
	.filters-section {
		margin-bottom: 3rem;
		text-align: left;
	}

	.filters-section h2 {
		font-size: 1.8rem;
		color: var(--light-muted);
		margin-bottom: 1.5rem;
		text-align: center;
	}

	.filters-container {
		background: rgba(255, 255, 255, 0.05);
		border-radius: 1rem;
		padding: 2rem;
		border: 1px solid rgba(255, 255, 255, 0.1);
		display: flex;
		flex-wrap: wrap;
		gap: 2rem;
		align-items: center;
		justify-content: center;
	}

	.filter-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		align-items: center;
	}

	.filter-group label {
		font-size: 0.9rem;
		color: var(--light-muted);
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.filter-group select,
	.filter-group input[type='number'] {
		background: rgba(255, 255, 255, 0.1);
		border: 1px solid rgba(255, 255, 255, 0.2);
		border-radius: 0.5rem;
		padding: 0.5rem;
		color: var(--light-muted);
		font-size: 0.9rem;
	}

	.filter-group select:disabled,
	.filter-group input[type='number']:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.range-inputs {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.range-inputs.disabled {
		opacity: 0.5;
	}

	.range-inputs span {
		color: var(--light-muted);
		font-size: 0.9rem;
	}

	.clear-filters {
		background: var(--vibrant);
		color: var(--dark-muted);
		border: none;
		border-radius: 0.5rem;
		padding: 0.5rem 1rem;
		cursor: pointer;
		font-size: 0.9rem;
		transition: all 0.3s ease;
	}

	.clear-filters:hover {
		transform: translateY(-2px);
		opacity: 0.9;
	}

	.active-filters {
		margin-top: 1rem;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 1rem;
		flex-wrap: wrap;
	}

	.filter-badge {
		background: var(--vibrant);
		color: var(--dark-muted);
		padding: 0.3rem 0.8rem;
		border-radius: 1rem;
		font-size: 0.8rem;
		font-weight: bold;
	}

	.results-count {
		color: var(--light-muted);
		font-size: 0.9rem;
		opacity: 0.8;
	}

	/* Métriques principales */
	.metrics-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: 1.5rem;
		margin-bottom: 3rem;
	}

	.metric-card {
		background: rgba(255, 255, 255, 0.05);
		border-radius: 1rem;
		padding: 1.5rem;
		border: 1px solid rgba(255, 255, 255, 0.1);
		transition: all 0.3s ease;
	}

	.metric-card:hover {
		transform: translateY(-5px);
		border-color: var(--vibrant);
	}

	.metric-icon {
		font-size: 2rem;
		margin-bottom: 0.5rem;
	}

	.metric-value {
		font-size: 2rem;
		font-weight: bold;
		color: var(--vibrant);
		margin-bottom: 0.5rem;
	}

	.metric-label {
		font-size: 0.9rem;
		color: var(--light-muted);
		opacity: 0.8;
	}

	/* Graphique */
	.chart-section {
		margin-bottom: 3rem;
		text-align: left;
	}

	.chart-section h2 {
		font-size: 1.8rem;
		color: var(--light-muted);
		margin-bottom: 1.5rem;
		text-align: center;
	}

	.chart-container {
		background: rgba(255, 255, 255, 0.05);
		border-radius: 1rem;
		padding: 2rem;
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.chart-bar {
		display: flex;
		align-items: center;
		gap: 1rem;
		margin-bottom: 1rem;
	}

	.bar-label {
		width: 150px;
		font-size: 0.9rem;
		color: var(--light-muted);
	}

	.bar-track {
		flex: 1;
		height: 20px;
		background: rgba(255, 255, 255, 0.1);
		border-radius: 10px;
		overflow: hidden;
	}

	.bar-fill {
		height: 100%;
		border-radius: 10px;
		transition: width 0.5s ease;
	}

	.bar-fill.high {
		background: linear-gradient(90deg, #ffd700, #ffed4e);
	}

	.bar-fill.medium {
		background: linear-gradient(90deg, #c0c0c0, #e5e5e5);
	}

	.bar-fill.low {
		background: linear-gradient(90deg, #cd7f32, #daa520);
	}

	.bar-value {
		width: 40px;
		font-weight: bold;
		color: var(--vibrant);
	}

	/* Top 10 */
	.top-section {
		margin-bottom: 3rem;
	}

	.top-section h2 {
		font-size: 1.8rem;
		color: var(--light-muted);
		margin-bottom: 1.5rem;
	}

	.top-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
		gap: 1rem;
	}

	.top-card {
		background: rgba(255, 255, 255, 0.05);
		border-radius: 1rem;
		padding: 1.5rem;
		display: flex;
		align-items: center;
		gap: 1rem;
		position: relative;
		transition: all 0.3s ease;
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.top-card:hover {
		transform: translateY(-3px);
		border-color: var(--vibrant);
	}

	.rank-number {
		font-size: 1.5rem;
		font-weight: bold;
		color: var(--vibrant);
		width: 40px;
	}

	.percentage {
		font-size: 0.9rem;
		color: var(--light-muted);
		opacity: 0.8;
	}

	.years-info {
		font-size: 0.8rem;
		color: var(--light-muted);
		opacity: 0.7;
		margin-top: 0.3rem;
	}

	.locations-info {
		font-size: 0.8rem;
		color: var(--light-muted);
		opacity: 0.6;
		margin-top: 0.2rem;
	}

	.continents-info {
		font-size: 0.8rem;
		color: var(--light-muted);
		opacity: 0.7;
		margin-top: 0.3rem;
		font-weight: 500;
	}

	.clickable {
		cursor: pointer;
		transition: all 0.3s ease;
	}

	.clickable:hover {
		transform: translateY(-5px);
		box-shadow: 0 8px 25px rgba(0, 0, 0, 0.3);
		border-color: var(--vibrant);
	}

	.clickable:active {
		transform: translateY(-2px);
	}

	/* Tous les artistes */
	.all-section h2 {
		font-size: 1.8rem;
		color: var(--light-muted);
		margin-bottom: 1.5rem;
	}

	.stats-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
		gap: 1.5rem;
	}

	.artist-card {
		background: rgba(255, 255, 255, 0.05);
		border-radius: 1rem;
		padding: 1.5rem;
		display: flex;
		align-items: center;
		gap: 1rem;
		position: relative;
		transition: all 0.3s ease;
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.artist-card:hover {
		transform: translateY(-5px);
		box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
		border-color: var(--vibrant);
	}

	.artist-image {
		width: 80px;
		height: 80px;
		border-radius: 50%;
		overflow: hidden;
		flex-shrink: 0;
	}

	.artist-image img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.artist-info {
		flex: 1;
		text-align: left;
	}

	.artist-info h3 {
		font-size: 1.3rem;
		color: var(--light-muted);
		margin: 0 0 0.5rem 0;
		font-weight: 600;
	}

	.concert-count {
		display: flex;
		align-items: baseline;
		gap: 0.5rem;
		margin-bottom: 0.3rem;
	}

	.count {
		font-size: 2rem;
		font-weight: bold;
		color: var(--vibrant);
	}

	.label {
		font-size: 1rem;
		color: var(--light-muted);
		opacity: 0.8;
	}

	.percentage-small {
		font-size: 0.8rem;
		color: var(--light-muted);
		opacity: 0.6;
	}

	.rank-badge {
		position: absolute;
		top: 1rem;
		right: 1rem;
		background: var(--vibrant);
		color: var(--dark-muted);
		border-radius: 50%;
		width: 30px;
		height: 30px;
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: bold;
		font-size: 0.9rem;
	}

	/* Couleurs spéciales pour les top 3 */
	.top-card:nth-child(1) {
		border-color: #ffd700;
		background: linear-gradient(135deg, rgba(255, 215, 0, 0.1), rgba(255, 215, 0, 0.05));
	}

	.top-card:nth-child(2) {
		border-color: #c0c0c0;
		background: linear-gradient(135deg, rgba(192, 192, 192, 0.1), rgba(192, 192, 192, 0.05));
	}

	.top-card:nth-child(3) {
		border-color: #cd7f32;
		background: linear-gradient(135deg, rgba(205, 127, 50, 0.1), rgba(205, 127, 50, 0.05));
	}

	@media (max-width: 768px) {
		.filters-container {
			flex-direction: column;
			align-items: stretch;
		}

		.metrics-grid {
			grid-template-columns: repeat(2, 1fr);
		}

		.top-grid {
			grid-template-columns: 1fr;
		}

		.stats-grid {
			grid-template-columns: 1fr;
		}

		.artist-card {
			padding: 1rem;
		}

		.artist-image {
			width: 60px;
			height: 60px;
		}

		.artist-info h3 {
			font-size: 1.1rem;
		}

		.count {
			font-size: 1.5rem;
		}
	}
</style>
