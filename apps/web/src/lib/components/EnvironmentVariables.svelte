<script lang="ts">
  /*
   * Copyright (c) 2025 Bouali Consulting Inc.
   * Author: Kaiss Bouali (kaissb)
   * Company: Bouali Consulting Inc.
   * GitHub: https://github.com/kaissb
   */
  import { createEventDispatcher } from 'svelte';

  export let variables: any[] = [];
  export let environmentId: number;

  const dispatch = createEventDispatcher();

  let newKey = '';
  let newValue = '';
  let revealedId: number | null = null;

  async function addVariable() {
    if (!newKey || !newValue) {
      alert('Both key and value are required.');
      return;
    }

    const response = await fetch(`http://localhost:8080/api/environments/${environmentId}/variables`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ key: newKey, value: newValue }),
    });

    if (response.ok) {
      newKey = '';
      newValue = '';
      dispatch('update');
    } else {
      const error = await response.json();
      alert(`Failed to add variable: ${error.error}`);
    }
  }

  async function deleteVariable(variableId: number) {
    if (!confirm('Are you sure you want to delete this variable?')) return;

    const response = await fetch(`http://localhost:8080/api/environments/${environmentId}/variables/${variableId}`, {
      method: 'DELETE',
    });

    if (response.ok) {
      dispatch('update');
    } else {
      alert('Failed to delete variable.');
    }
  }

  function toggleReveal(variableId: number) {
    if (revealedId === variableId) {
      revealedId = null;
    } else {
      revealedId = variableId;
    }
  }
</script>

<div class="variables-manager">
  <h4>Environment Variables</h4>
  <div class="variable-list">
    {#if variables && variables.length > 0}
      <table>
        <thead>
          <tr>
            <th>Key</th>
            <th>Value</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each variables as variable (variable.ID)}
            <tr>
              <td><code>{variable.key}</code></td>
              <td>
                {#if revealedId === variable.ID}
                  <code>{variable.value}</code>
                {:else}
                  <code>********</code>
                {/if}
              </td>
              <td class="actions">
                <button on:click={() => toggleReveal(variable.ID)}>{revealedId === variable.ID ? 'Hide' : 'Show'}</button>
                <button class="button-danger" on:click={() => deleteVariable(variable.ID)}>Delete</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    {:else}
      <p>No environment variables defined.</p>
    {/if}
  </div>

  <form on:submit|preventDefault={addVariable} class="add-variable-form">
    <input type="text" bind:value={newKey} placeholder="Variable Name" required />
    <input type="text" bind:value={newValue} placeholder="Variable Value" required />
    <button type="submit">Add Variable</button>
  </form>
</div>

<style>
  .variables-manager {
    margin-top: 1.5rem;
    padding-top: 1rem;
    border-top: 1px solid #eee;
  }
  table {
    width: 100%;
    border-collapse: collapse;
    margin-bottom: 1rem;
  }
  th, td {
    padding: 0.5rem;
    text-align: left;
    border-bottom: 1px solid #eee;
  }
  .add-variable-form {
    display: flex;
    gap: 0.5rem;
  }
  .add-variable-form input {
    flex-grow: 1;
  }
  .actions button {
    margin-right: 0.5rem;
    font-size: 0.8em;
    padding: 0.2rem 0.5rem;
  }
</style>
