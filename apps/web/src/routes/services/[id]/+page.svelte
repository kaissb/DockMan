<script lang="ts">
  /*
   * Copyright (c) 2025 Bouali Consulting Inc.
   * Author: Kaiss Bouali (kaissb)
   * Company: Bouali Consulting Inc.
   * GitHub: https://github.com/kaissb
   */
  import type { PageData } from './$types';

  export let data: PageData;

  async function runServiceAction(serviceId: number, action: 'up' | 'down') {
    const response = await fetch(`http://localhost:8080/api/services/${serviceId}/${action}`, {
      method: 'POST',
    });
    const result = await response.json();
    alert(result.output || result.message || result.error);
    if (response.ok) {
      location.reload();
    }
  }
</script>

<main>
  {#if data.service}
    <a href="/projects/{data.service.project_id}">&larr; Back to Project</a>
    <h1>{data.service.name} <span class="service-type">({data.service.type})</span></h1>

    {#if data.service.type === 'container'}
      <p><strong>Image:</strong> <code>{data.service.image}</code></p>
      <p><strong>Container ID:</strong> <code>{data.service.container_id?.slice(0, 12) || 'N/A'}</code></p>
    {:else if data.service.type === 'compose'}
      <p><strong>Compose File Path:</strong> <code>{data.service.compose_path}</code></p>
    {/if}

    <div class="actions">
        <button on:click={() => runServiceAction(data.service.ID, 'up')}>Up</button>
        <button class="button-danger" on:click={() => runServiceAction(data.service.ID, 'down')}>Down</button>
    </div>

    {#if data.service.type === 'compose' && data.service.SubServices?.length > 0}
      <div class="sub-services">
        <h2>Sub-Services (from Compose file)</h2>
        <ul>
          {#each data.service.SubServices as subService}
            <li>
              <strong>{subService.name}</strong> ({subService.type})
              <p>Image: <code>{subService.image || 'Not specified'}</code></p>
            </li>
          {/each}
        </ul>
      </div>
    {/if}

  {:else}
    <h1>Service not found</h1>
    <p>{data.error || 'Could not load the service.'}</p>
  {/if}
</main>

<style>
  .service-type {
    font-size: 0.8em;
    color: #666;
    font-weight: normal;
  }
  .actions {
    margin: 1rem 0;
  }
   .actions button {
    margin-right: 0.5rem;
  }
  .sub-services {
    margin-top: 2rem;
  }
  ul {
    list-style-type: none;
    padding-left: 0;
  }
  li {
    background-color: #f9f9f9;
    padding: 1rem;
    border-radius: 4px;
    margin-bottom: 1rem;
  }
</style>
