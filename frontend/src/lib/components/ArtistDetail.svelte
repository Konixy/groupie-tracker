<script lang="ts">
	import { fade, scale } from 'svelte/transition';
	import type { Artist } from '@/types';
	import { expoOut } from 'svelte/easing';

	let {
		artist = $bindable(),
		onClose,
		highlightedMember = null,
		highlightedDate = null,
		highlightedAlbum = null
	}: {
		artist: Artist;
		onClose: () => void;
		highlightedMember?: string | null;
		highlightedDate?: string | null;
		highlightedAlbum?: string | null;
	} = $props();

	let concerts = $state<Record<string, [string, string]>>({});
	let loadingConcerts = $state(false);

	$effect(() => {
		if (artist) {
			fetchConcerts();
		}
	});

	$effect(() => {
		const handleKeyDown = (e: KeyboardEvent) => {
			if (e.key === 'Escape') {
				onClose();
			}
		};

		document.addEventListener('keydown', handleKeyDown);
		return () => document.removeEventListener('keydown', handleKeyDown);
	});

	async function fetchConcerts() {
		try {
			loadingConcerts = true;
			const response = await fetch(`http://localhost:8080/artists/${artist.id}`);
			if (response.ok) {
				const data = await response.json();
				concerts = data.concerts || {};
			}
		} catch (err) {
			console.error('Erreur lors de la récupération des concerts:', err);
		} finally {
			loadingConcerts = false;
		}
	}

	function formatDate(dateStr: string) {
		try {
			const d = dateStr.split('-');
			const date = new Date(Number(d[2]), Number(d[1]) - 1, Number(d[0]));
			return date.toLocaleDateString('fr-FR', {
				year: 'numeric',
				month: 'short',
				day: 'numeric'
			});
		} catch {
			return dateStr;
		}
	}

	function formatCountryName(country: string) {
		return country.replaceAll('_', ' ');
	}
</script>

<div
	class="background-overlay"
	aria-hidden="true"
	transition:fade={{ duration: 400, easing: expoOut }}
	onclick={onClose}
></div>

<div class="artist-container">
	<div
		class="artist-card"
		transition:scale={{ duration: 400, easing: expoOut, start: 0.5 }}
		role="dialog"
		aria-modal="true"
		aria-labelledby="artist-name"
		tabindex="0"
	>
		<button class="close-button" onclick={onClose} aria-label="Fermer">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				width="16"
				height="16"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="3"
				stroke-linecap="round"
				stroke-linejoin="round"
				aria-hidden="true"
			>
				<path d="M18 6 6 18" />
				<path d="m6 6 12 12" />
			</svg>
		</button>

		<div class="artist-header">
			<img src={artist.image} alt="Photo de {artist.name}" class="artist-image" />
			<div class="artist-info">
				<h1 id="artist-name">{artist.name}</h1>
				<p><strong>Date de création :</strong> {artist.creationDate}</p>
				<p><strong>Premier album :</strong> {formatDate(artist.firstAlbum)}</p>
			</div>
		</div>

		<div class="artist-content">
			<div class="sticky-gradient"></div>

			<div class="members-section">
				<h2>Membres du groupe</h2>
				<div class="members-list">
					{#each artist.members as member}
						<span class="member-tag" class:highlighted={highlightedMember === member}>
							{member}
						</span>
					{/each}
				</div>
			</div>

			<div class="stats-grid">
				<div class="stat-card">
					<h3>Date de création</h3>
					<p class:highlighted={highlightedDate === artist.creationDate.toString()}>
						{artist.creationDate}
					</p>
				</div>
				<div class="stat-card">
					<h3>Premier album</h3>
					<p class:highlighted={highlightedAlbum === artist.firstAlbum}>
						{formatDate(artist.firstAlbum)}
					</p>
				</div>
				<div class="stat-card">
					<h3>Nombre de membres</h3>
					<p>{artist.members.length}</p>
				</div>
			</div>

			<div class="concerts-section">
				<h2>Concerts</h2>
				{#if loadingConcerts}
					<div class="loading" role="status" aria-live="polite">Chargement des concerts...</div>
				{:else if Object.keys(concerts).length > 0}
					<div class="concerts-list">
						{#each Object.entries(concerts) as [date, locations]}
							<div class="concert-item">
								<div class="concert-date">{formatDate(date)}</div>
								<div class="concert-locations">
									{#each locations.slice(0, 2) as location}
										<div class="location-item">📍 {formatCountryName(location)}</div>
									{/each}
									{#if locations.length > 2}
										<div class="location-item">... et {locations.length - 2} autres lieux</div>
									{/if}
								</div>
							</div>
						{/each}
					</div>
				{:else}
					<div class="no-concerts">Aucun concert trouvé</div>
				{/if}
			</div>
		</div>
	</div>
</div>

<style>
	.background-overlay {
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-color: rgba(0, 0, 0, 0.5);
		z-index: 999;
	}

	.artist-container {
		position: fixed;
		top: 50%;
		transform: translateY(-50%);
		width: 60vw;
		height: 90vh;
		display: flex;
		justify-content: center;
		align-items: center;
		z-index: 1000;
	}

	@media (max-width: 1000px) {
		.artist-container {
			width: 90vw;
			height: 90vh;
		}
	}

	.artist-card {
		position: relative;
		z-index: 1000;
		background: var(--light-vibrant);
		border-radius: 20px;
		padding: 1rem;
		width: 100%;
		height: 100%;
		box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
		display: flex;
		flex-direction: column;
	}

	.close-button {
		position: absolute;
		width: 2rem;
		height: 2rem;
		top: 2rem;
		right: 2rem;
		background: var(--light-muted);
		color: var(--dark-vibrant);
		border-radius: 50%;
		border: none;
		outline: none;
		padding: 0.5rem;
		font-size: 0.8rem;
		cursor: pointer;
		transition: background 0.3s ease;
	}

	.close-button:hover {
		background: var(--muted);
	}

	.artist-header {
		display: flex;
		align-items: center;
		gap: 2rem;
		padding: 2rem;
	}

	.artist-image {
		width: 200px;
		height: 200px;
		object-fit: cover;
		border-radius: 16px;
		box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
	}

	.artist-info h1 {
		margin: 0 0 1rem 0;
		font-size: 2.5rem;
		color: var(--dark-vibrant);
	}

	.artist-info p {
		margin: 0.5rem 0;
		font-size: 1.1rem;
		color: var(--dark-muted);
	}

	.artist-content {
		position: relative;
		display: flex;
		flex-direction: column;
		overflow-y: auto;
		padding: 0 2rem;
	}

	.sticky-gradient {
		position: sticky;
		top: 0;
		background: linear-gradient(to bottom, var(--light-vibrant), transparent);
		min-height: 4rem;
		margin: 0 -3rem;
		padding: 0 3rem;
		width: calc(100% + 6rem - 20px);
	}

	.members-section h2 {
		margin-top: -1rem;
		color: var(--dark-vibrant);
	}

	.members-list {
		display: flex;
		flex-wrap: wrap;
		margin-bottom: 1rem;
		gap: 0.5rem;
	}

	.member-tag {
		background: var(--light-muted);
		padding: 0.5rem 1rem;
		border-radius: 20px;
		font-size: 0.9rem;
		color: var(--dark-vibrant);
		transition: background 0.2s ease;
	}

	.member-tag:hover {
		background: var(--muted);
	}

	.member-tag.highlighted {
		background: var(--vibrant);
		color: var(--dark-muted);
		animation: blink 0.8s ease-in-out infinite;
		box-shadow: 0 0 20px var(--vibrant);
	}

	@keyframes blink {
		0%,
		100% {
			opacity: 1;
			transform: scale(1);
		}
		50% {
			opacity: 0.7;
			transform: scale(1.05);
		}
	}

	.stat-card p.highlighted {
		background: var(--vibrant);
		color: var(--light-vibrant);
		animation: blink 0.8s ease-in-out infinite;
		box-shadow: 0 0 20px var(--vibrant);
		border-radius: 8px;
		padding: 0.5rem;
		margin: 0.5rem 0;
	}

	.stats-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: 1.5rem;
		margin-bottom: 2rem;
	}

	.stat-card {
		background: var(--light-muted);
		padding: 1.5rem;
		border-radius: 12px;
		text-align: center;
		transition: transform 0.2s ease;
	}

	.stat-card:hover {
		transform: translateY(-2px);
	}

	.stat-card h3 {
		margin: 0 0 0.5rem 0;
		color: var(--dark-vibrant);
		font-size: 1.1rem;
	}

	.stat-card p {
		margin: 0;
		color: var(--dark-muted);
		font-size: 1.2rem;
		font-weight: bold;
	}

	.concerts-section {
		display: flex;
		flex-direction: column;
		margin-bottom: 2rem;
	}

	.concerts-section h2 {
		color: var(--dark-vibrant);
		margin-bottom: 1rem;
	}

	.concerts-list {
		display: grid;
		gap: 1rem;
	}

	.concert-item {
		background: #fff3e0;
		border-radius: 12px;
		padding: 1rem;
		border-left: 4px solid #ff9800;
		opacity: 0.9;
		transition: transform 0.2s ease;
	}

	.concert-item:hover {
		transform: translateY(-2px);
	}

	.concert-date {
		font-weight: bold;
		color: var(--dark-vibrant);
		font-size: 1rem;
		margin-bottom: 0.5rem;
	}

	.concert-locations {
		color: var(--dark-muted);
	}

	.location-item {
		margin: 0.2rem 0;
		padding-left: 0.5rem;
		border-left: 2px solid var(--light-muted);
		font-size: 0.9rem;
		text-transform: capitalize;
	}

	.loading {
		text-align: center;
		color: var(--dark-muted);
		font-style: italic;
		margin: 1rem 0;
	}

	.no-concerts {
		text-align: center;
		color: var(--dark-muted);
		font-style: italic;
		margin: 1rem 0;
	}

	@media (max-width: 768px) {
		.artist-header {
			flex-direction: column;
			text-align: center;
			gap: 1rem;
		}

		.artist-image {
			width: 150px;
			height: 150px;
		}

		.artist-info h1 {
			font-size: 2rem;
		}

		.stats-grid {
			grid-template-columns: 1fr;
		}
	}
</style>
