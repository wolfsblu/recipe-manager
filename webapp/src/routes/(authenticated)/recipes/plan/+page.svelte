<script lang="ts">
    import {Badge, Button, ButtonGroup, Card, Timeline, TimelineItem} from "flowbite-svelte";
    import {AngleLeftOutline, AngleRightOutline, EditOutline, TrashBinOutline} from "flowbite-svelte-icons";
    import {getDaysInMonth, getWeekdays, isToday} from "$lib/calendar/calendar";
    import Recipe from "$lib/components/recipes/Recipe.svelte";
    import PlannedRecipe from "$lib/components/recipes/PlannedRecipe.svelte";

    const now = new Date();
    const days = getDaysInMonth(now);
    const weekdays = getWeekdays()
</script>

{#snippet Col()}
    <div class="grid grid-cols-7 gap-3">
        {#each {length: 7} as _, i}
            {#if i !== 2}
                <PlannedRecipe title="My super duper recipe" />
            {:else}
                <span></span>
            {/if}
        {/each}
    </div>
{/snippet}

<Card size="none">
    <div>
        {#each weekdays as weekday}
            <span>{weekday}</span>
        {/each}
    </div>
    <div class="flex flex-col sm:flex-row justify-between items-center pb-4 gap-3">
        <h2 class="text-2xl font-semibold">January 2022</h2>
        <div class="flex space-x-2">
            <Button color="alternative">
                <AngleLeftOutline class="w-5 h-5"/>
            </Button>
            <Button color="alternative">
                <AngleRightOutline class="w-5 h-5"/>
            </Button>
            <Button>
                Today
            </Button>
        </div>
    </div>
    <Timeline>
        <TimelineItem date="Breakfast">
            {@render Col()}
        </TimelineItem>
        <TimelineItem date="Lunch">
            {@render Col()}
        </TimelineItem>
        <TimelineItem date="Dinner">
            {@render Col()}
        </TimelineItem>
    </Timeline>
</Card>