<script lang="ts">
    import { Button } from '$lib/components/ui/button';
    import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
    import * as Dialog from '$lib/components/ui/dialog';
    import { Badge } from '$lib/components/ui/badge';
    import { shoppingStore } from '$lib/stores/shopping.svelte';
    import { deleteShoppingList, updateShoppingList } from '$lib/api/shopping/shopping.svelte';
    import { Input } from "$lib/components/ui/input/index.js";
    import { toast } from 'svelte-sonner';
    import ChevronDownIcon from '@lucide/svelte/icons/chevron-down';
    import PlusIcon from '@lucide/svelte/icons/plus';
    import EditIcon from '@lucide/svelte/icons/edit-2';
    import TrashIcon from '@lucide/svelte/icons/trash-2';
    import CheckIcon from '@lucide/svelte/icons/check';
    import XIcon from '@lucide/svelte/icons/x';

    interface Props {
        onCreateNewList: () => void;
    }

    let { onCreateNewList }: Props = $props();

    let editingListId = $state<number | null>(null);
    let editingName = $state('');
    let dropdownOpen = $state(false);
    let showDeleteDialog = $state(false);
    let deleteListData = $state<{ id: number; name: string } | null>(null);
    let isDeleting = $state(false);

    const handleSelectList = (listId: number) => {
        shoppingStore.setCurrentListId(listId);
    };

    const startEdit = (listId: number, currentName: string) => {
        editingListId = listId;
        editingName = currentName;
        dropdownOpen = false;
    };

    const cancelEdit = () => {
        editingListId = null;
        editingName = '';
    };

    const saveEdit = async () => {
        if (!editingListId || !editingName.trim()) return;

        try {
            const updatedList = await updateShoppingList(editingListId, { name: editingName.trim() });
            shoppingStore.updateList(updatedList);
            toast.success('Shopping list renamed successfully');
            editingListId = null;
            editingName = '';
        } catch (error) {
            toast.error('Failed to rename shopping list');
        }
    };

    const showDeleteConfirmation = (listId: number, listName: string) => {
        dropdownOpen = false;
        deleteListData = { id: listId, name: listName };
        showDeleteDialog = true;
    };

    const confirmDelete = async () => {
        if (!deleteListData) return;

        isDeleting = true;
        try {
            await deleteShoppingList(deleteListData.id);
            shoppingStore.removeList(deleteListData.id);
            toast.success('Shopping list deleted successfully');
            showDeleteDialog = false;
            deleteListData = null;
        } catch (error) {
            toast.error('Failed to delete shopping list');
        } finally {
            isDeleting = false;
        }
    };

    const cancelDelete = () => {
        showDeleteDialog = false;
        deleteListData = null;
    };

    const handleKeyPress = (e: KeyboardEvent) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            saveEdit();
        } else if (e.key === 'Escape') {
            e.preventDefault();
            cancelEdit();
        }
    };

    const currentList = $derived(shoppingStore.currentList);
    const lists = $derived(shoppingStore.lists);
    const totalItems = $derived(currentList ? currentList.items.length : 0);
    const doneItems = $derived(currentList ? currentList.items.filter(item => item.done).length : 0);
</script>

<div class="flex items-center justify-between mb-6">
    <div class="flex items-center gap-4">
        <DropdownMenu.Root bind:open={dropdownOpen}>
            <DropdownMenu.Trigger>
                {#snippet child({ props })}
                    <Button {...props} variant="outline" class="justify-between min-w-[200px]">
                        <span class="truncate">
                            {currentList ? currentList.name : 'Select a list'}
                        </span>
                        <ChevronDownIcon class="h-4 w-4 ml-2" />
                    </Button>
                {/snippet}
            </DropdownMenu.Trigger>
            <DropdownMenu.Content>
                <DropdownMenu.Label>Your Shopping Lists</DropdownMenu.Label>
                <DropdownMenu.Separator />
                
                {#each lists as list (list.id)}
                    <DropdownMenu.Item 
                        class="flex items-center justify-between cursor-pointer"
                        onclick={() => handleSelectList(list.id)}
                    >
                        <div class="flex items-center gap-2">
                            <span class="truncate">{list.name}</span>
                            {#if list.id === currentList?.id}
                                <CheckIcon class="h-4 w-4 text-primary" />
                            {/if}
                        </div>
                        <div class="flex items-center gap-1">
                            <Button
                                size="sm"
                                variant="ghost"
                                onclick={(e) => {
                                    e.stopPropagation();
                                    startEdit(list.id, list.name);
                                }}
                                title="Rename list"
                            >
                                <EditIcon />
                            </Button>
                            <Button
                                size="sm"
                                variant="ghost"
                                onclick={(e) => {
                                    e.stopPropagation();
                                    showDeleteConfirmation(list.id, list.name);
                                }}
                                title="Delete list"
                            >
                                <TrashIcon />
                            </Button>
                        </div>
                    </DropdownMenu.Item>
                {/each}
                
                <DropdownMenu.Separator />
                <DropdownMenu.Item onclick={onCreateNewList} class="cursor-pointer">
                    <PlusIcon />
                    Create New List
                </DropdownMenu.Item>
            </DropdownMenu.Content>
        </DropdownMenu.Root>
    </div>

    <Button onclick={onCreateNewList}>
        <PlusIcon />
        New List
    </Button>
</div>

<!-- Edit Modal Overlay -->
{#if editingListId}
    <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" onclick={cancelEdit}>
        <div class="bg-background border rounded-lg p-4 max-w-sm w-full mx-4" onclick={(e) => e.stopPropagation()}>
            <h3 class="font-semibold mb-3">Rename Shopping List</h3>
            <div class="flex gap-2">
                <Input
                    type="text"
                    bind:value={editingName}
                    onkeypress={handleKeyPress}
                    placeholder="List name"
                />
                <Button onclick={saveEdit} disabled={!editingName.trim()}>
                    <CheckIcon />
                </Button>
                <Button variant="outline" onclick={cancelEdit}>
                    <XIcon />
                </Button>
            </div>
        </div>
    </div>
{/if}

<!-- Delete Confirmation Dialog -->
<Dialog.Root open={showDeleteDialog} onOpenChange={(open) => { if (!open) cancelDelete(); }}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>Delete Shopping List</Dialog.Title>
            <Dialog.Description>
                Are you sure you want to delete "{deleteListData?.name}"? This action cannot be undone.
            </Dialog.Description>
        </Dialog.Header>
        
        <Dialog.Footer>
            <Button
                type="button"
                variant="outline"
                onclick={cancelDelete}
                disabled={isDeleting}
            >
                Cancel
            </Button>
            <Button 
                variant="destructive" 
                onclick={confirmDelete}
                disabled={isDeleting}
            >
                {isDeleting ? 'Deleting...' : 'Delete'}
            </Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>