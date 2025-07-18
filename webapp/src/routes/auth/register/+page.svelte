<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
    import * as Card from "$lib/components/ui/card/index.js";
    import * as Form from "$lib/components/ui/form/index.js";
    import { formSchema, type FormSchema } from "./schema";
    import { Input } from "$lib/components/ui/input/index.js";
    import {register} from "$lib/api/auth/user.svelte";
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
                    await register(form.data)
                    toast.success("Account registered, please verify your email")
                    await goto("/")
                } catch (error) {
                    toast.error("Failed to register account")
                }
                setMessage(form, "Registration submitted successfully");
            }
        }
    })

    const { form: formData, enhance } = form;
</script>

<Card.Root class="m-auto w-full max-w-sm">
    <Card.Header>
        <Card.Title class="text-2xl">Register</Card.Title>
        <Card.Description>Enter your data below to create a new account</Card.Description>
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
                                </div>
                                <Input {...props} type="password" bind:value={$formData.password} placeholder="●●●●●●●" />
                            {/snippet}
                        </Form.Control>
                        <Form.Description />
                        <Form.FieldErrors />
                    </Form.Field>
                </div>
                <Button type="submit" class="w-full">Register</Button>
            </div>
            <div class="mt-4 text-center text-sm">
                Already have an account?
                <a href="/auth/login" class="underline"> Login </a>
            </div>
        </form>
    </Card.Content>
</Card.Root>