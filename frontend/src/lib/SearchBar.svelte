<script lang="ts">
  let query = $state('');

  function search() {
    // Émettre un événement Svelte
    const event = new CustomEvent('search', {
      detail: query
    });
    document.dispatchEvent(event);
  }

  // Recherche en temps réel
  $effect(() => {
    if (query !== undefined) {
      search();
    }
  });
</script>

<style>
  .search-bar {
    display: flex;
    justify-content: center;
    margin-top: 4rem;
  }
  input {
    border-radius: 20px 0 0 20px;
    border: none;
    padding: 1rem;
    width: 350px;
    font-size: 1rem;
    outline: none;
    background: #ffffff;
    color: #000000;
  }
  button {
    border-radius: 0 20px 20px 0;
    border: none;
    background: #727373;
    color: #fff;
    font-size: 1.5rem;
    padding: 0 1rem;
    cursor: pointer;
  }
</style>

<div class="search-bar">
  <input
    type="text"
    placeholder="Chercher un artiste ou un groupe…"
    bind:value={query}
    onkeydown={(e) => e.key === 'Enter' && search()}
  />
  <button onclick={search}>🔍</button>
</div>
