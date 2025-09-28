/*
 * Copyright (c) 2025 Bouali Consulting Inc.
 * Author: Kaiss Bouali (kaissb)
 * Company: Bouali Consulting Inc.
 * GitHub: https://github.com/kaissb
 */

import adapter from '@sveltejs/adapter-auto';
import { vitePreprocess } from '@sveltejs/kit/vite';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: vitePreprocess(),
  kit: {
    adapter: adapter(),
    csrf: {
      checkOrigin: false
    }
  }
};

export default config;
