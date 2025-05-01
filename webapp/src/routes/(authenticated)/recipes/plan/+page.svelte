<script lang="ts">
    import {Button, Card} from "flowbite-svelte";
    import {AngleLeftOutline, AngleRightOutline} from "flowbite-svelte-icons";
    import {getWeekdays, isToday} from "$lib/calendar/calendar";
    import PlannedRecipe from "$lib/components/recipes/PlannedRecipe.svelte";
    import {DateTime} from "luxon";

    let selectedDate = $state(DateTime.now())
    const weekdays = $derived(getWeekdays(selectedDate))

    const onNextClick = () => {
        selectedDate = selectedDate.plus({day: 7})
    }
    const onPrevClick = () => {
        selectedDate = selectedDate.minus({day: 7})
    }
</script>

<Card padding="none" size="none">
    <div class="flex flex-col sm:flex-row justify-between items-center pb-4 gap-3 p-6">
        <h2 class="text-2xl font-semibold">
            {selectedDate.toLocaleString({month: "long"})}
            {selectedDate.toLocaleString({year: "numeric"})}
        </h2>
        <div class="flex space-x-2">
            <Button color="alternative" onclick={() => onPrevClick()}>
                <AngleLeftOutline class="w-5 h-5"/>
            </Button>
            <Button color="alternative" onclick={() => onNextClick()}>
                <AngleRightOutline class="w-5 h-5"/>
            </Button>
            <Button onclick={() => selectedDate = DateTime.now()}>
                Today
            </Button>
        </div>
    </div>

    <div class="grid grid-cols-7">
        {#each weekdays as weekday}
            <div class="text-center py-3 dark:text-white
                        {isToday(weekday.date) ? 'bg-primary-200 dark:bg-primary-900' : ''}">
                <p class="font-thin mb-1">{weekday.name.toUpperCase()}</p>
                <p class="font-bold">{weekday.date.toLocaleString({day: "numeric"})}</p>

                <div class="p-3 pb-6 {isToday(weekday.date) ? 'bg-primary-200 dark:bg-primary-900' : ''}">
                    <PlannedRecipe title="My super duper recipe"/>
                </div>
            </div>
        {/each}
    </div>
</Card>