<script lang="ts">
  /*
   * Copyright (c) 2025 Bouali Consulting Inc.
   * Author: Kaiss Bouali (kaissb)
   * Company: Bouali Consulting Inc.
   * GitHub: https://github.com/kaissb
   */
  import type { PageData } from './$types';

  export let data: PageData;
  let replicas: { [key: string]: number } = {};

  async function scaleService(serviceId: number, subServiceName: string) {
    const numReplicas = replicas[subServiceName] || 1;
    const response = await fetch(`http://localhost:8080/api/services/${serviceId}/scale`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ sub_service_name: subServiceName, replicas: numReplicas }),
    });

    if (response.ok) {
      alert(`Service ${subServiceName} scaled to ${numReplicas} replicas.`);
    } else {
      const error = await response.json();
      alert(`Failed to scale service: ${error.error}`);
    }
  }

  async function runServiceAction(serviceId: number, action: 'up' | 'down') {
    const response = await fetch(`http://localhost:8080/api/services/${serviceId}/${action}`, {
      method: 'POST',
    });

    if (response.ok) {
      alert(`Service action '${action}' completed successfully.`);
    } else {
      const error = await response.json();
      alert(`Failed to perform action: ${error.error}`);
    }
  }
</script>

<main>
  {#if data.service}
    <h1>{data.service.name}</h1>
    <p class="service-type">Type: {data.service.type}</p>

    {#if data.service.type === 'compose'}
      <div class="actions">
        <button on:click={() => runServiceAction(data.service.ID, 'up')}>Up</button>
        <button class="button-danger" on:click={() => runServiceAction(data.service.ID, 'down')}>Down</button>
      </div>

      <div class="sub-services">
        <h2>Sub-Services (from Compose file)</h2>
        <ul>
          {#each data.service.SubServices as subService}
            <li>
              <strong>{subService.name}</strong> ({subService.type})
              <p>Image: <code>{subService.image}</code></p>
              <div class="scale-controls">
                <input type="number" min="0" bind:value={replicas[subService.name]} placeholder="Replicas" />
                <button on:click={() => scaleService(data.service.ID, subService.name)}>Scale</button>
              </div>
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
    text-transform: uppercase;
    margin-bottom: 1rem;
  }
  .actions {
    margin-bottom: 2rem;
  }
  .sub-services {
    margin-top: 2rem;
    border-top: 1px solid #eee;
    padding-top: 1rem;
  }
  ul {
    list-style: none;
    padding: 0;
  }
  li {
    background: #2a2a2a;
    padding: 1rem;
    border-radius: 4px;
    margin-bottom: 1rem;
  }
  .scale-controls {
    margin-top: 0.5rem;
    display: flex;
    gap: 0.5rem;
    align-items: center;
  }
  .scale-controls input {
    width: 80px;
  }
</style>
