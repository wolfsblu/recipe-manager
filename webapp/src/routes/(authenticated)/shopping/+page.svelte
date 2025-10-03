<script lang="ts">
    import type { PageData } from './$types';
    import { onMount } from 'svelte';
    import { Button } from '$lib/components/ui/button';
    import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
    import { Input } from '$lib/components/ui/input';
    import { Badge } from '$lib/components/ui/badge';
    import ShoppingCartIcon from '@lucide/svelte/icons/shopping-cart';
    import CheckIcon from '@lucide/svelte/icons/check';
    import CheckCheckIcon from '@lucide/svelte/icons/check-check';
    import PlusIcon from '@lucide/svelte/icons/plus';
    import TrashIcon from '@lucide/svelte/icons/trash-2';
    import { toast } from 'svelte-sonner';
    import * as m from "$lib/paraglide/messages.js";

    import { shoppingStore } from '$lib/stores/shopping.svelte';
    import { addShoppingListItem, updateShoppingListItem, deleteShoppingListItem, type ReadShoppingList, type ReadShoppingListItem } from '$lib/api/shopping/shopping.svelte';
    import ShoppingListSelector from '$lib/components/shopping/ShoppingListSelector.svelte';
    import CreateShoppingListDialog from '$lib/components/shopping/CreateShoppingListDialog.svelte';
    import EmptyShoppingListState from '$lib/components/shopping/EmptyShoppingListState.svelte';

    interface Props {
        data: PageData;
    }

    let { data }: Props = $props();

    let newItemName = $state('');
    let newItemQuantity = $state('');
    let newItemUnit = $state('');
    let isAddingItem = $state(false);
    let showCreateDialog = $state(false);

    // Initialize store with loaded data
    onMount(() => {
        if (data.shoppingLists) {
            shoppingStore.setLists(data.shoppingLists);
        }
        if (data.error) {
            toast.error(data.error);
        }
    });

    const addItem = async () => {
        if (!newItemName.trim() || !shoppingStore.currentListId) return;

        const quantityStr = newItemQuantity.trim() || null;
        const unitStr = newItemUnit.trim() || null;

        const newItem = {
            ingredient: newItemName.trim(),
            quantity: quantityStr,
            unit: unitStr,
            done: false
        };

        isAddingItem = true;
        try {
            const createdItem = await addShoppingListItem(shoppingStore.currentListId, newItem);
            shoppingStore.addItemToList(shoppingStore.currentListId, createdItem);
            
            newItemName = '';
            newItemQuantity = '';
            newItemUnit = '';
        } catch (error) {
            toast.error(m.shopping_errors_addItem());
        } finally {
            isAddingItem = false;
        }
    };

    const toggleItemStatus = async (item: ReadShoppingListItem) => {
        if (!shoppingStore.currentListId) return;

        try {
            const updatedItem = await updateShoppingListItem(
                shoppingStore.currentListId,
                item.id,
                {
                    ingredient: item.ingredient,
                    quantity: item.quantity,
                    unit: item.unit,
                    done: !item.done
                }
            );
            shoppingStore.updateItemInList(shoppingStore.currentListId, updatedItem);
        } catch (error) {
            toast.error(m.shopping_errors_updateItem());
        }
    };

    const removeItem = async (itemId: number) => {
        if (!shoppingStore.currentListId) return;

        try {
            await deleteShoppingListItem(shoppingStore.currentListId, itemId);
            shoppingStore.removeItemFromList(shoppingStore.currentListId, itemId);
        } catch (error) {
            toast.error(m.shopping_errors_deleteItem());
        }
    };

    const handleCreateList = () => {
        showCreateDialog = true;
    };

    const handleListCreated = (newList: ReadShoppingList) => {
        shoppingStore.addList(newList);
        shoppingStore.setCurrentListId(newList.id);
    };

    const handleMarkAllAsDone = async () => {
        if (!shoppingStore.currentListId || todoItems.length === 0) return;

        const listId = shoppingStore.currentListId;
        const count = todoItems.length;
        try {
            for (const item of todoItems) {
                const updatedItem = await updateShoppingListItem(listId, item.id, {
                    ingredient: item.ingredient,
                    quantity: item.quantity,
                    unit: item.unit,
                    done: true
                });
                shoppingStore.updateItemInList(listId, updatedItem);
            }

            toast.success(`Marked ${count} items as done`);
        } catch (error) {
            toast.error('Failed to mark items as done');
        }
    };

    const handleDeleteAllDone = async () => {
        if (!shoppingStore.currentListId || doneItems.length === 0) return;

        const listId = shoppingStore.currentListId;
        const count = doneItems.length;
        try {
            for (const item of doneItems) {
                await deleteShoppingListItem(listId, item.id);
                shoppingStore.removeItemFromList(listId, item.id);
            }

            toast.success(`Deleted ${count} done items`);
        } catch (error) {
            toast.error('Failed to delete done items');
        }
    };

    const formatQuantity = (quantity: string | null, unit: string | null) => {
        if (!quantity) return '';
        return unit ? `${quantity} ${unit}` : quantity;
    };

    const currentList = $derived(shoppingStore.currentList);
    const hasLists = $derived(shoppingStore.hasLists);
    const todoItems = $derived(currentList ? currentList.items.filter(item => !item.done) : []);
    const doneItems = $derived(currentList ? currentList.items.filter(item => item.done) : []);

    const getColumnConfig = (isDone: boolean) => {
        return isDone ? {
            title: m.shopping_list_done(),
            icon: CheckIcon,
            bgColor: 'bg-green-50 dark:bg-green-950/20',
            borderColor: 'border-green-200 dark:border-green-800',
        } : {
            title: currentList?.name || m.shopping_list_title(),
            icon: ShoppingCartIcon,
            bgColor: 'bg-orange-50 dark:bg-orange-950/20',
            borderColor: 'border-orange-200 dark:border-orange-800',
        };
    };
</script>


<svelte:head>
    <title>{m.shopping_title()}</title>
</svelte:head>

{#if !hasLists}
    <EmptyShoppingListState onCreateList={handleCreateList} />
{:else}
    <div class="flex flex-col gap-5 h-full p-5">
        <ShoppingListSelector onCreateNewList={handleCreateList} />
        
        {#if currentList}
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {#each [false, true] as isDone (isDone)}
                    {@const config = getColumnConfig(isDone)}
                    {@const items = isDone ? doneItems : todoItems}
                    {@const IconComponent = config.icon}

                    <Card class="h-full gap-0 py-0">
                        <CardHeader class="px-3 border-b pb-0!">
                            <CardTitle class="pt-4 pb-2 flex items-center gap-2 min-h-[52px]">
                                <IconComponent class="h-5 w-5"/>
                                {config.title}
                                <Badge variant="secondary" class="font-semibold text-sm rounded-full">
                                    {items.length}
                                </Badge>
                                <div class="ml-auto">
                                    {#if !isDone}
                                        <Button
                                            variant="ghost"
                                            size="sm"
                                            onclick={handleMarkAllAsDone}
                                            title="Mark all as done"
                                            disabled={items.length === 0}
                                            class={items.length === 0 ? 'invisible' : ''}
                                        >
                                            <CheckCheckIcon />
                                        </Button>
                                    {:else}
                                        <Button
                                            variant="ghost"
                                            size="sm"
                                            onclick={handleDeleteAllDone}
                                            title="Delete all done items"
                                            disabled={items.length === 0}
                                            class={items.length === 0 ? 'invisible' : ''}
                                        >
                                            <TrashIcon />
                                        </Button>
                                    {/if}
                                </div>
                            </CardTitle>
                        </CardHeader>
                        <CardContent class="h-full p-1">
                            <div
                                class="h-full rounded-lg border-2 border-dashed border-transparent"
                                role="region"
                                aria-label="{config.title}"
                            >
                                <div
                                    class="flex flex-col space-y-2 min-h-[100px] p-1 h-full"
                                    role="list"
                                    aria-label="{config.title} items"
                                >
                                    <!-- Always show add item form for shopping list -->
                                    {#if !isDone}
                                        <div class="border flex items-center gap-2 px-2 py-1.5 bg-muted/10 dark:bg-muted/5 rounded-lg border-dashed">
                                            <form onsubmit={(e) => { e.preventDefault(); addItem(); }}
                                                  class="flex items-center gap-2 w-full">
                                                <div class="w-full grid grid-cols-2 md:grid-cols-[2fr_4fr_6fr_min-content] items-center gap-2">
                                                    <Input
                                                        bind:value={newItemQuantity}
                                                        placeholder={m.shopping_list_addItemPlaceholder_quantity()}
                                                        disabled={isAddingItem}
                                                    />
                                                    <Input
                                                        bind:value={newItemUnit}
                                                        placeholder={m.shopping_list_addItemPlaceholder_unit()}
                                                        disabled={isAddingItem}
                                                    />
                                                    <Input
                                                        bind:value={newItemName}
                                                        placeholder={m.shopping_list_addItemPlaceholder_name()}
                                                        class="col-span-2 md:col-span-1"
                                                        disabled={isAddingItem}
                                                        required
                                                    />
                                                    <Button type="submit" variant="ghost" class="col-span-2 md:col-span-1" disabled={isAddingItem || !newItemName.trim()}>
                                                        <PlusIcon/>
                                                        <span class="md:hidden">{isAddingItem ? m.shopping_list_adding() : m.shopping_list_addButton()}</span>
                                                    </Button>
                                                </div>
                                            </form>
                                        </div>
                                    {/if}

                                    {#if items.length === 0}
                                        <div class="flex flex-col items-center justify-center flex-1 text-muted-foreground min-h-32">
                                            <IconComponent/>
                                            <p class="text-sm">{m.shopping_list_noItems()}</p>
                                        </div>
                                    {/if}
                                    {#each items as item (item.id)}
                                        <div
                                            role="listitem"
                                            aria-label="Shopping item: {item.ingredient}"
                                        >
                                            <div class="border flex items-center justify-between px-2 py-1.5 bg-muted/30 dark:bg-muted rounded-lg">
                                                <span class="font-medium text-foreground">{item.ingredient}</span>

                                                <div class="flex items-center gap-0">
                                                    {#if formatQuantity(item.quantity, item.unit)}
                                                        <Badge variant="outline" class="text-sm mr-1">{formatQuantity(item.quantity, item.unit)}</Badge>
                                                    {/if}
                                                    {#if !isDone}
                                                        <Button
                                                            size="sm"
                                                            variant="ghost"
                                                            onclick={() => toggleItemStatus(item)}
                                                            title={m.shopping_list_gotIt()}
                                                        >
                                                            <CheckIcon/>
                                                        </Button>
                                                    {:else}
                                                        <Button
                                                            size="sm"
                                                            variant="ghost"
                                                            onclick={() => toggleItemStatus(item)}
                                                            title={m.shopping_list_needMore()}
                                                        >
                                                            <ShoppingCartIcon/>
                                                        </Button>
                                                    {/if}
                                                    <Button
                                                        variant="ghost"
                                                        size="sm"
                                                        onclick={() => removeItem(item.id)}
                                                        class="text-destructive hover:text-destructive"
                                                        title={m.shopping_list_removeItem()}
                                                    >
                                                        <TrashIcon/>
                                                    </Button>
                                                </div>
                                            </div>
                                        </div>
                                    {/each}
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                {/each}
            </div>
        {/if}
    </div>
{/if}

<CreateShoppingListDialog
    open={showCreateDialog}
    onOpenChange={(open) => showCreateDialog = open}
    onSuccess={handleListCreated}
/>