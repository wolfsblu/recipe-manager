<script lang="ts">
    import { Card, Button, Label, Input } from 'flowbite-svelte';
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
            addToast({ message: "Account registered successfully", type: "success" })
            await goto("/")
        } catch (e) {
            addToast({ message: "Failed to register account", type: "error" })
        }
    }
</script>

<Card class="self-center my-auto">
    <form class="flex flex-col space-y-6" on:submit={onRegister}>
        <h3 class="text-xl font-medium text-gray-900 dark:text-white">Register</h3>
        <Label class="space-y-2">
            <span>Email</span>
            <Input type="email" name="email" bind:value={credentials.email} placeholder="name@company.com" required />
        </Label>
        <Label class="space-y-2">
            <span>Your password</span>
            <Input type="password" name="password" bind:value={credentials.password} placeholder="•••••" required />
        </Label>
        <Button type="submit" class="w-full">Create a new account</Button>
        <div class="text-sm font-medium text-gray-500 dark:text-gray-300">
            Already have an account? <a href="/auth/login" class="text-primary-700 hover:underline dark:text-primary-500">Login</a>
        </div>
    </form>
</Card>