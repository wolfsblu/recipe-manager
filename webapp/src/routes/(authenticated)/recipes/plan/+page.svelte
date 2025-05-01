<script lang="ts">
    import {Button, Card, Timeline, TimelineItem} from "flowbite-svelte";
    import {AngleLeftOutline, AngleRightOutline} from "flowbite-svelte-icons";
    import {getWeekdays, isToday} from "$lib/calendar/calendar";
    import PlannedRecipe from "$lib/components/recipes/PlannedRecipe.svelte";
    import {DateTime, Settings} from "luxon";

    const weekdays = getWeekdays(DateTime.now())
</script>

<Card size="none" padding="none">
    <div class="flex flex-col sm:flex-row justify-between items-center pb-4 gap-3 p-6">
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

    <div class="grid grid-cols-7 border-y">
        {#each weekdays as weekday}
            <div class="text-center py-3 dark:text-white {isToday(weekday.date) ? 'bg-primary-200 dark:bg-primary-900' : ''}">
                <p class="font-thin mb-1">{weekday.name.toUpperCase()}</p>
                <p class="font-bold">{weekday.date.toLocaleString({day: "numeric"})}</p>
            </div>
        {/each}
    </div>
    <div class="grid grid-cols-7">
        {#each weekdays as weekday}
            <div class="p-3 pb-6 {isToday(weekday.date) ? 'bg-primary-200 dark:bg-primary-900' : ''}">
                <PlannedRecipe title="My super duper recipe" />
            </div>
        {/each}
    </div>
</Card>