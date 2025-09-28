<script lang="ts">
  /*
   * Copyright (c) 2025 Bouali Consulting Inc.
   * Author: Kaiss Bouali (kaissb)
   * Company: Bouali Consulting Inc.
   * GitHub: https://github.com/kaissb
   */
  import { onMount } from 'svelte';

  let projects: any[] = [];
  let newProjectName = '';
  let newProjectDescription = '';

  async function fetchProjects() {
    try {
      const response = await fetch('http://localhost:8080/api/projects');
      if (response.ok) {
        projects = await response.json();
      } else {
        projects = [];
      }
    } catch (error) {
      console.error('Failed to fetch projects:', error);
      projects = [];
    }
  }

  async function createProject() {
    if (!newProjectName) {
      alert('Project name is required.');
      return;
    }
    try {
      const response = await fetch('http://localhost:8080/api/projects', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: newProjectName, description: newProjectDescription }),
      });
      if (response.ok) {
        newProjectName = '';
        newProjectDescription = '';
        await fetchProjects(); // Refresh the list
      } else {
        alert('Failed to create project.');
      }
    } catch (error) {
      console.error('Failed to create project:', error);
      alert('An error occurred while creating the project.');
    }
  }

  onMount(() => {
    fetchProjects();
  });
</script>

<main>
  <h1>Projects</h1>

  <div class="create-project-form">
    <h2>Create New Project</h2>
    <form on:submit|preventDefault={createProject}>
      <input type="text" bind:value={newProjectName} placeholder="Project Name" required />
      <textarea bind:value={newProjectDescription} placeholder="Project Description"></textarea>
      <button type="submit">Create Project</button>
    </form>
  </div>

  <div class="project-list">
    <h2>Existing Projects</h2>
    {#if projects && projects.length > 0}
      <ul>
        {#each projects as project}
          <li>
            <h3><a href={`/projects/${project.ID}`}>{project.name}</a></h3>
            <p>{project.description || 'No description'}</p>
          </li>
        {/each}
      </ul>
    {:else}
      <p>No projects found. Create one above!</p>
    {/if}
  </div>
</main>

<style>
  main {
    padding: 2rem;
  }
  .create-project-form, .project-list {
    margin-top: 2rem;
  }
  input, textarea {
    width: 100%;
    padding: 0.5rem;
    margin-bottom: 1rem;
  }
  ul {
    list-style: none;
    padding: 0;
  }
  li {
    border: 1px solid #ccc;
    padding: 1rem;
    margin-bottom: 1rem;
    border-radius: 4px;
  }
</style>
