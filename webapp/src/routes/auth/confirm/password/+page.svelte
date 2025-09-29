<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
    import * as Card from "$lib/components/ui/card/index.js";
    import * as Form from "$lib/components/ui/form/index.js";
    import { formSchema, type FormSchema } from "./schema";
    import { Input } from "$lib/components/ui/input/index.js";
    import {confirmPassword} from "$lib/api/auth/user.svelte";
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
                    const urlParams = new URLSearchParams(window.location.search);
                    const token = urlParams.get("token") ?? ""
                    await confirmPassword(form.data.password, token)
                    toast.success(m.auth_confirmPassword_successMessage())
                    await goto("/auth/login")
                } catch (e) {
                    toast.error(m.auth_confirmPassword_errorMessage())
                    await goto("/")
                }
                setMessage(form, "Reset performed successfully");
            }
        }
    })

    const { form: formData, enhance } = form;
</script>

<Card.Root class="m-auto w-full max-w-sm">
    <Card.Header>
        <Card.Title class="text-2xl">{m.auth_confirmPassword_title()}</Card.Title>
        <Card.Description>{m.auth_confirmPassword_description()}</Card.Description>
    </Card.Header>
    <Card.Content>
        <form method="POST" use:enhance>
            <div class="grid gap-4">
                <div class="grid gap-2">
                    <Form.Field {form} name="password">
                        <Form.Control>
                            {#snippet children({ props })}
                                <Form.Label>{m.auth_confirmPassword_password()}</Form.Label>
                                <Input {...props} bind:value={$formData.password} placeholder={m.auth_confirmPassword_passwordPlaceholder()} />
                            {/snippet}
                        </Form.Control>
                        <Form.Description />
                        <Form.FieldErrors />
                    </Form.Field>
                </div>
                <Button type="submit" class="w-full">{m.auth_confirmPassword_button()}</Button>
            </div>
        </form>
    </Card.Content>
</Card.Root>