<script>
    import * as Form from '$lib/components/ui/form/index.js'
    import {Input} from "$lib/components/ui/input/index.js";

    const {
        form,
        stepIndex,
        ingredientIndex,
    } = $props();

    const { form: formData } = form;

    const basePath = `steps[${stepIndex}].ingredients[${ingredientIndex}]`;
</script>

<div class="flex-grow grid grid-cols-[1fr_2fr_3fr] gap-2">
<Form.ElementField class="space-y-0" {form} name="{basePath}.amount">
    <Form.Control>
        {#snippet children({props})}
            <div class="flex flex-col md:flex-row gap-1">
                <Input {...props} class="w-full" type="number" bind:value={$formData.steps[stepIndex].ingredients[ingredientIndex].amount} placeholder="1" />
            </div>
        {/snippet}
    </Form.Control>
    <Form.Description class="sr-only" />
    <Form.FieldErrors/>
</Form.ElementField>
<Form.ElementField class="space-y-0" {form} name="{basePath}.unit">
    <Form.Control>
        {#snippet children({props})}
            <div class="flex flex-col md:flex-row gap-1">
                <Input {...props} class="w-full" bind:value={$formData.steps[stepIndex].ingredients[ingredientIndex].unit} placeholder="kg" />
            </div>
        {/snippet}
    </Form.Control>
    <Form.Description class="sr-only" />
    <Form.FieldErrors/>
</Form.ElementField>
<Form.ElementField class="space-y-0" {form} name="{basePath}.name">
    <Form.Control>
        {#snippet children({props})}
            <div class="flex flex-col md:flex-row gap-1">
                <Input {...props} class="w-full" bind:value={$formData.steps[stepIndex].ingredients[ingredientIndex].name} placeholder="Potatoes" />
            </div>
        {/snippet}
    </Form.Control>
    <Form.Description class="sr-only" />
    <Form.FieldErrors/>
</Form.ElementField>
</div>