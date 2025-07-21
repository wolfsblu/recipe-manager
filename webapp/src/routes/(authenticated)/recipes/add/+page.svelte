<script lang="ts">
    import * as Form from '$lib/components/ui/form/index.js'
    import {Input} from "$lib/components/ui/input/index.js";
    import {Button} from "$lib/components/ui/button/index.js";
    import {Textarea} from "$lib/components/ui/textarea/index.js";
    import {formSchema, type FormSchema} from "./schema";
    import {type Infer, superForm, type SuperValidated,} from "sveltekit-superforms";
    import {zodClient} from "sveltekit-superforms/adapters";

    let {data}: { data: { form: SuperValidated<Infer<FormSchema>> } } = $props();
    const form = superForm(data.form, {
        SPA: true,
        resetForm: false,
        validators: zodClient(formSchema),
        async onUpdate({form}) {
            if (form.valid) {

            }
        }
    })
    const {form: formData, enhance} = form;
</script>

<form class="flex flex-col flex-grow p-6" method="POST" use:enhance>
    <div class="flex-grow flex flex-col lg:flex-row gap-6">
        <div class="w-full flex flex-col gap-3">
            <p>
                Images
            </p>
            <Form.Field {form} name="title">
                <Form.Control>
                    {#snippet children({props})}
                        <Form.Label>Title</Form.Label>
                        <Input {...props} bind:value={$formData.title} placeholder="My super tasty recipe"/>
                    {/snippet}
                </Form.Control>
                <Form.Description/>
                <Form.FieldErrors/>
            </Form.Field>
            <div class="grid grid-cols-2 gap-3">
                <Form.Field {form} name="servings">
                    <Form.Control>
                        {#snippet children({props})}
                            <Form.Label>Servings</Form.Label>
                            <Input {...props} bind:value={$formData.servings} type="number" placeholder="4"/>
                        {/snippet}
                    </Form.Control>
                    <Form.Description/>
                    <Form.FieldErrors/>
                </Form.Field>
                <Form.Field {form} name="minutes">
                    <Form.Control>
                        {#snippet children({props})}
                            <Form.Label>Cooking Time (min)</Form.Label>
                            <Input {...props} bind:value={$formData.minutes} type="number" placeholder="4"/>
                        {/snippet}
                    </Form.Control>
                    <Form.Description/>
                    <Form.FieldErrors/>
                </Form.Field>
            </div>
            <p>
                Ingredients
            </p>
            <p>
                Tags
            </p>
        </div>
        <Form.Field class="w-full flex-grow flex flex-col" {form} name="description">
            <Form.Control>
                {#snippet children({props})}
                    <Form.Label>Description</Form.Label>
                    <Textarea {...props}
                              class="flex-grow resize-none"
                              bind:value={$formData.description}
                              placeholder="Type your message here."
                    />
                {/snippet}
            </Form.Control>
            <Form.Description/>
            <Form.FieldErrors/>
        </Form.Field>
    </div>
    <Button class="w-full" type="submit">Create</Button>

</form>
