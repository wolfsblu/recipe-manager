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
                    toast.success("Recipe updated successfully!");
                    await goto(`/recipes/${recipe.id}`);
                } catch (error) {
                    toast.error("Failed to update recipe");
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
    <title>Edit {recipe.name} - Recipe Manager</title>
</svelte:head>

<div class="container mx-auto px-4 py-6">
    <div class="mb-8">
        <h1 class="text-4xl font-bold text-foreground mb-2">Edit Recipe</h1>
        <p class="text-lg text-muted-foreground">Update your "{recipe.name}" recipe</p>
    </div>

    <RecipeForm 
        {form} 
        {ingredients} 
        {units} 
        {tags} 
        isEditing={true}
        submitText="Update Recipe"
        onCancel={handleCancel}
    />
</div>