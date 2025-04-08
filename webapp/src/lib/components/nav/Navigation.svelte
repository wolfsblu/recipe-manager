<script lang="ts">
    import {Drawer, Navbar, NavBrand, NavHamburger} from "flowbite-svelte";
    import Menu from "$lib/components/nav/Menu.svelte";
    import {sineIn} from "svelte/easing";

    interface Props {
        classes: string
    }
    let { classes }: Props = $props();

    let hideNavigation = $state(true)
    const toggleNavigation = () => {
        hideNavigation = !hideNavigation;
    }
</script>

<Navbar class={classes}>
    <NavBrand href="/">
        <span class="self-center whitespace-nowrap text-xl font-semibold dark:text-white">Go Chef</span>
    </NavBrand>
    <NavHamburger onClick={() => hideNavigation = !hideNavigation}  />
</Navbar>

<Drawer bind:hidden={hideNavigation}
        transitionType="fly"
        transitionParams={{x: -256, duration: 125, easing: sineIn}}
        divClass="overflow-y-auto z-50 bg-white dark:bg-gray-800"
        width="w-64">
    <Menu onNavigate={toggleNavigation} />
</Drawer>