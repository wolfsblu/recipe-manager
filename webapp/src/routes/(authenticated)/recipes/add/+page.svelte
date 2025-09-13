<script lang="ts">
    import type { PageProps } from './$types';
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
    import Ingredient from "$lib/components/recipes/Ingredient.svelte";

    let { data }: PageProps = $props();

    const form = superForm(data.form, {
        SPA: true,
        dataType: 'json',
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
            ingredients: [
                {amount: 0, ingredientId: 0, unitId: 0},
            ],
            instructions: '',
        }];
    }

    const moveStep = (stepIndex: number, direction: 'up' | 'down') => {
        const offset = direction === 'up' ? -1 : 1;
        const targetIndex = stepIndex + offset;

        if (targetIndex < 0 || targetIndex >= $formData.steps.length) {
            return
        }

        const newSteps = [...$formData.steps];
        [newSteps[stepIndex], newSteps[targetIndex]] = [newSteps[targetIndex], newSteps[stepIndex]];
        $formData.steps = newSteps;
    }

    const removeStepByIndex = (stepIndex: number) => {
        $formData.steps = $formData.steps.filter((_, i) => i !== stepIndex)
    }

    const addIngredient = (stepIndex: number) => {
        $formData.steps[stepIndex].ingredients = [
            ...$formData.steps[stepIndex].ingredients,
            {amount: 0, ingredientId: 0, unitId: 0},
        ]
    }
    const removeIngredientByIndex = (stepIndex: number, ingredientIndex: number) => {
        $formData.steps[stepIndex].ingredients = $formData.steps[stepIndex].ingredients.filter((_, i) => i !== ingredientIndex)
    }
</script>

<form class="p-6" method="POST" use:enhance>
    <ImageUpload />
    <div class="flex flex-col lg:flex-row gap-x-3 gap-y-1 mt-3">
        <div class="w-full flex flex-col">
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
                <Button disabled={stepIndex === 0} variant="ghost" size="icon" class="size-8" onclick={() => moveStep(stepIndex, 'up')}>
                    <ArrowUpIcon />
                </Button>
                <Button disabled={stepIndex === $formData.steps.length - 1} variant="ghost" size="icon" class="size-8" onclick={() => moveStep(stepIndex, 'down')}>
                    <ArrowDownIcon />
                </Button>
                <Button variant="ghost" size="icon" class="size-8" onclick={() => removeStepByIndex(stepIndex)}>
                    <TrashIcon class="stroke-red-600 dark:stroke-red-400" />
                </Button>
            </div>
        </div>
        <Separator class="mt-1 mb-2" orientation="horizontal" />
        <div class="grid lg:grid-cols-2 gap-x-3">
            <Form.Fieldset {form} name="steps[{stepIndex}].ingredients">
                <Form.Legend>Ingredients</Form.Legend>
                <div class="grid grid-cols-[1fr_1fr_2fr_max-content] md:grid-cols-[1fr_2fr_3fr_max-content] gap-1">
                {#each $formData.steps[stepIndex].ingredients as _, ingredientIndex}
                    <Ingredient {form} {stepIndex} {ingredientIndex} ingredients={data.ingredients} units={data.units} />
                    <div>
                        <Button variant="outline" onclick={() => removeIngredientByIndex(stepIndex, ingredientIndex)} type="button">
                            <TrashIcon class="stroke-red-600 dark:stroke-red-400" />
                        </Button>
                    </div>
                {/each}
                    <Button variant="outline" class="col-span-4" onclick={() => addIngredient(stepIndex)}>
                        <PlusIcon/> Add Ingredient
                    </Button>
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
