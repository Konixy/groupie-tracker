<script lang="ts">
	import { fade, scale } from 'svelte/transition';
	import type { Artist } from '@/types';
	import { expoOut } from 'svelte/easing';

	// Import des polices
	const link1 = document.createElement('link');
	link1.rel = 'stylesheet';
	link1.href = 'https://fonts.googleapis.com/css2?family=Russo+One&display=swap';
	document.head.appendChild(link1);

	const link2 = document.createElement('link');
	link2.rel = 'stylesheet';
	link2.href =
		'https://fonts.googleapis.com/css2?family=Montserrat:wght@400;500;600;700&display=swap';
	document.head.appendChild(link2);

	const link3 = document.createElement('link');
	link3.rel = 'stylesheet';
	link3.href = 'https://fonts.googleapis.com/css2?family=Jost:wght@400;500;600&display=swap';
	document.head.appendChild(link3);

	let { artist = $bindable(), onClose }: { artist: Artist; onClose: () => void } = $props();

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

	function handleViewConcerts() {
		onClose();

		setTimeout(() => {
			window.scrollTo({
				top: window.scrollY + 720,
				behavior: 'smooth'
			});
		}, 100);
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
				<div class="members-section">
					<h2>Membres du groupe</h2>
					<div class="members-list">
						{#each artist.members as member}
							<span class="member-tag">{member}</span>
						{/each}
					</div>
				</div>
			</div>
		</div>

		<div class="artist-content">
			<div class="concerts-banners-container">
				<div class="concerts-info-banner">
					<span>Nombre de concerts recensés par l'API : {Object.keys(concerts).length}</span>
				</div>
				<button class="view-concerts-btn" onclick={handleViewConcerts}>Voir les concerts</button>
			</div>

			<div class="stats-grid">
				<div class="stat-card">
					<h3>Date de création</h3>
					<p>{artist.creationDate}</p>
				</div>
				<div class="stat-card">
					<h3>Premier album</h3>
					<p>{formatDate(artist.firstAlbum)}</p>
				</div>
				<div class="stat-card">
					<h3>Nombre de membres</h3>
					<p>{artist.members.length}</p>
				</div>
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
		height: 95vh;
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
		border-radius: 80px;
		padding: 1rem;
		width: 100%;
		height: 100%;
		box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
		display: flex;
		flex-direction: column;
		background-image: url('/line1.svg');
		background-repeat: no-repeat;
		background-position: center;
		background-size: 100%;
		background-blend-mode: overlay;
	}

	.close-button {
		position: absolute;
		width: 3rem;
		height: 3rem;
		bottom: 0.5rem;
		left: 50%;
		transform: translateX(-50%);
		background: var(--light-muted);
		color: var(--dark-vibrant);
		border-radius: 50%;
		border: none;
		outline: none;
		padding: 0.8rem;
		font-size: 1rem;
		cursor: pointer;
		transition: background 0.3s ease;
		z-index: 10;
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
		width: 300px;
		height: 300px;
		object-fit: cover;
		border-radius: 60px;
	}

	.artist-info h1 {
		margin: 0 0 1rem 0;
		font-size: 3rem;
		color: var(--dark-vibrant);
		text-align: left;
		font-family: 'Russo One', sans-serif;
		line-height: 1;
	}

	.artist-content {
		position: relative;
		display: flex;
		flex-direction: column;
		padding: 0 2rem;
		height: 100%;
		justify-content: flex-start;
	}

	.concerts-banners-container {
		display: flex;
		gap: 1rem;
		align-items: center;
		justify-content: center;
		margin-top: -0.8rem;
		margin-bottom: 2.5rem;
	}

	.concerts-info-banner {
		background: var(--dark-vibrant);
		color: white;
		padding: 0.65rem 5rem;
		border-radius: 40px;
		font-size: 0.9rem;
		font-weight: 500;
	}

	.concerts-info-banner span {
		font-family: 'Montserrat', sans-serif;
		font-weight: 700;
	}

	.view-concerts-btn {
		background: var(--dark-vibrant);
		color: white;
		border: none;
		padding: 0.5rem 1.5rem;
		border-radius: 40px;
		font-size: 1rem;
		font-family: 'Jost', sans-serif;
		font-weight: 500;
		cursor: pointer;
		transition: background 0.2s ease;
	}

	.view-concerts-btn:hover {
		background: var(--dark-muted);
	}

	.members-section h2 {
		margin: 0 0 0.9rem 0;
		color: var(--dark-vibrant);
		font-size: 1.2rem;
		font-family: 'Montserrat', sans-serif;
		font-weight: 700;
	}

	.members-list {
		display: flex;
		flex-wrap: wrap;
		margin-bottom: 2rem;
		gap: 0.5rem;
	}

	.member-tag {
		background: var(--light-muted);
		padding: 0.5rem 1rem;
		border-radius: 20px;
		font-size: 0.9rem;
		color: white;
		transition: background 0.2s ease;
	}

	.member-tag:hover {
		background: var(--muted);
	}

	.stats-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
		gap: 1.5rem;
		margin-top: -0.8rem;
		margin-bottom: 1rem;
		justify-items: center;
	}

	.stat-card {
		background: var(--light-muted);
		padding: 2.5rem 1.5rem;
		border-radius: 50%;
		text-align: center;
		transition: all 0.3s ease;
		width: 160px;
		height: 160px;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
	}

	.stat-card h3 {
		margin: 0 0 0.5rem 0;
		color: var(--dark-vibrant);
		font-size: 0.9rem;
		font-family: 'Jost', sans-serif;
		font-weight: 600;
		line-height: 1.2;
	}

	.stat-card p {
		margin: 0;
		color: white;
		font-size: 1.1rem;
		font-weight: bold;
		line-height: 1.2;
		font-family: 'Montserrat', sans-serif;
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
			grid-template-columns: repeat(2, 1fr);
			gap: 1.5rem;
		}

		.stat-card {
			width: 130px;
			height: 130px;
			padding: 2rem 1rem;
		}

		.stat-card h3 {
			font-size: 0.8rem;
		}

		.stat-card p {
			font-size: 1rem;
		}

		.close-button {
			width: 2.5rem;
			height: 2.5rem;
			bottom: 1.5rem;
			padding: 0.5rem;
		}

		.concerts-banner {
			flex-direction: column;
			gap: 1rem;
			text-align: center;
			padding: 1rem;
		}
	}
</style>
