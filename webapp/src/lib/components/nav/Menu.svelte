<script lang="ts">
    import { Sidebar, SidebarGroup, SidebarItem, SidebarWrapper } from "flowbite-svelte";
    import {
        ArrowLeftToBracketOutline,
        ArrowRightToBracketOutline,
        EditOutline,
        ShoppingBagSolid,
        UserCircleSolid,
        ReceiptSolid,
        SearchOutline,
        RulerCombinedOutline,
        CartOutline,
        BellOutline
    } from "flowbite-svelte-icons";
    import { page } from "$app/state";
    import { isAuthenticated, logout } from "$lib/auth/user.svelte";

    interface Props {
        onNavigate?: () => void;
    }

    const { onNavigate }: Props = $props();

    const onLogout = (_: MouseEvent) => {
        logout()
        onNavigate?.()
    }

    let bubbleClass = 'inline-flex justify-center items-center px-2 ms-3 text-sm font-medium text-gray-800 bg-gray-200 rounded-full dark:bg-gray-700 dark:text-gray-300'
    let iconClass = 'w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white'
    let spanClass = 'flex-1 ms-3 whitespace-nowrap';
    let activeUrl = $derived(page.url.pathname);
</script>

<Sidebar {activeUrl} asideClass="w-full h-full">
    <SidebarWrapper class="h-full flex flex-col justify-between">
        <SidebarGroup>
            <SidebarItem label="Browse" href="/" on:click={() => onNavigate?.()}>
                <svelte:fragment slot="icon">
                    <SearchOutline class={iconClass} />
                </svelte:fragment>
            </SidebarItem>
            {#if isAuthenticated()}
                <SidebarItem label="Recipes" {spanClass} href="/recipes" on:click={() => onNavigate?.()}>
                    <svelte:fragment slot="icon">
                        <ReceiptSolid class={iconClass} />
                    </svelte:fragment>
                </SidebarItem>
                <SidebarItem label="Ingredients" {spanClass} href="/ingredients" on:click={() => onNavigate?.()}>
                    <svelte:fragment slot="icon">
                        <ShoppingBagSolid class={iconClass} />
                    </svelte:fragment>
                </SidebarItem>
                <SidebarItem label="Measurements" {spanClass} href="/measurements" on:click={() => onNavigate?.()}>
                    <svelte:fragment slot="icon">
                        <RulerCombinedOutline class={iconClass} />
                    </svelte:fragment>
                </SidebarItem>
                <SidebarItem label="Shopping List" {spanClass} href="/shopping" on:click={() => onNavigate?.()}>
                    <svelte:fragment slot="icon">
                        <CartOutline class={iconClass} />
                    </svelte:fragment>
                </SidebarItem>
            {/if}
        </SidebarGroup>
        <SidebarGroup border>
            {#if !isAuthenticated()}
                <SidebarItem label="Sign In" href="/auth/login" on:click={() => onNavigate?.()}>
                    <svelte:fragment slot="icon">
                        <ArrowRightToBracketOutline class={iconClass} />
                    </svelte:fragment>
                </SidebarItem>
                <SidebarItem label="Sign Up" href="/auth/register" on:click={() => onNavigate?.()}>
                    <svelte:fragment slot="icon">
                        <EditOutline class={iconClass} />
                    </svelte:fragment>
                </SidebarItem>
            {:else}
                <SidebarItem label="Account" href="/user/account" on:click={() => onNavigate?.()}>
                    <svelte:fragment slot="icon">
                        <UserCircleSolid class={iconClass} />
                    </svelte:fragment>
                </SidebarItem>
                <SidebarItem spanClass="flex-1 ms-3" label="Notifications" href="/user/notifications" on:click={() => onNavigate?.()}>
                    <svelte:fragment slot="icon">
                        <BellOutline class={iconClass} />
                    </svelte:fragment>
                    <svelte:fragment slot="subtext">
                        <span class="inline-flex justify-center items-center p-3 ms-3 w-3 h-3 text-sm font-medium text-primary-600 bg-primary-200 rounded-full dark:bg-primary-900 dark:text-primary-200">
                            3
                        </span>
                    </svelte:fragment>
                </SidebarItem>
                <SidebarItem label="Sign Out" href="/auth/login" on:click={onLogout}>
                    <svelte:fragment slot="icon">
                        <ArrowLeftToBracketOutline class={iconClass} />
                    </svelte:fragment>
                </SidebarItem>
            {/if}
        </SidebarGroup>
    </SidebarWrapper>
</Sidebar>