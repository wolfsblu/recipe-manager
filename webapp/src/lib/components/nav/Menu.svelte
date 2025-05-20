<script lang="ts">
    import {
        DarkMode, Indicator,
        Sidebar,
        SidebarDropdownItem,
        SidebarDropdownWrapper,
        SidebarGroup,
        SidebarItem,
        SidebarWrapper,
        Tooltip,
    } from "flowbite-svelte";
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
        BellOutline,
    } from "flowbite-svelte-icons";
    import { page } from "$app/state";
    import { isAuthenticated, logout } from "$lib/auth/user.svelte";
    import Notifications from "../notifications/Notifications.svelte";

    interface Props {
        onNavigate?: () => void;
    }

    const { onNavigate }: Props = $props();

    const onLogout = (_: MouseEvent) => {
        logout();
        onNavigate?.();
    };

    let bubbleClass =
        "inline-flex justify-center items-center px-2 ms-3 text-sm font-medium text-gray-800 bg-gray-200 rounded-full dark:bg-gray-700 dark:text-gray-300";
    let iconClass =
        "w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white";
    let spanClass = "flex-1 ms-3 whitespace-nowrap";
    let activeUrl = $derived(page.url.pathname);
    const isActive = (url: string) => {
        return activeUrl.endsWith(url);
    };

    let showNotifications = $state(false);
</script>

<Notifications bind:open={showNotifications} />

<Sidebar {activeUrl} class="h-full shadow" divClass="h-full" isSingle={false}>
    <SidebarWrapper class="h-full flex flex-col justify-between rounded-l-none">
        <SidebarGroup>
            <SidebarItem
                label="Browse"
                href="/"
                onclick={() => onNavigate?.()}
            >
                {#snippet icon()}
                    <SearchOutline class={iconClass} />
                {/snippet}
            </SidebarItem>
            {#if isAuthenticated()}
                <SidebarDropdownWrapper label="Recipes" isOpen>
                    {#snippet icon()}
                        <ReceiptSolid class={iconClass} />
                    {/snippet}
                    <SidebarDropdownItem
                        label="My Recipes"
                        href="/recipes"
                        onclick={() => onNavigate?.()}
                        active={isActive("/recipes")}
                    />
                    <SidebarDropdownItem
                        label="Planned Recipes"
                        href="/recipes/plan"
                        onclick={() => onNavigate?.()}
                        active={isActive("/recipes/plan")}
                    />
                    <SidebarDropdownItem
                        label="Saved Recipes"
                        href="/recipes/saved"
                        onclick={() => onNavigate?.()}
                        active={isActive("/recipes/saved")}
                    />
                </SidebarDropdownWrapper>
                <SidebarItem
                    label="Ingredients"
                    {spanClass}
                    href="/ingredients"
                    onclick={() => onNavigate?.()}
                >
                    {#snippet icon()}
                        <ShoppingBagSolid class={iconClass} />
                    {/snippet}
                </SidebarItem>
                <SidebarItem
                    label="Measurements"
                    {spanClass}
                    href="/measurements"
                    onclick={() => onNavigate?.()}
                >
                    {#snippet icon()}
                        <RulerCombinedOutline class={iconClass} />
                    {/snippet}
                </SidebarItem>
                <SidebarItem
                    label="Shopping List"
                    {spanClass}
                    href="/shopping"
                    onclick={() => onNavigate?.()}
                >
                    {#snippet icon()}
                        <CartOutline class={iconClass} />
                    {/snippet}
                </SidebarItem>
            {/if}
        </SidebarGroup>
        <SidebarGroup border>
            {#if !isAuthenticated()}
                <SidebarItem
                    label="Sign In"
                    href="/auth/login"
                    onclick={() => onNavigate?.()}
                >
                    {#snippet icon()}
                        <ArrowRightToBracketOutline class={iconClass} />
                    {/snippet}
                </SidebarItem>
                <SidebarItem
                    label="Sign Up"
                    href="/auth/register"
                    onclick={() => onNavigate?.()}
                >
                    {#snippet icon()}
                        <EditOutline class={iconClass} />
                    {/snippet}
                </SidebarItem>
            {:else}
                <SidebarItem
                    label="Account"
                    href="/user/account"
                    onclick={() => onNavigate?.()}
                >
                    {#snippet icon()}
                        <UserCircleSolid class={iconClass} />
                    {/snippet}
                </SidebarItem>
                <SidebarItem
                    spanClass="flex-1 ms-3"
                    label="Notifications"
                    href="/user/notifications"
                    onclick={(e) => { e.preventDefault(); showNotifications = true}}
                >
                    {#snippet icon()}
                        <div class="relative">
                            <BellOutline class={iconClass} />
                            <Indicator color="red" size="xs" placement="bottom-right" />
                        </div>
                    {/snippet}
                    {#snippet subtext()}
                        <span class={bubbleClass}> 3 </span>
                    {/snippet}
                </SidebarItem>

                <SidebarItem
                    label="Sign Out"
                    href="/auth/login"
                    onclick={onLogout}
                >
                    {#snippet icon()}
                        <ArrowLeftToBracketOutline class={iconClass} />
                    {/snippet}
                </SidebarItem>
            {/if}
            <section class="flex justify-center">
                <DarkMode class="cursor-pointer" />
                <Tooltip>Toggle dark mode</Tooltip>
            </section>
        </SidebarGroup>
    </SidebarWrapper>
</Sidebar>