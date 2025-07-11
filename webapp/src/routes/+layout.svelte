<script lang="ts">
	import '../app.css';
    import Languages from "$lib/components/i18n/Languages.svelte";
    import { ModeWatcher } from "mode-watcher";
    import * as Sidebar from "$lib/components/ui/sidebar/index.js";
    import AppSidebar from "$lib/components/nav/AppSidebar.svelte";

    let { children } = $props();
</script>

<Languages />
<ModeWatcher />

<Sidebar.Provider>
    <AppSidebar />
    <main>
        {@render children?.()}
    </main>
</Sidebar.Provider>
