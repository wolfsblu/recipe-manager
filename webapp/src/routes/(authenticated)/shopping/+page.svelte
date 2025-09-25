<script lang="ts">
    import {Button} from '$lib/components/ui/button';
    import {Card, CardContent, CardHeader, CardTitle} from '$lib/components/ui/card';
    import {Input} from '$lib/components/ui/input';
    import {Badge} from '$lib/components/ui/badge';
    import ShoppingCartIcon from '@lucide/svelte/icons/shopping-cart';
    import CheckIcon from '@lucide/svelte/icons/check';
    import PlusIcon from '@lucide/svelte/icons/plus';
    import TrashIcon from '@lucide/svelte/icons/trash-2';
    import {flip} from 'svelte/animate';
    import {fly, scale} from 'svelte/transition';
    import {elasticOut, quintOut} from 'svelte/easing';

    type ShoppingStatus = 'shopping' | 'pantry';

    interface ShoppingItem {
        id: string;
        name: string;
        quantity: string;
        status: ShoppingStatus;
    }

    let newItemName = $state('');
    let newItemQuantity = $state('');
    let newItemUnit = $state('');
    let shoppingItems = $state<ShoppingItem[]>([
        {id: '1', name: 'Milk', quantity: '1 gallon', status: 'shopping'},
        {id: '2', name: 'Bread', quantity: '1 loaf', status: 'shopping'},
        {id: '3', name: 'Apples', quantity: '2 lbs', status: 'shopping'},
        {id: '4', name: 'Chicken breast', quantity: '1 lb', status: 'shopping'},
        {id: '5', name: 'Rice', quantity: '2 cups', status: 'pantry'},
        {id: '6', name: 'Yogurt', quantity: '1 container', status: 'pantry'}
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
            status: 'shopping',
        };

        shoppingItems = [newItem, ...shoppingItems];
        newItemName = '';
        newItemQuantity = '';
        newItemUnit = '';
    }

    function moveItem(id: string, newStatus: ShoppingStatus) {
        shoppingItems = shoppingItems.map(item =>
            item.id === id ? {...item, status: newStatus} : item
        );
    }

    function removeItem(id: string) {
        shoppingItems = shoppingItems.filter(item => item.id !== id);
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

    const todoItems = $derived(shoppingItems.filter(item => item.status === 'shopping'));
    const doneItems = $derived(shoppingItems.filter(item => item.status === 'pantry'));

</script>


<svelte:head>
    <title>Shopping List</title>
</svelte:head>

<div class="container mx-auto p-5">
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {#each ['shopping', 'pantry'] as status (status)}
            {@const config = getColumnConfig(status)}
            {@const items = status === 'shopping' ? todoItems : doneItems}
            {@const IconComponent = config.icon}

            <Card class="h-fit gap-0 py-0">
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
                <CardContent class="p-1">
                    <div
                            class="rounded-lg border-2 border-dashed border-transparent transition-all duration-200"
                            role="region"
                            aria-label="{config.title}"
                    >
                        <div
                                class="space-y-2 min-h-[100px] p-1"
                                role="list"
                                aria-label="{config.title} items"
                        >
                            <!-- Always show add item form for shopping list -->
                            {#if status === 'shopping'}
                                <div class="border flex items-center gap-2 px-2 py-1.5 bg-muted/10 dark:bg-muted/5 rounded-lg border-dashed">
                                    <form onsubmit={(e) => { e.preventDefault(); addItem(); }}
                                          class="flex items-center gap-2 w-full">
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
                                                <PlusIcon/>
                                                <span class="md:hidden">Add Item</span>
                                            </Button>
                                        </div>
                                    </form>
                                </div>
                            {/if}

                            {#if items.length === 0}
                                <div class="flex flex-col items-center justify-center h-32 text-muted-foreground">
                                    <IconComponent/>
                                    <p class="text-sm">No items</p>
                                </div>
                            {/if}
                            {#each items as item (item.id)}
                                <div
                                        class="transition-all duration-200"
                                        role="listitem"
                                        aria-label="Shopping item: {item.name}"
                                        in:fly={{ x: status === 'shopping' ? -200 : 200, duration: 400, easing: quintOut }}
                                        out:fly={{ x: status === 'shopping' ? 200 : -200, duration: 300, easing: quintOut }}
                                        animate:flip={{ duration: 600, easing: quintOut }}
                                >
                                    <div class="border flex items-center justify-between px-2 py-1.5 bg-muted/30 dark:bg-muted rounded-lg transition-all duration-200">
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
                                                    <CheckIcon/>
                                                </Button>
                                            {:else if status === 'pantry'}
                                                <Button
                                                        size="sm"
                                                        variant="ghost"
                                                        onclick={() => moveItem(item.id, 'shopping')}
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
</div>