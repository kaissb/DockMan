<script>
  import { onMount } from 'svelte';
  
  let status = 'loading';
  
  onMount(async () => {
    try {
      // Try both localhost and the Docker service name
      const urls = [
        'http://localhost:8080/health',
        'http://api:8080/health'
      ];
      
      let lastError;
      
      for (const url of urls) {
        try {
          const response = await fetch(url);
          if (response.ok) {
            console.log(`Connected to API at ${url}`);
            status = 'connected';
            return;
          }
        } catch (err) {
          console.log(`Failed to connect to ${url}:`, err);
          lastError = err;
        }
      }
      
      // If we get here, all connection attempts failed
      console.error('All API connection attempts failed', lastError);
      status = 'error';
    } catch (err) {
      console.error('Error in health check:', err);
      status = 'error';
    }
  });
</script>

<main>
  <h1>Docker Manager</h1>
  <p>Status: {status}</p>
  {#if status === 'connected'}
    <p>API is connected and ready!</p>
    <a href="/containers">View Containers</a>
  {:else if status === 'error'}
    <p>Failed to connect to API server</p>
  {:else}
    <p>Connecting to API...</p>
  {/if}
</main>

<style>
  main {
    max-width: 800px;
    margin: 0 auto;
    padding: 2rem;
  }
  a {
    display: inline-block;
    margin-top: 1rem;
    padding: 0.5rem 1rem;
    background-color: #007bff;
    color: white;
    text-decoration: none;
    border-radius: 4px;
  }
</style>
