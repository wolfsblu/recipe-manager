import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

const config = {
	preprocess: vitePreprocess(),
	kit: {
		appDir: 'assets',
		adapter: adapter({
			pages: 'dist',
			assets: 'dist',
			precompress: false,
			strict: true
		}),
		env: {
			dir: './..'
		}
	}
};

export default config;
