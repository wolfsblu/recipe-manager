<script lang="ts">
    import Trash2 from '@lucide/svelte/icons/trash-2'
    import ShoppingCart from '@lucide/svelte/icons/shopping-cart'
    import ShoppingCartIcon from '@lucide/svelte/icons/shopping-cart'
    import Check from '@lucide/svelte/icons/check'
    import ChevronDown from '@lucide/svelte/icons/chevron-down'
    import PlusIcon from '@lucide/svelte/icons/plus'
    import { Button } from "$lib/components/ui/button/index.js";
    import { Input } from "$lib/components/ui/input/index.js";
    import { Label } from "$lib/components/ui/label/index.js";
    import * as Dialog from "$lib/components/ui/dialog/index.js";
    import * as Popover from "$lib/components/ui/popover/index.js";
    import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
    import { toast } from "svelte-sonner";
    import { deleteMealPlan } from "$lib/api/mealplan/mealplan.svelte.js";
    import { addShoppingListItem, type ReadShoppingList } from "$lib/api/shopping/shopping.svelte.js";
    import CreateShoppingListDialog from "$lib/components/shopping/CreateShoppingListDialog.svelte";
    import * as m from "$lib/paraglide/messages.js";

    let { recipe, date, shoppingLists = [], onDeleted } = $props<{
        recipe: any;
        date: string;
        shoppingLists?: ReadShoppingList[];
        onDeleted?: (recipeId: number) => void;
    }>();

    let dialogOpen = $state(false);
    let popoverOpen = $state(false);
    let selectedListId = $state<number | null>(null);
    let servings = $state<number>(recipe.servings || 1);
    let isAddingToList = $state(false);
    let localShoppingLists = $state<ReadShoppingList[]>([]);
    let showCreateDialog = $state(false);

    $effect(() => {
        localShoppingLists = shoppingLists;
    });

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

    const handleAddToShoppingListClick = () => {
        selectedListId = localShoppingLists.length > 0 ? localShoppingLists[0].id : null;
        servings = recipe.servings || 1;
        popoverOpen = true;
    }

    const handleConfirmAddToList = async () => {
        if (!selectedListId) {
            toast.error(m.shopping_mealPlan_errorNoList());
            return;
        }

        isAddingToList = true;
        try {
            const allIngredients = recipe.steps?.flatMap(step => step.ingredients || []) || [];
            const servingsMultiplier = servings / (recipe.servings || 1);
            for (const ingredient of allIngredients) {
                const adjustedAmount = ingredient.amount * servingsMultiplier;
                await addShoppingListItem(selectedListId, {
                    ingredient: ingredient.name || '',
                    quantity: adjustedAmount?.toString(),
                    unit: ingredient.unit?.name || '',
                    done: false,
                });
            }

            toast.success(m.shopping_mealPlan_successMessage({ count: allIngredients.length }));
            popoverOpen = false;
        } catch (error) {
            toast.error(m.shopping_mealPlan_errorMessage());
        } finally {
            isAddingToList = false;
        }
    }

    const handleCreateList = () => {
        popoverOpen = false;
        showCreateDialog = true;
    }

    const handleListCreated = (newList: ReadShoppingList) => {
        // Add the new list to local state
        localShoppingLists = [...localShoppingLists, newList];

        // Select the newly created list
        selectedListId = newList.id;

        showCreateDialog = false;
        popoverOpen = true;
    };

    const selectedList = $derived(localShoppingLists.find(list => list.id === selectedListId));
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

<Popover.Root bind:open={popoverOpen}>
    <Popover.Trigger>
        {#snippet child({ props })}
            <Button
                {...props}
                variant="ghost"
                size="sm"
                class="h-6 w-6 p-0 hover:bg-transparent"
                onclick={handleAddToShoppingListClick}
                title="Add to shopping list"
            >
                <ShoppingCart />
            </Button>
        {/snippet}
    </Popover.Trigger>
    <Popover.Content class="w-80">
        {#if localShoppingLists.length === 0}
            <div class="text-center py-6">
                <div class="p-3 bg-muted rounded-full mx-auto w-fit mb-3">
                    <ShoppingCartIcon class="h-6 w-6 text-muted-foreground" />
                </div>
                <h4 class="font-medium mb-1">{m.shopping_empty_title()}</h4>
                <p class="text-sm text-muted-foreground mb-4">
                    {m.shopping_empty_description()}
                </p>
                <Button onclick={handleCreateList} class="w-full">
                    <PlusIcon />
                    {m.shopping_empty_button()}
                </Button>
            </div>
        {:else}
            <div class="space-y-4">
                <div class="space-y-2">
                    <h4 class="font-medium leading-none">{m.shopping_mealPlan_addToList_title()}</h4>
                    <p class="text-sm text-muted-foreground">
                        {m.shopping_mealPlan_addToList_description()}
                    </p>
                </div>

                <div class="space-y-2">
                    <Label for="shopping-list">{m.shopping_list_title()}</Label>
                    <DropdownMenu.Root>
                        <DropdownMenu.Trigger>
                            {#snippet child({ props })}
                                <Button {...props} variant="outline" class="w-full justify-between">
                                    <span class="truncate">
                                        {selectedList ? selectedList.name : m.shopping_selector_selectList()}
                                    </span>
                                    <ChevronDown class="h-4 w-4 ml-2" />
                                </Button>
                            {/snippet}
                        </DropdownMenu.Trigger>
                        <DropdownMenu.Content class="w-72">
                            {#each localShoppingLists as list (list.id)}
                                <DropdownMenu.Item
                                    class="flex items-center justify-between cursor-pointer"
                                    onclick={() => selectedListId = list.id}
                                >
                                    <span class="truncate">{list.name}</span>
                                    {#if list.id === selectedListId}
                                        <Check class="h-4 w-4 text-primary" />
                                    {/if}
                                </DropdownMenu.Item>
                            {/each}
                        </DropdownMenu.Content>
                    </DropdownMenu.Root>
                </div>

                <div class="space-y-2">
                    <Label for="servings">{m.recipes_form_servings()}</Label>
                    <Input
                        id="servings"
                        type="number"
                        min="1"
                        bind:value={servings}
                        placeholder={m.recipes_form_servingsPlaceholder()}
                    />
                </div>

                <Button
                    class="w-full"
                    onclick={handleConfirmAddToList}
                    disabled={!selectedListId || isAddingToList}
                >
                    {isAddingToList ? m.shopping_list_adding() : m.shopping_mealPlan_addButton()}
                </Button>
            </div>
        {/if}
    </Popover.Content>
</Popover.Root>

<!-- Delete Meal Plan Dialog -->
<Dialog.Root bind:open={dialogOpen}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>{m.recipes_mealPlan_removeDialog_title()}</Dialog.Title>
            <Dialog.Description>
                {m.recipes_mealPlan_removeDialog_description({ recipeName: recipe.name, date })}
            </Dialog.Description>
        </Dialog.Header>
        <Dialog.Footer>
            <Dialog.Close>
                <Button variant="outline" class="w-full">
                    {m.recipes_mealPlan_removeDialog_cancel()}
                </Button>
            </Dialog.Close>
            <Button variant="destructive" onclick={handleConfirmDelete}>
                {m.recipes_mealPlan_removeDialog_remove()}
            </Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>

<!-- Create Shopping List Dialog -->
<CreateShoppingListDialog
    open={showCreateDialog}
    onOpenChange={(open) => showCreateDialog = open}
    onSuccess={handleListCreated}
/>
