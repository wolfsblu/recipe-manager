<script lang="ts">
    import type { PageData } from './$types';
    import RecipeForm from "$lib/components/recipes/RecipeForm.svelte";
    
    import { formSchema } from "./schema";
    import { setMessage, superForm } from "sveltekit-superforms";
    import { zodClient } from "sveltekit-superforms/adapters";
    import { toast } from "svelte-sonner";
    import { goto } from "$app/navigation";
    import { addRecipe } from "$lib/api/recipes/recipes.svelte";
    import * as m from "$lib/paraglide/messages.js";

    let { data }: { data: PageData } = $props();
    const { ingredients, units, tags } = data;

    const form = superForm(data.form, {
        SPA: true,
        dataType: 'json',
        resetForm: false,
        validators: zodClient(formSchema),
        async onUpdate({ form }) {
            if (form.valid) {
                try {
                    await addRecipe(form.data);
                    toast.success(m.recipes_add_successMessage());
                    await goto("/recipes");
                } catch (error) {
                    toast.error(m.recipes_add_errorMessage());
                }
            }
        }
    });
</script>

<svelte:head>
    <title>{m.recipes_add_pageTitle()}</title>
</svelte:head>

<div class="container mx-auto px-0 md:px-4 py-6">
    <!-- Form Header -->
    <div class="mb-8 px-4 md:px-0">
        <h1 class="text-4xl font-bold text-foreground mb-2">{m.recipes_add_title()}</h1>
        <p class="text-lg text-muted-foreground">{m.recipes_add_subtitle()}</p>
    </div>

    <RecipeForm
        {form}
        {ingredients}
        {units}
        {tags}
        isEditing={false}
        submitText={m.recipes_add_submitButton()}
    />
</div>