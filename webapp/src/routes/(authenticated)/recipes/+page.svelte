<script lang="ts">
    import * as Card from "$lib/components/ui/card/index.js";
    import * as Pagination from "$lib/components/ui/pagination/index.js";
    import RecipeCard from "$lib/components/recipes/RecipeCard.svelte";
    import type { PageProps } from './$types'

    const { data }: PageProps = $props()
</script>

<div class="flex flex-col gap-5 h-full justify-between p-5">
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-5 gap-5">
        {#each data.recipes as recipe, i}
            <RecipeCard {recipe} availableTags={data.tags} />
        {/each}
    </div>
</div>