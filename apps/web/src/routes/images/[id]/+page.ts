// Copyright (c) 2025 Bouali Consulting Inc.
// Author: Kaiss Bouali (kaissb)
// Company: Bouali Consulting Inc.
// GitHub: https://github.com/kaissb

import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
  try {
    const imageId = params.id;
    const response = await fetch(`http://localhost:8080/images/${imageId}`);

    if (!response.ok) {
      throw new Error(`Could not load image details for ${imageId}`);
    }

    const imageDetails = await response.json();
    return { imageDetails };

  } catch (error) {
    return {
      status: 500,
      error: error instanceof Error ? error.message : 'Unknown error',
    };
  }
};
