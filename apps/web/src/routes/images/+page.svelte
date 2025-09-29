<script lang="ts">
  /*
   * Copyright (c) 2025 Bouali Consulting Inc.
   * Author: Kaiss Bouali (kaissb)
   * Company: Bouali Consulting Inc.
   * GitHub: https://github.com/kaissb
   */
  import { onMount } from 'svelte';

  let images: any[] = [];
  let imageName = '';
  let pullLog = '';
  let isPulling = false;

  async function fetchImages() {
    try {
      const response = await fetch('http://localhost:8080/images');
      if (!response.ok) {
        throw new Error('Failed to fetch images');
      }
      images = await response.json();
    } catch (err) {
      alert(`Error: ${(err as Error).message}`);
    }
  }

  async function pullImage() {
    if (!imageName) {
      alert('Please enter an image name.');
      return;
    }
    isPulling = true;
    pullLog = `Pulling image: ${imageName}...\n`;
    try {
      const response = await fetch('http://localhost:8080/images/pull', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name: imageName }),
      });

      if (!response.body) {
        throw new Error('Response body is null');
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder();
      let done = false;

      while (!done) {
        const { value, done: readerDone } = await reader.read();
        done = readerDone;
        if (value) {
          const chunk = decoder.decode(value, { stream: true });
          pullLog += chunk;
        }
      }

      pullLog += `\nImage ${imageName} pulled successfully.`;
      imageName = '';
      await fetchImages();
    } catch (err) {
      pullLog += `\nError pulling image: ${(err as Error).message}`;
    } finally {
      isPulling = false;
    }
  }

  async function deleteImage(imageId: string) {
    if (!confirm('Are you sure you want to delete this image? This cannot be undone.')) {
      return;
    }
    try {
      const response = await fetch(`http://localhost:8080/images/${imageId}`, {
        method: 'DELETE',
      });
      const result = await response.json();
      if (!response.ok) {
        throw new Error(result.error || 'Failed to delete image');
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
      <input type="text" bind:value={imageName} placeholder="e.g., ubuntu:latest" required />
      <button type="submit" disabled={isPulling}>{isPulling ? 'Pulling...' : 'Pull Image'}</button>
    </form>
    {#if pullLog}
      <pre class="pull-log">{pullLog}</pre>
    {/if}
  </div>

  <table>
    <thead>
      <tr>
        <th>Tag</th>
        <th>Image ID</th>
        <th>Created</th>
        <th>Size</th>
        <th>Actions</th>
      </tr>
    </thead>
    <tbody>
      {#each images as image}
        <tr>
          <td><a href="/images/{image.ID}">{image.RepoTags && image.RepoTags.length > 0 ? image.RepoTags.join(', ') : 'N/A'}</a></td>
          <td>{image.ID.replace('sha256:', '').slice(0, 12)}</td>
          <td>{formatDate(image.Created)}</td>
          <td>{formatSize(image.Size)}</td>
          <td>
            <button class="button-danger" on:click={() => deleteImage(image.ID)} disabled={isProtectedImage(image)}>Delete</button>
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
</main>

<style>
  main {
    padding: 2rem;
  }
  .pull-image-form {
    margin-bottom: 2rem;
    background: #2a2a2a;
    padding: 1.5rem;
    border-radius: 8px;
  }
  .pull-image-form form {
    display: flex;
    gap: 1rem;
  }
  .pull-image-form input {
    flex-grow: 1;
  }
  .pull-log {
    margin-top: 1rem;
    background: #111;
    color: #eee;
    padding: 1rem;
    border-radius: 4px;
    max-height: 200px;
    overflow-y: auto;
    white-space: pre-wrap;
    word-wrap: break-word;
  }
  table {
    width: 100%;
    border-collapse: collapse;
  }
  th, td {
    padding: 0.75rem 1rem;
    text-align: left;
    border-bottom: 1px solid #333;
  }
  th {
    background-color: #333;
  }
</style>
