<script lang="ts">
    import type { PageData } from './$types';
    import RecipeForm from "$lib/components/recipes/RecipeForm.svelte";
    
    import { formSchema } from "./schema";
    import { setMessage, superForm } from "sveltekit-superforms";
    import { zodClient } from "sveltekit-superforms/adapters";
    import { toast } from "svelte-sonner";
    import { goto } from "$app/navigation";
    import { addRecipe } from "$lib/api/recipes/recipes.svelte";

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
                    toast.success("Recipe created successfully!");
                    await goto("/recipes");
                } catch (error) {
                    toast.error("Failed to create recipe");
                }
            }
        }
    });
</script>

<svelte:head>
    <title>Create New Recipe - Recipe Manager</title>
</svelte:head>

<div class="container mx-auto px-4 py-6">
    <!-- Form Header -->
    <div class="mb-8">
        <h1 class="text-4xl font-bold text-foreground mb-2">Create New Recipe</h1>
        <p class="text-lg text-muted-foreground">Share your culinary creation with the world</p>
    </div>

    <RecipeForm 
        {form} 
        {ingredients} 
        {units} 
        {tags} 
        isEditing={false}
        submitText="Create Recipe"
    />
</div>