import { invalidateAll } from '$app/navigation';
import { toast } from 'svelte-sonner';

interface CrudOperations<T> {
    add: (item: Omit<T, 'id'>) => Promise<any>;
    update: (id: number, item: Omit<T, 'id'>) => Promise<any>;
    delete: (id: number) => Promise<any>;
}

interface CrudItem {
    id: number;
    [key: string]: any;
}

export function useCrud<T extends CrudItem>(operations: CrudOperations<T>) {
    let showAddDialog = $state(false);
    let showEditDialog = $state(false);
    let showDeleteDialog = $state(false);
    let showBulkDeleteDialog = $state(false);
    let editingItem = $state<T | null>(null);
    let deletingItem = $state<T | null>(null);
    let isSubmitting = $state(false);

    const handleAdd = async (item: Omit<T, 'id'>) => {
        if (isSubmitting) return;
        
        isSubmitting = true;
        try {
            await operations.add(item);
            showAddDialog = false;
            await invalidateAll();
            return true;
        } catch (error) {
            toast.error('Failed to add item. Please try again.');
            return false;
        } finally {
            isSubmitting = false;
        }
    };

    const handleEdit = async (item: Omit<T, 'id'>) => {
        if (!editingItem || isSubmitting) return;
        
        isSubmitting = true;
        try {
            await operations.update(editingItem.id, item);
            showEditDialog = false;
            editingItem = null;
            await invalidateAll();
            return true;
        } catch (error) {
            toast.error('Failed to update item. Please try again.');
            return false;
        } finally {
            isSubmitting = false;
        }
    };

    const handleDelete = async () => {
        if (!deletingItem || isSubmitting) return;
        
        isSubmitting = true;
        try {
            await operations.delete(deletingItem.id);
            showDeleteDialog = false;
            deletingItem = null;
            await invalidateAll();
            return true;
        } catch (error) {
            toast.error('Failed to delete item. Please try again.');
            return false;
        } finally {
            isSubmitting = false;
        }
    };

    const handleBulkDelete = async (items: T[]) => {
        if (items.length === 0 || isSubmitting) return;
        
        isSubmitting = true;
        try {
            const deletePromises = items.map(item => operations.delete(item.id));
            const results = await Promise.allSettled(deletePromises);
            
            const failures = results.filter(result => result.status === 'rejected');
            if (failures.length > 0) {
                if (failures.length === items.length) {
                    toast.error('Failed to delete all selected items. Please try again.');
                } else {
                    toast.error(`Failed to delete ${failures.length} out of ${items.length} items. Please check and try again.`);
                }
            } else {
                toast.success(`Successfully deleted ${items.length} item${items.length > 1 ? 's' : ''}.`);
            }
            
            showBulkDeleteDialog = false;
            await invalidateAll();
            return failures.length === 0; // Return true only if all deletions succeeded
        } catch (error) {
            toast.error('Failed to delete items. Please try again.');
            return false;
        } finally {
            isSubmitting = false;
        }
    };

    const openAddDialog = () => {
        showAddDialog = true;
    };

    const openEditDialog = (item: T) => {
        editingItem = item;
        showEditDialog = true;
    };

    const openDeleteDialog = (item: T) => {
        deletingItem = item;
        showDeleteDialog = true;
    };

    const openBulkDeleteDialog = (items: T[]) => {
        showBulkDeleteDialog = true;
        return handleBulkDelete(items);
    };

    const closeDialogs = () => {
        showAddDialog = false;
        showEditDialog = false;
        showDeleteDialog = false;
        showBulkDeleteDialog = false;
        editingItem = null;
        deletingItem = null;
    };

    return {
        // State (getters)
        get showAddDialog() { return showAddDialog; },
        get showEditDialog() { return showEditDialog; },
        get showDeleteDialog() { return showDeleteDialog; },
        get showBulkDeleteDialog() { return showBulkDeleteDialog; },
        get editingItem() { return editingItem; },
        get deletingItem() { return deletingItem; },
        get isSubmitting() { return isSubmitting; },

        // State (setters)
        set showAddDialog(value: boolean) { showAddDialog = value; },
        set showEditDialog(value: boolean) { showEditDialog = value; },
        set showDeleteDialog(value: boolean) { showDeleteDialog = value; },
        set showBulkDeleteDialog(value: boolean) { showBulkDeleteDialog = value; },

        // Actions
        handleAdd,
        handleEdit,
        handleDelete,
        handleBulkDelete,
        openAddDialog,
        openEditDialog,
        openDeleteDialog,
        openBulkDeleteDialog,
        closeDialogs
    };
}