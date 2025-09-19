<script lang="ts">
    import { createRawSnippet } from "svelte";
    import * as Dialog from "$lib/components/ui/dialog/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import { Input } from "$lib/components/ui/input/index.js";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import { addIngredient, updateIngredient, deleteIngredient } from "$lib/api/ingredients/ingredients.svelte";
    import { useCrud } from "$lib/hooks/useCrud.svelte";
    import { renderSnippet } from "$lib/components/ui/data-table/index.js";
    import { dialogStore } from "$lib/stores/dialog.svelte";
    import type { PageProps } from './$types';

    type Ingredient = {
        id: number;
        name: string;
    };

    const { data }: PageProps = $props();

    let newIngredientName = $state('');
    let editIngredientName = $state('');

    const crud = useCrud<Ingredient>({
        add: (ingredient) => addIngredient(ingredient),
        update: (id, ingredient) => updateIngredient(id, ingredient),
        delete: (id) => deleteIngredient(id)
    });

    const columns = [
        {
            accessorKey: "name",
            header: "Name",
            cell: ({ row }: any) => {
                const nameSnippet = createRawSnippet<[string]>((getName) => {
                    const name = getName();
                    return {
                        render: () => `<div>${name}</div>`
                    };
                });
                return renderSnippet(nameSnippet, row.getValue("name"));
            },
            enableSorting: false
        }
    ];

    const handleAddSubmit = async () => {
        if (!newIngredientName.trim()) return;
        
        const success = await crud.handleAdd({ name: newIngredientName.trim() });
        if (success) {
            newIngredientName = '';
        }
    };

    const handleEditSubmit = async () => {
        if (!editIngredientName.trim()) return;
        
        const success = await crud.handleEdit({ name: editIngredientName.trim() });
        if (success) {
            editIngredientName = '';
        }
    };

    const openEditDialog = (ingredient: Ingredient) => {
        editIngredientName = ingredient.name;
        crud.openEditDialog(ingredient);
    };

    $effect(() => {
        if (crud.editingItem) {
            editIngredientName = crud.editingItem?.name || '';
        }
    });

    $effect(() => {
        if (dialogStore.addIngredientDialogOpen) {
            crud.openAddDialog();
            dialogStore.addIngredientDialogOpen = false;
        }
    });
</script>

<DataTable
    data={data.ingredients}
    searchColumn="name"
    searchPlaceholder="Filter ingredients"
    entityName="Ingredient"
    {columns}
    onAdd={crud.openAddDialog}
    onEdit={openEditDialog}
    onDelete={crud.openDeleteDialog}
    onBulkDelete={crud.openBulkDeleteDialog}
/>

<!-- Add Ingredient Dialog -->
<Dialog.Root bind:open={crud.showAddDialog}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>Add New Ingredient</Dialog.Title>
        </Dialog.Header>
        <div class="grid gap-4 py-4">
            <div class="grid grid-cols-4 items-center gap-4">
                <label for="name" class="text-right">Name</label>
                <Input
                    id="name"
                    class="col-span-3"
                    bind:value={newIngredientName}
                    placeholder="Enter ingredient name"
                    onkeydown={(e) => e.key === 'Enter' && handleAddSubmit()}
                />
            </div>
        </div>
        <Dialog.Footer>
            <Button variant="outline" onclick={() => crud.showAddDialog = false} disabled={crud.isSubmitting}>
                Cancel
            </Button>
            <Button onclick={handleAddSubmit} disabled={!newIngredientName.trim() || crud.isSubmitting}>
                {crud.isSubmitting ? 'Adding...' : 'Add Ingredient'}
            </Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>

<!-- Edit Ingredient Dialog -->
<Dialog.Root bind:open={crud.showEditDialog}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>Edit Ingredient</Dialog.Title>
        </Dialog.Header>
        <div class="grid gap-4 py-4">
            <div class="grid grid-cols-4 items-center gap-4">
                <label for="edit-name" class="text-right">Name</label>
                <Input
                    id="edit-name"
                    class="col-span-3"
                    bind:value={editIngredientName}
                    placeholder="Enter ingredient name"
                    onkeydown={(e) => e.key === 'Enter' && handleEditSubmit()}
                />
            </div>
        </div>
        <Dialog.Footer>
            <Button variant="outline" onclick={() => crud.showEditDialog = false} disabled={crud.isSubmitting}>
                Cancel
            </Button>
            <Button onclick={handleEditSubmit} disabled={!editIngredientName.trim() || crud.isSubmitting}>
                {crud.isSubmitting ? 'Saving...' : 'Save Changes'}
            </Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>

<!-- Delete Ingredient Dialog -->
<Dialog.Root bind:open={crud.showDeleteDialog}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>Delete Ingredient</Dialog.Title>
        </Dialog.Header>
        <div class="py-4">
            <p>Are you sure you want to delete "{crud.deletingItem?.name}"? This action cannot be undone.</p>
        </div>
        <Dialog.Footer>
            <Button variant="outline" onclick={() => crud.showDeleteDialog = false} disabled={crud.isSubmitting}>
                Cancel
            </Button>
            <Button variant="destructive" onclick={crud.handleDelete} disabled={crud.isSubmitting}>
                {crud.isSubmitting ? 'Deleting...' : 'Delete'}
            </Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>