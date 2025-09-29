<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
    import * as Card from "$lib/components/ui/card/index.js";
    import * as Form from "$lib/components/ui/form/index.js";
    import { formSchema, type FormSchema } from "./schema";
    import { Input } from "$lib/components/ui/input/index.js";
    import {register, resetPassword} from "$lib/api/auth/user.svelte";
    import {goto} from "$app/navigation";
    import { toast } from "svelte-sonner";
    import {
        type SuperValidated,
        type Infer,
        superForm,
        setMessage,
    } from "sveltekit-superforms";
    import { zodClient } from "sveltekit-superforms/adapters";
    import * as m from "$lib/paraglide/messages.js";

    let { data }: { data: { form: SuperValidated<Infer<FormSchema>> } } = $props();

    const form = superForm(data.form, {
        SPA: true,
        resetForm: false,
        validators: zodClient(formSchema),
        async onUpdate({ form }) {
            if (form.valid) {
                try {
                    await resetPassword(email)
                } finally {
                    toast.success(m.auth_reset_successMessage())
                    await goto("/auth/login")
                }
                setMessage(form, "Reset requested successfully");
            }
        }
    })

    const { form: formData, enhance } = form;
</script>

<Card.Root class="m-auto w-full max-w-sm">
    <Card.Header>
        <Card.Title class="text-2xl">{m.auth_reset_title()}</Card.Title>
        <Card.Description>{m.auth_reset_description()}</Card.Description>
    </Card.Header>
    <Card.Content>
        <form method="POST" use:enhance>
            <div class="grid gap-4">
                <div class="grid gap-2">
                    <Form.Field {form} name="email">
                        <Form.Control>
                            {#snippet children({ props })}
                                <Form.Label>{m.auth_reset_email()}</Form.Label>
                                <Input {...props} bind:value={$formData.email} placeholder={m.auth_reset_emailPlaceholder()} />
                            {/snippet}
                        </Form.Control>
                        <Form.Description />
                        <Form.FieldErrors />
                    </Form.Field>
                </div>
                <Button type="submit" class="w-full">{m.auth_reset_button()}</Button>
            </div>
            <div class="mt-4 text-center text-sm">
                {m.auth_reset_haveAccount()}
                <a href="/auth/login" class="underline">{m.auth_register_loginLink()}</a> {m.auth_reset_orText()}
                <a href="/auth/register" class="underline">{m.auth_reset_createAccount()}</a>
            </div>
        </form>
    </Card.Content>
</Card.Root>