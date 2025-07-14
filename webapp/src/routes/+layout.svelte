<script lang="ts">
	import '../app.css';
    import Languages from "$lib/components/i18n/Languages.svelte";
    import { ModeWatcher } from "mode-watcher";
    import { commandContext } from '$lib/context';
    import * as Breadcrumb from "$lib/components/ui/breadcrumb/index.js";
    import * as Sidebar from "$lib/components/ui/sidebar/index.js";
    import AppSidebar from "$lib/components/nav/AppSidebar.svelte";
    import { Separator } from "$lib/components/ui/separator";
    import { Toaster } from "$lib/components/ui/sonner/index.js";
    import SearchButton from "$lib/components/nav/SearchButton.svelte";
    import {UseBoolean} from "$lib/hooks/use-boolean.svelte";
    import AppCommand from "$lib/components/nav/AppCommand.svelte";
    import ModeToggle from "$lib/components/theme/ModeToggle.svelte";
    import AppShortcuts from "$lib/components/nav/AppShortcuts.svelte";
    import type {LayoutProps} from './$types';
    import { page } from '$app/state'

    let { data, children }: LayoutProps = $props();
    commandContext.set(new UseBoolean(false))

    const breadcrumbs = $derived(page.data.breadcrumbs || data.breadcrumbs)
</script>

<AppCommand />
<AppShortcuts />
<Languages />
<ModeWatcher />
<Toaster />

<Sidebar.Provider>
    <AppSidebar />
    <Sidebar.Inset>
        <main class="flex flex-col flex-grow">
            <header class="shrink-0 items-center gap-2 transition-[width,height] ease-linear p-3 border-b">
                <div class="flex items-center justify-between">
                    <div class="flex items-center gap-2 px-2">
                        <Sidebar.Trigger class="-ml-1" />
                        <Separator orientation="vertical" class="mr-2 data-[orientation=vertical]:h-4" />
                        <Breadcrumb.Root>
                            <Breadcrumb.List>
                                {#each breadcrumbs as breadcrumb, i (breadcrumb.link)}
                                    {#if i > 0}
                                        <Breadcrumb.Separator class="hidden md:block" />
                                    {/if}
                                    <Breadcrumb.Item class="hidden md:block">
                                        <Breadcrumb.Link href={breadcrumb.link}>{breadcrumb.name}</Breadcrumb.Link>
                                    </Breadcrumb.Item>
                                {/each}
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
    </Sidebar.Inset>
</Sidebar.Provider>
