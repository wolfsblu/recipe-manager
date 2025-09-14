<script lang="ts">
    import fruits from './fruits.jpg'
    import ClockIcon from '@lucide/svelte/icons/clock'
    import { Badge } from "$lib/components/ui/badge/index.js";
    import { ScrollArea } from "$lib/components/ui/scroll-area/index.js";
    import { Separator } from "$lib/components/ui/separator/index.js";

    let { recipe } = $props()

    const formatMinutesAsHours = (minutes: number) => {
        if (!minutes || minutes < 0) {
            return "N/A";
        }

        if (minutes < 60) {
            return `${minutes}m`;
        }

        const hours = minutes / 60;

        // Round to a maximum of 2 decimal places.
        // We multiply by 100, round the result, then divide by 100.
        // This avoids trailing zeros for numbers like 1.5 or 2.0.
        const roundedHours = Math.round(hours * 100) / 100;

        return `${roundedHours}h`;
    }
</script>

<div class="bg-card rounded-lg shadow-sm overflow-hidden transition-all duration-300">
    <a href="/recipes/{recipe.id}" class="block relative">
        {#if recipe.minutes}
            <Badge variant="secondary" class="absolute top-2 left-2">
                <ClockIcon /> {formatMinutesAsHours(recipe.minutes)}
            </Badge>
        {/if}
        {#if recipe.images && recipe.images.length > 0}
            <img loading="lazy" src={recipe.images[0]} alt="Recipe" class="w-full h-52 object-cover" />
        {:else}
            <img loading="lazy" src={fruits} alt="Fruits" class="w-full h-52 object-cover" />
        {/if}
    </a>
    <div class="pt-3 px-4 pb-2 space-y-2">
        <a href="/recipes/{recipe.id}" class="inline-block font-semibold text-base">
            {recipe.name}
        </a>
        <div class="text-sm text-muted-foreground grid space-y-1">
            <div class="flex justify-between">
                <span>Servings</span>
                <Badge>{recipe.servings}</Badge>
            </div>
            <Separator />
            <div class="flex justify-between">
                <span>Ingredients</span>
                <Badge>12</Badge>
            </div>
            {#if recipe.tags}
                <Separator />
            {/if}
        </div>
    </div>
    <ScrollArea class="whitespace-nowrap pb-4 px-4" orientation="horizontal">
        <div class="flex w-max gap-x-1">
            {#each recipe.tags as tag}
                <Badge variant="secondary" class="text-muted-foreground">
                    {tag}
                </Badge>
            {/each}
        </div>
    </ScrollArea>
</div>
