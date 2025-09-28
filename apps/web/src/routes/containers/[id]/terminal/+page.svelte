<script lang="ts">
  /*
   * Copyright (c) 2025 Bouali Consulting Inc.
   * Author: Kaiss Bouali (kaissb)
   * Company: Bouali Consulting Inc.
   * GitHub: https://github.com/kaissb
   */
  import '@xterm/xterm/css/xterm.css';
  import { page } from '$app/stores';
  import { onMount, onDestroy } from 'svelte';
  import { browser } from '$app/environment';

  let terminalEl: HTMLElement;
  const containerId = $page.params.id;
  const socketURL = `ws://localhost:8080/ws/terminal/${containerId}`;

  onMount(async () => {
    if (browser) {
      const { Terminal } = await import('@xterm/xterm');
      const { AttachAddon } = await import('@xterm/addon-attach');
      const { FitAddon } = await import('@xterm/addon-fit');

      const term = new Terminal({ cursorBlink: true, convertEol: true });
      const fitAddon = new FitAddon();
      term.loadAddon(fitAddon);

      term.open(terminalEl);

      const socket = new WebSocket(socketURL);
      const attachAddon = new AttachAddon(socket, { bidirectional: true });

      term.loadAddon(attachAddon);

      socket.onopen = () => {
        fitAddon.fit();
        term.focus();
      };

      socket.onerror = (err) => {
        console.error('WebSocket error:', err);
        term.write('--- WebSocket Connection Error ---');
      };

      socket.onclose = () => {
        term.write('--- WebSocket Connection Closed ---');
      };

      const resizeListener = () => {
        fitAddon.fit();
      };
      window.addEventListener('resize', resizeListener);

      onDestroy(() => {
        window.removeEventListener('resize', resizeListener);
        socket.close();
        term.dispose();
      });
    }
  });
</script>

<main>
  <h1>Terminal for Container {containerId.slice(0, 12)}</h1>
  <a href="/containers">&larr; Back to Containers</a>

  <div bind:this={terminalEl} class="terminal-container"></div>
</main>

<style>
  .terminal-container {
    width: 100%;
    height: 80vh;
    margin-top: 1rem;
  }
</style>
