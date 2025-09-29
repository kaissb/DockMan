<script lang="ts">
  /*
   * Copyright (c) 2025 Bouali Consulting Inc.
   * Author: Kaiss Bouali (kaissb)
   * Company: Bouali Consulting Inc.
   * GitHub: https://github.com/kaissb
   */
  import type { PageData } from './$types';
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import EnvironmentVariables from '$lib/components/EnvironmentVariables.svelte';

  export let data: PageData;

  let newServiceName = '';
  let newServiceType = 'container';
  let newServiceImage = '';
  let newServiceComposePath = '';
  let selectedEnvironmentId: number | null = null;

  async function loadProject() {
    const projectId = $page.params.id;
    const projectRes = await fetch(`http://localhost:8080/api/projects/${projectId}`);
    if (projectRes.ok) {
      data.project = await projectRes.json();
      const environmentsRes = await fetch(`http://localhost:8080/api/projects/${projectId}/environments`);
      if (environmentsRes.ok) {
        const environments = await environmentsRes.json();
        for (const env of environments) {
            const variablesRes = await fetch(`http://localhost:8080/api/environments/${env.ID}/variables`);
            if (variablesRes.ok) {
                env.Variables = await variablesRes.json();
            }
        }
        data.project.Environments = environments;
      }
    }
  }

  async function addService() {
    if (!selectedEnvironmentId) {
      alert('Please select an environment.');
      return;
    }

    const serviceData: any = {
      name: newServiceName,
      type: newServiceType,
    };

    if (newServiceType === 'container') {
      serviceData.image = newServiceImage;
    } else if (newServiceType === 'compose') {
      serviceData.compose_path = newServiceComposePath;
    }

    const response = await fetch(`http://localhost:8080/api/environments/${selectedEnvironmentId}/services`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(serviceData),
    });

    if (response.ok) {
      // For simplicity, we'll just alert and reload. A better UX would update the list directly.
      alert('Service added!');
      location.reload();
    } else {
      const error = await response.json();
      alert(`Failed to add service: ${error.error}`);
    }
  }

  async function runServiceAction(environmentId: number, serviceId: number, action: 'up' | 'down') {
    const response = await fetch(`http://localhost:8080/api/environments/${environmentId}/services/${serviceId}/${action}`, {
      method: 'POST',
    });
    const result = await response.json();
    alert(result.output || result.message || result.error);
  }

</script>

<main>
  <a href="/projects">&larr; Back to All Projects</a>
  {#if data.project}
    <h1>{data.project.name}</h1>
    <p>{data.project.description}</p>

    <div class="environments">
      <h2>Environments</h2>
      {#each data.project.Environments as environment}
        <div class="environment-card">
          <h3>{environment.name}</h3>

          <EnvironmentVariables variables={environment.Variables} environmentId={environment.ID} on:update={loadProject} />
          
          <h4>Services</h4>
          {#if environment.Services && environment.Services.length > 0}
            <ul>
              {#each environment.Services as service}
                <li>
                  <a href="/services/{service.ID}"><strong>{service.name}</strong> ({service.type})</a>
                  {#if service.type === 'container'}
                    <p>Image: <code>{service.image}</code></p>
                  {:else if service.type === 'compose'}
                    <p>Path: <code>{service.compose_path}</code></p>
                  {/if}
                  <div class="actions">
                    <button on:click={() => runServiceAction(environment.ID, service.ID, 'up')}>Up</button>
                    <button class="button-danger" on:click={() => runServiceAction(environment.ID, service.ID, 'down')}>Down</button>
                  </div>
                </li>
              {/each}
            </ul>
          {:else}
            <p>No services in this environment.</p>
          {/if}

          <div class="add-service-form">
            <h4>Add New Service to {environment.name}</h4>
            <form on:submit|preventDefault={() => { selectedEnvironmentId = environment.ID; addService(); }}>
              <input type="text" bind:value={newServiceName} placeholder="Service Name" required>
              <select bind:value={newServiceType}>
                <option value="container">Container</option>
                <option value="compose">Compose</option>
              </select>
              {#if newServiceType === 'container'}
                <input type="text" bind:value={newServiceImage} placeholder="Image (e.g., nginx:latest)" required>
              {:else if newServiceType === 'compose'}
                <input type="text" bind:value={newServiceComposePath} placeholder="/path/to/docker-compose.yml" required>
              {/if}
              <button type="submit">Add Service</button>
            </form>
          </div>
        </div>
      {:else}
        <p>No environments found for this project.</p>
      {/each}
    </div>
  {:else}
    <h1>Project not found</h1>
    <p>{data.error || 'Could not load the project.'}</p>
  {/if}
</main>

<style>
  .environments { margin-top: 2rem; }
  .environment-card { border: 1px solid #ccc; border-radius: 8px; padding: 1rem; margin-bottom: 1rem; }
  ul { list-style-type: none; padding-left: 0; }
  li { background-color: #f9f9f9; padding: 1rem; border-radius: 4px; margin-bottom: 1rem; }
  .actions button { margin-right: 0.5rem; }
  .add-service-form { margin-top: 1.5rem; padding-top: 1rem; border-top: 1px solid #eee; }
  input, select { display: block; margin-bottom: 0.5rem; width: 100%; max-width: 400px; }
</style>
