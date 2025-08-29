<script lang="ts">
	import type { Artist } from '@/types';
	import { Vibrant } from 'node-vibrant/browser';

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

	let {
		artists = [],
		selectedArtist = $bindable(),
		loading = $bindable(),
		current = $bindable(),
		highlightedMember = null,
		highlightedDate = null,
		highlightedAlbum = null
	}: {
		artists: Artist[];
		selectedArtist: Artist | null;
		loading: boolean;
		current: number;
		highlightedMember?: string | null;
		highlightedDate?: string | null;
		highlightedAlbum?: string | null;
	} = $props();
	// let current = $state(Math.floor(Math.random() * 52));
	// let current = $state(0);
	let isAnimating = $state(false);

	$effect(() => {
		if (selectedArtist) {
			document.body.style.overflow = 'hidden';
		} else {
			document.body.style.overflow = '';
		}
	});

	$effect(() => {
		current;
		if (artists.length > 0) {
			(async () => {
				const img = document.querySelector(`.img-center`) as HTMLImageElement | null;

				if (img) {
					img.crossOrigin = 'anonymous';
					await img.decode();
					const v = await new Vibrant(img).getPalette();
					const vibrant = v.Vibrant?.rgb || [255, 255, 255];
					const lightVibrant = v.LightVibrant?.rgb || [255, 255, 255];
					const darkVibrant = v.DarkVibrant?.rgb || [100, 100, 100];
					const muted = v.Muted?.rgb || [255, 255, 255];
					const lightMuted = v.LightMuted?.rgb || [255, 255, 255];
					const darkMuted = v.DarkMuted?.rgb || [100, 100, 100];

					document.body.style.setProperty(
						'--vibrant',
						`rgb(${vibrant[0]}, ${vibrant[1]}, ${vibrant[2]})`
					);
					document.body.style.setProperty(
						'--light-vibrant',
						`rgb(${lightVibrant[0]}, ${lightVibrant[1]}, ${lightVibrant[2]})`
					);
					document.body.style.setProperty(
						'--dark-vibrant',
						`rgb(${darkVibrant[0]}, ${darkVibrant[1]}, ${darkVibrant[2]})`
					);
					document.body.style.setProperty('--muted', `rgb(${muted[0]}, ${muted[1]}, ${muted[2]})`);
					document.body.style.setProperty(
						'--light-muted',
						`rgb(${lightMuted[0]}, ${lightMuted[1]}, ${lightMuted[2]})`
					);
					document.body.style.setProperty(
						'--dark-muted',
						`rgb(${darkMuted[0]}, ${darkMuted[1]}, ${darkMuted[2]})`
					);

					document.body.style.transition = `background-color 1s ease, color 1s ease, -webkit-text-stroke 1s ease, text-stroke 1s ease`;
					loading = false;
				}
			})();
		}
	});

	$effect(() => {
		const fn = (e: KeyboardEvent) => {
			if (e.key === 'ArrowLeft') {
				prev();
			} else if (e.key === 'ArrowRight') {
				next();
			}
		};
		window.addEventListener('keydown', fn);
		return () => window.removeEventListener('keydown', fn);
	});

	// Calculer la position de chaque carte en fonction de sa position relative au centre
	function getCardPosition(index: number) {
		const total = artists.length;
		if (total === 0)
			return {
				transform: '',
				opacity: 0,
				zIndex: 0,
				width: 0,
				height: 0,
				relativePos: 0,
				fontSize: '1.2rem'
			};

		// Calculer la position relative au centre
		let relativePos = index - current;

		// Gérer la circularité
		if (relativePos > total / 2) relativePos -= total;
		if (relativePos < -total / 2) relativePos += total;

		// Positions pour 1, 2, ou 3+ artistes
		if (total === 1) {
			return {
				transform: 'translate(-50%, -50%) translateZ(0px) rotateY(0deg)',
				opacity: 1,
				zIndex: 4,
				width: 'clamp(250px, 30vw, 400px)',
				height: 'clamp(250px, 30vw, 400px)',
				relativePos,
				fontSize: 'clamp(1.2rem, 4vw, 2.8rem)'
			};
		}

		if (total === 2) {
			if (relativePos === -1) {
				return {
					transform: 'translate(-50%, -50%) translateX(-250px) translateZ(-80px) rotateY(15deg)',
					opacity: 0.7,
					zIndex: 2,
					width: 'clamp(200px, 25vw, 320px)',
					height: 'clamp(200px, 25vw, 320px)',
					relativePos,
					fontSize: 'clamp(1rem, 3vw, 2.2rem)'
				};
			} else {
				return {
					transform: 'translate(-50%, -50%) translateX(250px) translateZ(-80px) rotateY(-15deg)',
					opacity: 0.7,
					zIndex: 2,
					width: 'clamp(200px, 25vw, 320px)',
					height: 'clamp(200px, 25vw, 320px)',
					relativePos,
					fontSize: 'clamp(1rem, 3vw, 2.2rem)'
				};
			}
		}

		// Pour 3+ artistes
		switch (relativePos) {
			case -2:
				return {
					transform: 'translate(-50%, -50%) translateX(-200px) translateZ(-150px) rotateY(20deg)',
					opacity: 0,
					zIndex: 1,
					width: 'clamp(150px, 20vw, 250px)',
					height: 'clamp(150px, 20vw, 250px)',
					relativePos,
					fontSize: 'clamp(0.8rem, 2.5vw, 1.8rem)'
				};
			case -1:
				return {
					transform: 'translate(-50%, -50%) translateX(-250px) translateZ(-80px) rotateY(15deg)',
					opacity: 1,
					zIndex: 2,
					width: 'clamp(200px, 25vw, 320px)',
					height: 'clamp(200px, 25vw, 320px)',
					relativePos,
					fontSize: 'clamp(1rem, 3vw, 2.2rem)'
				};
			case 0:
				return {
					transform: 'translate(-50%, -50%) translateZ(0px) rotateY(0deg)',
					opacity: 1,
					zIndex: 4,
					width: 'clamp(250px, 30vw, 400px)',
					height: 'clamp(250px, 30vw, 400px)',
					relativePos,
					fontSize: 'clamp(1.2rem, 4vw, 2.8rem)'
				};
			case 1:
				return {
					transform: 'translate(-50%, -50%) translateX(250px) translateZ(-80px) rotateY(-15deg)',
					opacity: 1,
					zIndex: 2,
					width: 'clamp(200px, 25vw, 320px)',
					height: 'clamp(200px, 25vw, 320px)',
					relativePos,
					fontSize: 'clamp(1rem, 3vw, 2.2rem)'
				};
			case 2:
				return {
					transform: 'translate(-50%, -50%) translateX(200px) translateZ(-150px) rotateY(-20deg)',
					opacity: 0,
					zIndex: 1,
					width: 'clamp(150px, 20vw, 250px)',
					height: 'clamp(150px, 20vw, 250px)',
					relativePos,
					fontSize: 'clamp(0.8rem, 2.5vw, 1.8rem)'
				};
			default:
				return {
					transform: 'translate(-50%, -50%) translateX(400px) translateZ(-200px) rotateY(-30deg)',
					opacity: 0,
					zIndex: 0,
					width: 'clamp(120px, 15vw, 200px)',
					height: 'clamp(120px, 15vw, 200px)',
					relativePos,
					fontSize: 'clamp(0.6rem, 2vw, 1.5rem)'
				};
		}
	}

	async function prev() {
		if (isAnimating) return;
		isAnimating = true;
		current = (current - 1 + artists.length) % artists.length;
		if (selectedArtist) selectedArtist = artists[current];
		await new Promise((resolve) => setTimeout(resolve, 300));
		isAnimating = false;
	}

	async function next() {
		if (isAnimating) return;
		isAnimating = true;
		current = (current + 1) % artists.length;
		if (selectedArtist) selectedArtist = artists[current];
		await new Promise((resolve) => setTimeout(resolve, 300));
		isAnimating = false;
	}

	function handleCardClick(artistIndex: number, relativePos: number) {
		if (relativePos === 1) {
			next();
		} else if (relativePos === -1) {
			prev();
		} else {
			selectedArtist = artists[artistIndex];
		}
	}

	function isInRange(index: number) {
		const range = 3;
		const distance = Math.abs(index - current);
		const wrapDistance = artists.length - distance;
		return Math.min(distance, wrapDistance) < range;
	}
</script>

{#if artists.length}
	<div class="carousel-container">
		<button class="arrow" onclick={prev} disabled={isAnimating}>&larr;</button>
		<div class="cards">
			{#each artists as artist, index}
				{@const position = getCardPosition(index)}
				<button
					class="clickable-card"
					style="
          transform: {position.transform};
          opacity: {position.opacity};
          z-index: {position.zIndex};
          width: {position.width};
          height: {position.height};
        "
					onclick={() => handleCardClick(index, position.relativePos)}
					aria-label="Artiste {index + 1}"
				>
					<!-- Optimization -->
					<!-- Slow 4g: before: 51.20s, after: 9.74s -->

					{#if isInRange(index)}
						<img
							src={artist.customImage}
							alt={artist.name}
							class="artist-img-{index} {position.relativePos === 0
								? 'img-center'
								: position.relativePos === 1
									? 'img-right'
									: 'img-left'}"
						/>
						<h2 style="font-size: {position.fontSize}">{artist.name}</h2>
					{/if}
				</button>
			{/each}
		</div>
		<button class="arrow" onclick={next} disabled={isAnimating}>&rarr;</button>
	</div>

	<button class="view-more-button" onclick={() => (selectedArtist = artists[current])}>
		Découvrir
	</button>
{:else}
	<div class="error-message">
		Une erreur est survenue lors de la récupération des artistes.
		<br />
		<button class="retry-button" onclick={() => window.location.reload()}>Réessayer</button>
	</div>
{/if}

<style>
	.carousel-container {
		display: flex;
		align-items: center;
		justify-content: center;
		perspective: 1200px;
	}

	.arrow {
		background: #fff3;
		border: none;
		border-radius: 50%;
		color: #fff;
		font-size: 2rem;
		width: 50px;
		height: 50px;
		cursor: pointer;
		margin: 0 2rem;
		display: flex;
		align-items: center;
		justify-content: center;
		opacity: 0.7;
		transition: opacity 0.2s;
		z-index: 2000;
	}

	.arrow:hover {
		opacity: 1;
		background: #fff6;
	}

	.arrow:disabled {
		opacity: 0.3;
		cursor: not-allowed;
	}

	.cards {
		display: flex;
		align-items: center;
		justify-content: center;
		position: relative;
		width: 70vw;
		height: 60vh;
		transform-style: preserve-3d;
	}

	.clickable-card {
		background: #ddd;
		border-radius: 20px;
		box-shadow: 0 8px 32px #0006;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		position: absolute;
		top: 50%;
		left: 50%;
		transition: all 0.6s ease;
		cursor: pointer;
		border: none;
		padding: 0;
		transform-style: preserve-3d;
		outline: none;
	}

	.clickable-card img {
		width: 100%;
		height: 100%;
		object-fit: cover;
		border-radius: 16px;
		box-shadow: 0 4px 24px #0005;
		transition: transform 0.3s ease;
	}

	.clickable-card h2 {
		margin: 0;
		font-family: 'Russo One', sans-serif;
		font-weight: 400;
		color: white;
		letter-spacing: -0.05em;
		position: absolute;
		bottom: 18px;
		left: 50%;
		transform: translateX(-50%);
		z-index: 3;
		white-space: normal;
		overflow: visible;
		max-width: 90%;
		width: 90%;
		transition: all 0.3s ease;
		text-align: center;
		line-height: 1.1;
	}

	.clickable-card:hover img {
		transform: scale(1.05);
	}

	.view-more-button {
		background: var(--light-vibrant);
		color: var(--dark-vibrant);
		border: none;
		padding: 1rem 2rem;
		border-radius: 20px;
		font-size: 1.1rem;
		font-weight: bold;
		cursor: pointer;
		margin-top: 0.5rem;
		margin-bottom: 2rem;
		transition: all 0.3s ease;
		box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
		display: block;
		font-family: 'Montserrat', sans-serif;
	}

	.view-more-button:hover {
		transform: translateY(-2px);
		box-shadow: 0 6px 20px rgba(0, 0, 0, 0.3);
	}

	.view-more-button:active {
		transform: translateY(0);
	}

	/* Responsive */
	@media (max-width: 768px) {
		.carousel-container {
			perspective: 1000px;
		}
	}

	.error-message {
		text-align: center;
		color: var(--light-vibrant);
		text-align: center;
		margin: 2rem;
		margin-top: 10rem;
		margin-bottom: 10rem;
		z-index: 1000;
	}

	.retry-button {
		background: var(--light-vibrant);
		color: var(--dark-vibrant);
		border: none;
		padding: 1rem 2rem;
		border-radius: 20px;
		cursor: pointer;
		transition: all 0.3s ease;
		font-family: inherit;
		font-size: 1rem;
		box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
		display: block;
		margin: 0 auto;
		margin-top: 2rem;
	}

	.retry-button:hover {
		transform: translateY(-2px);
		box-shadow: 0 6px 20px rgba(0, 0, 0, 0.3);
	}
</style>
