import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

const config = {
	preprocess: vitePreprocess(),
	kit: {
		appDir: 'assets',
		adapter: adapter({
			fallback: 'index.html',
			pages: 'dist',
			assets: 'dist',
			precompress: false,
			strict: true
		}),
		alias: {
			$lib: 'src/lib',
		}
	},
	trailingSlash: 'always'
};

export default config;
