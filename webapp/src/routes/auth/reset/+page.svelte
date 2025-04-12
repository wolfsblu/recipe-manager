<script lang="ts">
    import { Card, Button, Label, Input } from 'flowbite-svelte';
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

<Card class="self-center my-auto">
    <form class="flex flex-col space-y-6" onsubmit={onSubmit}>
        <h3 class="text-xl font-medium text-gray-900 dark:text-white">
            Reset your password
        </h3>
        <Label class="space-y-2">
            <span>Email</span>
            <Input type="email" name="email" placeholder="name@company.com" required bind:value={email} />
        </Label>
        <Button type="submit" class="w-full">Reset your password</Button>
    </form>
</Card>