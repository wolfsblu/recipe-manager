<script lang="ts">
    import { Badge } from "$lib/components/ui/badge/index.js";
    import RecipeCard from '../recipes/RecipeCard.svelte';
    import CalendarIcon from '@lucide/svelte/icons/calendar';
    import type { components } from '$lib/api/schema';

    type MealPlanDay = components['schemas']['ReadMealPlan'];

    let { mealPlanDay, availableTags = [] }: { mealPlanDay: MealPlanDay; availableTags?: any[] } = $props();

    const formatDate = (dateString: string) => {
        const date = new Date(dateString);
        const today = new Date();
        const tomorrow = new Date(today);
        tomorrow.setDate(tomorrow.getDate() + 1);
        const yesterday = new Date(today);
        yesterday.setDate(yesterday.getDate() - 1);

        // Reset time for comparison
        const dateOnly = new Date(date.getFullYear(), date.getMonth(), date.getDate());
        const todayOnly = new Date(today.getFullYear(), today.getMonth(), today.getDate());
        const tomorrowOnly = new Date(tomorrow.getFullYear(), tomorrow.getMonth(), tomorrow.getDate());
        const yesterdayOnly = new Date(yesterday.getFullYear(), yesterday.getMonth(), yesterday.getDate());

        if (dateOnly.getTime() === todayOnly.getTime()) {
            return 'Today';
        } else if (dateOnly.getTime() === tomorrowOnly.getTime()) {
            return 'Tomorrow';
        } else if (dateOnly.getTime() === yesterdayOnly.getTime()) {
            return 'Yesterday';
        } else {
            return date.toLocaleDateString('en-US', {
                weekday: 'long',
                month: 'short',
                day: 'numeric'
            });
        }
    };

    const getDayOfWeek = (dateString: string) => {
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', { weekday: 'short' });
    };

    const getFormattedDate = (dateString: string) => {
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric'
        });
    };
</script>

<div class="space-y-4">
    <!-- Day Header -->
    <div class="flex items-center justify-between pb-2 border-b">
        <div class="flex items-center gap-3">
            <CalendarIcon class="h-5 w-5 text-muted-foreground" />
            <div class="flex flex-col">
                <h2 class="text-xl font-semibold">{formatDate(mealPlanDay.date)}</h2>
                <div class="text-sm text-muted-foreground">
                    {getDayOfWeek(mealPlanDay.date)} • {getFormattedDate(mealPlanDay.date)}
                </div>
            </div>
        </div>
        {#if mealPlanDay.recipes && mealPlanDay.recipes.length > 0}
            <Badge variant="secondary" class="text-sm">
                {mealPlanDay.recipes.length} {mealPlanDay.recipes.length === 1 ? 'recipe' : 'recipes'}
            </Badge>
        {/if}
    </div>

    <!-- Recipes Grid -->
    {#if mealPlanDay.recipes && mealPlanDay.recipes.length > 0}
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-4">
            {#each mealPlanDay.recipes as recipe}
                <RecipeCard {recipe} {availableTags} />
            {/each}
        </div>
    {:else}
        <div class="text-center py-12 text-muted-foreground">
            <CalendarIcon class="h-12 w-12 mx-auto mb-4 opacity-30" />
            <p class="text-lg mb-1">No recipes planned for this day</p>
            <p class="text-sm">Add some recipes to your meal plan!</p>
        </div>
    {/if}
</div>