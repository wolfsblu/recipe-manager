<script lang="ts">
    import Trash2 from '@lucide/svelte/icons/trash-2'
    import ShoppingCart from '@lucide/svelte/icons/shopping-cart'
    import { Button } from "$lib/components/ui/button/index.js";
    import * as Dialog from "$lib/components/ui/dialog/index.js";
    import { toast } from "svelte-sonner";
    import { deleteMealPlan } from "$lib/api/mealplan/mealplan.svelte.js";
    import * as m from "$lib/paraglide/messages.js";

    let { recipe, date, onDeleted } = $props()

    let dialogOpen = $state(false);

    const handleDeleteClick = () => {
        dialogOpen = true;
    }

    const handleConfirmDelete = async () => {
        try {
            await deleteMealPlan(recipe.id, date);
            toast.success(m.recipes_mealPlan_removedSuccess());
            onDeleted?.(recipe.id);
            dialogOpen = false;
        } catch (error) {
            toast.error(m.recipes_mealPlan_removedError());
        }
    }

    const handleAddToShoppingList = async (e: Event) => {
        // TODO: Implement add to shopping list functionality
        toast.success(m.shopping_addToList());
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
            <Dialog.Title>{m.recipes_mealPlan_removeDialog_title()}</Dialog.Title>
            <Dialog.Description>
                {m.recipes_mealPlan_removeDialog_description({ recipeName: recipe.name, date })}
            </Dialog.Description>
        </Dialog.Header>
        <Dialog.Footer>
            <Dialog.Close asChild>
                <Button variant="outline">
                    {m.recipes_mealPlan_removeDialog_cancel()}
                </Button>
            </Dialog.Close>
            <Button variant="destructive" onclick={handleConfirmDelete}>
                {m.recipes_mealPlan_removeDialog_remove()}
            </Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>
