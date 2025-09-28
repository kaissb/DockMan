<script lang="ts">
  /*
   * Copyright (c) 2025 Bouali Consulting Inc.
   * Author: Kaiss Bouali (kaissb)
   * Company: Bouali Consulting Inc.
   * GitHub: https://github.com/kaissb
   */
  import { page } from '$app/stores';
  import { onMount, onDestroy } from 'svelte';

  let logs: string[] = [];
  let socket: WebSocket;
  const containerId = $page.params.id;

  onMount(() => {
    // Connect to the WebSocket endpoint
    socket = new WebSocket(`ws://localhost:8080/ws/logs/${containerId}`);

    socket.onopen = () => {
      console.log('WebSocket connection established');
      logs = [...logs, '--- Connection established ---'];
    };

    socket.onmessage = (event) => {
      // The raw log from Docker might contain non-printable characters or headers.
      // For now, we'll just append it.
      logs = [...logs, event.data];
    };

    socket.onclose = () => {
      console.log('WebSocket connection closed');
      logs = [...logs, '--- Connection closed ---'];
    };

    socket.onerror = (error) => {
      console.error('WebSocket error:', error);
      logs = [...logs, '--- WebSocket Error ---'];
    };
  });

  onDestroy(() => {
    // Clean up the connection when the component is destroyed
    if (socket) {
      socket.close();
    }
  });
</script>

<main>
  <h1>Logs for Container {containerId.slice(0, 12)}</h1>
  <a href="/containers">&larr; Back to Containers</a>

  <div class="log-container">
    {#each logs as log}
      <pre>{log}</pre>
    {/each}
  </div>
</main>

<style>
  .log-container {
    background-color: #1a1a1a;
    color: #f0f0f0;
    padding: 1rem;
    border-radius: 4px;
    margin-top: 1rem;
    height: 60vh;
    overflow-y: scroll;
    font-family: 'Courier New', Courier, monospace;
  }

  pre {
    margin: 0;
    white-space: pre-wrap;
    word-wrap: break-word;
  }

  a {
    display: inline-block;
    margin-bottom: 1rem;
  }
</style>
