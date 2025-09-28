<script lang="ts">
    import Trash2 from '@lucide/svelte/icons/trash-2'
    import ShoppingCart from '@lucide/svelte/icons/shopping-cart'
    import { Button } from "$lib/components/ui/button/index.js";
    import { toast } from "svelte-sonner";

    let { recipe } = $props()

    const handleDeleteFromMealPlan = async (e: Event) => {
        // TODO: Implement delete from meal plan functionality
        toast.success("Removed from meal plan");
    }

    const handleAddToShoppingList = async (e: Event) => {
        // TODO: Implement add to shopping list functionality
        toast.success("Added to shopping list");
    }
</script>

<Button
    variant="ghost"
    size="sm"
    class="h-6 w-6 p-0 hover:bg-transparent"
    onclick={handleDeleteFromMealPlan}
    title="Remove from meal plan"
>
    <Trash2 />
</Button>
<Button
    variant="ghost"
    size="sm"
    class="h-6 w-6 p-0 hover:bg-transparent"
    onclick={handleAddToShoppingList}
    title="Add to shopping list"
>
    <ShoppingCart />
</Button>
