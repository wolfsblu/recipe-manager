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
    import {TrashBinOutline} from "flowbite-svelte-icons";

    const listItemClass = "px-0 flex justify-between"

    interface Props {
        title: string
    }

    let {
        title
    }: Props = $props()
</script>

<div class="flex flex-col justify-between w-full h-36 bg-cover rounded-md shadow-md" style="background-image: url({placeholder})">
    <div class="flex justify-end p-1">
        <Button class="h-7 w-7 p-2" color="red" outline pill>
            <TrashBinOutline class="w-4 h-4"/>
        </Button>
    </div>
    <div class="flex flex-col justify-end py-3 h-1/2 text-white bg-gradient-to-b from-transparent to-gray-800 rounded-b-md">
        <h1 class="font-semibold">
            {title}
        </h1>
    </div>
</div>
<!--
<Card class="card" img={placeholder} imgClass="max-h-42 object-cover" padding="sm" size="none">
    <div class="flex flex-col flex-grow justify-between">
        <section class="mb-2 flex gap-1 justify-between items-top">
            <h1 class="text-lg font-bold tracking-tight text-gray-900 dark:text-white">
                {title}
            </h1>
            <Button class="h-7 w-7 p-2" outline pill>
                <DotsVerticalOutline class="w-6 h-6"/>
            </Button>
            <Dropdown placement="left">
                <DropdownItem>Edit</DropdownItem>
                <DropdownDivider/>
                <DropdownItem>Delete</DropdownItem>
            </Dropdown>
        </section>
    </div>
</Card>

<style>
    :global(.card > div) {
        display: flex;
        flex-grow: 1;
    }
</style>
-->