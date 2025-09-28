/*
 * Copyright (c) 2025 Bouali Consulting Inc.
 * Author: Kaiss Bouali (kaissb)
 * Company: Bouali Consulting Inc.
 * GitHub: https://github.com/kaissb
 */

import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
  const projectId = params.id;
  try {
    const projectRes = await fetch(`http://localhost:8080/api/projects/${projectId}`);
    if (!projectRes.ok) {
      throw new Error(`Could not load project ${projectId}`);
    }
    const project = await projectRes.json();

    // Fetch environments and their services
    const environmentsRes = await fetch(`http://localhost:8080/api/projects/${projectId}/environments`);
    if (!environmentsRes.ok) {
      throw new Error('Could not load environments');
    }
    project.Environments = await environmentsRes.json();

    return { project };

  } catch (error) {
    return {
      status: 500,
      error: error instanceof Error ? error.message : 'Unknown error',
    };
  }
};
