<script lang="ts">
    import {
        SpeedDial,
        SpeedDialButton,
        SpeedDialTrigger,
        Search,
    } from "flowbite-svelte";
    import {
        ShoppingBagSolid,
        ReceiptSolid,
        RulerCombinedOutline,
        CartOutline,
    } from "flowbite-svelte-icons";
    import Recipe from "$lib/components/recipes/Recipe.svelte";

    const iconClass = "w-6 h-6";
</script>

<form class="flex gap-2 mb-4">
    <Search placeholder="Find a recipe" />
</form>

{#snippet DummyRecipe(
    title: string,
    servings: number,
    ingredients: number,
    minutes: number
)}
    <Recipe {title} {servings} {ingredients} {minutes} />
{/snippet}

<section
    class="grid gap-4 grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4"
>
    {#each { length: 4 } as _, i}
        {@render DummyRecipe("My super super tasty recipe", 4, 12, 45)}
    {/each}
</section>

<SpeedDialTrigger class="absolute end-6 bottom-6" />
<SpeedDial class="fixed end-6 bottom-6">
    <SpeedDialButton name="Shopping Item">
        <CartOutline class={iconClass} />
    </SpeedDialButton>
    <SpeedDialButton name="Measurement">
        <RulerCombinedOutline class={iconClass} />
    </SpeedDialButton>
    <SpeedDialButton name="Ingredient">
        <ShoppingBagSolid class={iconClass} />
    </SpeedDialButton>
    <SpeedDialButton name="Recipe" href="/recipes/create">
        <ReceiptSolid class={iconClass} />
    </SpeedDialButton>
</SpeedDial>
