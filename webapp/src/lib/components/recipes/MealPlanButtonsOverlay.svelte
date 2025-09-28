<script lang="ts">
    import Trash2 from '@lucide/svelte/icons/trash-2'
    import ShoppingCart from '@lucide/svelte/icons/shopping-cart'
    import { Button } from "$lib/components/ui/button/index.js";
    import * as Dialog from "$lib/components/ui/dialog/index.js";
    import { toast } from "svelte-sonner";
    import { deleteMealPlan } from "$lib/api/mealplan/mealplan.svelte.js";

    let { recipe, date, onDeleted } = $props()

    let dialogOpen = $state(false);

    const handleDeleteClick = () => {
        dialogOpen = true;
    }

    const handleConfirmDelete = async () => {
        try {
            await deleteMealPlan(recipe.id, date);
            toast.success("Removed from meal plan");
            onDeleted?.(recipe.id);
            dialogOpen = false;
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
    onclick={handleDeleteClick}
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

<Dialog.Root bind:open={dialogOpen}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>Remove Recipe from Meal Plan</Dialog.Title>
            <Dialog.Description>
                Are you sure you want to remove "{recipe.name}" from your meal plan? This action cannot be undone.
            </Dialog.Description>
        </Dialog.Header>
        <Dialog.Footer>
            <Dialog.Close asChild>
                <Button variant="outline">
                    Cancel
                </Button>
            </Dialog.Close>
            <Button variant="destructive" onclick={handleConfirmDelete}>
                Remove
            </Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>
