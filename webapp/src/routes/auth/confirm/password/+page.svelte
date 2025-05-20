<script lang="ts">
    import { Card, Button, Label, Input } from 'flowbite-svelte';
    import {confirmPassword} from "$lib/auth/user.svelte.js";
    import {addToast} from "$lib/components/notifications/toasts";
    import {goto} from "$app/navigation";

    let newPassword = $state("")

    const onReset = async (e: SubmitEvent) => {
        e.preventDefault()
        try {
            const urlParams = new URLSearchParams(window.location.search);
            const token = urlParams.get("token") ?? ""
            await confirmPassword(newPassword, token)
            addToast({ message: "Password was reset successfully", type: "success", group: "password" })
            await goto("/auth/login")
        } catch (e) {
            addToast({ message: "Failed to reset password", type: "error", group: "password" })
        }
    }
</script>

<Card class="self-center my-auto p-6">
    <form class="flex flex-col space-y-6" onsubmit={onReset}>
        <h3 class="text-xl font-medium text-gray-900 dark:text-white">Register</h3>
        <Label class="space-y-2">
            <span>New password</span>
            <Input type="password" name="password" bind:value={newPassword} placeholder="•••••" required />
        </Label>
        <Button type="submit" class="w-full">Reset password</Button>
    </form>
</Card>