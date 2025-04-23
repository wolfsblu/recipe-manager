<script lang="ts">
    import {Button, Card} from "flowbite-svelte";
    import {AngleLeftOutline, AngleRightOutline} from "flowbite-svelte-icons";
    import {getDaysInMonth, getWeekdays, isToday} from "$lib/calendar/calendar";

    const now = new Date();
    const days = getDaysInMonth(now);
    const weekdays = getWeekdays()
</script>

<Card size="none">
    <div class="border-b flex flex-col sm:flex-row justify-between items-center pb-4 gap-3">
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
    <div class="grid grid-cols-7 text-center text-gray-600 dark:text-gray-300 font-semibold border-b">
        {#each weekdays as weekday, i}
            <div class="block sm:hidden py-2">{weekday.substring(0, 1)}</div>
            <div class="hidden sm:block py-2">{weekday}</div>
        {/each}
    </div>
    <!-- Calendar grid -->
    <div class="grid grid-cols-7 gap-px bg-gray-200 dark:bg-gray-700">
        <!-- Week 1 -->
        {#each days as day, i}
            <div class="overflow-y-auto overflow-x-hidden flex-col gap-2 bg-white dark:bg-gray-800 h-12 sm:h-30 text-base p-1">
                <time class="sticky top-0 mb-1 text-gray-600 dark:text-gray-300 w-8 h-8 flex items-center justify-center rounded-full font-semibold
                      {isToday(day) ? 'bg-primary-300 text-primary-600 dark:bg-primary-600 dark:text-primary-300' : ''}"
                      datetime="2022-01-01">
                    {day.getDate()}
                </time>
                {#if i === 23}
                <div class="overflow-y-auto">
                    {#each Array(10).fill(0) as item, i}
                        <Card size="none" class="mb-2" padding="none">
                            <div class="text-center text-nowrap p-1">
                                <p class="overflow-hidden text-ellipsis">Schnitzel Braten mit Sose und Ketchup</p>
                            </div>
                        </Card>
                    {/each}
                </div>
                {/if}
            </div>
        {/each}
    </div>
</Card>