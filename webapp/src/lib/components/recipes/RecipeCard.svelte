<script lang="ts">
    import fruits from './fruits.jpg'
    import ClockIcon from '@lucide/svelte/icons/clock'
    import TagIcon from "@lucide/svelte/icons/tag";
    import PlusIcon from "@lucide/svelte/icons/plus";
    import { Badge } from "$lib/components/ui/badge/index.js";
    import { ScrollArea } from "$lib/components/ui/scroll-area/index.js";
    import { Separator } from "$lib/components/ui/separator/index.js";
    import Button from '../ui/button/button.svelte';

    let { recipe } = $props()
</script>

<div class="bg-card rounded-lg shadow-sm overflow-hidden transition-all duration-300">
    <a href="/recipes/1" class="block relative">
        <Badge variant="secondary" class="absolute top-2 left-2">
            <ClockIcon /> {recipe.minutes}m
        </Badge>
        <img src={fruits} alt="fruits" class="w-full h-52 object-cover" />
    </a>
    <div class="pt-3 px-4 pb-2 space-y-2">
        <a href="/recipes/1" class="inline-block font-semibold text-base">
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
                <Separator />
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
