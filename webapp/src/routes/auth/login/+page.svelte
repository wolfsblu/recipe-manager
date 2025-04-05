<script lang="ts">
    import { Card, Button, Label, Input, Checkbox, Toast, Progressbar } from 'flowbite-svelte';
    import { ExclamationCircleOutline } from "flowbite-svelte-icons";
    import { login } from "$lib/auth/user.svelte";
    import { blur } from 'svelte/transition'
    import { writable } from "svelte/store";
    import { linear } from "svelte/easing";

    interface Toast {
        id: number
        message: string
        timeout: number
        visible: boolean
        progress: number
    }

    const toasts = writable<Toast[]>([]);

    const dismissToast = (id: number) => {
        toasts.update((all) => {
            const idx = all.findIndex(t => t.id === id)
            if (idx !== -1) {
                all[idx].visible = false;
            }
            return [...all];
        })
    }
    const addToast = (toast: Toast) => {
        toasts.update((all) => [toast, ...all])
        setTimeout(() => dismissToast(toast.id), toast.timeout)
    }

    const credentials = {
        email: '',
        password: '',
    }

    const onSubmit = async (e: SubmitEvent) => {
        e.preventDefault()
        try {
            await login(credentials)
        } catch (error) {
            addToast({
                id: Math.floor(Math.random() * 1000),
                message: "Failed to login, please verify your credentials.",
                timeout: 5000,
                progress: 100,
                visible: true
            })
        }
    }
</script>

<section class="absolute top-5 end-5">
    {#each $toasts as toast (toast.id)}
        <Toast class="min-w-80" color="red" transition={blur} params={{ amount: 10 }} bind:toastStatus={toast.visible}>
            <svelte:fragment slot="icon">
                <ExclamationCircleOutline class="w-5 h-5" />
                <span class="sr-only">Error icon</span>
            </svelte:fragment>
            {toast.message}
        </Toast>
        {#if toast.visible}
            <Progressbar animate class="mb-3" bind:progress={toast.progress} size="h-0.5" easing={linear} tweenDuration={toast.timeout} />
        {/if}
    {/each}
</section>

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