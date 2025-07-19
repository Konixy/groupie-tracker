<script lang="ts">
  let { artist, onBack } = $props();

  function formatDate(timestamp: number) {
    return new Date(timestamp * 1000).getFullYear();
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
    </div>
  {/if}
</div> 