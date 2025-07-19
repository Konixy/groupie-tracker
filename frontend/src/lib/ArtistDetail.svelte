<script lang="ts">
  let { artist, onBack } = $props();
  let concerts = $state<Record<string, string[]>>({});
  let loadingConcerts = $state(false);

  $effect(() => {
    if (artist) {
      fetchConcerts();
    }
  });

  async function fetchConcerts() {
    try {
      loadingConcerts = true;
      const response = await fetch(`http://localhost:8080/artists/${artist.id}`);
      if (response.ok) {
        const data = await response.json();
        concerts = data.concerts || {};
      }
    } catch (err) {
      console.error('Erreur lors de la récupération des concerts:', err);
    } finally {
      loadingConcerts = false;
    }
  }

  function formatDate(timestamp: number) {
    return new Date(timestamp * 1000).getFullYear();
  }

  function formatConcertDate(dateStr: string) {
    try {
      const date = new Date(dateStr);
      return date.toLocaleDateString('fr-FR', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      });
    } catch {
      return dateStr;
    }
  }

  function isPastDate(dateStr: string): boolean {
    try {
      const date = new Date(dateStr);
      const now = new Date();
      return date < now;
    } catch {
      return false;
    }
  }

  function getPastConcerts() {
    const entries = Object.entries(concerts);
    return entries
      .filter(([date]) => isPastDate(date))
      .sort(([dateA], [dateB]) => {
        const dateObjA = new Date(dateA);
        const dateObjB = new Date(dateB);
        return dateObjB.getTime() - dateObjA.getTime(); // Plus récent en premier
      }); // Afficher TOUS les concerts passés
  }
</script>

<style>
  .detail-container {
    min-height: 100vh;
    background: #4e4e4e;
    padding: 2rem;
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .back-button {
    position: absolute;
    top: 2rem;
    left: 2rem;
    background: #333;
    color: white;
    border: none;
    padding: 0.8rem 1.5rem;
    border-radius: 8px;
    cursor: pointer;
    font-size: 1rem;
    transition: background 0.3s;
  }

  .back-button:hover {
    background: #555;
  }

  .artist-card {
    background: white;
    border-radius: 20px;
    padding: 3rem;
    max-width: 800px;
    width: 100%;
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
    margin-top: 2rem;
  }

  .artist-header {
    display: flex;
    align-items: center;
    gap: 2rem;
    margin-bottom: 2rem;
  }

  .artist-image {
    width: 200px;
    height: 200px;
    object-fit: cover;
    border-radius: 16px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
  }

  .artist-info h1 {
    margin: 0 0 1rem 0;
    font-size: 2.5rem;
    color: #333;
  }

  .artist-info p {
    margin: 0.5rem 0;
    font-size: 1.1rem;
    color: #666;
  }

  .members-section {
    margin-top: 2rem;
  }

  .members-section h2 {
    color: #333;
    margin-bottom: 1rem;
  }

  .members-list {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .member-tag {
    background: #f0f0f0;
    padding: 0.5rem 1rem;
    border-radius: 20px;
    font-size: 0.9rem;
    color: #333;
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1.5rem;
    margin-top: 2rem;
  }

  .stat-card {
    background: #f8f8f8;
    padding: 1.5rem;
    border-radius: 12px;
    text-align: center;
  }

  .stat-card h3 {
    margin: 0 0 0.5rem 0;
    color: #333;
    font-size: 1.1rem;
  }

  .stat-card p {
    margin: 0;
    color: #666;
    font-size: 1.2rem;
    font-weight: bold;
  }

  .concerts-section {
    margin-top: 2rem;
  }

  .concerts-section h2 {
    color: #333;
    margin-bottom: 1rem;
  }

  .concerts-list {
    display: grid;
    gap: 1rem;
  }

  .concert-item {
    background: #fff3e0;
    border-radius: 12px;
    padding: 1rem;
    border-left: 4px solid #FF9800;
    opacity: 0.9;
  }

  .concert-date {
    font-weight: bold;
    color: #333;
    font-size: 1rem;
    margin-bottom: 0.5rem;
  }

  .concert-locations {
    color: #666;
  }

  .location-item {
    margin: 0.2rem 0;
    padding-left: 0.5rem;
    border-left: 2px solid #ddd;
    font-size: 0.9rem;
  }

  .loading {
    text-align: center;
    color: #666;
    font-style: italic;
    margin: 1rem 0;
  }


</style>

<div class="detail-container">
  <button class="back-button" onclick={onBack}>
    ← Retour
  </button>

  {#if artist}
    <div class="artist-card">
      <div class="artist-header">
        <img src={artist.image} alt={artist.name} class="artist-image" />
        <div class="artist-info">
          <h1>{artist.name}</h1>
          <p><strong>Date de création :</strong> {formatDate(artist.creationDate)}</p>
          <p><strong>Premier album :</strong> {artist.firstAlbum}</p>
        </div>
      </div>

      <div class="members-section">
        <h2>Membres du groupe</h2>
        <div class="members-list">
          {#each artist.members as member}
            <span class="member-tag">{member}</span>
          {/each}
        </div>
      </div>

      <div class="stats-grid">
        <div class="stat-card">
          <h3>Date de création</h3>
          <p>{formatDate(artist.creationDate)}</p>
        </div>
        <div class="stat-card">
          <h3>Premier album</h3>
          <p>{artist.firstAlbum}</p>
        </div>
        <div class="stat-card">
          <h3>Nombre de membres</h3>
          <p>{artist.members.length}</p>
        </div>
      </div>

      {#if loadingConcerts}
        <div class="concerts-section">
          <h2>Concerts passés</h2>
          <div class="loading">Chargement des concerts...</div>
        </div>
      {:else if getPastConcerts().length > 0}
        <div class="concerts-section">
          <h2>Tous les concerts passés</h2>
          <div class="concerts-list">
            {#each getPastConcerts() as [date, locations]}
              <div class="concert-item">
                <div class="concert-date">{formatConcertDate(date)}</div>
                <div class="concert-locations">
                  {#each locations.slice(0, 2) as location}
                    <div class="location-item">📍 {location}</div>
                  {/each}
                  {#if locations.length > 2}
                    <div class="location-item">... et {locations.length - 2} autres lieux</div>
                  {/if}
                </div>
              </div>
            {/each}
          </div>

        </div>
      {/if}
    </div>
  {/if}
</div> 