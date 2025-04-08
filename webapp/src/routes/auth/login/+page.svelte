<script lang="ts">
    import { Card, Button, Label, Input, Checkbox } from 'flowbite-svelte';
    import { login } from "$lib/auth/user.svelte";
    import { addToast } from "$lib/components/notifications/toasts";

    const credentials = {
        email: '',
        password: '',
    }

    const onSubmit = async (e: SubmitEvent) => {
        e.preventDefault()
        try {
            await login(credentials)
        } catch (error) {
            addToast({ message: "Failed to login, please verify your credentials.", type: "error" })
        }
    }
</script>

<Card class="self-center my-auto">
    <form class="flex flex-col space-y-6" on:submit="{onSubmit}">
        <h3 class="text-xl font-medium text-gray-900 dark:text-white">Sign in to our platform</h3>
        <Label class="space-y-2">
            <span>Email</span>
            <Input type="email" name="email" placeholder="name@company.com" required bind:value={credentials.email} />
        </Label>
        <Label class="space-y-2">
            <span>Your password</span>
            <Input type="password" name="password" placeholder="•••••" required bind:value={credentials.password} />
        </Label>
        <div class="flex items-start">
            <Checkbox>Remember me</Checkbox>
            <a href="/" class="ms-auto text-sm text-primary-700 hover:underline dark:text-primary-500"> Lost password? </a>
        </div>
        <Button type="submit" class="w-full">Login to your account</Button>
        <div class="text-sm font-medium text-gray-500 dark:text-gray-300">
            Not registered? <a href="/" class="text-primary-700 hover:underline dark:text-primary-500"> Create account </a>
        </div>
    </form>
</Card>