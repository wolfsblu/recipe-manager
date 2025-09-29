<script lang="ts">
    import type { PageData } from './$types';
    import RecipeForm from "$lib/components/recipes/RecipeForm.svelte";
    
    import { formSchema } from "../../add/schema";
    import { setMessage, superForm } from "sveltekit-superforms";
    import { zodClient } from "sveltekit-superforms/adapters";
    import { toast } from "svelte-sonner";
    import { goto } from "$app/navigation";
    import { updateRecipe } from "$lib/api/recipes/recipes.svelte";
    import { commitImageDeletions, clearImageDeletions } from "$lib/stores/imageDeleteQueue";
    import * as m from "$lib/paraglide/messages.js";

    let { data }: { data: PageData } = $props();
    const { recipe, ingredients, units, tags } = data;

    const form = superForm(data.form, {
        SPA: true,
        dataType: 'json',
        resetForm: false,
        validators: zodClient(formSchema),
        async onUpdate({ form }) {
            if (form.valid) {
                try {
                    await updateRecipe(recipe.id, form.data);
                    await commitImageDeletions();
                    toast.success(m.recipes_edit_successMessage());
                    await goto(`/recipes/${recipe.id}`);
                } catch (error) {
                    toast.error(m.recipes_edit_errorMessage());
                }
            }
        }
    });

    const handleCancel = () => {
        clearImageDeletions();
        goto(`/recipes/${recipe.id}`);
    };

</script>

<svelte:head>
    <title>{m.recipes_edit_pageTitle({ recipeName: recipe.name })}</title>
</svelte:head>

<div class="container mx-auto px-0 md:px-4 py-6">
    <div class="mb-8 px-4 md:px-0">
        <h1 class="text-4xl font-bold text-foreground mb-2">{m.recipes_edit_title()}</h1>
        <p class="text-lg text-muted-foreground">{m.recipes_edit_subtitle({ recipeName: recipe.name })}</p>
    </div>

    <RecipeForm
        {form}
        {ingredients}
        {units}
        {tags}
        isEditing={true}
        submitText={m.recipes_edit_submitButton()}
        onCancel={handleCancel}
    />
</div>