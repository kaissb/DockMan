/*
 * Copyright (c) 2025 Bouali Consulting Inc.
 * Author: Kaiss Bouali (kaissb)
 * Company: Bouali Consulting Inc.
 * GitHub: https://github.com/kaissb
 */

import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
  const serviceId = params.id;
  try {
    const response = await fetch(`http://localhost:8080/api/services/${serviceId}`);
    if (!response.ok) {
      throw new Error(`Could not load service ${serviceId}`);
    }
    const service = await response.json();
    return { service };
  } catch (error) {
    return {
      status: 500,
      error: error instanceof Error ? error.message : 'Unknown error',
    };
  }
};
