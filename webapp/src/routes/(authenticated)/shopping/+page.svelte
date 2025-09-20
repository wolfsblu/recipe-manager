<script lang="ts">
    import { Button } from '$lib/components/ui/button';
    import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
    import { Input } from '$lib/components/ui/input';
    import { Badge } from '$lib/components/ui/badge';
    import PlusIcon from '@lucide/svelte/icons/plus';
    import ShoppingCartIcon from '@lucide/svelte/icons/shopping-cart';
    import TrashIcon from '@lucide/svelte/icons/trash-2';
    import CheckIcon from '@lucide/svelte/icons/check';
    import PackageIcon from '@lucide/svelte/icons/package';
    import MoveRightIcon from '@lucide/svelte/icons/move-right';
    import MoveLeftIcon from '@lucide/svelte/icons/move-left';
    import ListIcon from '@lucide/svelte/icons/list';
    import ChefHatIcon from '@lucide/svelte/icons/chef-hat';
    import { flip } from 'svelte/animate';
    import { fade, fly, scale, slide } from 'svelte/transition';
    import { quintOut, elasticOut } from 'svelte/easing';

    type ShoppingStatus = 'shopping' | 'pantry';

    interface ShoppingItem {
        id: string;
        name: string;
        quantity: string;
        category: string;
        status: ShoppingStatus;
        recipe?: string;
        urgent: boolean;
    }

    let newItemName = $state('');
    let newItemQuantity = $state('');
    let newItemUnit = $state('');
    let newItemCategory = $state('Other');
    let dragState = $state({ isDragging: false, draggedItem: null as ShoppingItem | null });
    let draggedItemId = $state<string | null>(null);
    let dragOverItemId = $state<string | null>(null);
    let dragOverEndZone = $state<string | null>(null);
    let shoppingItems = $state<ShoppingItem[]>([
        { id: '1', name: 'Milk', quantity: '1 gallon', category: 'Dairy', status: 'shopping', recipe: 'Pancakes', urgent: false },
        { id: '2', name: 'Bread', quantity: '1 loaf', category: 'Bakery', status: 'shopping', urgent: true },
        { id: '3', name: 'Apples', quantity: '2 lbs', category: 'Produce', status: 'shopping', recipe: 'Apple Pie', urgent: false },
        { id: '4', name: 'Chicken breast', quantity: '1 lb', category: 'Meat', status: 'shopping', recipe: 'Chicken Teriyaki', urgent: false },
        { id: '5', name: 'Rice', quantity: '2 cups', category: 'Pantry', status: 'pantry', recipe: 'Fried Rice', urgent: false },
        { id: '6', name: 'Yogurt', quantity: '1 container', category: 'Dairy', status: 'pantry', urgent: false }
    ]);

    function addItem() {
        if (!newItemName.trim()) return;

        // Combine quantity and unit for display
        let quantityDisplay = '';
        if (newItemQuantity.trim()) {
            quantityDisplay = newItemQuantity.trim();
            if (newItemUnit.trim()) {
                quantityDisplay += ' ' + newItemUnit.trim();
            }
        }

        const newItem: ShoppingItem = {
            id: Date.now().toString(),
            name: newItemName.trim(),
            quantity: quantityDisplay,
            category: newItemCategory,
            status: 'shopping',
            urgent: false
        };

        shoppingItems = [newItem, ...shoppingItems];
        newItemName = '';
        newItemQuantity = '';
        newItemUnit = '';
        newItemCategory = 'Other';
    }

    function moveItem(id: string, newStatus: ShoppingStatus) {
        shoppingItems = shoppingItems.map(item =>
            item.id === id ? { ...item, status: newStatus } : item
        );
    }

    function removeItem(id: string) {
        shoppingItems = shoppingItems.filter(item => item.id !== id);
    }

    function reorderItemsById(sourceId: string, targetId: string, status: ShoppingStatus) {
        const items = shoppingItems.filter(item => item.status === status);
        const otherItems = shoppingItems.filter(item => item.status !== status);

        const sourceIndex = items.findIndex(item => item.id === sourceId);
        const targetIndex = items.findIndex(item => item.id === targetId);

        if (sourceIndex === -1 || targetIndex === -1) return;

        // Remove the source item
        const [movedItem] = items.splice(sourceIndex, 1);

        // Adjust target index if we removed an item before it
        const adjustedTargetIndex = sourceIndex < targetIndex ? targetIndex - 1 : targetIndex;

        // Insert before the target item (where the indicator shows)
        items.splice(adjustedTargetIndex, 0, movedItem);

        shoppingItems = [...otherItems, ...items];
    }

    function insertItemAtPosition(sourceId: string, targetStatus: ShoppingStatus, targetId: string) {
        const sourceItem = shoppingItems.find(item => item.id === sourceId);
        if (!sourceItem) return;

        // Remove source item from original position
        shoppingItems = shoppingItems.filter(item => item.id !== sourceId);

        // Get items in target column
        const targetItems = shoppingItems.filter(item => item.status === targetStatus);
        const otherItems = shoppingItems.filter(item => item.status !== targetStatus);

        // Find position of target item
        const targetIndex = targetItems.findIndex(item => item.id === targetId);

        // Update source item status
        const updatedSourceItem = { ...sourceItem, status: targetStatus };

        // Insert at target position
        if (targetIndex !== -1) {
            targetItems.splice(targetIndex, 0, updatedSourceItem);
        } else {
            // If target not found, add to end
            targetItems.push(updatedSourceItem);
        }

        shoppingItems = [...otherItems, ...targetItems];
    }


    function handleDragStart(e: DragEvent, item: ShoppingItem) {
        if (!e.dataTransfer) return;

        e.dataTransfer.setData('text/plain', item.id);
        e.dataTransfer.effectAllowed = 'move';

        draggedItemId = item.id;
        dragState.isDragging = true;
        dragState.draggedItem = item;
    }

    function handleDragEnd() {
        draggedItemId = null;
        dragOverItemId = null;
        dragOverEndZone = null;
        dragState.isDragging = false;
        dragState.draggedItem = null;
    }

    function handleDragOver(e: DragEvent, targetItem?: ShoppingItem, isEndZone?: boolean, endZoneStatus?: string) {
        e.preventDefault();
        if (e.dataTransfer) {
            e.dataTransfer.dropEffect = 'move';
        }

        if (isEndZone && endZoneStatus) {
            dragOverItemId = null;
            dragOverEndZone = endZoneStatus;
        } else if (targetItem && targetItem.id !== draggedItemId) {
            dragOverEndZone = null;
            dragOverItemId = targetItem.id;
        }
    }

    function handleDragLeave(e: DragEvent) {
        // Only clear if we're leaving the item entirely
        const rect = (e.currentTarget as HTMLElement).getBoundingClientRect();
        const isLeavingItem = (
            e.clientX < rect.left || e.clientX > rect.right ||
            e.clientY < rect.top || e.clientY > rect.bottom
        );

        if (isLeavingItem) {
            dragOverItemId = null;
        }
    }

    function handleDrop(e: DragEvent, targetType: 'column' | 'item', target: ShoppingStatus | ShoppingItem) {
        e.preventDefault();
        if (!e.dataTransfer) return;

        const draggedId = e.dataTransfer.getData('text/plain');
        const draggedItem = shoppingItems.find(item => item.id === draggedId);
        if (!draggedItem) return;

        if (targetType === 'column') {
            // Moving between columns (drop on empty space)
            const targetStatus = target as ShoppingStatus;
            if (draggedItem.status !== targetStatus) {
                moveItem(draggedId, targetStatus);
            }
        } else {
            // Dropping on specific item
            const targetItem = target as ShoppingItem;
            if (draggedItem.status === targetItem.status && draggedId !== targetItem.id) {
                // Same column reordering
                reorderItemsById(draggedId, targetItem.id, draggedItem.status);
            } else if (draggedItem.status !== targetItem.status) {
                // Cross-column move with positioning
                insertItemAtPosition(draggedId, targetItem.status, targetItem.id);
            }
        }

        dragOverItemId = null;
        dragOverEndZone = null;
        handleDragEnd();
    }

    function handleEndZoneDrop(e: DragEvent, status: ShoppingStatus) {
        e.preventDefault();
        if (!e.dataTransfer) return;

        const draggedId = e.dataTransfer.getData('text/plain');
        const draggedItem = shoppingItems.find(item => item.id === draggedId);
        if (!draggedItem) return;

        if (draggedItem.status === status) {
            // Move to end of the same column
            const items = shoppingItems.filter(item => item.status === status);
            const otherItems = shoppingItems.filter(item => item.status !== status);
            const filteredItems = items.filter(item => item.id !== draggedId);
            filteredItems.push(draggedItem);
            shoppingItems = [...otherItems, ...filteredItems];
        } else {
            // Move to different column (add to end)
            moveItem(draggedId, status);
        }

        dragOverEndZone = null;
        handleDragEnd();
    }

    function isValidDropTarget(targetStatus: ShoppingStatus): boolean {
        return dragState.draggedItem !== null;
    }

    function getDropZoneClass(status: ShoppingStatus): string {
        if (!dragState.isDragging || !dragState.draggedItem) return 'border-transparent';

        const isSameColumn = dragState.draggedItem.status === status;
        const isDifferentColumn = dragState.draggedItem.status !== status;

        if (isSameColumn) {
            return status === 'shopping' ? 'border-orange-300 bg-orange-50/30 border-opacity-60' : 'border-green-300 bg-green-50/30 border-opacity-60';
        } else if (isDifferentColumn) {
            return 'border-primary bg-primary/5 border-opacity-50';
        }

        return 'border-transparent';
    }

    function getColumnConfig(status: ShoppingStatus) {
        const configs = {
            'shopping': {
                title: 'Shopping List',
                icon: ShoppingCartIcon,
                bgColor: 'bg-orange-50 dark:bg-orange-950/20',
                borderColor: 'border-orange-200 dark:border-orange-800',
            },
            'pantry': {
                title: 'Done',
                icon: CheckIcon,
                bgColor: 'bg-green-50 dark:bg-green-950/20',
                borderColor: 'border-green-200 dark:border-green-800',
            }
        };
        return configs[status];
    }

    const shoppingListItems = $derived(shoppingItems.filter(item => item.status === 'shopping'));
    const pantryItems = $derived(shoppingItems.filter(item => item.status === 'pantry'));
</script>

<svelte:head>
    <title>Shopping List</title>
</svelte:head>

<div class="container mx-auto p-5">


    <!-- Kanban Board -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {#each ['shopping', 'pantry'] as status (status)}
            {@const config = getColumnConfig(status)}
            {@const items = status === 'shopping' ? shoppingListItems : pantryItems}
            {@const IconComponent = config.icon}

            <Card class="h-fit gap-0 py-0">
                <CardHeader class="px-3 border-b pb-0!">
                    <CardTitle class="pt-4 pb-2 flex items-center gap-2">
                        <IconComponent class="h-5 w-5" />
                        {config.title}
                        <Badge variant="secondary" class="font-semibold text-sm rounded-full">
                            {#key items.length}
                                <span in:scale={{ duration: 200, easing: elasticOut }}>{items.length}</span>
                            {/key}
                        </Badge>
                    </CardTitle>
                </CardHeader>
                <CardContent class="p-2">
                    <div
                        class="rounded-lg border-2 border-dashed transition-all duration-200 {getDropZoneClass(status)}"
                        role="region"
                        aria-label="{config.title} drop zone"
                        ondragover={handleDragOver}
                        ondrop={(e) => handleDrop(e, 'column', status)}
                    >
                    {#if items.length === 0}
                        <div
                            class="flex flex-col items-center justify-center h-32 text-muted-foreground"
                            in:fade={{ duration: 300, delay: 200 }}
                            out:fade={{ duration: 200 }}
                        >
                            <IconComponent
                                class="h-8 w-8 mb-2 opacity-50 transition-all duration-300"
                            />
                            <p class="text-sm">No items</p>
                            {#if dragState.isDragging && isValidDropTarget(status)}
                                <p
                                    class="text-xs mt-2 text-primary font-medium"
                                    in:scale={{ duration: 200, easing: elasticOut }}
                                >
                                    Drop here
                                </p>
                            {/if}
                        </div>
                    {:else}
                        <div class="space-y-2" role="list" aria-label="{config.title} items">
                            {#if status === 'shopping'}
                                <div class="border flex items-center gap-3 px-3 py-2 bg-muted/10 dark:bg-muted/5 rounded-lg border-dashed">
                                    <form onsubmit={(e) => { e.preventDefault(); addItem(); }} class="flex items-center gap-2 w-full">
                                        <div class="w-full grid grid-cols-2 md:grid-cols-[2fr_4fr_6fr_min-content] items-center gap-2">
                                            <Input
                                                bind:value={newItemQuantity}
                                                placeholder="1"
                                            />
                                            <Input
                                                bind:value={newItemUnit}
                                                placeholder="kg"
                                            />
                                            <Input
                                                bind:value={newItemName}
                                                placeholder="Potatoes"
                                                class="col-span-2 md:col-span-1"
                                                required
                                            />
                                            <Button type="submit" variant="ghost" class="col-span-2 md:col-span-1">
                                                <PlusIcon />
                                                <span class="md:hidden">Add Item</span>
                                            </Button>
                                        </div>
                                    </form>
                                </div>
                            {/if}

                            {#each items as item (item.id)}
                                <div
                                    class="cursor-move transition-all duration-200 {draggedItemId === item.id ? 'opacity-50 scale-95' : 'opacity-100 scale-100'}"
                                    role="listitem"
                                    aria-label="Draggable item: {item.name}"
                                    draggable="true"
                                    ondragstart={(e) => handleDragStart(e, item)}
                                    ondragend={handleDragEnd}
                                    ondragover={(e) => handleDragOver(e, item)}
                                    ondragleave={handleDragLeave}
                                    ondrop={(e) => handleDrop(e, 'item', item)}
                                    in:fly={{ x: status === 'shopping' ? -200 : 200, duration: 400, easing: quintOut }}
                                    out:fly={{ x: status === 'shopping' ? 200 : -200, duration: 300, easing: quintOut }}
                                    animate:flip={{ duration: 600, easing: quintOut }}
                                >
                                    <!-- Drop indicator -->
                                    {#if dragOverItemId === item.id && draggedItemId !== item.id && dragState.isDragging}
                                        <div class="h-0.5 bg-primary rounded-full mx-2 mb-2 transition-all duration-200" in:scale={{ duration: 200 }}></div>
                                    {/if}

                                    <div class="border flex items-center justify-between px-3 py-2 bg-muted/30 dark:bg-muted hover:bg-muted/50 dark:hover:bg-muted/30 rounded-lg transition-all duration-200 {draggedItemId === item.id ? 'opacity-50 scale-95 border-primary' : ''}">
                                        <span class="font-medium text-foreground">{item.name}</span>

                                        <div class="flex items-center gap-0">
                                            {#if item.quantity}
                                                <Badge variant="outline" class="text-sm mr-1">{item.quantity}</Badge>
                                            {/if}
                                            {#if status === 'shopping'}
                                                <Button
                                                    size="sm"
                                                    variant="ghost"
                                                    onclick={() => moveItem(item.id, 'pantry')}
                                                    title="Got it"
                                                >
                                                    <CheckIcon />
                                                </Button>
                                            {:else if status === 'pantry'}
                                                <Button
                                                    size="sm"
                                                    variant="ghost"
                                                    onclick={() => moveItem(item.id, 'shopping')}
                                                    title="Need more"
                                                >
                                                    <ShoppingCartIcon />
                                                </Button>
                                            {/if}
                                            <Button
                                                variant="ghost"
                                                size="sm"
                                                onclick={() => removeItem(item.id)}
                                                class="text-destructive hover:text-destructive"
                                                title="Remove item"
                                            >
                                                <TrashIcon />
                                            </Button>
                                        </div>
                                    </div>
                                </div>
                            {/each}

                            <!-- End zone drop area -->
                            {#if dragState.isDragging}
                                <div
                                    class="h-8 rounded-lg transition-all duration-200"
                                    role="region"
                                    aria-label="Drop zone for end of {config.title}"
                                    ondragover={(e) => handleDragOver(e, undefined, true, status)}
                                    ondrop={(e) => handleEndZoneDrop(e, status)}
                                >
                                    {#if dragOverEndZone === status}
                                        <div class="h-0.5 bg-primary rounded-full mx-2 transition-all duration-200" in:scale={{ duration: 200 }}></div>
                                    {/if}
                                </div>
                            {/if}
                        </div>
                    {/if}
                    </div>
                </CardContent>

            </Card>
        {/each}
    </div>
</div>