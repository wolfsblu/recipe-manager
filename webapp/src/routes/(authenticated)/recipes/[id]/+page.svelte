<script lang="ts">
    import type { PageData } from './$types';
    import pluralize from 'pluralize'
    import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card/index.js";
    import { Badge } from "$lib/components/ui/badge/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import { Separator } from "$lib/components/ui/separator/index.js";
    import { ScrollArea } from "$lib/components/ui/scroll-area/index.js";
    import * as Breadcrumb from "$lib/components/ui/breadcrumb/index.js";
    import * as AlertDialog from "$lib/components/ui/alert-dialog/index.js";
    import ClockIcon from '@lucide/svelte/icons/clock';
    import StepsIcon from '@lucide/svelte/icons/footprints';
    import UsersIcon from '@lucide/svelte/icons/users';
    import ChefHatIcon from '@lucide/svelte/icons/chef-hat';
    import ListIcon from '@lucide/svelte/icons/list';
    import ImageIcon from '@lucide/svelte/icons/image';
    import TagIcon from '@lucide/svelte/icons/tag';
    import IngredientIcon from "@lucide/svelte/icons/salad";
    import ChevronUpIcon from '@lucide/svelte/icons/chevron-up';
    import ChevronDownIcon from '@lucide/svelte/icons/chevron-down';
    import EditIcon from '@lucide/svelte/icons/edit';
    import TrashIcon from '@lucide/svelte/icons/trash-2';
    import {goto} from "$app/navigation";
    import { deleteRecipe } from "$lib/api/recipes/recipes.svelte";
    import { toast } from "svelte-sonner";
    import VoteButtons from "$lib/components/recipes/VoteButtons.svelte";
    import AddToMealPlanButton from "$lib/components/recipes/AddToMealPlanButton.svelte";
    import * as m from "$lib/paraglide/messages.js";

    let { data }: { data: PageData } = $props();
    const { recipe } = data;

    let showDeleteDialog = $state(false)

    const handleEdit = async () => {
        await goto(`/recipes/${recipe.id}/edit`)
    }

    const handleDelete = () => {
        showDeleteDialog = true;
    }

    const confirmDelete = async () => {
        showDeleteDialog = false;
        try {
            await deleteRecipe(recipe.id);
            toast.success(m.recipes_detail_deleteSuccess());
            await goto("/recipes");
        } catch (error) {
            toast.error(m.recipes_detail_deleteError());
        }
    }

    const formatMinutesAsHours = (minutes: number) => {
        if (!minutes || minutes < 0) {
            return m.common_time_notAvailable();
        }

        if (minutes < 60) {
            return m.common_time_minutes({ count: minutes });
        }

        const hours = minutes / 60;
        const roundedHours = Math.round(hours * 100) / 100;
        return m.common_time_hours({ count: roundedHours });
    }

    const ingredients = Array.from(
        recipe.steps
            .flatMap(step => step.ingredients || [])
            .reduce((ingredientMap, ingredient) => {
                const key = `${ingredient.ingredient?.id}-${ingredient.unit.id}`;
                
                if (ingredientMap.has(key)) {
                    // If same ingredient with same unit exists, add the amounts
                    ingredientMap.get(key).amount += ingredient.amount;
                } else {
                    // Create new entry
                    ingredientMap.set(key, {
                        ingredient: ingredient.ingredient,
                        unit: ingredient.unit,
                        amount: ingredient.amount
                    });
                }
                
                return ingredientMap;
            }, new Map())
            .values()
    ).sort((a, b) => 
        (a.ingredient?.name || '').localeCompare(b.ingredient?.name || '')
    )
</script>

<svelte:head>
    <title>{m.recipes_detail_pageTitle({ recipeName: recipe.name })}</title>
</svelte:head>

<div class="container mx-auto px-4 py-6">
    <div class="mb-8">
        <div class="flex flex-col lg:flex-row gap-8">
            <div class="lg:w-1/2">
                <div class="flex flex-row-reverse md:flex-row gap-3">
                    <div class="flex flex-col items-center gap-2">
                        <VoteButtons recipeId={recipe.id} votes={recipe.votes} />

                        <Separator orientation="horizontal" class="my-1 w-full" />

                        <Button
                            variant="ghost"
                            size="sm"
                            onclick={handleEdit}
                            title={m.recipes_detail_editButton()}
                        >
                            <EditIcon class="w-5 h-5" />
                        </Button>

                        <AddToMealPlanButton recipeId={recipe.id} />

                        <Button
                            variant="ghost"
                            size="sm"
                            onclick={handleDelete}
                            title={m.recipes_detail_deleteButton()}
                        >
                            <TrashIcon class="w-5 h-5" />
                        </Button>
                    </div>

                    <!-- Recipe Image -->
                    <div class="flex-1 relative aspect-square rounded-lg overflow-hidden shadow-lg">
                    {#if recipe.images && recipe.images.length > 0}
                        <img
                            src={recipe.images[0]}
                            alt={recipe.name}
                            class="w-full h-full object-cover object-center"
                        />
                    {:else}
                        <div class="w-full h-full bg-muted flex items-center justify-center">
                            <div class="text-center text-muted-foreground">
                                <ImageIcon class="w-16 h-16 mx-auto mb-2" />
                                <p>{m.recipes_detail_noImage()}</p>
                            </div>
                        </div>
                    {/if}
                    </div>
                </div>
            </div>

            <div class="lg:w-1/2 space-y-6">
                <div>
                    <h1 class="text-4xl font-bold text-foreground mb-2">{recipe.name}</h1>

                    {#if recipe.tags && recipe.tags.length > 0}
                        <ScrollArea class="whitespace-nowrap" orientation="horizontal">
                            <div class="flex gap-1">
                                {#each recipe.tags as tag}
                                    <Badge variant="secondary" class="whitespace-nowrap text-sm">
                                        <TagIcon class="w-4 h-4 text-muted-foreground" />
                                        {tag.name}
                                    </Badge>
                                {/each}
                            </div>
                        </ScrollArea>
                    {/if}

                    {#if recipe.description}
                        <p class="text-lg text-muted-foreground leading-relaxed">{recipe.description}</p>
                    {/if}
                </div>

                <div class="grid grid-cols-2 gap-4">
                    <Card class="p-0">
                        <CardContent class="p-4 flex items-center gap-3">
                            <ClockIcon class="w-5 h-5 text-primary" />
                            <div>
                                <p class="text-sm text-muted-foreground">{m.recipes_detail_cookTime()}</p>
                                <p class="text-lg font-semibold">{formatMinutesAsHours(recipe.minutes)}</p>
                            </div>
                        </CardContent>
                    </Card>

                    <Card class="p-0">
                        <CardContent class="p-4 flex items-center gap-3">
                            <UsersIcon class="w-5 h-5 text-primary" />
                            <div>
                                <p class="text-sm text-muted-foreground">{m.recipes_detail_servings()}</p>
                                <p class="text-lg font-semibold">{recipe.servings}</p>
                            </div>
                        </CardContent>
                    </Card>

                    <Card class="p-0">
                        <CardContent class="p-4 flex items-center gap-3">
                            <StepsIcon class="w-5 h-5 text-primary" />
                            <div>
                                <p class="text-sm text-muted-foreground">{m.recipes_detail_steps()}</p>
                                <p class="text-lg font-semibold">{recipe.steps.length}</p>
                            </div>
                        </CardContent>
                    </Card>

                    <Card class="p-0">
                        <CardContent class="p-4 flex items-center gap-3">
                            <IngredientIcon class="w-5 h-5 text-primary" />
                            <div>
                                <p class="text-sm text-muted-foreground">{m.recipes_detail_ingredients()}</p>
                                <p class="text-lg font-semibold">
                                    {recipe.steps.reduce((total, step) => total + (step.ingredients?.length || 0), 0)}
                                </p>
                            </div>
                        </CardContent>
                    </Card>
                </div>

                {#if recipe.images && recipe.images.length > 1}
                    <ScrollArea class="whitespace-nowrap lg:whitespace-normal" orientation="horizontal">
                        <div class="flex lg:grid lg:grid-cols-3 gap-4 pb-2 lg:pb-0">
                            {#each recipe.images.slice(1) as image}
                                <div class="aspect-square w-40 md:w-48 lg:w-auto rounded-lg overflow-hidden shadow-sm flex-shrink-0">
                                    <img
                                            src={image}
                                            alt="{recipe.name} - Additional view"
                                            class="w-full h-full object-cover hover:scale-105 transition-transform duration-200"
                                    />
                                </div>
                            {/each}
                        </div>
                    </ScrollArea>
                {/if}
            </div>
        </div>
    </div>

    {#if recipe.steps && recipe.steps.length > 0 && ingredients.length > 0}
        <h2 class="font-bold text-2xl text-foreground mb-3">
            {m.recipes_detail_ingredientsTitle()}
        </h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-2">
            {#each ingredients as ingredient}
                <div class="border flex items-center justify-between px-3 py-2 bg-muted/30 rounded-lg">
                    <span class="font-medium text-muted-foreground">
                        {ingredient.ingredient?.name || m.recipes_detail_unknownIngredient()}
                    </span>
                    <Badge variant="outline" class="text-sm text-muted-foreground">
                        {ingredient.amount}
                        {pluralize(ingredient.unit.name, ingredient.amount)}
                    </Badge>
                </div>
            {/each}
        </div>
    {/if}

    <Separator class="my-6" />

    <h2 class="font-bold text-2xl text-foreground mb-3">{m.recipes_detail_instructionsTitle()}</h2>
    {#if recipe.steps && recipe.steps.length > 0}
        <!-- Instructions Section -->
        <div class="space-y-6">
            <div class="space-y-8">
                {#each recipe.steps as step, stepIndex}
                    <Card class="overflow-hidden">
                        <CardHeader class="bg-muted/50 gap-0 py-1">
                            <CardTitle class="flex items-center gap-3">
                                <div class="w-8 h-8 bg-primary text-primary-foreground rounded-full flex items-center justify-center text-sm font-bold">
                                    {stepIndex + 1}
                                </div>
                                {m.recipes_detail_stepLabel({ number: stepIndex + 1 })}
                            </CardTitle>
                        </CardHeader>
                        <CardContent class="py-0">
                            <div class="grid lg:grid-cols-2 gap-6">
                                {#if step.ingredients && step.ingredients.length > 0}
                                    <div>
                                        <h4 class="font-semibold mb-3 flex items-center gap-2">
                                            <IngredientIcon class="w-4 h-4" />
                                            {m.recipes_form_ingredients()}
                                        </h4>
                                        <div class="space-y-1">
                                            {#each step.ingredients as ingredient}
                                                <div class="flex items-center justify-between px-3 py-2 bg-muted/30 rounded-lg">
                                                    <span class="font-medium text-muted-foreground">
                                                        {ingredient.ingredient?.name || m.recipes_detail_unknownIngredient()}
                                                    </span>
                                                    <Badge variant="outline" class="text-sm text-muted-foreground">
                                                        {ingredient.amount}
                                                        {pluralize(ingredient.unit.name, ingredient.amount)}
                                                    </Badge>
                                                </div>
                                            {/each}
                                        </div>
                                    </div>
                                {/if}

                                <div>
                                    <h4 class="font-semibold mb-3 flex items-center gap-2">
                                        <ChefHatIcon class="w-4 h-4" />
                                        {m.recipes_form_instructions()}
                                    </h4>
                                    {#if step.instructions}
                                        <div class="prose prose-sm max-w-none">
                                            <p class="text-muted-foreground leading-relaxed whitespace-pre-wrap">{step.instructions}</p>
                                        </div>
                                    {:else}
                                        <p class="text-muted-foreground italic">{m.recipes_detail_noInstructions()}</p>
                                    {/if}
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                {/each}
            </div>
        </div>
    {:else}
        <Card>
            <CardContent class="p-8 text-center">
                <ChefHatIcon class="w-16 h-16 mx-auto mb-4 text-muted-foreground" />
                <p class="text-lg text-muted-foreground">{m.recipes_detail_noCookingInstructions()}</p>
            </CardContent>
        </Card>
    {/if}

</div>

<AlertDialog.Root bind:open={showDeleteDialog}>
    <AlertDialog.Content>
        <AlertDialog.Header>
            <AlertDialog.Title>{m.recipes_detail_deleteDialog_title()}</AlertDialog.Title>
            <AlertDialog.Description>
                {m.recipes_detail_deleteDialog_description({ recipeName: recipe.name })}
            </AlertDialog.Description>
        </AlertDialog.Header>
        <AlertDialog.Footer>
            <AlertDialog.Cancel>{m.recipes_detail_deleteDialog_cancel()}</AlertDialog.Cancel>
            <AlertDialog.Action onclick={confirmDelete}>
                {m.recipes_detail_deleteDialog_delete()}
            </AlertDialog.Action>
        </AlertDialog.Footer>
    </AlertDialog.Content>
</AlertDialog.Root>