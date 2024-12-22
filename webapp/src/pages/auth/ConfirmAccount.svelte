<script lang="ts">
    import Layout from "../../Layout.svelte";
    import Navbar from "../../lib/components/navigation/Navbar.svelte";
    import {onMount} from "svelte";
    import {createRouter} from "../../lib/router.svelte";
    import {createUser} from "../../lib/auth/user.svelte";
    import {te} from '../../lib/i18n/i18n.svelte'

    const user = createUser()
    const router = createRouter()

    let countdown = $state(5)
    let error: APIError | null = $state(null)

    onMount(async () => {
        const urlParams = new URLSearchParams(window.location.search);
        const token = urlParams.get("token") ?? ""
        try {
            await user.confirmAccount(token)
            let handle = setInterval(() => {
                if (--countdown <= 0) {
                    clearInterval(handle)
                    router.redirectToLogin()
                }
            }, 1000)
        } catch (err: unknown) {
            error = err as APIError
        }
    })
</script>

<Layout Header={Navbar}>
    <div class="bg-gray-50 flex flex-col flex-wrap md:h-full items-center justify-center">
        <div class="bg-white md:shadow-lg md:w-1/3 rounded p-6">
            <h1 class="font-light mb-3 text-3xl">Account Confirmation</h1>
            {#if error}
                <p>
                    {@html te(error.message)}
                </p>
            {:else}
                <p>
                    Thank you for confirming your account, you will be redirected to the
                    <a href="/login">login</a> automatically in {countdown} seconds.
                </p>
            {/if}
        </div>
    </div>
</Layout>
