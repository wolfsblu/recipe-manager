<script lang="ts">
    import NavButton from "./NavButton.svelte";
    import Plus from "../../icons/Plus.svelte";
    import User from "../../icons/User.svelte";
    import logo from "../../../assets/images/logo.svg"
    import {t} from "../../i18n/i18n.svelte"
    import Dropdown from "../dropdown/Dropdown.svelte";
    import type {Component} from "svelte";
</script>

<div class="flex items-center justify-between bg-orange-200 px-3 py-2">
    <a class="flex items-center gap-x-3" href="/">
        <span class="block bg-orange-50 rounded-full w-11 h-11 p-2">
            <img src={logo} alt="Logo">
        </span>
    </a>

    <div class="flex items-center gap-2">
        <NavButton text={t("navigation.actions.newRecipe")} href="/recipes/create" icon={Plus} />

        {#snippet navButton(icon: Component, href?: string)}
            <NavButton {href} {icon} />
        {/snippet}

        <Dropdown>
            {#snippet button()}
                {@render navButton( User)}
            {/snippet}
            <!-- Active: "bg-gray-100 text-gray-900 outline-none", Not Active: "text-gray-700" -->
            <a href="/" class="block px-4 py-2 text-sm text-gray-700" role="menuitem" tabindex="-1" id="menu-item-0">Account settings</a>
            <a href="/" class="block px-4 py-2 text-sm text-gray-700" role="menuitem" tabindex="-1" id="menu-item-1">Support</a>
            <a href="/" class="block px-4 py-2 text-sm text-gray-700" role="menuitem" tabindex="-1" id="menu-item-2">License</a>
            <form method="POST" action="#" role="none">
                <button type="submit" class="block w-full px-4 py-2 text-left text-sm text-gray-700" role="menuitem" tabindex="-1" id="menu-item-3">Sign out</button>
            </form>
        </Dropdown>
    </div>
</div>