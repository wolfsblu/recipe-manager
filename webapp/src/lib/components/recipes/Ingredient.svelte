<script>
    import * as Form from '$lib/components/ui/form/index.js'
    import {Input} from "$lib/components/ui/input/index.js";
    import Combobox from "$lib/components/recipes/Combobox.svelte";

    let {
        form,
        stepIndex,
        ingredients,
        ingredientIndex,
        units
    } = $props();

    const { form: formData } = form;

    const ingredientOptions = ingredients.map(
        ingredient => ({
            value: ingredient.id,
            label: ingredient.name,
        })
    )
    const unitOptions = units.map(
        unit => ({
            value: unit.id,
            label: unit.name,
        })
    )

    const basePath = `steps[${stepIndex}].ingredients[${ingredientIndex}]`;
</script>

<Form.ElementField class="space-y-0" {form} name="{basePath}.amount">
    <Form.Control>
        {#snippet children({props})}
            <Input {...props} class="w-full" type="number" bind:value={$formData.steps[stepIndex].ingredients[ingredientIndex].amount} placeholder="1" />
        {/snippet}
    </Form.Control>
    <Form.Description class="sr-only" />
    <Form.FieldErrors/>
</Form.ElementField>

<Combobox {form}
          name="{basePath}.unitId"
          bind:value={$formData.steps[stepIndex].ingredients[ingredientIndex].unitId}
/>

<Combobox {form}
          name="{basePath}.ingredientId"
          bind:value={$formData.steps[stepIndex].ingredients[ingredientIndex].ingredientId}
/>
