<script lang="ts">
  /*
   * Copyright (c) 2025 Bouali Consulting Inc.
   * Author: Kaiss Bouali (kaissb)
   * Company: Bouali Consulting Inc.
   * GitHub: https://github.com/kaissb
   */
  import { onMount } from 'svelte';
  import StatsChart from '$lib/components/StatsChart.svelte';

  let containers: any[] = [];
  let error: string | null = null;
  let loading = false;

  async function fetchContainers() {
    error = null;
    loading = true;
    try {
      const response = await fetch('http://localhost:8080/containers');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      containers = await response.json();
      if (!Array.isArray(containers)) {
        containers = [];
      }
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }

  async function handleAction(containerId: string, action: 'start' | 'stop' | 'restart' | 'delete') {
    if (action === 'delete') {
      if (!confirm('Are you sure you want to permanently delete this container? This cannot be undone.')) {
        return;
      }
    }

    try {
      const method = action === 'delete' ? 'DELETE' : 'POST';
      const url = action === 'delete' ? `http://localhost:8080/containers/${containerId}` : `http://localhost:8080/containers/${containerId}/${action}`;

      const response = await fetch(url, {
        method,
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ error: 'An unknown error occurred' }));
        throw new Error(errorData.error || `Failed to ${action} container`);
      }

      await fetchContainers(); // Refresh the list
    } catch (err) {
      alert(`Error: ${err.message}`);
    }
  }

  onMount(() => {
    fetchContainers();
  });
</script>

<main>
  <h1>Docker Containers</h1>
  
  {#if error}
    <p class="error">{error}</p>
  {:else}
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Image</th>
          <th>Status</th>
          <th>Names</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        {#each containers as container}
          <tr>
            <td>{container.Id.slice(0, 12)}</td>
            <td>{container.Image}</td>
            <td>{container.Status}</td>
            <td>{container.Names.join(', ')}</td>
            <td>
              <button on:click={() => handleAction(container.Id, 'start')} disabled={container.State === 'running'}>Start</button>
              <button on:click={() => handleAction(container.Id, 'stop')} disabled={container.State !== 'running'}>Stop</button>
              <button on:click={() => handleAction(container.Id, 'restart')}>Restart</button>
              <a href={`/containers/${container.Id}/logs`} class="button">Logs</a>
              <a href={`/containers/${container.Id}/terminal`} class="button" class:disabled={container.State !== 'running'}>Terminal</a>
              <button on:click={() => handleAction(container.Id, 'delete')} class="button-danger">Delete</button>
            </td>
          </tr>
          {#if container.State === 'running'}
            <tr class="stats-row">
              <td colspan="5">
                <StatsChart containerId={container.Id} />
              </td>
            </tr>
          {/if}
        {/each}
      </tbody>
    </table>
  {/if}
</main>

<style>
  .button.disabled {
    pointer-events: none;
    opacity: 0.5;
    background-color: #ccc;
  }
  table {
    width: 100%;
    border-collapse: collapse;
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
  .stats-row td {
    padding: 0;
    border-top: none;
  }
</style>
