<script lang="ts">
    import { Button } from '$lib/components/ui/button';
    import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
    import { Input } from '$lib/components/ui/input';
    import { Badge } from '$lib/components/ui/badge';
    import ShoppingCartIcon from '@lucide/svelte/icons/shopping-cart';
    import PlusIcon from '@lucide/svelte/icons/plus';
    import TrashIcon from '@lucide/svelte/icons/trash-2';
    import CheckIcon from '@lucide/svelte/icons/check';
    import { flip } from 'svelte/animate';
    import { fade, fly, scale } from 'svelte/transition';
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
    // Drag state
    let isDragging = $state(false);
    let draggedItem = $state<ShoppingItem | null>(null);
    let dragOverItem = $state<ShoppingItem | null>(null);
    let dragOverZone = $state<ShoppingStatus | null>(null);
    let ghostElement = $state<HTMLElement | null>(null);
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

    function moveItemToEnd(id: string, status: ShoppingStatus) {
        const items = shoppingItems.filter(item => item.status === status);
        const otherItems = shoppingItems.filter(item => item.status !== status);
        const draggedItem = items.find(item => item.id === id);

        if (!draggedItem) return;

        const filteredItems = items.filter(item => item.id !== id);
        filteredItems.push(draggedItem); // Add to end

        shoppingItems = [...otherItems, ...filteredItems];
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


    // Native pointer event handlers
    let currentPointerId: number | null = null;

    function handlePointerDown(e: PointerEvent, item: ShoppingItem) {
        if ((e.target as HTMLElement).closest('button')) return; // Don't drag if clicking buttons

        // Only handle primary pointer (first touch/mouse)
        if (currentPointerId !== null) return;

        e.preventDefault();
        e.stopPropagation();

        currentPointerId = e.pointerId;
        isDragging = true;
        draggedItem = item;

        const element = e.currentTarget as HTMLElement;
        element.setPointerCapture(e.pointerId);

        // Create ghost element
        createGhost(element, e.clientX, e.clientY);
    }

    function handlePointerMove(e: PointerEvent, element: HTMLElement) {
        if (!isDragging || !ghostElement || e.pointerId !== currentPointerId) return;

        e.preventDefault();

        // Update ghost position (minimize work in handler)
        requestAnimationFrame(() => {
            if (!ghostElement) return;
            ghostElement.style.left = `${e.clientX - 50}px`;
            ghostElement.style.top = `${e.clientY - 25}px`;
        });

        // Only update drop targets if we've actually moved the pointer
        // This prevents immediate indicator showing when drag starts
        const elementBelow = document.elementFromPoint(e.clientX, e.clientY);
        const currentElement = element.closest('[data-item-id]');

        // Only update if we're no longer over the original dragged element
        if (elementBelow && !currentElement?.contains(elementBelow)) {
            updateDropTargets(elementBelow);
        } else {
            // Clear targets when still over original element
            dragOverItem = null;
            dragOverZone = null;
        }
    }

    function handlePointerUp(e: PointerEvent, element: HTMLElement) {
        if (!isDragging || !draggedItem || e.pointerId !== currentPointerId) return;

        e.preventDefault();

        // Release pointer capture
        element.releasePointerCapture(e.pointerId);

        // Perform the drop
        const elementBelow = document.elementFromPoint(e.clientX, e.clientY);
        performDrop(elementBelow);

        // Cleanup
        cleanup();
    }

    function handlePointerCancel(e: PointerEvent, element: HTMLElement) {
        if (e.pointerId !== currentPointerId) return;

        // Release capture and cleanup
        element.releasePointerCapture(e.pointerId);
        cleanup();
    }

    function createGhost(originalElement: HTMLElement, x: number, y: number) {
        const ghost = originalElement.cloneNode(true) as HTMLElement;
        ghost.style.position = 'fixed';
        ghost.style.left = `${x - 50}px`;
        ghost.style.top = `${y - 25}px`;
        ghost.style.width = `${originalElement.offsetWidth}px`;
        ghost.style.height = `${originalElement.offsetHeight}px`;
        ghost.style.pointerEvents = 'none';
        ghost.style.zIndex = '9999';
        ghost.style.opacity = '0.8';
        ghost.style.transform = 'rotate(3deg)';
        ghost.style.boxShadow = '0 10px 20px rgba(0,0,0,0.3)';
        ghost.style.borderRadius = '8px';

        document.body.appendChild(ghost);
        ghostElement = ghost;
    }

    function updateDropTargets(element: Element | null) {
        dragOverItem = null;
        dragOverZone = null;

        if (!element) return;

        // Check if hovering over an item first (highest priority for showing between-item indicators)
        const itemElement = element.closest('[data-item-id]');
        if (itemElement && itemElement !== element.closest(`[data-item-id="${draggedItem?.id}"]`)) {
            const itemId = itemElement.getAttribute('data-item-id');
            const item = shoppingItems.find(i => i.id === itemId);
            if (item) {
                dragOverItem = item;
                return;
            }
        }

        // Check if hovering over an explicit end zone (second priority)
        const endZoneElement = element.closest('[data-end-zone="true"]');
        if (endZoneElement) {
            const zone = endZoneElement.getAttribute('data-zone') as ShoppingStatus;
            if (zone) {
                dragOverZone = zone;
                return;
            }
        }

        // Check if hovering over a zone with different status (for cross-lane moves)
        const zoneElement = element.closest('[data-zone]');
        if (zoneElement && !zoneElement.hasAttribute('data-end-zone')) {
            const zone = zoneElement.getAttribute('data-zone') as ShoppingStatus;
            if (zone && zone !== draggedItem?.status) {
                dragOverZone = zone;
                return;
            }
        }

        // Don't show zone indicators for same-lane general container areas
        // This prevents the indicator from jumping to end when between items in same lane
    }

    function performDrop(element: Element | null) {
        if (!draggedItem || !element) return;

        // Check if dropping on an item
        const itemElement = element.closest('[data-item-id]');
        if (itemElement && itemElement !== element.closest(`[data-item-id="${draggedItem.id}"]`)) {
            const itemId = itemElement.getAttribute('data-item-id');
            const targetItem = shoppingItems.find(i => i.id === itemId);
            if (targetItem) {
                if (draggedItem.status === targetItem.status) {
                    // Reorder within same column
                    reorderItemsById(draggedItem.id, targetItem.id, draggedItem.status);
                } else {
                    // Move to different column at specific position
                    insertItemAtPosition(draggedItem.id, targetItem.status, targetItem.id);
                }
                return;
            }
        }

        // Check if dropping on an end zone
        const endZoneElement = element.closest('[data-end-zone="true"]');
        if (endZoneElement) {
            const zone = endZoneElement.getAttribute('data-zone') as ShoppingStatus;
            if (zone) {
                if (draggedItem.status !== zone) {
                    // Moving to different column at end
                    moveItem(draggedItem.id, zone);
                } else {
                    // Moving to end of same column
                    moveItemToEnd(draggedItem.id, zone);
                }
                return;
            }
        }

        // Check if dropping on a general zone (fallback for spaces between items)
        const zoneElement = element.closest('[data-zone]');
        if (zoneElement && !zoneElement.hasAttribute('data-end-zone')) {
            const zone = zoneElement.getAttribute('data-zone') as ShoppingStatus;
            if (zone && draggedItem.status !== zone) {
                // Moving to different column (will be placed at end)
                moveItem(draggedItem.id, zone);
            }
        }
    }

    function cleanup() {
        currentPointerId = null;
        isDragging = false;
        draggedItem = null;
        dragOverItem = null;
        dragOverZone = null;

        if (ghostElement) {
            document.body.removeChild(ghostElement);
            ghostElement = null;
        }
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
                <CardContent class="p-1">
                    <div
                        class="rounded-lg border-2 border-dashed border-transparent transition-all duration-200"
                        role="region"
                        aria-label="{config.title} drop zone"
                    >
                    <div
                        class="space-y-2 min-h-[100px] p-1"
                        role="list"
                        aria-label="{config.title} items"
                        data-zone={status}
                    >
                        <!-- Always show add item form for shopping list -->
                        {#if status === 'shopping'}
                            <div class="border flex items-center gap-2 px-2 py-1.5 bg-muted/10 dark:bg-muted/5 rounded-lg border-dashed">
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

                        {#if items.length === 0 && status !== 'shopping'}
                            <div
                                class="flex flex-col items-center justify-center h-32 text-muted-foreground"
                                data-zone={status}
                                in:fade={{ duration: 300, delay: 200 }}
                                out:fade={{ duration: 200 }}
                            >
                                <IconComponent
                                    class="h-8 w-8 mb-2 opacity-50 transition-all duration-300"
                                />
                                <p class="text-sm">No items</p>
                            </div>
                        {/if}

                            {#each items as item (item.id)}
                                <div
                                    class="select-none cursor-grab active:cursor-grabbing transition-all duration-200 {isDragging && draggedItem?.id === item.id ? 'opacity-50 scale-95' : 'opacity-100 scale-100'}"
                                    style="touch-action: none;"
                                    role="listitem"
                                    aria-label="Draggable item: {item.name}"
                                    data-item-id={item.id}
                                    onpointerdown={(e) => handlePointerDown(e, item)}
                                    onpointermove={(e) => handlePointerMove(e, e.currentTarget)}
                                    onpointerup={(e) => handlePointerUp(e, e.currentTarget)}
                                    onpointercancel={(e) => handlePointerCancel(e, e.currentTarget)}
                                    in:fly={{ x: status === 'shopping' ? -200 : 200, duration: 400, easing: quintOut }}
                                    out:fly={{ x: status === 'shopping' ? 200 : -200, duration: 300, easing: quintOut }}
                                    animate:flip={{ duration: 600, easing: quintOut }}
                                >
                                    <!-- Drop indicator -->
                                    {#if dragOverItem?.id === item.id && draggedItem?.id !== item.id}
                                        <div class="h-0.5 bg-primary rounded-full mx-2 mb-2 transition-all duration-200" in:scale={{ duration: 200 }}></div>
                                    {/if}
                                    <div class="border flex items-center justify-between px-2 py-1.5 bg-muted/30 dark:bg-muted hover:bg-muted/50 dark:hover:bg-muted/30 rounded-lg transition-all duration-200">
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

                            <!-- End drop zone for dragging to end of list -->
                            {#if isDragging && items.length > 0}
                                <div
                                    class="h-8"
                                    data-zone={status}
                                    data-end-zone="true"
                                >
                                    {#if dragOverZone === status}
                                        <div class="h-0.5 bg-primary rounded-full mx-2 mt-3"></div>
                                    {/if}
                                </div>
                            {/if}

                    </div>
                    </div>
                </CardContent>

            </Card>
        {/each}
    </div>
</div>