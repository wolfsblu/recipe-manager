<script lang="ts">
    import Trash2 from '@lucide/svelte/icons/trash-2'
    import ShoppingCart from '@lucide/svelte/icons/shopping-cart'
    import { Button } from "$lib/components/ui/button/index.js";
    import { toast } from "svelte-sonner";
    import { deleteMealPlan } from "$lib/api/mealplan/mealplan.svelte.js";

    let { recipe, date, onDeleted } = $props()

    const handleDeleteFromMealPlan = async (e: Event) => {
        try {
            await deleteMealPlan(recipe.id, date);
            toast.success("Removed from meal plan");
            onDeleted?.(recipe.id);
        } catch (error) {
            toast.error("Failed to remove from meal plan. Please try again.");
        }
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
