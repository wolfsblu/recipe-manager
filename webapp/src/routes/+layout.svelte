<script lang="ts">
	import '../app.css';
    import Languages from "$lib/components/i18n/Languages.svelte";
    import NavDrawer from "$lib/components/nav/NavDrawer.svelte";

    let { children } = $props();
</script>

<Languages />

<div class="flex h-screen">
    <NavDrawer />

    <main class="flex flex-col grow">
        {@render children()}
    </main>
</div>
