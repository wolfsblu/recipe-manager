<script lang="ts">
    import RecipeCard from "$lib/components/recipes/RecipeCard.svelte";
    import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
    import { useInfiniteScroll } from "$lib/utils/infiniteScroll.svelte";
    import { getRecipes } from "$lib/api/recipes/recipes.svelte";
    import type { PageProps } from './$types';
    import { onMount } from "svelte";

    const { data }: PageProps = $props();

    const scroll = useInfiniteScroll({
        fetchPage: async (cursor) => {
            const response = await getRecipes({ cursor });
            return {
                data: response?.data ?? [],
                nextCursor: response?.nextCursor,
                hasMore: response?.hasMore ?? false
            };
        }
    });

    onMount(() => {
        // Load initial page
        scroll.loadMore();
    });
</script>

<div class="flex flex-col gap-5 h-full p-5">
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-5 gap-5">
        {#each scroll.items as recipe}
            <RecipeCard {recipe} availableTags={data.tags} />
        {/each}
    </div>

    <InfiniteScroll
        hasMore={scroll.hasMore}
        loading={scroll.loading}
        onLoadMore={scroll.loadMore}
    />
</div>