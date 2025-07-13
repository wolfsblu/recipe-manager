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
    import NavUser from "$lib/components/nav/user/NavUser.svelte";

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
            <NavUser />
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
