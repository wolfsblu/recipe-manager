<script lang="ts">
    import {dndzone} from "svelte-dnd-action";
    import {Button, Card, SpeedDialTrigger} from "flowbite-svelte";
    import {
        AngleLeftOutline,
        AngleRightOutline,
        PlusOutline
    } from "flowbite-svelte-icons";
    import {isToday} from "$lib/calendar/calendar";
    import PlannedRecipe from "$lib/components/recipes/PlannedRecipe.svelte";
    import {DateTime} from "luxon";
    import {flip} from "svelte/animate";

    let selectedDate = $state(DateTime.now())

    const flipDuration = 100

    const onNextClick = () => {
        selectedDate = selectedDate.plus({day: 7})
    }
    const onPrevClick = () => {
        selectedDate = selectedDate.minus({day: 7})
    }

    let recipes = $state([
        {
            id: 1,
            name: 'Monday',
            date: new Date(2025, 3, 27),
            items: [
                {id: 10, name: 'My super duper recipe'},
            ]
        },
        {
            id: 2,
            name: 'Tuesday',
            date: new Date(2025, 3, 28),
            items: [
                {id: 20, name: 'My super duper recipe'},
            ]
        },
        {
            id: 3,
            name: 'Wednesday',
            date: new Date(2025, 3, 29),
            items: [
                {id: 30, name: 'My super duper recipe'},
            ]
        },
        {
            id: 4,
            name: 'Thursday',
            date: new Date(2025, 3, 30),
            items: [
                {id: 40, name: 'My super duper recipe'},
            ]
        },
        {
            id: 5,
            name: 'Friday',
            date: new Date(2025, 4, 1),
            items: [
                {id: 50, name: 'My super duper recipe'},
            ]
        },
        {
            id: 6,
            name: 'Saturday',
            date: new Date(2025, 4, 2),
            items: [
                {id: 60, name: 'My super duper recipe'},
            ]
        },
        {
            id: 7,
            name: 'Sunday',
            date: new Date(2025, 4, 3),
            items: [
                {id: 70, name: 'My super duper recipe'},
            ]
        },
    ])

    const handleDndConsiderCards = (cid: number, e) => {
        const colIdx = recipes.findIndex(c => c.id === cid);
        recipes[colIdx].items = e.detail.items;
        recipes = [...recipes];
    }
    const handleDndFinalizeCards = (cid: number, e) => {
        const colIdx = recipes.findIndex(c => c.id === cid);
        recipes[colIdx].items = e.detail.items;
        recipes = [...recipes];
    };
</script>

<Card class="max-w-full" >
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
        {#each recipes as recipe (recipe.id)}
            <div class="text-center dark:text-white flex flex-col
                        {isToday(DateTime.fromJSDate(recipe.date)) ? 'bg-gray-100 dark:bg-gray-900' : ''}">
                <div class="border-y py-1 border-gray-200 dark:border-gray-700">
                    <p class="font-thin mb-1">{DateTime.fromJSDate(recipe.date).weekdayShort}</p>
                    <p class="font-bold">{DateTime.fromJSDate(recipe.date).toLocaleString({day: "numeric"})}</p>
                </div>

                <div class="flex flex-col flex-grow space-y-3 p-3"
                     onconsider={(e) => handleDndConsiderCards(recipe.id, e)}
                     onfinalize={(e) => handleDndFinalizeCards(recipe.id, e)}
                     use:dndzone={{items: recipe.items, flipDurationMs: flipDuration, dropTargetStyle: {}}}
                >
                    {#each recipe.items as item (item.id)}
                        <div animate:flip={{duration: flipDuration}}>
                            <PlannedRecipe title={item.name}/>
                        </div>
                    {/each}
                </div>
            </div>
        {/each}
    </div>
</Card>

<SpeedDialTrigger class="absolute end-6 bottom-6">
    {#snippet icon()}
        <PlusOutline class="h-8 w-8" />
    {/snippet}
</SpeedDialTrigger>
