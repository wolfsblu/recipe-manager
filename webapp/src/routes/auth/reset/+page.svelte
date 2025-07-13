<script lang="ts">
    import { resetPassword } from "$lib/auth/user.svelte";
    import { goto } from "$app/navigation";
    import { addToast } from "$lib/components/notifications/toasts";

    let email = $state('')

    const onSubmit = async (e: SubmitEvent) => {
        e.preventDefault()
        try {
            await resetPassword(email)
        } finally {
            addToast({ message: "Successfully requested password reset", type: "success" })
            await goto("/auth/login")
        }
    }
</script>

<p>Reset</p>