<script lang="ts">
  import { onMount } from 'svelte';
  import Logo from './lib/logo.svelte';
  import Carousel from './lib/Carousel.svelte';
  import SearchBar from './lib/SearchBar.svelte';
  import ArtistDetail from './lib/ArtistDetail.svelte';

  type Artist = {
    id: number;
    name: string;
    image: string;
    members: string[];
    creationDate: number;
    firstAlbum: string;
    locations: string;
    concertDates: string;
    relations: string;
  };

  let artists: Artist[] = [];
  let filteredArtists = $state<Artist[]>([]);
  let searchResults = $state<Artist[]>([]);
  let selectedArtist = $state<Artist | null>(null);

  onMount(async () => {
    const res = await fetch('http://localhost:8080/artists');
    const data = await res.json();
    console.log(data);
    // Supporte à la fois l'API locale (data.artists) et distante (data)
    artists = Array.isArray(data) ? data : data.artists;
    filteredArtists = artists;

    // Écouter les événements de clic sur les artistes
    document.addEventListener('artistClick', (event: any) => {
      selectedArtist = event.detail;
    });

    // Écouter les événements de recherche
    document.addEventListener('search', (event: any) => {
      handleSearch(event.detail);
    });
  });

  function handleSearch(query: string) {
    if (!query.trim()) {
      filteredArtists = artists;
      searchResults = [];
      // Réinitialiser le carrousel à la position 0
      const event = new CustomEvent('centerOnArtist', {
        detail: { index: 0 }
      });
      document.dispatchEvent(event);
      return;
    }

    // Trouver tous les artistes qui commencent par la recherche
    const results = artists.filter(artist =>
      artist.name.toLowerCase().startsWith(query.toLowerCase())
    );
    
    searchResults = results;
    
    if (results.length > 0) {
      // Centrer le carrousel sur le premier artiste trouvé
      const event = new CustomEvent('centerOnArtist', {
        detail: { index: 0 }
      });
      document.dispatchEvent(event);
    }
  }

  function handleBack() {
    selectedArtist = null;
  }
</script>

<style>
  main {
    min-height: 100vh;
    background: #4e4e4e;
    box-sizing: border-box;
    padding: 2rem 0 4rem 0;
    display: flex;
    flex-direction: column;
    align-items: center;
  }
</style>

{#if selectedArtist}
  <ArtistDetail artist={selectedArtist} onBack={handleBack} />
{:else}
  <main>
    <Logo />
    <Carousel artists={filteredArtists} searchResults={searchResults} />
    <SearchBar on:search={e => handleSearch(e.detail)} />
  </main>
{/if}
