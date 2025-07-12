<script lang="ts">
	import '../app.css';
    import Languages from "$lib/components/i18n/Languages.svelte";
    import { ModeWatcher } from "mode-watcher";
    import { commandContext } from '$lib/context';
    import * as Breadcrumb from "$lib/components/ui/breadcrumb/index.js";
    import * as Sidebar from "$lib/components/ui/sidebar/index.js";
    import AppSidebar from "$lib/components/nav/AppSidebar.svelte";
    import {Separator} from "$lib/components/ui/separator";
    import SearchButton from "$lib/components/nav/SearchButton.svelte";
    import {UseBoolean} from "$lib/hooks/use-boolean.svelte";
    import AppCommand from "$lib/components/nav/AppCommand.svelte";
    import ModeToggle from "$lib/components/theme/ModeToggle.svelte";

    let { children } = $props();
    commandContext.set(new UseBoolean(false))
</script>

<AppCommand />
<Languages />
<ModeWatcher />

<Sidebar.Provider>
    <AppSidebar />
    <main class="flex flex-col flex-grow">
        <header class="bg-background shrink-0 items-center gap-2 transition-[width,height] ease-linear p-3">
            <div class="flex items-center justify-between">
                <div class="flex items-center gap-2 px-2">
                    <Sidebar.Trigger class="-ml-1" />
                    <Separator orientation="vertical" class="mr-2 data-[orientation=vertical]:h-4" />
                    <Breadcrumb.Root>
                        <Breadcrumb.List>
                            <Breadcrumb.Item class="hidden md:block">
                                <Breadcrumb.Link href="#">Building Your Application</Breadcrumb.Link>
                            </Breadcrumb.Item>
                            <Breadcrumb.Separator class="hidden md:block" />
                            <Breadcrumb.Item>
                                <Breadcrumb.Page>Data Fetching</Breadcrumb.Page>
                            </Breadcrumb.Item>
                        </Breadcrumb.List>
                    </Breadcrumb.Root>
                </div>
                <div class="flex gap-2">
                    <SearchButton />
                    <ModeToggle />
                </div>
            </div>
        </header>

        {@render children?.()}
    </main>
</Sidebar.Provider>
