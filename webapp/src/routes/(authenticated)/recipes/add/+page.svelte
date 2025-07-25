<script lang="ts">
    import * as Form from '$lib/components/ui/form/index.js'
    import {Input} from "$lib/components/ui/input/index.js";
    import {Button} from "$lib/components/ui/button/index.js";
    import {Textarea} from "$lib/components/ui/textarea/index.js";
    import {formSchema, type FormSchema} from "./schema";
    import {type Infer, setMessage, superForm, type SuperValidated,} from "sveltekit-superforms";
    import {zodClient} from "sveltekit-superforms/adapters";
    import {login} from "$lib/api/auth/user.svelte";
    import {toast} from "svelte-sonner";
    import {goto} from "$app/navigation";
    import {addRecipe} from "$lib/api/recipes/recipes.svelte";
    import RecipeImages from "$lib/components/recipes/RecipeImages.svelte";
    import { TagsInput } from '$lib/components/ui/tags-input';
    import ImageUpload from "$lib/components/recipes/ImageUpload.svelte";
    import UnitCombobox from "$lib/components/recipes/UnitCombobox.svelte";

    let {data}: { data: { form: SuperValidated<Infer<FormSchema>> } } = $props();
    const form = superForm(data.form, {
        SPA: true,
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


    const removeIngredientByIndex = (index: number) => {
        $formData.ingredients = $formData.ingredients.filter((_, i) => i !== index);
    }

    const addIngredient = () => {
        $formData.ingredients = [...$formData.ingredients, ""];
    }
</script>

<form class="flex flex-col flex-grow p-6" method="POST" use:enhance>

    <ImageUpload />

    <div class="flex-grow flex flex-col lg:flex-row gap-6">
        <div class="w-full flex flex-col gap-3">
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
                            <Input {...props} bind:value={$formData.minutes} type="number" placeholder="4"/>
                        {/snippet}
                    </Form.Control>
                    <Form.Description/>
                    <Form.FieldErrors/>
                </Form.Field>
            </div>
            <Form.Fieldset {form} name="ingredients">
                <Form.Legend>Ingredients</Form.Legend>
                {#each $formData.ingredients as _, i}
                    <Form.ElementField {form} name="ingredients[{i}]">
                        <Form.Control>
                            {#snippet children({props})}
                                <Form.Label>Ingredient {i + 1}</Form.Label>
                                <Input {...props} bind:value={$formData.ingredients[i]} />
                                <Button onclick={() => removeIngredientByIndex(i)} type="button">
                                    Remove
                                </Button>
                            {/snippet}
                        </Form.Control>
                        <Form.Description/>
                        <Form.FieldErrors/>
                    </Form.ElementField>
                {/each}
                <Form.Description/>
                <Form.FieldErrors/>
                <Button onclick={addIngredient} type="button">
                    Add Recipe
                </Button>
            </Form.Fieldset>
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
        <Form.Field class="w-full flex-grow flex flex-col" {form} name="description">
            <Form.Control>
                {#snippet children({props})}
                    <Form.Label>Description</Form.Label>
                    <Textarea {...props}
                              class="flex-grow resize-none"
                              bind:value={$formData.description}
                              placeholder="Type your message here."
                    />
                {/snippet}
            </Form.Control>
            <Form.Description/>
            <Form.FieldErrors/>
        </Form.Field>
    </div>
    <Button class="w-full" type="submit">Create</Button>
</form>
