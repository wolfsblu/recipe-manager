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

    let { data }: PageProps = $props();

    let mealPlan = $state(data.mealPlan);

    let dateRange = $derived({
        from: data.from,
        to: data.until
    });

    $effect(() => {
        mealPlan = data.mealPlan;
    });

    let daysWithRecipes = $derived((mealPlan || []).filter(day => day.recipes && day.recipes.length > 0));

    const handleDateRangeChange = async (newRange: { from: string; to: string }) => {
        if (newRange.from && newRange.to) {
            goto(`/mealplan?from=${newRange.from}&until=${newRange.to}`, { replaceState: true });
        }
    };
</script>

<div class="flex flex-col gap-5 p-5 h-full">
    <!-- Meal Plan Content -->
    {#if daysWithRecipes.length > 0}
        <!-- Date Range Picker -->
        <div class="flex justify-between items-center">
            <DateRangePicker
                value={dateRange}
                onValueChange={handleDateRangeChange}
                placeholder={m.mealPlan_dateRangePlaceholder()}
            />
        </div>

        <div class="space-y-6">
            {#each mealPlan as mealPlanDay, i}
                <MealPlanDay bind:mealPlanDay={mealPlan[i]} availableTags={data.tags} shoppingLists={data.shoppingLists} />
            {/each}
        </div>
    {:else}
        <!-- Empty State -->
        <div class="flex items-center justify-center flex-grow">
            <Card class="max-w-md w-full">
                <CardContent class="pt-6">
                    <div class="text-center">
                        <div class="p-4 bg-muted rounded-full mx-auto w-fit mb-4">
                            <CalendarIcon class="h-8 w-8 text-muted-foreground" />
                        </div>
                        <h3 class="text-lg font-semibold mb-2">{m.mealPlan_empty_title()}</h3>
                        <p class="text-muted-foreground mb-4">
                            {m.mealPlan_empty_description()}
                        </p>
                        <Button href="/recipes" class="w-full">
                            <ChefHatIcon />
                            {m.mealPlan_empty_button()}
                        </Button>
                        <div class="relative my-4">
                            <div class="absolute inset-0 flex items-center">
                                <span class="w-full border-t"></span>
                            </div>
                            <div class="relative flex justify-center text-xs">
                                <span class="bg-card px-2 text-muted-foreground">{m.mealPlan_empty_orPickPeriod()}</span>
                            </div>
                        </div>
                        <DateRangePicker
                            value={dateRange}
                            onValueChange={handleDateRangeChange}
                            placeholder={m.mealPlan_dateRangePlaceholder()}
                        />
                    </div>
                </CardContent>
            </Card>
        </div>
    {/if}
</div>