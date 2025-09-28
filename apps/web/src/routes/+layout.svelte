<script lang="ts">
  /*
   * Copyright (c) 2025 Bouali Consulting Inc.
   * Author: Kaiss Bouali (kaissb)
   * Company: Bouali Consulting Inc.
   * GitHub: https://github.com/kaissb
   */
  import '../app.css';
</script>

<header>
  <nav>
    <a href="/">Home</a>
    <a href="/containers">Containers</a>
    <a href="/projects">Projects</a>
    <a href="/images">Images</a>
  </nav>
</header>

<main>
  <slot />
</main>

<style>
  header {
    background-color: #f8f9fa;
    padding: 1rem;
    border-bottom: 1px solid #e0e0e0;
  }
  nav a {
    margin-right: 1rem;
    font-weight: bold;
    color: #333;
  }
  main {
    padding: 1rem;
  }
  :global(body) {
    margin: 0;
  }
  :global(body) {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
    line-height: 1.5;
    color: #333;
  }

  :global(a) {
    color: #007bff;
    text-decoration: none;
  }

  :global(a:hover) {
    text-decoration: underline;
  }
</style>
