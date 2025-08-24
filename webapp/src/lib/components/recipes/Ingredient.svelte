<script>
    import * as Form from '$lib/components/ui/form/index.js'
    import {Input} from "$lib/components/ui/input/index.js";

    const {
        form,
        value = $bindable(),
        path
    } = $props();
</script>

<Form.ElementField class="space-y-0" {form} name="{path}.amount">
    <Form.Control>
        {#snippet children({props})}
            <div class="flex flex-col md:flex-row gap-1">
                <Input {...props} class="w-auto" type="number" bind:value={value.amount} placeholder="1" />
            </div>
        {/snippet}
    </Form.Control>
    <Form.Description class="sr-only" />
    <Form.FieldErrors/>
</Form.ElementField>
<Form.ElementField class="space-y-0" {form} name="{path}.unit">
    <Form.Control>
        {#snippet children({props})}
            <div class="flex flex-col md:flex-row gap-1">
                <Input {...props} class="w-auto" bind:value={value.unit} placeholder="kg" />
            </div>
        {/snippet}
    </Form.Control>
    <Form.Description class="sr-only" />
    <Form.FieldErrors/>
</Form.ElementField>
<Form.ElementField class="space-y-0 w-full" {form} name="{path}.name">
    <Form.Control>
        {#snippet children({props})}
            <div class="flex flex-col md:flex-row gap-1">
                <Input {...props} class="w-full" bind:value={value.name} placeholder="Potatoes" />
            </div>
        {/snippet}
    </Form.Control>
    <Form.Description class="sr-only" />
    <Form.FieldErrors/>
</Form.ElementField>