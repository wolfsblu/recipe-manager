<script lang="ts">
    import {register} from "$lib/auth/user.svelte.js";
    import {addToast} from "$lib/components/notifications/toasts";
    import {goto} from "$app/navigation";

    const credentials = {
        email: '',
        password: '',
    }

    const onRegister = async (e: SubmitEvent) => {
        e.preventDefault()
        try {
            await register(credentials)
            addToast({ message: "Account registered, please verify your email", type: "success", group: "register" })
            await goto("/")
        } catch (e) {
            addToast({ message: "Failed to register account", type: "error", group: "register" })
        }
    }
</script>


<p>Register</p>