<script lang="ts">
    import fruits from './fruits.jpg'
    import ClockIcon from '@lucide/svelte/icons/clock'
    import { Badge } from "$lib/components/ui/badge/index.js";
    import { ScrollArea } from "$lib/components/ui/scroll-area/index.js";
    import { Separator } from "$lib/components/ui/separator/index.js";
</script>

<div class="bg-card rounded-lg shadow-sm overflow-hidden transition-all duration-300">
    <div class="relative">
        <Badge variant="secondary" class="absolute top-2 left-2">
            <ClockIcon /> 45m
        </Badge>
        <img src={fruits} alt="fruits" class="w-full h-52 object-cover" />
    </div>
    <div class="pt-4 px-4 pb-2 space-y-2">
        <h3 class="font-semibold text-base">
            My super duper tasty recipe
        </h3>
        <div class="text-sm text-muted-foreground grid space-y-1">
            <div class="flex justify-between">
                <span>Servings</span>
                <Badge>4</Badge>
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