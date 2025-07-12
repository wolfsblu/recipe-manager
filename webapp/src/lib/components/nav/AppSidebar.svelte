<script lang="ts">
    import CalendarIcon from "@lucide/svelte/icons/calendar";
    import ChevronUpIcon from "@lucide/svelte/icons/chevron-up";
    import HouseIcon from "@lucide/svelte/icons/house";
    import InboxIcon from "@lucide/svelte/icons/inbox";
    import SearchIcon from "@lucide/svelte/icons/search";
    import SettingsIcon from "@lucide/svelte/icons/settings";
    import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
    import * as Sidebar from "$lib/components/ui/sidebar/index";
    import ModeToggle from "$lib/components/theme/ModeToggle.svelte";

    // Menu items.
    const items = [
        {
            title: "Home",
            url: "#",
            icon: HouseIcon,
        },
        {
            title: "Inbox",
            url: "#",
            icon: InboxIcon,
        },
        {
            title: "Calendar",
            url: "#",
            icon: CalendarIcon,
        },
        {
            title: "Search",
            url: "#",
            icon: SearchIcon,
        },
        {
            title: "Settings",
            url: "#",
            icon: SettingsIcon,
        },
    ];
</script>

<Sidebar.Root>
    <Sidebar.Content>
        <Sidebar.Group>
            <Sidebar.GroupLabel>Application</Sidebar.GroupLabel>
            <Sidebar.GroupContent>
                <Sidebar.Menu>
                    {#each items as item (item.title)}
                        <Sidebar.MenuItem>
                            <Sidebar.MenuButton>
                                {#snippet child({ props })}
                                    <a href={item.url} {...props}>
                                        <item.icon />
                                        <span>{item.title}</span>
                                    </a>
                                {/snippet}
                            </Sidebar.MenuButton>
                        </Sidebar.MenuItem>
                    {/each}
                </Sidebar.Menu>
            </Sidebar.GroupContent>
        </Sidebar.Group>
    </Sidebar.Content>
    <Sidebar.Footer>
        <Sidebar.Menu>
            <Sidebar.MenuItem>
                <DropdownMenu.Root>
                    <DropdownMenu.Trigger>
                        {#snippet child({ props })}
                            <Sidebar.MenuButton
                                    {...props}
                                    class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                            >
                                Username
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
                        <DropdownMenu.Item>
                            <span>Sign out</span>
                        </DropdownMenu.Item>
                    </DropdownMenu.Content>
                </DropdownMenu.Root>
            </Sidebar.MenuItem>
        </Sidebar.Menu>
    </Sidebar.Footer>
</Sidebar.Root>
