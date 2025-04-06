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
		})
	},
	paths: {
		base: '',
		assets: '/'
	},
	trailingSlash: 'always'
};

export default config;
