<script lang="ts">
	import { fade, scale } from 'svelte/transition';
	import type { Artist } from '@/types';
	import { expoOut } from 'svelte/easing';
	import { config } from '@/config/config';

	// Import des polices
	const link1 = document.createElement('link');
	link1.rel = 'stylesheet';
	link1.href = 'https://fonts.googleapis.com/css2?family=Jost:wght@400;500;600&display=swap';
	document.head.appendChild(link1);

	// Import local de la police Hanson
	const hansonStyle = document.createElement('style');
	hansonStyle.textContent = `
		@font-face {
			font-family: 'Hanson';
			src: url('/font/Hanson-Bold.ttf') format('truetype');
			font-weight: bold;
			font-style: normal;
		}
	`;
	document.head.appendChild(hansonStyle);

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
			const response = await fetch(`${config.apiBaseUrl}/artists/${artist.id}`);
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
				width="40px"
				height="40px"
				viewBox="0 0 24 24"
				fill="none"
				xmlns="http://www.w3.org/2000/svg"
			>
				<path
					d="M12 2C6.49 2 2 6.49 2 12C2 17.51 6.49 22 12 22C17.51 22 22 17.51 22 12C22 6.49 17.51 2 12 2ZM15.36 14.3C15.65 14.59 15.65 15.07 15.36 15.36C15.21 15.51 15.02 15.58 14.83 15.58C14.64 15.58 14.45 15.51 14.3 15.36L12 13.06L9.7 15.36C9.55 15.51 9.36 15.58 9.17 15.58C8.98 15.58 8.79 15.51 8.64 15.36C8.35 15.07 8.35 14.59 8.64 14.3L10.94 12L8.64 9.7C8.35 9.41 8.35 8.93 8.64 8.64C8.93 8.35 9.41 8.35 9.7 8.64L12 10.94L14.3 8.64C14.59 8.35 15.07 8.35 15.36 8.64C15.65 8.93 15.65 9.41 15.36 9.7L13.06 12L15.36 14.3Z"
					fill="currentColor"
				/>
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
		background: var(--dark-vibrant);
		border-radius: 80px;
		padding: 1rem;
		width: 100%;
		height: 100%;
		box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
		display: flex;
		flex-direction: column;
	}

	.artist-card::before {
		content: '';
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background-image: url('/lines1.svg');
		background-repeat: no-repeat;
		background-position: 90% 85%;
		background-size: 180%;
		opacity: 0.15;
		filter: brightness(0) saturate(100%) invert(1);
		pointer-events: none;
		z-index: -1;
	}

	.close-button {
		position: absolute;
		width: 5rem;
		height: 5rem;
		top: 2rem;
		right: 2rem;
		background: transparent;
		color: var(--light-vibrant);
		border: none;
		outline: none;
		padding: 1.2rem;
		font-size: 1rem;
		cursor: pointer;
		transition: opacity 0.3s ease;
		z-index: 10;
	}

	.close-button:hover {
		opacity: 0.7;
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
		filter: brightness(1.1) contrast(0.8);
	}

	.artist-info h1 {
		margin: 0 0 1rem 0;
		font-size: 3rem;
		color: var(--light-vibrant);
		text-align: left;
		font-family: 'Bukhari Script', cursive;
		line-height: 1;
	}

	.artist-content {
		position: relative;
		display: flex;
		flex-direction: column;
		padding: 0 2rem;
		height: 100%;
		justify-content: space-between;
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
		font-family: 'Hanson', sans-serif;
		font-weight: bold;
	}

	.view-concerts-btn {
		background: var(--dark-vibrant);
		color: white;
		border: none;
		padding: 0.5rem 1.5rem;
		border-radius: 40px;
		font-size: 1rem;
		font-family: 'Hanson', sans-serif;
		font-weight: bold;
		cursor: pointer;
		transition: background 0.2s ease;
	}

	.view-concerts-btn:hover {
		background: var(--dark-muted);
	}

	.members-section h2 {
		margin: 0 0 0.9rem 0;
		color: var(--light-vibrant);
		font-size: 1.2rem;
		font-family: 'Hanson', sans-serif;
		font-weight: bold;
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
		display: flex;
		gap: 0;
		margin-top: auto;
		margin-bottom: 0;
		justify-content: center;
		width: 100%;
		flex: 1;
	}

	.stat-card {
		background: var(--dark-muted);
		padding: 1.5rem 1rem;
		border-radius: 0;
		text-align: center;
		transition: all 0.3s ease;
		flex: 1;
		height: 100%;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		border: 2px solid color-mix(in srgb, var(--light-vibrant), transparent 70%);
		border-right: none;
	}

	.stat-card:first-child {
		border-radius: 0.5rem 0 0 0.5rem;
	}

	.stat-card:last-child {
		border-radius: 0 0.5rem 0.5rem 0;
		border-right: 2px solid color-mix(in srgb, var(--light-vibrant), transparent 70%);
	}

	.stat-card h3 {
		margin: 0 0 0.5rem 0;
		color: color-mix(in srgb, var(--light-vibrant), transparent 30%);
		font-size: 0.8rem;
		font-family: 'Jost', sans-serif;
		font-weight: 400;
		line-height: 1.2;
	}

	.stat-card p {
		margin: 0;
		color: var(--light-vibrant);
		font-size: 1.2rem;
		font-weight: bold;
		line-height: 1.2;
		font-family: 'Hanson', sans-serif;
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
			width: 4.5rem;
			height: 4.5rem;
			top: 2rem;
			right: 2rem;
			padding: 1rem;
			background: transparent;
		}
	}
</style>
