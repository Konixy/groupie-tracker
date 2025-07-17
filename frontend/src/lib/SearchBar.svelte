<script lang="ts">
  export let dispatch: (type: string, detail?: any) => void = () => {};
  let query = '';

  function search() {
    dispatch('search', { detail: query });
  }
</script>

<style>
  .search-bar {
    display: flex;
    justify-content: center;
    margin-top: 2rem;
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
    on:keydown={(e) => e.key === 'Enter' && search()}
  />
  <button on:click={search}>🔍</button>
</div>
