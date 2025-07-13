<script lang="ts">
    import CalendarIcon from "@lucide/svelte/icons/calendar";
    import ChevronUpIcon from "@lucide/svelte/icons/chevron-up";
    import ChefHat from "@lucide/svelte/icons/chef-hat";
    import InboxIcon from "@lucide/svelte/icons/inbox";
    import PlusIcon from "@lucide/svelte/icons/plus";
    import LoginIcon from "@lucide/svelte/icons/log-in";
    import SearchIcon from "@lucide/svelte/icons/search";
    import UtensilsIcon from "@lucide/svelte/icons/utensils";
    import SettingsIcon from "@lucide/svelte/icons/settings";
    import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
    import * as Sidebar from "$lib/components/ui/sidebar/index";
    import {isAuthenticated, logout, user} from "$lib/auth/user.svelte";
    import {toast} from "svelte-sonner";

    const onLogout = async () => {
        await logout()
        toast.success("Logged out")
    }
</script>

<Sidebar.Root variant="inset">
    <Sidebar.Header>
        <Sidebar.Menu>
            <Sidebar.MenuItem>
                <Sidebar.MenuButton>
                    {#snippet child({ props })}
                        <a href="/" {...props}>
                            <UtensilsIcon class="!size-5" />
                            <span class="text-base font-semibold">Grecipes</span>
                        </a>
                    {/snippet}
                </Sidebar.MenuButton>
            </Sidebar.MenuItem>
        </Sidebar.Menu>
    </Sidebar.Header>
    <Sidebar.Content>
        <Sidebar.Group>
            <Sidebar.GroupLabel>Application</Sidebar.GroupLabel>
            <Sidebar.GroupContent>
                <Sidebar.Menu>
                    <Sidebar.MenuItem>
                        <Sidebar.MenuButton>
                            {#snippet child({ props })}
                                <a href="/recipes" {...props}>
                                    <ChefHat />
                                    <span>Recipes</span>
                                </a>
                            {/snippet}
                        </Sidebar.MenuButton>
                        <Sidebar.MenuAction showOnHover>
                            <PlusIcon />
                        </Sidebar.MenuAction>
                    </Sidebar.MenuItem>
                </Sidebar.Menu>
            </Sidebar.GroupContent>
        </Sidebar.Group>
    </Sidebar.Content>
    <Sidebar.Footer>
        {#if isAuthenticated()}
        <Sidebar.Menu>
            <Sidebar.MenuItem>
                <DropdownMenu.Root>
                    <DropdownMenu.Trigger>
                        {#snippet child({ props })}
                            <Sidebar.MenuButton
                                    {...props}
                                    class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                            >
                                {user.email}
                                <ChevronUpIcon class="ml-auto" />
                            </Sidebar.MenuButton>
                        {/snippet}
                    </DropdownMenu.Trigger>
                    <DropdownMenu.Content
                            side="top"
                            class="w-(--bits-dropdown-menu-anchor-width)"
                    >
                        <DropdownMenu.Item>
                            <span>Account</span>
                        </DropdownMenu.Item>
                        <DropdownMenu.Item>
                            <span>Billing</span>
                        </DropdownMenu.Item>
                        <DropdownMenu.Item onclick={onLogout}>
                            <span>Sign out</span>
                        </DropdownMenu.Item>
                    </DropdownMenu.Content>
                </DropdownMenu.Root>
            </Sidebar.MenuItem>
        </Sidebar.Menu>
        {:else}
            <Sidebar.Menu>
                <Sidebar.MenuItem>
                    <Sidebar.MenuButton>
                        {#snippet child({ props })}
                            <a href="/auth/login" {...props}>
                                <LoginIcon />
                                <span>Login</span>
                            </a>
                        {/snippet}
                    </Sidebar.MenuButton>
                </Sidebar.MenuItem>
            </Sidebar.Menu>
        {/if}
    </Sidebar.Footer>
</Sidebar.Root>
