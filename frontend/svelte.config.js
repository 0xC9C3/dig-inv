import adapter from '@sveltejs/adapter-static';
import {vitePreprocess} from '@sveltejs/vite-plugin-svelte';

const config = {
	preprocess: vitePreprocess(),
	kit: { adapter:
			adapter({
				fallback: 'index.html' // may differ from host to host
			})
	}
};

export default config;
