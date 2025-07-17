<script lang="ts">
  let { artists = [] } = $props();
  let current = $state(0);

  function prev() {
    current = (current - 1 + artists.length) % artists.length;
  }
  function next() {
    current = (current + 1) % artists.length;
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
    background: radial-gradient(circle, #444 60%, #222 100%);
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
  }
  .cards {
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    width: 70vw;
    height: 60vh;
  }
  .card {
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
    transform: translate(-50%, -50%) scale(1);
    transition: transform 0.5s cubic-bezier(.4,2,.6,1), box-shadow 0.5s, opacity 0.5s, filter 0.5s;
    z-index: 2;
    width: 400px;
    height: 400px;
    opacity: 1;
  }
  .card.side {
    width: 200px;
    height: 200px;
    z-index: 1;
    box-shadow: 0 2px 8px #0002;
    filter: grayscale(0.7);
    transform: translateY(-50%) scale(0.7);
    pointer-events: none;
  }
  .card.leftleft  { transform: translate(-50%, -50%) translateX(-235px) scale(1); z-index: 1; }
  .card.left      { transform: translate(-50%, -50%) translateX(-150px) scale(1.4); z-index: 2; }
  .card.center    { transform: translate(-50%, -50%) scale(0.9); z-index: 3; }
  .card.right     { transform: translate(-50%, -50%) translateX(150px) scale(1.4); z-index: 2; }
  .card.rightright{ transform: translate(-50%, -50%) translateX(235px) scale(1); z-index: 1; }

  .card img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-radius: 16px;
    box-shadow: 0 4px 24px #0005;
  }
  .card h2 {
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
  }
</style>

{#if !artists.length}
  <div style="color: #fff; text-align: center; margin: 2rem;">Aucun artiste à afficher.</div>
{/if}

<div class="carousel-container">
  <button class="arrow" onclick={prev}>&larr;</button>
  <div class="cards">
    {#if artists.length}
      <!-- Left++ side image (hidden behind, blurred) -->
      <div class="card side leftleft">
        <img src={artists[(current - 2 + artists.length) % artists.length].image} alt="artist" />
        <h2>{artists[(current - 2 + artists.length) % artists.length].name}</h2>
      </div>
      <!-- Left side image (hidden behind, blurred) -->
      <div class="card side left">
        <img src={artists[(current - 1 + artists.length) % artists.length].image} alt="artist" />
        <h2>{artists[(current - 1 + artists.length) % artists.length].name}</h2>
      </div>
      <!-- Center image (large, always present) -->
      <div class="card center">
        <img src={artists[current].image} alt={artists[current].name} />
        <h2>{artists[current].name}</h2>
      </div>
      <!-- Right side image (hidden behind, blurred) -->
      <div class="card side right">
        <img src={artists[(current + 1) % artists.length].image} alt="artist" />
        <h2>{artists[(current + 1) % artists.length].name}</h2>
      </div>
      <!-- Right++ side image (hidden behind, blurred) -->
      <div class="card side rightright">
        <img src={artists[(current + 2) % artists.length].image} alt="artist" />
        <h2>{artists[(current + 2) % artists.length].name}</h2>
      </div>
    {/if}
  </div>
  <button class="arrow" onclick={next}>&rarr;</button>
</div>
