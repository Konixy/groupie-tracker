<script lang="ts">
  let { artists = [], searchResults = [] } = $props();
  let current = $state(0);

  // Écouter les événements de centrage sur un artiste
  $effect(() => {
    const handleCenterOnArtist = (event: any) => {
      current = event.detail.index;
    };
    
    document.addEventListener('centerOnArtist', handleCenterOnArtist);
    
    return () => {
      document.removeEventListener('centerOnArtist', handleCenterOnArtist);
    };
  });

  // Utiliser les résultats de recherche s'ils existent, sinon tous les artistes
  let displayArtists = $derived(searchResults.length > 0 ? searchResults : artists);

  function prev() {
    current = (current - 1 + displayArtists.length) % displayArtists.length;
  }
  function next() {
    current = (current + 1) % displayArtists.length;
  }

  function handleCardClick(artistIndex: number) {
    // Émettre un événement pour ouvrir les détails de l'artiste
    const event = new CustomEvent('artistClick', {
      detail: displayArtists[artistIndex]
    });
    document.dispatchEvent(event);
  }

  function handleViewMore() {
    // Émettre un événement pour ouvrir les détails de l'artiste central
    const event = new CustomEvent('artistClick', {
      detail: displayArtists[current]
    });
    document.dispatchEvent(event);
  }

  function handleKeyDown(event: KeyboardEvent, artistIndex: number) {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      handleCardClick(artistIndex);
    }
  }
</script>

<style>
  .carousel-container {
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 2rem 0;
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
    z-index: 10;
  }
  .arrow:hover {
    opacity: 1;
    background: #fff6;
  }
  .cards {
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    width: 70vw;
    height: 60vh;
  }
  .clickable-card img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-radius: 16px;
    box-shadow: 0 4px 24px #0005;
  }
  .clickable-card h2 {
    margin: 0;
    font-size: 1.5rem;
    color: #222;
    background: rgba(255,255,255,0.7);
    border-radius: 8px;
    padding: 0.3em 1em;
    position: absolute;
    bottom: 18px;
    left: 50%;
    transform: translateX(-50%);
    z-index: 3;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 90%;
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
    transform: translate(-50%, -50%);
    transition: transform 0.5s cubic-bezier(.4,2,.6,1), box-shadow 0.5s, opacity 0.5s, filter 0.5s;
    z-index: 2;
    width: 400px;
    height: 400px;
    opacity: 1;
    cursor: pointer;
    border: none;
    padding: 0;
  }
  .clickable-card:hover {
    box-shadow: 0 12px 40px #0008;
    transform: translate(-50%, -50%) scale(1.02);
  }
  .clickable-card:focus {
    outline: 2px solid #667eea;
    outline-offset: 2px;
  }
  .clickable-card.left {
    transform: translate(-50%, -50%) translateX(-220px);
    z-index: 2;
    width: 300px;
    height: 300px;
  }
  .clickable-card.center {
    transform: translate(-50%, -50%);
    z-index: 3;
    width: clamp(250px, 30vw, 400px);
    height: clamp(250px, 30vw, 400px);
  }
  .clickable-card.right {
    transform: translate(-50%, -50%) translateX(220px);
    z-index: 2;
    width: 300px;
    height: 300px;
  }
  .clickable-card.leftleft {
    transform: translate(-50%, -50%) translateX(-180px);
    z-index: 1;
    width: 250px;
    height: 250px;
  }
  .clickable-card.rightright {
    transform: translate(-50%, -50%) translateX(180px);
    z-index: 1;
    width: 250px;
    height: 250px;
  }
  .view-more-button {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: none;
    padding: 1rem 2rem;
    border-radius: 25px;
    font-size: 1.1rem;
    font-weight: bold;
    cursor: pointer;
    margin-top: 2rem;
    transition: all 0.3s ease;
    box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
    display: block;
    margin: 2rem auto 0 auto;
  }
  .view-more-button:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.3);
  }
  .view-more-button:active {
    transform: translateY(0);
  }
</style>

{#if !displayArtists.length}
  <div style="color: #fff; text-align: center; margin: 2rem;">Aucun artiste à afficher.</div>
{/if}

<div class="carousel-container">
  <button class="arrow" onclick={prev}>&larr;</button>
  <div class="cards">
    {#if displayArtists.length === 1}
      <!-- Si un seul artiste, afficher seulement la carte centrale -->
      <button class="clickable-card center" onclick={() => handleCardClick(0)} onkeydown={(event) => handleKeyDown(event, 0)} aria-label="Artiste actuel">
        <img src={displayArtists[0].image} alt={displayArtists[0].name} />
        <h2>{displayArtists[0].name}</h2>
      </button>
    {:else if displayArtists.length === 2}
      <!-- Si deux artistes, afficher seulement les cartes left et right -->
      <button class="clickable-card left" onclick={() => handleCardClick(0)} onkeydown={(event) => handleKeyDown(event, 0)} aria-label="Artiste précédent">
        <img src={displayArtists[0].image} alt="artist" />
        <h2>{displayArtists[0].name}</h2>
      </button>
      <button class="clickable-card right" onclick={() => handleCardClick(1)} onkeydown={(event) => handleKeyDown(event, 1)} aria-label="Artiste suivant">
        <img src={displayArtists[1].image} alt="artist" />
        <h2>{displayArtists[1].name}</h2>
      </button>
    {:else if displayArtists.length >= 3}
      <!-- Si 3 artistes ou plus, afficher le carousel complet -->
      <!-- Left++ side image (hidden behind, blurred) -->
      <button class="clickable-card leftleft" onclick={() => handleCardClick((current - 2 + displayArtists.length) % displayArtists.length)} onkeydown={(event) => handleKeyDown(event, (current - 2 + displayArtists.length) % displayArtists.length)} aria-label="Artiste précédent">
        <img src={displayArtists[(current - 2 + displayArtists.length) % displayArtists.length].image} alt="artist" />
        <h2>{displayArtists[(current - 2 + displayArtists.length) % displayArtists.length].name}</h2>
      </button>
      <!-- Left side image (hidden behind, blurred) -->
      <button class="clickable-card left" onclick={() => handleCardClick((current - 1 + displayArtists.length) % displayArtists.length)} onkeydown={(event) => handleKeyDown(event, (current - 1 + displayArtists.length) % displayArtists.length)} aria-label="Artiste précédent">
        <img src={displayArtists[(current - 1 + displayArtists.length) % displayArtists.length].image} alt="artist" />
        <h2>{displayArtists[(current - 1 + displayArtists.length) % displayArtists.length].name}</h2>
      </button>
      <!-- Center image (large, always present) -->
      <button class="clickable-card center" onclick={() => handleCardClick(current)} onkeydown={(event) => handleKeyDown(event, current)} aria-label="Artiste actuel">
        <img src={displayArtists[current].image} alt={displayArtists[current].name} />
        <h2>{displayArtists[current].name}</h2>
      </button>
      <!-- Right side image (hidden behind, blurred) -->
      <button class="clickable-card right" onclick={() => handleCardClick((current + 1) % displayArtists.length)} onkeydown={(event) => handleKeyDown(event, (current + 1) % displayArtists.length)} aria-label="Artiste suivant">
        <img src={displayArtists[(current + 1) % displayArtists.length].image} alt="artist" />
        <h2>{displayArtists[(current + 1) % displayArtists.length].name}</h2>
      </button>
      <!-- Right++ side image (hidden behind, blurred) -->
      <button class="clickable-card rightright" onclick={() => handleCardClick((current + 2) % displayArtists.length)} onkeydown={(event) => handleKeyDown(event, (current + 2) % displayArtists.length)} aria-label="Artiste suivant">
        <img src={displayArtists[(current + 2) % displayArtists.length].image} alt="artist" />
        <h2>{displayArtists[(current + 2) % displayArtists.length].name}</h2>
      </button>
    {/if}
  </div>
  <button class="arrow" onclick={next}>&rarr;</button>
</div>

{#if displayArtists.length}
  <button class="view-more-button" onclick={handleViewMore}>
    Voir plus
  </button>
{/if}
