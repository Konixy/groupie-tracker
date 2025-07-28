<script lang="ts">
	let { artist, onBack } = $props();
	let concerts = $state<Record<string, string[]>>({});
	let loading = $state(true);
	let error = $state<string | null>(null);

	$effect(() => {
		if (artist) {
			fetchConcerts();
		}
	});

	async function fetchConcerts() {
		try {
			loading = true;
			error = null;

			const response = await fetch(`http://localhost:8080/artists/${artist.id}`);
			if (!response.ok) {
				throw new Error('Erreur lors de la récupération des concerts');
			}

			const data = await response.json();
			concerts = data.concerts || {};
		} catch (err) {
			error = err instanceof Error ? err.message : 'Erreur inconnue';
		} finally {
			loading = false;
		}
	}

	function formatDate(dateStr: string) {
		try {
			const date = new Date(dateStr);
			return date.toLocaleDateString('fr-FR', {
				weekday: 'long',
				year: 'numeric',
				month: 'long',
				day: 'numeric'
			});
		} catch {
			return dateStr;
		}
	}

	function isPastDate(dateStr: string): boolean {
		try {
			const date = new Date(dateStr);
			const now = new Date();
			return date < now;
		} catch {
			return false;
		}
	}

	function sortConcerts() {
		const entries = Object.entries(concerts);
		return entries.sort(([dateA], [dateB]) => {
			const dateObjA = new Date(dateA);
			const dateObjB = new Date(dateB);
			return dateObjA.getTime() - dateObjB.getTime();
		});
	}

	function getPastConcerts() {
		return sortConcerts().filter(([date]) => isPastDate(date));
	}

	function getUpcomingConcerts() {
		return sortConcerts().filter(([date]) => !isPastDate(date));
	}
</script>

<div class="concerts-container">
	<button class="back-button" onclick={onBack}> ← Retour </button>

	{#if artist}
		<div class="concerts-card">
			<div class="artist-header">
				<img src={artist.image} alt={artist.name} class="artist-image" />
				<div class="artist-info">
					<h1>{artist.name}</h1>
					<p>Consultez les dates et lieux de concerts</p>
				</div>
			</div>

			<div class="concerts-section">
				<h2>Concerts à venir</h2>

				{#if loading}
					<div class="loading">Chargement des concerts...</div>
				{:else if error}
					<div class="error">Erreur: {error}</div>
				{:else if Object.keys(concerts).length === 0}
					<div class="no-concerts">Aucun concert programmé pour le moment</div>
				{:else}
					{#if getUpcomingConcerts().length > 0}
						<div class="concerts-list">
							{#each getUpcomingConcerts() as [date, locations]}
								<div class="concert-item upcoming">
									<div class="concert-date">{formatDate(date)}</div>
									<div class="concert-locations">
										{#each locations as location}
											<div class="location-item">📍 {location}</div>
										{/each}
									</div>
								</div>
							{/each}
						</div>
					{:else}
						<div class="no-concerts">Aucun concert à venir</div>
					{/if}

					{#if getPastConcerts().length > 0}
						<div class="past-concerts-section">
							<h3>Concerts passés</h3>
							<div class="concerts-list past">
								{#each getPastConcerts() as [date, locations]}
									<div class="concert-item past">
										<div class="concert-date">{formatDate(date)}</div>
										<div class="concert-locations">
											{#each locations as location}
												<div class="location-item">📍 {location}</div>
											{/each}
										</div>
									</div>
								{/each}
							</div>
						</div>
					{/if}
				{/if}
			</div>
		</div>
	{/if}
</div>

<style>
	.concerts-container {
		min-height: 100vh;
		background: #4e4e4e;
		padding: 2rem;
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	.back-button {
		position: absolute;
		top: 2rem;
		left: 2rem;
		background: #333;
		color: white;
		border: none;
		padding: 0.8rem 1.5rem;
		border-radius: 8px;
		cursor: pointer;
		font-size: 1rem;
		transition: background 0.3s;
	}

	.back-button:hover {
		background: #555;
	}

	.concerts-card {
		background: white;
		border-radius: 20px;
		padding: 3rem;
		max-width: 900px;
		width: 100%;
		box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
		margin-top: 2rem;
	}

	.artist-header {
		display: flex;
		align-items: center;
		gap: 2rem;
		margin-bottom: 2rem;
		border-bottom: 2px solid #f0f0f0;
		padding-bottom: 2rem;
	}

	.artist-image {
		width: 150px;
		height: 150px;
		object-fit: cover;
		border-radius: 16px;
		box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
	}

	.artist-info h1 {
		margin: 0 0 1rem 0;
		font-size: 2.5rem;
		color: #333;
	}

	.concerts-section h2 {
		color: #333;
		margin-bottom: 1.5rem;
		font-size: 1.8rem;
	}

	.concerts-list {
		display: grid;
		gap: 1rem;
	}

	.concert-item {
		background: #f8f8f8;
		border-radius: 12px;
		padding: 1.5rem;
		border-left: 4px solid #667eea;
	}

	.concert-item.upcoming {
		border-left-color: #4caf50;
		background: #f1f8e9;
	}

	.concert-item.past {
		border-left-color: #ff9800;
		background: #fff3e0;
		opacity: 0.8;
	}

	.concert-date {
		font-weight: bold;
		color: #333;
		font-size: 1.1rem;
		margin-bottom: 0.5rem;
	}

	.concert-locations {
		color: #666;
	}

	.location-item {
		margin: 0.3rem 0;
		padding-left: 1rem;
		border-left: 2px solid #ddd;
	}

	.loading {
		text-align: center;
		color: white;
		font-size: 1.2rem;
		margin: 2rem 0;
	}

	.error {
		text-align: center;
		color: #ff6b6b;
		background: rgba(255, 107, 107, 0.1);
		padding: 1rem;
		border-radius: 8px;
		margin: 2rem 0;
	}

	.no-concerts {
		text-align: center;
		color: #666;
		font-style: italic;
		margin: 2rem 0;
	}

	.past-concerts-section {
		margin-top: 3rem;
		padding-top: 2rem;
		border-top: 2px solid #f0f0f0;
	}

	.past-concerts-section h3 {
		color: #666;
		margin-bottom: 1.5rem;
		font-size: 1.5rem;
	}

	.past-concerts-section .concerts-list {
		gap: 0.8rem;
	}

	.past-concerts-section .concert-item {
		padding: 1rem;
		font-size: 0.95rem;
	}
</style>
