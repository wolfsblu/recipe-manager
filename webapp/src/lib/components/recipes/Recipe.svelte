<script lang="ts">
    import placeholder from '$lib/components/recipes/fruits.jpg'
    import {
        Badge,
        Button,
        Card,
        Dropdown,
        DropdownDivider,
        DropdownItem,
        Listgroup,
        ListgroupItem
    } from "flowbite-svelte";
    import {DotsVerticalOutline} from "flowbite-svelte-icons";

    const listItemClass = "px-0 flex justify-between"

    interface Props {
        title: string
        ingredients: number
        servings: number
        minutes: number
    }

    let {
        title,
        ingredients,
        servings,
        minutes,
    }: Props = $props()
</script>

<Card class="card" img={placeholder} imgClass="max-h-42 object-cover" padding="sm" size="none">
    <div class="flex flex-col flex-grow justify-between">
        <section class="mb-2 flex gap-1 justify-between items-top">
            <h1 class="text-lg font-bold tracking-tight text-gray-900 dark:text-white">
                {title}
            </h1>
            <Button class="h-7 w-7 p-2" outline pill>
                <DotsVerticalOutline class="w-6 h-6"/>
            </Button>
            <Dropdown>
                <DropdownItem>Edit</DropdownItem>
                <DropdownDivider/>
                <DropdownItem>Delete</DropdownItem>
            </Dropdown>
        </section>
        <Listgroup class="border-0 dark:bg-transparent!">
            <ListgroupItem normalClass={listItemClass}>
                <span>Time</span>
                <span>{minutes} min</span>
            </ListgroupItem>
            <ListgroupItem normalClass={listItemClass}>
                <span>Servings</span>
                <Badge>{servings}</Badge>
            </ListgroupItem>
            <ListgroupItem normalClass={listItemClass}>
                <span>Ingredients</span>
                <Badge>{ingredients}</Badge>
            </ListgroupItem>
        </Listgroup>
    </div>
</Card>

<style>
    :global(.card > div) {
        display: flex;
        flex-grow: 1;
    }
</style>