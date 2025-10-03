<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
    import { Card, CardContent } from "$lib/components/ui/card/index.js";
    import MealPlanDay from "$lib/components/mealplan/MealPlanDay.svelte";
    import DateRangePicker from "$lib/components/ui/date-range-picker.svelte";
    import CalendarIcon from '@lucide/svelte/icons/calendar';
    import ChefHatIcon from '@lucide/svelte/icons/chef-hat';
    import { goto } from '$app/navigation';
    import type { PageProps } from './$types';
    import * as m from "$lib/paraglide/messages.js";

    const { data }: PageProps = $props();

    let mealPlan = $state(data.mealPlan || []);

    let dateRange = $state({
        from: data.from,
        to: data.until
    });

    // Filter out days with no recipes
    let daysWithRecipes = $derived(mealPlan.filter(day => day.recipes && day.recipes.length > 0));

    const handleDateRangeChange = async (newRange: { from: string; to: string }) => {
        if (newRange.from && newRange.to) {
            dateRange = newRange;
            // Navigate to the same page with new query parameters
            const url = new URL(window.location.href);
            url.searchParams.set('from', newRange.from);
            url.searchParams.set('until', newRange.to);
            goto(url.toString());
        }
    };
</script>

<div class="flex flex-col gap-5 p-5 h-full">
    <!-- Date Range Picker -->
    <div class="flex justify-between items-center">
        <DateRangePicker
            value={dateRange}
            onValueChange={handleDateRangeChange}
            placeholder={m.mealPlan_dateRangePlaceholder()}
        />
    </div>

    <!-- Meal Plan Content -->
    {#if daysWithRecipes.length > 0}
        <div class="space-y-6">
            {#each mealPlan as mealPlanDay, i}
                <MealPlanDay bind:mealPlanDay={mealPlan[i]} availableTags={data.tags} shoppingLists={data.shoppingLists} />
            {/each}
        </div>
    {:else}
        <!-- Empty State -->
        <div class="flex items-center justify-center flex-grow">
            <Card class="max-w-md w-full">
                <CardContent>
                    <div class="text-center">
                        <div class="p-4 bg-muted rounded-full mx-auto w-fit mb-4">
                            <CalendarIcon class="h-8 w-8 text-muted-foreground" />
                        </div>
                        <h3 class="text-lg font-semibold mb-2">{m.mealPlan_empty_title()}</h3>
                        <p class="text-muted-foreground mb-4">
                            {m.mealPlan_empty_description()}
                        </p>
                        <Button href="/recipes">
                            <ChefHatIcon />
                            {m.mealPlan_empty_button()}
                        </Button>
                    </div>
                </CardContent>
            </Card>
        </div>
    {/if}
</div>