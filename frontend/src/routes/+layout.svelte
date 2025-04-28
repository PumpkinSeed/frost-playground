<script lang="ts">
	import '../app.css';
	import Navigation from '$lib/components/Navigation.svelte';
</script>

<div class="min-h-screen bg-gray-50">
	<Navigation />
	<slot />
</div>
