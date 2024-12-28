<script lang="ts">
    import NavButton from "./NavButton.svelte";
    import Plus from "../../icons/Plus.svelte";
    import User from "../../icons/User.svelte";
    import logo from "../../../assets/images/logo.svg"
    import {t} from "../../i18n/i18n.svelte"
    import Dropdown from "../dropdown/Dropdown.svelte";
    import type {Component} from "svelte";
    import {createUser} from "../../auth/user.svelte";
    import {createRouter} from "../../router.svelte";
    import type {MenuItem} from "../dropdown/types";
    import Logout from "../../icons/Logout.svelte";
    import Login from "../../icons/Login.svelte";

    const router = createRouter()
    const user = createUser()

    const onLogout = async () => {
        console.log("HI");
        await user.logout()
        router.redirectToNext()
    }

    const profileMenu: MenuItem[] = $derived(user.profile ? [
        { class: "text-xs font-semibold", label: user.profile.email },
        { icon: User, label: "Account", href: "/" },
        { icon: Logout, label: "Sign Out", onClick: onLogout },
    ] : [
        { icon: Login, label: "Sign In", href: "/login" },
    ])
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

        <Dropdown menu={profileMenu}>
            {#snippet button()}
                {@render navButton(User)}
            {/snippet}
        </Dropdown>
    </div>
</div>