<script lang="ts">
  let { query = $bindable() }: {
    query: string;
  } = $props();
</script>

<style>
  .search-bar {
    display: flex;
    justify-content: center;
    align-items: center;
    margin-top: 4rem;
  }
  input {
    border-radius: 20px;
    border: none;
    padding: 1rem;
    width: 350px;
    font-size: 1rem;
    font-family: "Arial Rounded MT Bold", Arial, sans-serif;
    outline: none;
    background: #ffffff;
    color: #000000;
    padding-left: 2.8rem;
  }

  .lucide-search-icon {
    margin-right: -2.2rem;
    z-index: 100;
    color: rgb(101, 101, 101);
  }
</style>

<div class="search-bar">
  <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-search-icon lucide-search"><path d="m21 21-4.34-4.34"/><circle cx="11" cy="11" r="8"/></svg>
  <input
    type="text"
    placeholder="Chercher un artiste ou un groupe…"
    bind:value={query}
  />
</div>
