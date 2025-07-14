<script lang="ts">
    import * as Card from "$lib/components/ui/card/index.js";
    import * as Pagination from "$lib/components/ui/pagination/index.js";
    import RecipeCard from "$lib/components/recipes/RecipeCard.svelte";
</script>

<div class="flex flex-col gap-5 h-full justify-between p-5">
<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-5 gap-5">
    {#each Array(10) as _, i}
        <RecipeCard />
    {/each}
</div>
<Pagination.Root count={100} perPage={10}>
    {#snippet children({ pages, currentPage })}
        <Pagination.Content>
            <Pagination.Item>
                <Pagination.PrevButton />
            </Pagination.Item>
            {#each pages as page (page.key)}
                {#if page.type === "ellipsis"}
                    <Pagination.Item>
                        <Pagination.Ellipsis />
                    </Pagination.Item>
                {:else}
                    <Pagination.Item>
                        <Pagination.Link {page} isActive={currentPage === page.value}>
                            {page.value}
                        </Pagination.Link>
                    </Pagination.Item>
                {/if}
            {/each}
            <Pagination.Item>
                <Pagination.NextButton />
            </Pagination.Item>
        </Pagination.Content>
    {/snippet}
</Pagination.Root>
</div>