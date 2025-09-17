<script lang="ts">
    import type { SuperForm, SuperValidated } from 'sveltekit-superforms';
    import { Card, CardContent, CardHeader, CardTitle } from "$lib/components/ui/card/index.js";
    import { Badge } from "$lib/components/ui/badge/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import { Input } from "$lib/components/ui/input/index.js";
    import { Textarea } from "$lib/components/ui/textarea/index.js";
    import { Separator } from "$lib/components/ui/separator/index.js";
    import * as Form from '$lib/components/ui/form/index.js';
    import { TagsInput } from '$lib/components/ui/tags-input';
    import ImageUpload from "$lib/components/recipes/ImageUpload.svelte";
    import Combobox from "$lib/components/recipes/Combobox.svelte";

    import ClockIcon from '@lucide/svelte/icons/clock';
    import UsersIcon from '@lucide/svelte/icons/users';
    import ChefHatIcon from '@lucide/svelte/icons/chef-hat';
    import ListIcon from '@lucide/svelte/icons/list';
    import SaladIcon from '@lucide/svelte/icons/salad';
    import PlusIcon from '@lucide/svelte/icons/plus';
    import ArrowUpIcon from '@lucide/svelte/icons/arrow-up';
    import ArrowDownIcon from '@lucide/svelte/icons/arrow-down';
    import TrashIcon from '@lucide/svelte/icons/trash';
    import SaveIcon from '@lucide/svelte/icons/save';

    import type { formSchema } from "../../routes/(authenticated)/recipes/add/schema";

    // Props
    interface Props {
        form: SuperForm<typeof formSchema>;
        ingredients: Array<{ id: number; name: string }>;
        units: Array<{ id: number; name: string }>;
        tags: Array<{ id: number; name: string }>;
        isEditing?: boolean;
        submitText?: string;
        onCancel?: () => void;
    }

    let { 
        form, 
        ingredients, 
        units, 
        tags, 
        isEditing = false, 
        submitText = isEditing ? "Update Recipe" : "Create Recipe",
        onCancel
    }: Props = $props();

    const { form: formData, enhance } = form;

    // Step management
    const addStep = () => {
        $formData.steps = [...$formData.steps, {
            ingredients: [{ amount: 0, ingredientId: 0, unitId: 0 }],
            instructions: '',
        }];
    };

    const moveStep = (stepIndex: number, direction: 'up' | 'down') => {
        const offset = direction === 'up' ? -1 : 1;
        const targetIndex = stepIndex + offset;

        if (targetIndex < 0 || targetIndex >= $formData.steps.length) {
            return;
        }

        const newSteps = [...$formData.steps];
        [newSteps[stepIndex], newSteps[targetIndex]] = [newSteps[targetIndex], newSteps[stepIndex]];
        $formData.steps = newSteps;
    };

    const removeStep = (stepIndex: number) => {
        $formData.steps = $formData.steps.filter((_, i) => i !== stepIndex);
    };

    // Ingredient management
    const addIngredient = (stepIndex: number) => {
        $formData.steps[stepIndex].ingredients = [
            ...$formData.steps[stepIndex].ingredients,
            { amount: 0, ingredientId: 0, unitId: 0 },
        ];
    };

    const removeIngredient = (stepIndex: number, ingredientIndex: number) => {
        $formData.steps[stepIndex].ingredients = $formData.steps[stepIndex].ingredients.filter((_, i) => i !== ingredientIndex);
    };

</script>

<form method="POST" use:enhance class="space-y-8">
    <!-- Recipe Header Section -->
    <Card class="border-0 md:border shadow-none md:shadow-sm">
        <CardHeader class="px-4 md:px-6">
            <CardTitle class="flex items-center gap-2">
                <ChefHatIcon class="w-5 h-5" />
                Recipe Details
            </CardTitle>
        </CardHeader>
        <CardContent class="space-y-6 px-4 md:px-6">
            <!-- Image Upload -->
            <div>
                <label class="text-sm font-medium text-foreground mb-3 block">Recipe Images</label>
                <ImageUpload bind:value={$formData.images} deferDeletion={isEditing} />
            </div>

            <div class="grid lg:grid-cols-2 gap-6">
                <!-- Left Column -->
                <div class="space-y-4">
                    <!-- Recipe Name -->
                    <Form.Field {form} name="name">
                        <Form.Control>
                            {#snippet children({ props })}
                                <Form.Label>Recipe Name</Form.Label>
                                <Input
                                    {...props}
                                    bind:value={$formData.name}
                                    placeholder="My amazing recipe"
                                />
                            {/snippet}
                        </Form.Control>
                        <Form.FieldErrors />
                    </Form.Field>

                    <!-- Recipe Stats -->
                    <div class="grid grid-cols-2 gap-4">
                        <Form.Field {form} name="servings">
                            <Form.Control>
                                {#snippet children({ props })}
                                    <Form.Label class="flex items-center gap-2">
                                        <UsersIcon class="w-4 h-4" />
                                        Servings
                                    </Form.Label>
                                    <Input
                                        {...props}
                                        bind:value={$formData.servings}
                                        type="number"
                                        placeholder="4"
                                        min="1"
                                    />
                                {/snippet}
                            </Form.Control>
                            <Form.FieldErrors />
                        </Form.Field>

                        <Form.Field {form} name="minutes">
                            <Form.Control>
                                {#snippet children({ props })}
                                    <Form.Label class="flex items-center gap-2">
                                        <ClockIcon class="w-4 h-4" />
                                        Cook Time (min)
                                    </Form.Label>
                                    <Input
                                        {...props}
                                        bind:value={$formData.minutes}
                                        type="number"
                                        placeholder="30"
                                        min="1"
                                    />
                                {/snippet}
                            </Form.Control>
                            <Form.FieldErrors />
                        </Form.Field>
                    </div>

                    <!-- Tags -->
                    <Form.Field {form} name="tags">
                        <Form.Control>
                            {#snippet children({ props })}
                                <Form.Label>Tags</Form.Label>
                                <TagsInput
                                    {...props}
                                    bind:value={$formData.tags}
                                    placeholder="Add a tag"
                                    suggestions={tags?.map(tag => ({ id: tag.id, label: tag.name })) || []}
                                />
                            {/snippet}
                        </Form.Control>
                        <Form.FieldErrors />
                    </Form.Field>
                </div>

                <!-- Right Column -->
                <div>
                    <Form.Field {form} name="description">
                        <Form.Control>
                            {#snippet children({ props })}
                                <Form.Label>Description</Form.Label>
                                <Textarea
                                    {...props}
                                    bind:value={$formData.description}
                                    placeholder="Tell us about your recipe..."
                                    class="min-h-[200px] resize-none"
                                />
                            {/snippet}
                        </Form.Control>
                        <Form.FieldErrors />
                    </Form.Field>
                </div>
            </div>

        </CardContent>
    </Card>

    <!-- Instructions Section -->
    {#if $formData.steps.length === 0}
        <!-- Empty State Card -->
        <Card class="border-0 md:border shadow-none md:shadow-sm">
            <CardHeader class="px-4 md:px-6">
                <CardTitle class="flex items-center justify-between">
                    <div class="flex items-center gap-2">
                        <ListIcon class="w-5 h-5" />
                        Cooking Instructions
                    </div>
                    <Button type="button" onclick={addStep} variant="outline" size="sm">
                        <PlusIcon class="w-4 h-4" />
                        Add Step
                    </Button>
                </CardTitle>
            </CardHeader>
            <CardContent class="px-4 md:px-6">
                <div class="text-center py-12 text-muted-foreground">
                    <ChefHatIcon class="w-12 h-12 mx-auto mb-4" />
                    <p class="text-lg mb-2">No cooking steps yet</p>
                    <p class="text-sm">Add your first step to get started</p>
                </div>
            </CardContent>
        </Card>
    {:else}
        <!-- Steps as Individual Cards -->
        <div class="space-y-6">
            <div class="flex items-center justify-between px-4 md:px-0">
                <h2 class="text-2xl font-bold flex items-center gap-2">
                    <ListIcon class="w-6 h-6" />
                    Cooking Instructions
                </h2>
            </div>

            {#each $formData.steps as _, stepIndex}
                <Card class="overflow-hidden border-0 md:border shadow-none md:shadow-sm">
                    <CardHeader class="bg-muted/50 gap-0 py-1 px-4 md:px-6">
                        <div class="flex items-center justify-between">
                            <CardTitle class="flex items-center gap-3">
                                <div class="w-8 h-8 bg-primary text-primary-foreground rounded-full flex items-center justify-center text-sm font-bold">
                                    {stepIndex + 1}
                                </div>
                                Step {stepIndex + 1}
                            </CardTitle>
                            <div class="flex items-center gap-1">
                                <Button
                                    type="button"
                                    variant="ghost"
                                    size="sm"
                                    onclick={() => moveStep(stepIndex, 'up')}
                                    disabled={stepIndex === 0}
                                >
                                    <ArrowUpIcon class="w-4 h-4" />
                                </Button>
                                <Button
                                    type="button"
                                    variant="ghost"
                                    size="sm"
                                    onclick={() => moveStep(stepIndex, 'down')}
                                    disabled={stepIndex === $formData.steps.length - 1}
                                >
                                    <ArrowDownIcon class="w-4 h-4" />
                                </Button>
                                <Button
                                    type="button"
                                    variant="ghost"
                                    size="sm"
                                    onclick={() => removeStep(stepIndex)}
                                >
                                    <TrashIcon class="w-4 h-4 text-destructive" />
                                </Button>
                            </div>
                        </div>
                    </CardHeader>
                    <CardContent class="px-4 md:px-6">
                        <div class="grid lg:grid-cols-2 gap-6">
                            <!-- Ingredients -->
                            <div>
                                <h4 class="font-semibold mb-3 flex items-center gap-2">
                                    <SaladIcon class="w-4 h-4" />
                                    Ingredients
                                </h4>
                                <Form.Fieldset {form} name="steps[{stepIndex}].ingredients">
                                    <div class="space-y-3">
                                        {#each $formData.steps[stepIndex].ingredients as _, ingredientIndex}
                                            <div class="space-y-2 md:space-y-0 md:grid md:grid-cols-[2fr_2fr_3fr_auto] md:gap-2 md:items-start">
                                                <div class="grid grid-cols-[1fr_2fr] gap-2 md:contents">
                                                    <Form.ElementField class="space-y-0" {form} name="steps[{stepIndex}].ingredients[{ingredientIndex}].amount">
                                                        <Form.Control>
                                                            {#snippet children({ props })}
                                                                <Input
                                                                    {...props}
                                                                    type="number"
                                                                    bind:value={$formData.steps[stepIndex].ingredients[ingredientIndex].amount}
                                                                    placeholder="1"
                                                                    class="text-sm"
                                                                />
                                                            {/snippet}
                                                        </Form.Control>
                                                        <Form.FieldErrors />
                                                    </Form.ElementField>

                                                    <Combobox
                                                        {form}
                                                        name="steps[{stepIndex}].ingredients[{ingredientIndex}].unitId"
                                                        options={units?.map(unit => ({ value: unit.id, label: unit.name })) || []}
                                                        bind:value={$formData.steps[stepIndex].ingredients[ingredientIndex].unitId}
                                                    />
                                                </div>

                                                <Combobox
                                                    {form}
                                                    name="steps[{stepIndex}].ingredients[{ingredientIndex}].ingredientId"
                                                    options={ingredients?.map(ingredient => ({ value: ingredient.id, label: ingredient.name })) || []}
                                                    bind:value={$formData.steps[stepIndex].ingredients[ingredientIndex].ingredientId}
                                                />

                                                <div class="md:flex md:justify-center md:items-start">
                                                    <Button
                                                        type="button"
                                                        variant="ghost"
                                                        onclick={() => removeIngredient(stepIndex, ingredientIndex)}
                                                        class="w-full md:w-auto md:h-9 md:px-3"
                                                    >
                                                        <TrashIcon class="w-4 h-4 text-destructive md:mr-0 mr-2" />
                                                        <span class="md:hidden">Delete</span>
                                                    </Button>
                                                </div>
                                            </div>
                                        {/each}

                                        <Button
                                            type="button"
                                            variant="outline"
                                            size="sm"
                                            onclick={() => addIngredient(stepIndex)}
                                            class="w-full"
                                        >
                                            <PlusIcon class="w-4 h-4" />
                                            Add Ingredient
                                        </Button>
                                    </div>
                                    <Form.FieldErrors />
                                </Form.Fieldset>
                            </div>

                            <!-- Instructions -->
                            <div>
                                <h4 class="font-semibold mb-3 flex items-center gap-2">
                                    <ChefHatIcon class="w-4 h-4" />
                                    Instructions
                                </h4>
                                <Form.Field {form} name="steps[{stepIndex}].instructions">
                                    <Form.Control>
                                        {#snippet children({ props })}
                                            <Textarea
                                                {...props}
                                                bind:value={$formData.steps[stepIndex].instructions}
                                                placeholder="Describe what to do in this step..."
                                                class="min-h-[120px] resize-none"
                                            />
                                        {/snippet}
                                    </Form.Control>
                                    <Form.FieldErrors />
                                </Form.Field>
                            </div>
                        </div>
                    </CardContent>
                </Card>
            {/each}

            <!-- Add Step Button at Bottom -->
            <div class="flex justify-center pt-2">
                <Button type="button" onclick={addStep} variant="outline">
                    <PlusIcon class="w-4 h-4" />
                    Add Another Step
                </Button>
            </div>
        </div>
    {/if}

    <!-- Form Actions -->
    <div class="flex justify-between items-center px-4 md:px-0">
        {#if onCancel}
            <Button type="button" variant="outline" onclick={onCancel}>
                Cancel
            </Button>
        {:else}
            <Button type="reset" variant="outline">
                Reset
            </Button>
        {/if}

        <Button type="submit" disabled={$formData.steps.length === 0}>
            {#if isEditing}
                <SaveIcon class="w-4 h-4" />
            {:else}
                <ChefHatIcon class="w-4 h-4" />
            {/if}
            {submitText}
        </Button>
    </div>
</form>