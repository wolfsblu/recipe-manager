<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
    import * as Popover from "$lib/components/ui/popover/index.js";
    import { Calendar } from "$lib/components/ui/calendar/index.js";
    import CalendarIcon from '@lucide/svelte/icons/calendar';
    import PlusIcon from '@lucide/svelte/icons/plus';
    import { toast } from "svelte-sonner";
    import {
        DateFormatter,
        type DateValue,
        getLocalTimeZone,
        today,
        CalendarDate
    } from "@internationalized/date";
    import { cn } from "$lib/utils.js";
    import { buttonVariants } from "$lib/components/ui/button/index.js";
    import { createMealPlan } from "$lib/api/mealplan/mealplan.svelte";

    let { recipeId }: { recipeId: number } = $props();

    let open = $state(false);
    let selectedDate = $state<DateValue | undefined>(today(getLocalTimeZone()));

    const df = new DateFormatter("en-US", {
        dateStyle: "long"
    });

    const handleAddToMealPlan = async () => {
        if (!selectedDate) {
            toast.error("Please select a date");
            return;
        }

        try {
            // Convert DateValue to string format (YYYY-MM-DD)
            const dateString = `${selectedDate.year}-${selectedDate.month.toString().padStart(2, '0')}-${selectedDate.day.toString().padStart(2, '0')}`;

            await createMealPlan(recipeId, dateString);

            const formattedDate = df.format(selectedDate.toDate(getLocalTimeZone()));
            toast.success(`Recipe added to meal plan for ${formattedDate}`);
            open = false;
        } catch (error) {
            toast.error("Failed to add recipe to meal plan");
        }
    };
</script>

<Popover.Root bind:open>
    <Popover.Trigger>
        <Button
            variant="ghost"
            size="sm"
            title="Add to meal plan"
        >
            <CalendarIcon class="w-5 h-5" />
        </Button>
    </Popover.Trigger>
    <Popover.Content class="w-80 p-4" align="center">
        <div class="space-y-4">
            <div class="space-y-2">
                <h4 class="font-medium">Add to Meal Plan</h4>
                <p class="text-sm text-muted-foreground">
                    Select a date to add this recipe to your meal plan.
                </p>
            </div>

            <div class="space-y-2">
                <Popover.Root>
                    <Popover.Trigger
                        class={cn(
                            buttonVariants({ variant: "outline" }),
                            "w-full justify-start text-left font-normal",
                            !selectedDate && "text-muted-foreground"
                        )}
                    >
                        <CalendarIcon class="mr-2 h-4 w-4" />
                        {selectedDate ? df.format(selectedDate.toDate(getLocalTimeZone())) : "Pick a date"}
                    </Popover.Trigger>
                    <Popover.Content class="w-auto p-0">
                        <Calendar
                            type="single"
                            bind:value={selectedDate}
                            class="rounded-md border"
                        />
                    </Popover.Content>
                </Popover.Root>
            </div>

            <div class="flex justify-end gap-2">
                <Button variant="outline" size="sm" onclick={() => open = false}>
                    Cancel
                </Button>
                <Button size="sm" onclick={handleAddToMealPlan} disabled={!selectedDate}>
                    <PlusIcon />
                    Add
                </Button>
            </div>
        </div>
    </Popover.Content>
</Popover.Root>