<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
    import * as Alert from "$lib/components/ui/alert/index.js";
    import SuccessIcon from "@lucide/svelte/icons/check-circle-2";
    import {confirmEmail} from "$lib/api/auth/user.svelte";
    import {goto} from "$app/navigation";
    import {onDestroy, onMount} from "svelte";
    import {toast} from "svelte-sonner";
    import * as m from "$lib/paraglide/messages.js";

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
            toast.error(m.auth_confirmEmail_errorMessage())
            await goto("/")
        }
    })
    onDestroy(() => {
        clearInterval(redirectHandle)
    })
</script>

<Alert.Root class="max-w-xl m-auto">
    <SuccessIcon />
    <Alert.Title>{m.auth_confirmEmail_successTitle()}</Alert.Title>
    <Alert.Description>
        {m.auth_confirmEmail_successDescription({ seconds: redirectDelay })}

        <Button href="/auth/login">{m.auth_confirmEmail_loginButton()}</Button>
    </Alert.Description>
</Alert.Root>