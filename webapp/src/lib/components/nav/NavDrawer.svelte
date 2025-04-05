<script lang="ts">

    import {Sidebar, SidebarGroup, SidebarItem, SidebarWrapper} from "flowbite-svelte";
    import {
        ArrowRightToBracketOutline,
        ChartPieSolid,
        EditOutline,
        GridSolid,
    } from "flowbite-svelte-icons";
    import {page} from "$app/state";

    let spanClass = 'flex-1 ms-3 whitespace-nowrap';
    let activeUrl = $derived(page.url.pathname);
</script>

<Sidebar {activeUrl}>
    <SidebarWrapper class="h-full flex flex-col justify-between">
        <SidebarGroup>
            <SidebarItem label="Dashboard" href="/">
                <svelte:fragment slot="icon">
                    <ChartPieSolid class="w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white" />
                </svelte:fragment>
            </SidebarItem>
            <SidebarItem label="About" {spanClass} href="/about">
                <svelte:fragment slot="icon">
                    <GridSolid class="w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white" />
                </svelte:fragment>
                <svelte:fragment slot="subtext">
                    <span class="inline-flex justify-center items-center px-2 ms-3 text-sm font-medium text-gray-800 bg-gray-200 rounded-full dark:bg-gray-700 dark:text-gray-300">Pro</span>
                </svelte:fragment>
            </SidebarItem>
        </SidebarGroup>
        <SidebarGroup border>
            <SidebarItem label="Sign In" href="/auth/login">
                <svelte:fragment slot="icon">
                    <ArrowRightToBracketOutline class="w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white" />
                </svelte:fragment>
            </SidebarItem>
            <SidebarItem label="Sign Up" href="/auth/register">
                <svelte:fragment slot="icon">
                    <EditOutline class="w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white" />
                </svelte:fragment>
            </SidebarItem>
        </SidebarGroup>
    </SidebarWrapper>
</Sidebar>