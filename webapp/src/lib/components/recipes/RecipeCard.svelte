<script lang="ts">
    import fruits from './fruits.jpg'
    import ClockIcon from '@lucide/svelte/icons/clock'
    import { Badge } from "$lib/components/ui/badge/index.js";
    import { ScrollArea } from "$lib/components/ui/scroll-area/index.js";
    import { Separator } from "$lib/components/ui/separator/index.js";

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
            <Badge variant="secondary" class="text-muted-foreground">Spicy</Badge>
            <Badge variant="secondary" class="text-muted-foreground">Mediterranean</Badge>
            <Badge variant="secondary" class="text-muted-foreground">Lunch</Badge>
            <Badge variant="secondary" class="text-muted-foreground">Pasta</Badge>
            <Badge variant="secondary" class="text-muted-foreground">Organic</Badge>
            <Badge variant="secondary" class="text-muted-foreground">Healthy</Badge>
            <Badge variant="secondary" class="text-muted-foreground">Quick</Badge>
        </div>
    </ScrollArea>
</div>
