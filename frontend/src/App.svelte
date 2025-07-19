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
  });

  function handleSearch(query: string) {
    filteredArtists = artists.filter(artist =>
      artist.name.toLowerCase().includes(query.toLowerCase())
    );
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
    <Carousel artists={filteredArtists} />
    <SearchBar on:search={e => handleSearch(e.detail)} />
  </main>
{/if}
