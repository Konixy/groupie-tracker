<script lang="ts">
  let { artists = [], searchResults = [] } = $props();
  let current = $state(0);
  let isAnimating = $state(false);

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

  // Calculer la position de chaque carte en fonction de sa position relative au centre
  function getCardPosition(index: number) {
    const total = displayArtists.length;
    if (total === 0) return { transform: '', opacity: 0, zIndex: 0, width: 0, height: 0 };
    
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
        height: 'clamp(250px, 30vw, 400px)'
      };
    }
    
    if (total === 2) {
      if (relativePos === -1) {
        return {
          transform: 'translate(-50%, -50%) translateX(-250px) translateZ(-80px) rotateY(15deg)',
          opacity: 0.7,
          zIndex: 2,
          width: '320px',
          height: '320px'
        };
      } else {
        return {
          transform: 'translate(-50%, -50%) translateX(250px) translateZ(-80px) rotateY(-15deg)',
          opacity: 0.7,
          zIndex: 2,
          width: '320px',
          height: '320px'
        };
      }
    }
    
    // Pour 3+ artistes
    switch (relativePos) {
      case -2:
        return {
          transform: 'translate(-50%, -50%) translateX(-200px) translateZ(-150px) rotateY(20deg)',
          opacity: 1,
          zIndex: 1,
          width: '250px',
          height: '250px'
        };
      case -1:
        return {
          transform: 'translate(-50%, -50%) translateX(-250px) translateZ(-80px) rotateY(15deg)',
          opacity: 1,
          zIndex: 2,
          width: '320px',
          height: '320px'
        };
      case 0:
        return {
          transform: 'translate(-50%, -50%) translateZ(0px) rotateY(0deg)',
          opacity: 1,
          zIndex: 4,
          width: 'clamp(250px, 30vw, 400px)',
          height: 'clamp(250px, 30vw, 400px)'
        };
      case 1:
        return {
          transform: 'translate(-50%, -50%) translateX(250px) translateZ(-80px) rotateY(-15deg)',
          opacity: 1,
          zIndex: 2,
          width: '320px',
          height: '320px'
        };
      case 2:
        return {
          transform: 'translate(-50%, -50%) translateX(200px) translateZ(-150px) rotateY(-20deg)',
          opacity: 1,
          zIndex: 1,
          width: '250px',
          height: '250px'
        };
      default:
        return {
          transform: 'translate(-50%, -50%) translateX(400px) translateZ(-200px) rotateY(-30deg)',
          opacity: 0,
          zIndex: 0,
          width: '200px',
          height: '200px'
        };
    }
  }

  async function prev() {
    if (isAnimating) return;
    isAnimating = true;
    current = (current - 1 + displayArtists.length) % displayArtists.length;
    await new Promise(resolve => setTimeout(resolve, 600));
    isAnimating = false;
  }

  async function next() {
    if (isAnimating) return;
    isAnimating = true;
    current = (current + 1) % displayArtists.length;
    await new Promise(resolve => setTimeout(resolve, 600));
    isAnimating = false;
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
    z-index: 10;
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
    transition: all 0.6s cubic-bezier(0.25, 0.46, 0.45, 0.94);
    cursor: pointer;
    border: none;
    padding: 0;
    transform-style: preserve-3d;
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
    transition: all 0.3s ease;
  }

  .clickable-card:hover {
    box-shadow: 0 12px 40px #0008;
  }

  .clickable-card:hover img {
    transform: scale(1.05);
  }

  .clickable-card:focus {
    outline: 2px solid #667eea;
    outline-offset: 2px;
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

  /* Responsive */
  @media (max-width: 768px) {
    .carousel-container {
      perspective: 1000px;
    }
  }
</style>

{#if !displayArtists.length}
  <div style="color: #fff; text-align: center; margin: 2rem;">Aucun artiste à afficher.</div>
{/if}

<div class="carousel-container">
  <button class="arrow" onclick={prev} disabled={isAnimating}>&larr;</button>
  <div class="cards">
    {#each displayArtists as artist, index}
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
        onclick={() => handleCardClick(index)} 
        onkeydown={(event) => handleKeyDown(event, index)} 
        aria-label="Artiste {index + 1}"
      >
        <img src={artist.image} alt={artist.name} />
        <h2>{artist.name}</h2>
      </button>
    {/each}
  </div>
  <button class="arrow" onclick={next} disabled={isAnimating}>&rarr;</button>
</div>

{#if displayArtists.length}
  <button class="view-more-button" onclick={handleViewMore}>
    Voir plus
  </button>
{/if}
