<script lang="ts">
    import { Label } from "$lib/components/ui/label/index.js";
    import { Input } from "$lib/components/ui/input/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import { Textarea } from "$lib/components/ui/textarea/index.js";
    import { formSchema, type FormSchema } from "./schema";
    import * as Form from "$lib/components/ui/form/index.js";
    import {
        type SuperValidated,
        type Infer,
        superForm,
        setMessage,
    } from "sveltekit-superforms";
    import {zodClient} from "sveltekit-superforms/adapters";

    let { data }: { data: { form: SuperValidated<Infer<FormSchema>> } } = $props();
    const form = superForm(data.form, {
        SPA: true,
        resetForm: false,
        validators: zodClient(formSchema),
        async onUpdate({ form }) {
            if (form.valid) {

            }
        }
    })
    const { form: formData, enhance } = form;
</script>

<div class="flex-grow grid lg:grid-cols-2">
    <div class="flex flex-col gap-4 p-6 md:p-10">
        <div class="flex flex-1 justify-center">
            <form class="w-full flex flex-col gap-6" method="POST" use:enhance>
                <div class="flex flex-col flex-grow gap-3">
                    <Form.Field {form} name="title">
                        <Form.Control>
                            {#snippet children({ props })}
                                <Form.Label>Title</Form.Label>
                                <Input {...props} bind:value={$formData.title} placeholder="My super tasty recipe" />
                            {/snippet}
                        </Form.Control>
                        <Form.Description />
                        <Form.FieldErrors />
                    </Form.Field>
                    <div class="grid grid-cols-2 gap-3">
                        <Form.Field {form} name="servings">
                            <Form.Control>
                                {#snippet children({ props })}
                                    <Form.Label>Servings</Form.Label>
                                    <Input {...props} bind:value={$formData.servings} type="number" placeholder="4" />
                                {/snippet}
                            </Form.Control>
                            <Form.Description />
                            <Form.FieldErrors />
                        </Form.Field>
                        <Form.Field {form} name="minutes">
                            <Form.Control>
                                {#snippet children({ props })}
                                    <Form.Label>Cooking Time (min)</Form.Label>
                                    <Input {...props} bind:value={$formData.minutes} type="number" placeholder="4" />
                                {/snippet}
                            </Form.Control>
                            <Form.Description />
                            <Form.FieldErrors />
                        </Form.Field>
                    </div>
                    <Form.Field {form} name="description" class="flex-grow flex flex-col">
                        <Form.Control>
                            {#snippet children({ props })}
                                <Form.Label>Description</Form.Label>
                                <Textarea {...props}
                                          class="flex-grow"
                                          bind:value={$formData.description}
                                          placeholder="Type your message here."
                                />
                            {/snippet}
                        </Form.Control>
                        <Form.Description />
                        <Form.FieldErrors />
                    </Form.Field>
                    <Button type="submit" class="w-full">Create</Button>
                </div>
            </form>
        </div>
    </div>
    <div class="bg-muted relative hidden lg:block">
        <img
                src="/placeholder.svg"
                alt="placeholder"
                class="absolute inset-0 h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"
        />
    </div>
</div>