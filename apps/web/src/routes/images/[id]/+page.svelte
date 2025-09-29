<script lang="ts">
  /*
   * Copyright (c) 2025 Bouali Consulting Inc.
   * Author: Kaiss Bouali (kaissb)
   * Company: Bouali Consulting Inc.
   * GitHub: https://github.com/kaissb
   */
  import type { PageData } from './$types';

  export let data: PageData;

  function formatBytes(bytes: number, decimals = 2) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
  }
</script>

<div class="container">
  <a href="/images">&larr; Back to All Images</a>
  {#if data.imageDetails}
    <h1>Image Details: {data.imageDetails.RepoTags[0] || data.imageDetails.Id.substring(7, 19)}</h1>

    <div class="details-grid">
      <div><strong>ID:</strong> <code>{data.imageDetails.Id}</code></div>
      <div><strong>Created:</strong> {new Date(data.imageDetails.Created).toLocaleString()}</div>
      <div><strong>Size:</strong> {formatBytes(data.imageDetails.Size)}</div>
      <div><strong>Architecture:</strong> {data.imageDetails.Architecture}</div>
    </div>

    <h2>Tags</h2>
    <ul>
      {#each data.imageDetails.RepoTags as tag}
        <li><code>{tag}</code></li>
      {/each}
    </ul>

    <h2>Layers</h2>
    <div class="layers-list">
      {#each data.imageDetails.RootFS.Layers as layer, i}
        <div class="layer-item">Layer {i + 1}: <code>{layer}</code></div>
      {/each}
    </div>

  {:else if data.error}
    <p class="error">Error: {data.error}</p>
  {/if}
</div>

<style>
  .container {
    padding: 2rem;
  }
  .details-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 1rem;
    margin: 1.5rem 0;
    background: #2a2a2a;
    padding: 1rem;
    border-radius: 8px;
  }
  .layers-list {
    font-family: monospace;
    background: #2a2a2a;
    padding: 1rem;
    border-radius: 8px;
    max-height: 400px;
    overflow-y: auto;
  }
  .layer-item {
    padding: 0.25rem 0;
  }
  .error {
    color: red;
  }
</style>
