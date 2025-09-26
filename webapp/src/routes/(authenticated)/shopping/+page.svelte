<script lang="ts">
    import type { PageData } from './$types';
    import { onMount } from 'svelte';
    import { Button } from '$lib/components/ui/button';
    import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
    import { Input } from '$lib/components/ui/input';
    import { Badge } from '$lib/components/ui/badge';
    import ShoppingCartIcon from '@lucide/svelte/icons/shopping-cart';
    import CheckIcon from '@lucide/svelte/icons/check';
    import PlusIcon from '@lucide/svelte/icons/plus';
    import TrashIcon from '@lucide/svelte/icons/trash-2';
    import { flip } from 'svelte/animate';
    import { fly, scale } from 'svelte/transition';
    import { elasticOut, quintOut } from 'svelte/easing';
    import { toast } from 'svelte-sonner';
    
    import { shoppingStore } from '$lib/stores/shopping.svelte';
    import { addShoppingListItem, updateShoppingListItem, deleteShoppingListItem } from '$lib/api/shopping/shopping.svelte';
    import ShoppingListSelector from '$lib/components/shopping/ShoppingListSelector.svelte';
    import CreateShoppingListDialog from '$lib/components/shopping/CreateShoppingListDialog.svelte';
    import EmptyShoppingListState from '$lib/components/shopping/EmptyShoppingListState.svelte';
    import type { ReadShoppingListItem } from '$lib/api/shopping/shopping.svelte';

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
            toast.error('Failed to add item');
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
            toast.error('Failed to update item');
        }
    };

    const removeItem = async (itemId: number) => {
        if (!shoppingStore.currentListId) return;

        try {
            await deleteShoppingListItem(shoppingStore.currentListId, itemId);
            shoppingStore.removeItemFromList(shoppingStore.currentListId, itemId);
        } catch (error) {
            toast.error('Failed to delete item');
        }
    };

    const handleCreateList = () => {
        showCreateDialog = true;
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
            title: 'Done',
            icon: CheckIcon,
            bgColor: 'bg-green-50 dark:bg-green-950/20',
            borderColor: 'border-green-200 dark:border-green-800',
        } : {
            title: currentList?.name || 'Shopping List',
            icon: ShoppingCartIcon,
            bgColor: 'bg-orange-50 dark:bg-orange-950/20',
            borderColor: 'border-orange-200 dark:border-orange-800',
        };
    };
</script>


<svelte:head>
    <title>Shopping Lists</title>
</svelte:head>

{#if !hasLists}
    <EmptyShoppingListState onCreateList={handleCreateList} />
{:else}
    <div class="container mx-auto p-5">
        <ShoppingListSelector onCreateNewList={handleCreateList} />
        
        {#if currentList}
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {#each [false, true] as isDone (isDone)}
                    {@const config = getColumnConfig(isDone)}
                    {@const items = isDone ? doneItems : todoItems}
                    {@const IconComponent = config.icon}

                    <Card class="h-full gap-0 py-0">
                        <CardHeader class="px-3 border-b pb-0!">
                            <CardTitle class="pt-4 pb-2 flex items-center gap-2">
                                <IconComponent class="h-5 w-5"/>
                                {config.title}
                                <Badge variant="secondary" class="font-semibold text-sm rounded-full">
                                    {#key items.length}
                                        <span in:scale={{ duration: 200, easing: elasticOut }}>{items.length}</span>
                                    {/key}
                                </Badge>
                            </CardTitle>
                        </CardHeader>
                        <CardContent class="h-full p-1">
                            <div
                                class="h-full rounded-lg border-2 border-dashed border-transparent transition-all duration-200"
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
                                                        placeholder="1"
                                                        disabled={isAddingItem}
                                                    />
                                                    <Input
                                                        bind:value={newItemUnit}
                                                        placeholder="kg"
                                                        disabled={isAddingItem}
                                                    />
                                                    <Input
                                                        bind:value={newItemName}
                                                        placeholder="Potatoes"
                                                        class="col-span-2 md:col-span-1"
                                                        disabled={isAddingItem}
                                                        required
                                                    />
                                                    <Button type="submit" variant="ghost" class="col-span-2 md:col-span-1" disabled={isAddingItem || !newItemName.trim()}>
                                                        <PlusIcon/>
                                                        <span class="md:hidden">{isAddingItem ? 'Adding...' : 'Add Item'}</span>
                                                    </Button>
                                                </div>
                                            </form>
                                        </div>
                                    {/if}

                                    {#if items.length === 0}
                                        <div class="flex flex-col items-center justify-center flex-1 text-muted-foreground min-h-32">
                                            <IconComponent/>
                                            <p class="text-sm">No items</p>
                                        </div>
                                    {/if}
                                    {#each items as item (item.id)}
                                        <div
                                            class="transition-all duration-200"
                                            role="listitem"
                                            aria-label="Shopping item: {item.ingredient}"
                                            in:fly={{ x: !isDone ? -200 : 200, duration: 400, easing: quintOut }}
                                            out:fly={{ x: !isDone ? 200 : -200, duration: 300, easing: quintOut }}
                                            animate:flip={{ duration: 600, easing: quintOut }}
                                        >
                                            <div class="border flex items-center justify-between px-2 py-1.5 bg-muted/30 dark:bg-muted rounded-lg transition-all duration-200">
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
                                                            title="Got it"
                                                        >
                                                            <CheckIcon/>
                                                        </Button>
                                                    {:else}
                                                        <Button
                                                            size="sm"
                                                            variant="ghost"
                                                            onclick={() => toggleItemStatus(item)}
                                                            title="Need more"
                                                        >
                                                            <ShoppingCartIcon/>
                                                        </Button>
                                                    {/if}
                                                    <Button
                                                        variant="ghost"
                                                        size="sm"
                                                        onclick={() => removeItem(item.id)}
                                                        class="text-destructive hover:text-destructive"
                                                        title="Remove item"
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
/>