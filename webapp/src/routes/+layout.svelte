<script lang="ts">
	import '../app.css';
    import Languages from "$lib/components/i18n/Languages.svelte";
    import NavDrawer from "$lib/components/nav/Navigation.svelte";
    import Menu from '$lib/components/nav/Menu.svelte';
    import Toasts from "$lib/components/notifications/Toasts.svelte";

    let { children } = $props();
</script>

<Languages />
<Toasts />

<div class="flex flex-col h-screen">
    <NavDrawer classes="md:hidden" />

    <div class="flex grow">
        <div class="hidden md:block w-64">
            <Menu />
        </div>

        <main class="flex flex-col grow p-3">
            {@render children()}
        </main>
    </div>
</div>