<script lang="ts">
  import { onMount } from 'svelte';
  import Logo from './lib/logo.svelte';
  import Carousel from './lib/Carousel.svelte';
  import SearchBar from './lib/SearchBar.svelte';
  import ArtistDetail from './lib/ArtistDetail.svelte';
  import ConcertsMap from './lib/ConcertsMap.svelte';
  import type { Artist } from './types';

  let artists: Artist[] = [];
  let filteredArtists = $state<Artist[]>([]);
  let searchResults = $state<Artist[]>([]);
  let selectedArtist = $state<Artist | null>(null);
  let searchQuery = $state('');

  onMount(async () => {
    const res = await fetch('http://localhost:8080/artists');
    const data = await res.json();
    // Supporte à la fois l'API locale (data.artists) et distante (data)
    artists = Array.isArray(data) ? data : data.artists;
    filteredArtists = artists;
  });

  $effect(() => {
    if (searchQuery !== null) {
      handleSearch();
    }
  });

  function handleSearch() {
    if (!searchQuery.trim()) {
      filteredArtists = artists;
      searchResults = [];
      return;
    }

    const searchTerm = searchQuery.toLowerCase();
    
    // Créer un tableau avec les artistes et leur score de pertinence
    const scoredResults = artists.map(artist => {
      const artistName = artist.name.toLowerCase();
      let score = 0;
      
      // Score le plus élevé : commence par la recherche
      if (artistName.startsWith(searchTerm)) {
        score = 1000;
      }
      // Score élevé : contient la recherche quelque part dans le nom
      else if (artistName.includes(searchTerm)) {
        const index = artistName.indexOf(searchTerm);
        score = 1000 - index; // Plus tôt dans le nom = score plus élevé
      }
      
      return { artist, score };
    });
    
    // Filtrer les artistes qui correspondent et les trier par score
    let results = scoredResults
      .filter(item => item.score > 0)
      .sort((a, b) => {
        // Tri par score décroissant
        if (b.score !== a.score) {
          return b.score - a.score;
        }
        // Si même score, trier par ordre alphabétique
        return a.artist.name.localeCompare(b.artist.name);
      })
      .map(item => item.artist)
      // Éliminer les doublons en gardant seulement la première occurrence
      .filter((artist, index, self) => 
        index === self.findIndex(a => a.id === artist.id)
      );

    // Si on a moins de 3 résultats, ajouter des artistes similaires
    // if (results.length < 3 && results.length > 0) {
    //   const exactMatches = results;
    //   const remainingArtists = artists.filter(artist => 
    //     !exactMatches.some(exact => exact.id === artist.id)
    //   );
      
    //   // Trouver des artistes qui ont au moins 2 lettres en commun avec la recherche
    //   const similarResults = remainingArtists.map(artist => {
    //     const artistName = artist.name.toLowerCase();
    //     let commonLetters = 0;
        
    //     // Compter les lettres communes
    //     for (let i = 0; i < searchTerm.length; i++) {
    //       if (artistName.includes(searchTerm[i])) {
    //         commonLetters++;
    //       }
    //     }
        
    //     return { artist, commonLetters };
    //   })
    //   .filter(item => item.commonLetters >= 2) // Au moins 2 lettres en commun
    //   .sort((a, b) => b.commonLetters - a.commonLetters) // Plus de lettres communes = meilleur score
    //   .map(item => item.artist)
    //   .slice(0, 5 - results.length); // Compléter jusqu'à 5 artistes max
      
    //   results = [...exactMatches, ...similarResults];
    // }
    
    searchResults = results;

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
  <SearchBar bind:query={searchQuery} />
  <Carousel bind:selectedArtist={selectedArtist} artists={filteredArtists} searchResults={searchResults} searchTerm={searchQuery.trim()} />
  {#if selectedArtist}
    <ArtistDetail bind:artist={selectedArtist} onClose={() => selectedArtist = null} />
  {/if}
  <ConcertsMap />
</main>
