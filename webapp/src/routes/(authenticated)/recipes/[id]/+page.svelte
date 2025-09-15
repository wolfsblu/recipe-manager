<script lang="ts">
    import type { PageData } from './$types';
    import pluralize from 'pluralize'
    import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card/index.js";
    import { Badge } from "$lib/components/ui/badge/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import { Separator } from "$lib/components/ui/separator/index.js";
    import { ScrollArea } from "$lib/components/ui/scroll-area/index.js";
    import * as Breadcrumb from "$lib/components/ui/breadcrumb/index.js";
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
    import {goto} from "$app/navigation";

    let { data }: { data: PageData } = $props();
    const { recipe } = data;

    let votes = $state(0)
    let userVote = $state<'up' | 'down' | null>(null)

    const handleVote = (voteType: 'up' | 'down') => {
        if (userVote === voteType) {
            // Remove vote if clicking the same button
            votes -= voteType === 'up' ? 1 : -1
            userVote = null;
        } else if (userVote === null) {
            // Add new vote
            votes += voteType === 'up' ? 1 : -1
            userVote = voteType;
        } else {
            // Switch vote
            votes += voteType === 'up' ? 2 : -2
            userVote = voteType;
        }
    }

    const handleEdit = async () => {
        await goto(`/recipes/${recipe.id}/edit`)
    }

    const formatMinutesAsHours = (minutes: number) => {
        if (!minutes || minutes < 0) {
            return "N/A";
        }

        if (minutes < 60) {
            return `${minutes}m`;
        }

        const hours = minutes / 60;
        const roundedHours = Math.round(hours * 100) / 100;
        return `${roundedHours}h`;
    }
</script>

<svelte:head>
    <title>{recipe.name} - Recipe Manager</title>
</svelte:head>

<div class="container mx-auto px-4 py-6">
    <div class="mb-8">
        <div class="flex flex-col lg:flex-row gap-8">
            <div class="lg:w-1/2">
                <div class="flex gap-3">
                    <div class="flex flex-col items-center gap-2">
                        <Button
                            variant="ghost"
                            size="sm"
                            onclick={() => handleVote('up')}
                            class={userVote === 'up' ? 'text-green-600 hover:text-green-700' : ''}
                            title="Upvote recipe"
                        >
                            <ChevronUpIcon class="w-5 h-5" />
                        </Button>

                        <span class="text-base font-bold text-foreground">
                            {votes}
                        </span>

                        <Button
                            variant="ghost"
                            size="sm"
                            onclick={() => handleVote('down')}
                            class={userVote === 'down' ? 'text-red-600 hover:text-red-700' : ''}
                            title="Downvote recipe"
                        >
                            <ChevronDownIcon class="w-5 h-5" />
                        </Button>

                        <Separator orientation="horizontal" class="my-1 w-full" />

                        <Button
                            variant="ghost"
                            size="sm"
                            onclick={handleEdit}
                            title="Edit recipe"
                        >
                            <EditIcon class="w-5 h-5" />
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
                                <p>No image available</p>
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
                                <p class="text-sm text-muted-foreground">Cook Time</p>
                                <p class="text-lg font-semibold">{formatMinutesAsHours(recipe.minutes)}</p>
                            </div>
                        </CardContent>
                    </Card>

                    <Card class="p-0">
                        <CardContent class="p-4 flex items-center gap-3">
                            <UsersIcon class="w-5 h-5 text-primary" />
                            <div>
                                <p class="text-sm text-muted-foreground">Servings</p>
                                <p class="text-lg font-semibold">{recipe.servings}</p>
                            </div>
                        </CardContent>
                    </Card>

                    <Card class="p-0">
                        <CardContent class="p-4 flex items-center gap-3">
                            <StepsIcon class="w-5 h-5 text-primary" />
                            <div>
                                <p class="text-sm text-muted-foreground">Steps</p>
                                <p class="text-lg font-semibold">{recipe.steps.length}</p>
                            </div>
                        </CardContent>
                    </Card>

                    <Card class="p-0">
                        <CardContent class="p-4 flex items-center gap-3">
                            <IngredientIcon class="w-5 h-5 text-primary" />
                            <div>
                                <p class="text-sm text-muted-foreground">Ingredients</p>
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

    {#if recipe.steps && recipe.steps.length > 0}
        <div class="space-y-6">
            <h2 class="text-2xl font-bold text-foreground">Instructions</h2>

            <div class="space-y-8">
                {#each recipe.steps as step, stepIndex}
                    <Card class="overflow-hidden">
                        <CardHeader class="bg-muted/50 gap-0 py-1">
                            <CardTitle class="flex items-center gap-3">
                                <div class="w-8 h-8 bg-primary text-primary-foreground rounded-full flex items-center justify-center text-sm font-bold">
                                    {stepIndex + 1}
                                </div>
                                Step {stepIndex + 1}
                            </CardTitle>
                        </CardHeader>
                        <CardContent class="py-0">
                            <div class="grid lg:grid-cols-2 gap-6">
                                {#if step.ingredients && step.ingredients.length > 0}
                                    <div>
                                        <h4 class="font-semibold mb-3 flex items-center gap-2">
                                            <IngredientIcon class="w-4 h-4" />
                                            Ingredients
                                        </h4>
                                        <div class="space-y-1">
                                            {#each step.ingredients as ingredient}
                                                <div class="flex items-center justify-between px-3 py-2 bg-muted/30 rounded-lg">
                                                    <span class="font-medium text-muted-foreground">
                                                        {ingredient.ingredient?.name || 'Unknown ingredient'}
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
                                        Instructions
                                    </h4>
                                    {#if step.instructions}
                                        <div class="prose prose-sm max-w-none">
                                            <p class="text-muted-foreground leading-relaxed whitespace-pre-wrap">{step.instructions}</p>
                                        </div>
                                    {:else}
                                        <p class="text-muted-foreground italic">No instructions provided for this step.</p>
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
                <p class="text-lg text-muted-foreground">No cooking instructions available for this recipe.</p>
            </CardContent>
        </Card>
    {/if}

</div>