<script lang="ts">
    import {createRouter} from "../../router.svelte";
    import {createUser} from "../../auth/user.svelte";
    import User from "../../icons/User.svelte";
    import Login from "../../icons/Login.svelte";
    import Logout from "../../icons/Logout.svelte";

    const router = createRouter()
    const user = createUser()

    const onLogout = async () => {
        console.log("HI");
        await user.logout()
        router.redirectToNext()
    }

    const iconClass = "relative h-5 w-5"
    const itemClass = "flex flex-1 gap-x-3 items-center justify-self-stretch px-4 py-2 text-gray-700 text-sm"
</script>

<div class="py-1" role="none">
    {#if user.profile}
        <span class={`text-xs font-semibold ${itemClass}`}>
            {user.profile.email}
        </span>
        <a href="/user/account" tabindex="-1" class={`cursor-pointer hover:bg-orange-50 ${itemClass}`}>
            <span class={iconClass}>
                <User/>
            </span>
            Account
        </a>
    {:else}
        <a href="/login" tabindex="-1" class={`cursor-pointer hover:bg-orange-50 ${itemClass}`}>
            <span class={iconClass}>
                <Login/>
            </span>
            Sign In
        </a>
    {/if}
</div>
{#if user.profile}
    <div class="py-1" role="none">
        <button tabindex="-1" class={`hover:bg-orange-50 ${itemClass}`} onclick={onLogout}>
                    <span class={iconClass}>
                        <Logout/>
                    </span>
            Sign Out
        </button>
    </div>
{/if}