<script lang="ts">
  /*
   * Copyright (c) 2025 Bouali Consulting Inc.
   * Author: Kaiss Bouali (kaissb)
   * Company: Bouali Consulting Inc.
   * GitHub: https://github.com/kaissb
   */
  import { onMount } from 'svelte';

  let images: any[] = [];
  let error: string | null = null;
  let loading = true;
  let newImageName = '';
  let pullLog = '';
  let isPulling = false;

  async function fetchImages() {
    error = null;
    loading = true;
    try {
      const response = await fetch('http://localhost:8080/images');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const rawImages = await response.json();
      images = rawImages.filter((img: any) => img.RepoTags && img.RepoTags[0] !== '<none>:<none>');
    } catch (err) {
      error = (err as Error).message;
    } finally {
      loading = false;
    }
  }

  async function pullImage() {
    if (!newImageName) return;
    isPulling = true;
    pullLog = `Pulling image: ${newImageName}...\n`;
    try {
      const response = await fetch('http://localhost:8080/images/pull', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: newImageName }),
      });

      const reader = response.body?.getReader();
      const decoder = new TextDecoder();

      while (true) {
        const { done, value } = await reader?.read() || { done: true };
        if (done) break;
        pullLog += decoder.decode(value, { stream: true });
      }

      newImageName = '';
      await fetchImages();
    } catch (err) {
      pullLog += `\nError pulling image: ${(err as Error).message}`;
    } finally {
      isPulling = false;
    }
  }

  async function deleteImage(imageId: string) {
    if (!confirm('Are you sure you want to delete this image?')) return;
    try {
      const response = await fetch(`http://localhost:8080/images/${imageId}`, {
        method: 'DELETE',
      });
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Failed to delete image');
      }
      await fetchImages();
    } catch (err) {
      alert(`Error: ${(err as Error).message}`);
    }
  }

  function formatSize(bytes: number): string {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  }

  function isProtectedImage(image: any): boolean {
    if (!image.RepoTags) return false;
    return image.RepoTags.some((tag: string) => tag.includes('dockman'));
  }

  function formatDate(timestamp: number): string {
    return new Date(timestamp * 1000).toLocaleString();
  }

  onMount(() => {
    fetchImages();
  });
</script>

<main>
  <h1>Docker Images</h1>

  <div class="pull-image-form">
    <h2>Pull New Image</h2>
    <form on:submit|preventDefault={pullImage}>
      <input type="text" bind:value={newImageName} placeholder="e.g., ubuntu:latest" required disabled={isPulling} />
      <button type="submit" disabled={isPulling}>{isPulling ? 'Pulling...' : 'Pull Image'}</button>
    </form>
    {#if pullLog}
      <pre class="pull-log">{pullLog}</pre>
    {/if}
  </div>

  {#if loading}
    <p>Loading images...</p>
  {:else if error}
    <p class="error">{error}</p>
  {:else}
    <table>
      <thead>
        <tr>
          <th>Tag</th>
          <th>ID</th>
          <th>Created</th>
          <th>Size</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        {#each images as image}
          <tr>
            <td>{image.RepoTags[0]}</td>
            <td>{image.Id.replace('sha256:', '').slice(0, 12)}</td>
            <td>{formatDate(image.Created)}</td>
            <td>{formatSize(image.Size)}</td>
            <td>
              <button class="button-danger" on:click={() => deleteImage(image.Id)} disabled={isProtectedImage(image)}>Delete</button>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</main>

<style>
  table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 1rem;
  }
  th, td {
    border: 1px solid #ddd;
    padding: 8px;
    text-align: left;
  }
  th {
    background-color: #f2f2f2;
  }
  .error {
    color: red;
  }
  .pull-image-form {
    margin-bottom: 2rem;
  }
  .pull-log {
    background-color: #222;
    color: #eee;
    padding: 1rem;
    border-radius: 4px;
    max-height: 300px;
    overflow-y: auto;
    white-space: pre-wrap;
    word-wrap: break-word;
  }
</style>
