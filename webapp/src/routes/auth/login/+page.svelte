<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
    import * as Card from "$lib/components/ui/card/index.js";
    import * as Form from "$lib/components/ui/form/index.js";
    import { formSchema, type FormSchema } from "./schema";
    import { Input } from "$lib/components/ui/input/index.js";
    import {login} from "$lib/api/auth/user.svelte";
    import {goto} from "$app/navigation";
    import { toast } from "svelte-sonner";
    import {
        type SuperValidated,
        type Infer,
        superForm,
        setMessage,
    } from "sveltekit-superforms";
    import { zodClient } from "sveltekit-superforms/adapters";

    let { data }: { data: { form: SuperValidated<Infer<FormSchema>> } } = $props();

    const form = superForm(data.form, {
        SPA: true,
        resetForm: false,
        validators: zodClient(formSchema),
        async onUpdate({ form }) {
            if (form.valid) {
                try {
                    await login(form.data)
                    toast.success("Signed in successfully")
                    await goto("/")
                } catch (error) {
                    toast.error("Failed to login, please verify your credentials")
                }
                setMessage(form, "Login successful");
            }
        }
    })

    const { form: formData, enhance } = form;
</script>

<Card.Root class="m-auto w-full max-w-sm">
    <Card.Header>
        <Card.Title class="text-2xl">Login</Card.Title>
        <Card.Description>Enter your email below to login to your account</Card.Description>
    </Card.Header>
    <Card.Content>
        <form method="POST" use:enhance>
            <div class="grid gap-4">
                <div class="grid gap-2">
                    <Form.Field {form} name="email">
                        <Form.Control>
                            {#snippet children({ props })}
                                <Form.Label>Email</Form.Label>
                                <Input {...props} bind:value={$formData.email} placeholder="mail@example.com" />
                            {/snippet}
                        </Form.Control>
                        <Form.Description />
                        <Form.FieldErrors />
                    </Form.Field>
                </div>
                <div class="grid gap-2">
                    <Form.Field {form} name="password">
                        <Form.Control>
                            {#snippet children({ props })}
                                <div class="flex items-center">
                                    <Form.Label>Password</Form.Label>
                                    <a href="/auth/reset" class="ml-auto inline-block text-sm underline">
                                        Forgot password?
                                    </a>
                                </div>
                                <Input {...props} type="password" bind:value={$formData.password} placeholder="●●●●●●●" />
                            {/snippet}
                        </Form.Control>
                        <Form.Description />
                        <Form.FieldErrors />
                    </Form.Field>
                </div>
                <Button type="submit" class="w-full">Login</Button>
            </div>
            <div class="mt-4 text-center text-sm">
                Don't have an account?
                <a href="/auth/register" class="underline"> Sign up </a>
            </div>
        </form>
    </Card.Content>
</Card.Root>