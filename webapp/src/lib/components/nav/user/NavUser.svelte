<script lang="ts">
    import * as Avatar from "$lib/components/ui/avatar";
    import * as DropdownMenu from "$lib/components/ui/dropdown-menu";
    import * as Sidebar from "$lib/components/ui/sidebar";
    import { useSidebar } from "$lib/components/ui/sidebar";
    import BadgeCheckIcon from "@lucide/svelte/icons/badge-check";
    import BellIcon from "@lucide/svelte/icons/bell";
    import ChevronsUpDownIcon from "@lucide/svelte/icons/chevrons-up-down";
    import CreditCardIcon from "@lucide/svelte/icons/credit-card";
    import LogOutIcon from "@lucide/svelte/icons/log-out";
    import SparklesIcon from "@lucide/svelte/icons/sparkles";
    import {logout, user} from "$lib/api/auth/user.svelte.js";
    import avatarFallback from "$lib/components/nav/user/avatar.png";
    import {toast} from "svelte-sonner";
    import {goto} from "$app/navigation";

    const sidebar = useSidebar();

    const handleMobileClose = () => {
        if (sidebar.isMobile) {
            sidebar.setOpenMobile(false);
        }
    };

    const onLogout = async () => {
        await logout()
        toast.success("Logged out")
        await goto("/")
        handleMobileClose();
    }
</script>

{#snippet avatar({ id, email })}
    <Avatar.Root class="size-8 rounded-lg">
        <Avatar.Image src={avatarFallback} alt={id}  />
    </Avatar.Root>
    <div class="grid flex-1 text-left text-sm leading-tight">
        <span class="truncate">{email}</span>
    </div>
{/snippet}

<Sidebar.Menu>
    <Sidebar.MenuItem>
        <DropdownMenu.Root>
            <DropdownMenu.Trigger>
                {#snippet child({ props })}
                    <Sidebar.MenuButton
                            size="lg"
                            class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                            {...props}
                    >
                        {@render avatar(user)}
                        <ChevronsUpDownIcon class="ml-auto size-4" />
                    </Sidebar.MenuButton>
                {/snippet}
            </DropdownMenu.Trigger>
            <DropdownMenu.Content
                    class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
                    side={sidebar.isMobile ? "bottom" : "right"}
                    align="end"
                    sideOffset={4}
            >
                <DropdownMenu.Label class="p-0 font-normal">
                    <div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                        {@render avatar(user)}
                    </div>
                </DropdownMenu.Label>
                <DropdownMenu.Separator />
                <DropdownMenu.Group>
                    <DropdownMenu.Item onclick={handleMobileClose}>
                        <BadgeCheckIcon />
                        Account
                    </DropdownMenu.Item>
                    <DropdownMenu.Item onclick={handleMobileClose}>
                        <BellIcon />
                        Notifications
                    </DropdownMenu.Item>
                </DropdownMenu.Group>
                <DropdownMenu.Separator />
                <DropdownMenu.Item onclick={onLogout}>
                    <LogOutIcon />
                    Log out
                </DropdownMenu.Item>
            </DropdownMenu.Content>
        </DropdownMenu.Root>
    </Sidebar.MenuItem>
</Sidebar.Menu>
