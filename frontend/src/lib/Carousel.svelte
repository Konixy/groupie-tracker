<script lang="ts">
  export let artists: Array<{
    id: number;
    name: string;
    image: string;
    members: string[];
    creationDate: number;
    firstAlbum: string;
    locations: string;
    concertDates: string;
    relations: string;
  }> = [];
  let current = 0;

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
  }
  .arrow:hover {
    opacity: 1;
  }
  .cards {
    display: flex;
    align-items: center;
    justify-content: center;
    perspective: 1000px;
  }
  .card {
    background: #ddd;
    border-radius: 16px;
    width: 300px;
    height: 300px;
    margin: 0 1rem;
    box-shadow: 0 8px 32px #0004;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    transition: transform 0.4s, box-shadow 0.4s;
    position: relative;
    z-index: 2;
  }
  .card.side {
    width: 220px;
    height: 220px;
    opacity: 0.6;
    z-index: 1;
    box-shadow: 0 4px 16px #0002;
    transform: scale(0.85) translateY(30px);
  }
  .card img {
    width: 80%;
    border-radius: 12px;
    margin-bottom: 1rem;
  }
  .card h2 {
    margin: 0;
    font-size: 1.3rem;
    color: #222;
  }
</style>

{#if !artists.length}
  <div style="color: #fff; text-align: center; margin: 2rem;">Aucun artiste à afficher.</div>
{/if}

<div class="carousel-container">
  <button class="arrow" on:click={prev}>&larr;</button>
  <div class="cards">
    {#if artists.length}
      {#if artists[current - 1]}
        <div class="card side">
          <img src={artists[(current - 1 + artists.length) % artists.length].image} alt="artist" />
          <h2>{artists[(current - 1 + artists.length) % artists.length].name}</h2>
        </div>
      {/if}
      <div class="card">
        <img src={artists[current].image} alt="artist" />
        <h2>{artists[current].name}</h2>
      </div>
      {#if artists[current + 1]}
        <div class="card side">
          <img src={artists[(current + 1) % artists.length].image} alt="artist" />
          <h2>{artists[(current + 1) % artists.length].name}</h2>
        </div>
      {/if}
    {/if}
  </div>
  <button class="arrow" on:click={next}>&rarr;</button>
</div>
