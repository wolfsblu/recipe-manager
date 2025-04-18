<script lang="ts">
    import {
        SpeedDial,
        SpeedDialButton,
        Search,
        Button
    } from "flowbite-svelte";
    import {
        ShoppingBagSolid,
        ReceiptSolid,
        RulerCombinedOutline,
        SearchOutline
    } from "flowbite-svelte-icons";
    import Recipe from "$lib/components/recipes/Recipe.svelte";

    const iconClass = "w-6 h-6"
</script>

<form class="flex gap-2 mb-4">
    <Search wrapperClass="relative w-full" class="h-full" placeholder="Find a recipe" />
    <Button class="cursor-pointer">
        <SearchOutline class="w-6 h-6 me-2 -ms-1" />
        Search
    </Button>
</form>

{#snippet DummyRecipe(title: string, servings: number, ingredients: number, minutes: number)}
    <Recipe title={title} servings={servings} ingredients={ingredients} minutes={minutes} />
{/snippet}

<section class="grid gap-4 grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
    {#each {length: 4} as _, i}
        {@render DummyRecipe("My super super tasty recipe", 4, 12, 45)}
    {/each}
</section>

<SpeedDial>
    <SpeedDialButton name="Measurement">
        <RulerCombinedOutline class={iconClass} />
    </SpeedDialButton>
    <SpeedDialButton name="Ingredient">
        <ShoppingBagSolid class={iconClass} />
    </SpeedDialButton>
    <SpeedDialButton name="Recipe" href="/recipes/create">
        <ReceiptSolid class={iconClass}  />
    </SpeedDialButton>
</SpeedDial>