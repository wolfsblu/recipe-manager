<script lang="ts">
    import * as Form from '$lib/components/ui/form/index.js'
    import {Input} from "$lib/components/ui/input/index.js";
    import {Button} from "$lib/components/ui/button/index.js";
    import {Textarea} from "$lib/components/ui/textarea/index.js";
    import {formSchema, type FormSchema} from "./schema";
    import {type Infer, setMessage, superForm, type SuperValidated,} from "sveltekit-superforms";
    import {zodClient} from "sveltekit-superforms/adapters";
    import {toast} from "svelte-sonner";
    import {goto} from "$app/navigation";
    import {addRecipe} from "$lib/api/recipes/recipes.svelte";
    import { TagsInput } from '$lib/components/ui/tags-input';
    import ImageUpload from "$lib/components/recipes/ImageUpload.svelte";
    import UnitCombobox from "$lib/components/recipes/UnitCombobox.svelte";
    import { Separator } from "$lib/components/ui/separator/index.js";
    import PlusIcon from "@lucide/svelte/icons/plus";
    import TrashIcon from "@lucide/svelte/icons/trash";
    import IngredientCombobox from "$lib/components/recipes/IngredientCombobox.svelte";

    let {data}: { data: { form: SuperValidated<Infer<FormSchema>> } } = $props();
    const form = superForm(data.form, {
        SPA: true,
        resetForm: false,
        validators: zodClient(formSchema),
        async onUpdate({form}) {
            if (form.valid) {
                try {
                    await addRecipe(form.data)
                    toast.success("Recipe added successfully")
                    await goto("/recipes")
                } catch (error) {
                    toast.error("Failed to save recipe")
                }
                setMessage(form, "Recipe added successfully");
            }
        }
    })
    const {form: formData, enhance} = form;

    const addStep = () => {
        $formData.steps = [...$formData.steps, {
            ingredients: [{
                amount: 100,
                unit: 'kg',
                name: 'Potatoes',
            }]
        }];
    }
</script>

<form class="p-6" method="POST" use:enhance>
    <ImageUpload />
    <div class="flex flex-col lg:flex-row gap-x-6 gap-y-3 mt-3">
        <div class="w-full flex flex-col gap-3">
            <Form.Field {form} name="name">
                <Form.Control>
                    {#snippet children({props})}
                        <Form.Label>Name</Form.Label>
                        <Input {...props} bind:value={$formData.name} placeholder="My super tasty recipe"/>
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
                            <Form.Label>Time (min.)</Form.Label>
                            <Input {...props} bind:value={$formData.minutes} type="number" placeholder="45"/>
                        {/snippet}
                    </Form.Control>
                    <Form.Description/>
                    <Form.FieldErrors/>
                </Form.Field>
            </div>
            <Form.Field {form} name="tags">
                <Form.Control>
                    {#snippet children({props})}
                        <Form.Label>Tags</Form.Label>
                        <TagsInput {...props} bind:value={$formData.tags} placeholder="Add a tag" />
                    {/snippet}
                </Form.Control>
                <Form.Description/>
                <Form.FieldErrors/>
            </Form.Field>
        </div>

        <Form.Field class="w-full flex flex-col" {form} name="description">
            <Form.Control>
                {#snippet children({props})}
                    <Form.Label>Description</Form.Label>
                    <Textarea {...props}
                              class="flex-grow resize-none"
                              bind:value={$formData.description}
                              placeholder="Grandma's original pasta recipe handed down for generations."
                    />
                {/snippet}
            </Form.Control>
            <Form.Description class="sr-only">Test</Form.Description>
            <Form.FieldErrors/>
        </Form.Field>
    </div>

    <Separator class="my-3" orientation="horizontal" />

    <h1>Steps</h1>
        {#each $formData.steps as _, stepIndex}
            <Form.Fieldset {form} name="steps[{stepIndex}].ingredients">
                <Form.Legend>Ingredients</Form.Legend>
                {#each $formData.steps[stepIndex].ingredients as _, ingredientIndex}
                    <Form.ElementField {form} name="steps[{stepIndex}].ingredients[{ingredientIndex}].amount">
                        <Form.Control>
                            {#snippet children({props})}
                                <div class="flex flex-col md:flex-row gap-1">
                                    <Input {...props} class="w-auto" type="number" bind:value={$formData.steps[stepIndex].ingredients[ingredientIndex].amount} placeholder="1" />
                                </div>
                            {/snippet}
                        </Form.Control>
                        <Form.Description class="sr-only" />
                        <Form.FieldErrors/>
                    </Form.ElementField>
                    <Form.ElementField {form} name="steps[{stepIndex}].ingredients[{ingredientIndex}].unit">
                        <Form.Control>
                            {#snippet children({props})}
                                <div class="flex flex-col md:flex-row gap-1">
                                    <Input {...props} class="w-auto" bind:value={$formData.steps[stepIndex].ingredients[ingredientIndex].unit} placeholder="kg" />
                                </div>
                            {/snippet}
                        </Form.Control>
                        <Form.Description class="sr-only" />
                        <Form.FieldErrors/>
                    </Form.ElementField>
                    <Form.ElementField {form} name="steps[{stepIndex}].ingredients[{ingredientIndex}].name">
                        <Form.Control>
                            {#snippet children({props})}
                                <div class="flex flex-col md:flex-row gap-1">
                                    <Input {...props} class="w-auto" bind:value={$formData.steps[stepIndex].ingredients[ingredientIndex].name} placeholder="Potatoes" />
                                </div>
                            {/snippet}
                        </Form.Control>
                        <Form.Description class="sr-only" />
                        <Form.FieldErrors/>
                    </Form.ElementField>
                {/each}
                <Form.Description/>
                <Form.FieldErrors/>
            </Form.Fieldset>

        {/each}

    <div class="flex justify-between">
        <Button type="button" onclick={addStep}>Add Step</Button>
        <Button type="submit">Create</Button>
    </div>
</form>
