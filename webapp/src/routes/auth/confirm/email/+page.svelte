<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
    import * as Alert from "$lib/components/ui/alert/index.js";
    import SuccessIcon from "@lucide/svelte/icons/check-circle-2";
    import {confirmEmail} from "$lib/auth/user.svelte";
    import {goto} from "$app/navigation";
    import {onDestroy, onMount} from "svelte";
    import {toast} from "svelte-sonner";

    let redirectHandle = 0
    let redirectDelay = $state(5)

    onMount(async () => {
        const urlParams = new URLSearchParams(window.location.search);
        const token = urlParams.get("token") ?? ""
        try {
            await confirmEmail(token)
            redirectHandle = setInterval(() => {
                if (--redirectDelay <= 0) {
                    clearInterval(redirectHandle)
                    goto("/auth/login")
                }
            }, 1000)
        } catch (err) {
            toast.error("Failed to verify your email, please try again")
            await goto("/")
        }
    })
    onDestroy(() => {
        clearInterval(redirectHandle)
    })
</script>

<Alert.Root class="max-w-xl m-auto">
    <SuccessIcon />
    <Alert.Title>Success! Your account has been confirmed</Alert.Title>
    <Alert.Description>
        Your email address has been successfully verified, you will be redirected
        to the login page in {redirectDelay} seconds, alternatively click on the
        button below.

        <Button href="/auth/login">Login</Button>
    </Alert.Description>
</Alert.Root>