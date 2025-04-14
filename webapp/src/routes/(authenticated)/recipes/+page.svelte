<script lang="ts">
    import {
        SpeedDial,
        SpeedDialButton,
        Search,
        Button,
        Card,
        Listgroup,
        ListgroupItem,
        Badge,
        Dropdown,
        DropdownItem,
        DropdownDivider
    } from "flowbite-svelte";
    import {
        ShoppingBagSolid,
        ReceiptSolid,
        RulerCombinedOutline,
        SearchOutline,
        DotsVerticalOutline
    } from "flowbite-svelte-icons";

    const iconClass = "w-6 h-6"
    const listItemClass = "px-0 flex justify-between"
</script>

<form class="flex gap-2 mb-4">
    <Search wrapperClass="relative w-full" class="h-full" placeholder="Find a recipe" />
    <Button class="cursor-pointer">
        <SearchOutline class="w-6 h-6 me-2 -ms-1" />
        Search
    </Button>
</form>

{#snippet DummyCard(index: number)}
    <Card img="/recipes/fruits.jpg" imgClass="max-h-42 object-cover" padding="sm" size="none" class="card">
        <div class="flex flex-col flex-grow justify-between">
            <section class="mb-2 flex gap-1 justify-between items-top">
                <h1 class="text-lg font-bold tracking-tight text-gray-900 dark:text-white">
                    {#if index === 2}
                        My super duper super duper super duper super tasty recipe
                    {:else}
                        My super duper tasty recipe
                    {/if}
                </h1>
                <Button class="h-7 w-7 p-2" outline pill>
                    <DotsVerticalOutline class="w-6 h-6" />
                </Button>
                <Dropdown>
                    <DropdownItem>Edit</DropdownItem>
                    <DropdownDivider />
                    <DropdownItem>Delete</DropdownItem>
                </Dropdown>
            </section>
            <Listgroup class="border-0 dark:bg-transparent!">
                <ListgroupItem normalClass={listItemClass}>
                    <span>Time</span>
                    <span>45 min</span>
                </ListgroupItem>
                <ListgroupItem normalClass={listItemClass}>
                    <span>Servings</span>
                    <Badge>4</Badge>
                </ListgroupItem>
                <ListgroupItem normalClass={listItemClass}>
                    <span>Ingredients</span>
                    <Badge>12</Badge>
                </ListgroupItem>
            </Listgroup>
        </div>
    </Card>
{/snippet}

<section class="grid gap-4 grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
    {#each {length: 4} as _, i}
        {@render DummyCard(i)}
    {/each}
</section>

<SpeedDial>
    <SpeedDialButton name="Measurement">
        <RulerCombinedOutline class={iconClass} />
    </SpeedDialButton>
    <SpeedDialButton name="Ingredient">
        <ShoppingBagSolid class={iconClass} />
    </SpeedDialButton>
    <SpeedDialButton name="Recipe">
        <ReceiptSolid class={iconClass} />
    </SpeedDialButton>
</SpeedDial>

<style>
    :global(.card > div) {
        display: flex;
        flex-grow: 1;
    }
</style>