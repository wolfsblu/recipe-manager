<script>
    import * as Form from '$lib/components/ui/form/index.js'
    import {Input} from "$lib/components/ui/input/index.js";
    import Combobox from "$lib/components/recipes/Combobox.svelte";

    let {
        form,
        stepIndex,
        ingredientIndex,
    } = $props();

    const { form: formData } = form;

    const ingredients = [
        { value: "greensalad", label: "Green Salad" },
        { value: "carrots", label: "Carrots" },
        { value: "chicken", label: "Chicken" },
        { value: "salt", label: "Salt" },
        { value: "pepper", label: "Pepper" }
    ];

    const units = [
        { value: 'kg', label: 'Kilogram' },
        { value: 'mg', label: 'Milligram' },
        { value: 'l', label: 'Litre' },
        { value: 'pcs', label: 'Pieces' },
    ]

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
          name="{basePath}.unit"
          options={units}
          bind:value={$formData.steps[stepIndex].ingredients[ingredientIndex].unit}
/>

<Combobox {form}
          name="{basePath}.name"
          options={ingredients}
          bind:value={$formData.steps[stepIndex].ingredients[ingredientIndex].name}
/>
