<script lang="ts">
    import CalendarIcon from "@lucide/svelte/icons/calendar";
    import ChevronUpIcon from "@lucide/svelte/icons/chevron-up";
    import ChefIcon from "@lucide/svelte/icons/chef-hat";
    import InboxIcon from "@lucide/svelte/icons/inbox";
    import LoginIcon from "@lucide/svelte/icons/log-in";
    import PlusIcon from "@lucide/svelte/icons/plus";
    import RulerIcon from "@lucide/svelte/icons/ruler";
    import SaladIcon from "@lucide/svelte/icons/salad";
    import SearchIcon from "@lucide/svelte/icons/search";
    import ShoppingIcon from "@lucide/svelte/icons/shopping-cart";
    import UtensilsIcon from "@lucide/svelte/icons/utensils";
    import SettingsIcon from "@lucide/svelte/icons/settings";
    import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
    import * as Sidebar from "$lib/components/ui/sidebar/index";
    import NavUser from "$lib/components/nav/user/NavUser.svelte";
    import { goto } from '$app/navigation'
    import { isAuthenticated } from "$lib/api/auth/user.svelte";

    const onAddRecipe = () => goto("/recipes/add")
</script>

<Sidebar.Root variant="inset" collapsible="icon">
    <Sidebar.Header>
        <Sidebar.Menu>
            <Sidebar.MenuItem>
                <Sidebar.MenuButton>
                    {#snippet child({ props })}
                        <a href="/" {...props}>
                            <UtensilsIcon class="!size-5" />
                            <span class="text-base font-semibold">GRM</span>
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
                                    <ChefIcon />
                                    <span>Recipes</span>
                                </a>
                            {/snippet}
                        </Sidebar.MenuButton>
                        <Sidebar.MenuAction onclick={onAddRecipe}>
                            <PlusIcon />
                        </Sidebar.MenuAction>
                    </Sidebar.MenuItem>
                    <Sidebar.MenuItem>
                        <Sidebar.MenuButton>
                            {#snippet child({ props })}
                                <a href="/ingredients" {...props}>
                                    <SaladIcon />
                                    <span>Ingredients</span>
                                </a>
                            {/snippet}
                        </Sidebar.MenuButton>
                        <Sidebar.MenuAction>
                            <PlusIcon />
                        </Sidebar.MenuAction>
                    </Sidebar.MenuItem>
                    <Sidebar.MenuItem>
                        <Sidebar.MenuButton>
                            {#snippet child({ props })}
                                <a href="/units" {...props}>
                                    <RulerIcon />
                                    <span>Units</span>
                                </a>
                            {/snippet}
                        </Sidebar.MenuButton>
                        <Sidebar.MenuAction>
                            <PlusIcon />
                        </Sidebar.MenuAction>
                    </Sidebar.MenuItem>
                    <Sidebar.MenuItem>
                        <Sidebar.MenuButton>
                            {#snippet child({ props })}
                                <a href="/mealplan" {...props}>
                                    <CalendarIcon />
                                    <span>Meal Plan</span>
                                </a>
                            {/snippet}
                        </Sidebar.MenuButton>
                    </Sidebar.MenuItem>
                    <Sidebar.MenuItem>
                        <Sidebar.MenuButton>
                            {#snippet child({ props })}
                                <a href="/shopping" {...props}>
                                    <ShoppingIcon />
                                    <span>Shopping List</span>
                                </a>
                            {/snippet}
                        </Sidebar.MenuButton>
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
