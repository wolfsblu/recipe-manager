<script lang="ts">
    import fruits from './fruits.jpg'
    import { Badge } from "$lib/components/ui/badge/index.js";
    import { Separator } from "$lib/components/ui/separator/index.js";
</script>

<a href="/recipes/1">
    <div class="bg-card rounded-lg shadow-sm overflow-hidden transition-all duration-300">
            <img src={fruits} alt="fruits" class="w-full h-52 object-cover" />
        <div class="p-4 space-y-2">
            <h3 class="font-semibold text-base">
                My super duper tasty recipe
            </h3>
            <div class="text-sm text-muted-foreground grid space-y-1">
                <div class="flex justify-between">
                    <span>Time</span>
                    <span>45 min</span>
                </div>
                <Separator />
                <div class="flex justify-between">
                    <span>Servings</span>
                    <Badge>4</Badge>
                </div>
                <Separator />
                <div class="flex justify-between">
                    <span>Ingredients</span>
                    <Badge>12</Badge>
                </div>
            </div>
        </div>
    </div>
</a>