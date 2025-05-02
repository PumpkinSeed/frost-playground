import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		fs: {
			// Allow serving files from the entire project directory
			allow: ['..']
		}
	},
	// Add base configuration for static builds
	base: '',
	build: {
		// Ensure proper path resolution for static files
		rollupOptions: {
			output: {
				manualChunks: () => 'everything.js', // forces all code into a single file
				inlineDynamicImports: true // disables dynamic import chunks
			}
		}
	}
});
