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
    import ArrowDownIcon from "@lucide/svelte/icons/arrow-down";
    import ArrowUpIcon from "@lucide/svelte/icons/arrow-up";
    import TrashIcon from "@lucide/svelte/icons/trash";
    import IngredientCombobox from "$lib/components/recipes/IngredientCombobox.svelte";
    import Ingredient from "$lib/components/recipes/Ingredient.svelte";

    let {data}: { data: { form: SuperValidated<Infer<FormSchema>> } } = $props();
    const form = superForm(data.form, {
        SPA: true,
        dataType: 'json',
        resetForm: false,
        validators: zodClient(formSchema),
        async onUpdate({form}) {
            console.log(form)
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
            ingredients: [
                {name: '', amount: '', unit: ''},
            ],
            instructions: '',
        }];
    }
    const moveStepUp = (stepIndex: number) => {
        if (stepIndex <= 0) {
            return;
        }
        const newSteps = [...$formData.steps];
        [newSteps[stepIndex - 1], newSteps[stepIndex]] = [newSteps[stepIndex], newSteps[stepIndex - 1]];
        $formData.steps = newSteps;
    }
    const moveStepDown = (stepIndex: number) => {
        if (stepIndex >= $formData.steps.length - 1) {
            return;
        }
        const newSteps = [...$formData.steps];
        [newSteps[stepIndex], newSteps[stepIndex + 1]] = [newSteps[stepIndex + 1], newSteps[stepIndex]];
        $formData.steps = newSteps;
    }
    const removeStepByIndex = (stepIndex: number) => {
        $formData.steps = $formData.steps.filter((_, i) => i !== stepIndex)
    }

    const addIngredient = (stepIndex: number) => {
        $formData.steps[stepIndex].ingredients = [
            ...$formData.steps[stepIndex].ingredients,
            {name: '', unit: '', amount: ''},
        ]
    }
    const removeIngredientByIndex = (stepIndex: number, ingredientIndex: number) => {
        $formData.steps[stepIndex].ingredients = $formData.steps[stepIndex].ingredients.filter((_, i) => i !== ingredientIndex)
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

    {#each $formData.steps as _, stepIndex}
        <div class="flex items-center justify-between">
            <h1>Step {stepIndex + 1}</h1>
            <div>
                <Button disabled={stepIndex === 0} variant="ghost" size="icon" class="size-8" onclick={() => moveStepUp(stepIndex)}>
                    <ArrowUpIcon />
                </Button>
                <Button disabled={stepIndex === $formData.steps.length - 1} variant="ghost" size="icon" class="size-8" onclick={() => moveStepDown(stepIndex)}>
                    <ArrowDownIcon />
                </Button>
                <Button variant="ghost" size="icon" class="size-8" onclick={() => removeStepByIndex(stepIndex)}>
                    <TrashIcon />
                </Button>
            </div>
        </div>
        <Separator class="mt-1 mb-2" orientation="horizontal" />
        <div class="grid lg:grid-cols-2 gap-x-6">
            <Form.Fieldset {form} name="steps[{stepIndex}].ingredients">
                {#if $formData.steps[stepIndex].ingredients.length === 0}
                    <Button onclick={() => addIngredient(stepIndex)}>
                        Add Ingredients
                    </Button>
                {:else}
                    <Form.Legend>Ingredients</Form.Legend>
                {/if}
                <div class="grid grid-cols-[1fr_2fr_3fr_max-content] gap-2">
                {#each $formData.steps[stepIndex].ingredients as _, ingredientIndex}
                    <Ingredient {form} {stepIndex} {ingredientIndex} />
                    <div>
                        <Button variant="destructive" onclick={() => removeIngredientByIndex(stepIndex, ingredientIndex)} type="button">
                            <TrashIcon />
                        </Button>
                        {#if ingredientIndex === $formData.steps[stepIndex].ingredients.length - 1}
                            <Button onclick={() => addIngredient(stepIndex)} type="button">
                                <PlusIcon />
                            </Button>
                        {/if}
                    </div>
                {/each}
                </div>
                <Form.FieldErrors/>
            </Form.Fieldset>
            <Form.Fieldset class="flex flex-col" {form} name="steps[{stepIndex}].instructions">
                <Form.Control>
                    {#snippet children({props})}
                        <Form.Label>Instructions</Form.Label>
                        <Textarea {...props}
                                  class="flex-grow"
                                  bind:value={$formData.steps[stepIndex].instructions}
                                  placeholder="Stir the potatoes."
                        />
                    {/snippet}
                </Form.Control>
                <Form.Description class="sr-only">Test</Form.Description>
                <Form.FieldErrors/>
            </Form.Fieldset>
        </div>
    {/each}

    <div class="flex justify-between">
        <Button type="reset">Cancel</Button>
        <div>
            <Button type="button" onclick={addStep}>Add Step</Button>
            {#if $formData.steps.length > 0}
                <Button type="submit">Create</Button>
            {/if}
        </div>
    </div>
</form>
