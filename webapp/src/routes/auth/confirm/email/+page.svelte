<script lang="ts">
    import {Alert, Button, Card} from "flowbite-svelte";
    import {CheckOutline} from "flowbite-svelte-icons";
    import {onDestroy, onMount} from "svelte";
    import {goto} from "$app/navigation";
    import {confirmEmail} from "$lib/auth/user.svelte.js";

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
        } catch (err) {}
    })

    onDestroy(() => {
        clearInterval(redirectHandle)
    })
</script>

<Card class="self-center my-auto" padding="none" size="lg">
    <Alert color="green">
        <div class="flex items-center gap-3">
            <CheckOutline class="w-5 h-5" />
            <span class="text-lg font-medium">
                Your account has been confirmed
            </span>
        </div>
        <p class="mt-2">
            Your email address has been successfully verified, you will be redirected
            to the login page in {redirectDelay} seconds, alternatively click on the button below
        </p>
        <div class="mt-4">
            <Button class="w-full" size="xs" color="green" href="/auth/login">Go to login</Button>
        </div>
    </Alert>
</Card>