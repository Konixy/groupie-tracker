<script lang="ts">
  import { onMount } from 'svelte';
  import Logo from './lib/logo.svelte';
  import Carousel from './lib/Carousel.svelte';
  import SearchBar from './lib/SearchBar.svelte';

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
  let filteredArtists: Artist[] = [];

  onMount(async () => {
    const res = await fetch('http://localhost:8080/artists');
    const data = await res.json();
    console.log(data);
    // Supporte à la fois l'API locale (data.artists) et distante (data)
    artists = Array.isArray(data) ? data : data.artists;
    filteredArtists = artists;
  });

  function handleSearch(query: string) {
    filteredArtists = artists.filter(artist =>
      artist.name.toLowerCase().includes(query.toLowerCase())
    );
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

<main>
  <Logo />
  <Carousel artists={filteredArtists} />
  <SearchBar on:search={e => handleSearch(e.detail)} />
</main>
