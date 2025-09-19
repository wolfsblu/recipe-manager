<script lang="ts">
    import { createRawSnippet } from "svelte";
    import * as Dialog from "$lib/components/ui/dialog/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import { Input } from "$lib/components/ui/input/index.js";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import { addUnit, updateUnit, deleteUnit } from "$lib/api/units/units.svelte";
    import { useCrud } from "$lib/hooks/useCrud.svelte";
    import { renderSnippet } from "$lib/components/ui/data-table/index.js";
    import { dialogStore, addUnitDialogOpen } from "$lib/stores/dialog.svelte";
    import type { PageProps } from './$types';

    type Unit = {
        id: number;
        name: string;
        symbol?: string;
    };

    const { data }: PageProps = $props();

    let newUnitName = $state('');
    let newUnitSymbol = $state('');
    let editUnitName = $state('');
    let editUnitSymbol = $state('');

    const crud = useCrud<Unit>({
        add: (unit) => addUnit(unit),
        update: (id, unit) => updateUnit(id, unit),
        delete: (id) => deleteUnit(id)
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
        },
        {
            accessorKey: "symbol",
            header: "Symbol",
            cell: ({ row }: any) => {
                const symbolSnippet = createRawSnippet<[string]>((getSymbol) => {
                    const symbol = getSymbol();
                    return {
                        render: () => `<div class="text-muted-foreground">${symbol || '-'}</div>`
                    };
                });
                return renderSnippet(symbolSnippet, row.getValue("symbol"));
            },
            enableSorting: false
        }
    ];

    const handleAddSubmit = async () => {
        if (!newUnitName.trim()) return;
        
        const unitData: Omit<Unit, 'id'> = { name: newUnitName.trim() };
        if (newUnitSymbol.trim()) {
            unitData.symbol = newUnitSymbol.trim();
        }
        
        const success = await crud.handleAdd(unitData);
        if (success) {
            newUnitName = '';
            newUnitSymbol = '';
        }
    };

    const handleEditSubmit = async () => {
        if (!editUnitName.trim()) return;
        
        const unitData: Omit<Unit, 'id'> = { name: editUnitName.trim() };
        if (editUnitSymbol.trim()) {
            unitData.symbol = editUnitSymbol.trim();
        }
        
        const success = await crud.handleEdit(unitData);
        if (success) {
            editUnitName = '';
            editUnitSymbol = '';
        }
    };

    const openEditDialog = (unit: Unit) => {
        editUnitName = unit.name;
        editUnitSymbol = unit.symbol || '';
        crud.openEditDialog(unit);
    };

    $effect(() => {
        if (crud.editingItem) {
            editUnitName = crud.editingItem?.name || '';
            editUnitSymbol = crud.editingItem?.symbol || '';
        }
    });

    $effect(() => {
        if ($addUnitDialogOpen) {
            crud.openAddDialog();
            addUnitDialogOpen.set(false);
        }
    });
</script>

<DataTable
    data={data.units}
    searchColumn="name"
    searchPlaceholder="Filter units"
    entityName="Unit"
    {columns}
    onAdd={crud.openAddDialog}
    onEdit={openEditDialog}
    onDelete={crud.openDeleteDialog}
    onBulkDelete={crud.openBulkDeleteDialog}
/>

<!-- Add Unit Dialog -->
<Dialog.Root bind:open={crud.showAddDialog}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>Add New Unit</Dialog.Title>
        </Dialog.Header>
        <div class="grid gap-4 py-4">
            <div class="grid grid-cols-4 items-center gap-4">
                <label for="name" class="text-right">Name</label>
                <Input
                    id="name"
                    class="col-span-3"
                    bind:value={newUnitName}
                    placeholder="Enter unit name"
                    onkeydown={(e) => e.key === 'Enter' && handleAddSubmit()}
                />
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
                <label for="symbol" class="text-right">Symbol</label>
                <Input
                    id="symbol"
                    class="col-span-3"
                    bind:value={newUnitSymbol}
                    placeholder="Enter unit symbol (optional)"
                    onkeydown={(e) => e.key === 'Enter' && handleAddSubmit()}
                />
            </div>
        </div>
        <Dialog.Footer>
            <Button variant="outline" onclick={() => crud.showAddDialog = false} disabled={crud.isSubmitting}>
                Cancel
            </Button>
            <Button onclick={handleAddSubmit} disabled={!newUnitName.trim() || crud.isSubmitting}>
                {crud.isSubmitting ? 'Adding...' : 'Add Unit'}
            </Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>

<!-- Edit Unit Dialog -->
<Dialog.Root bind:open={crud.showEditDialog}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>Edit Unit</Dialog.Title>
        </Dialog.Header>
        <div class="grid gap-4 py-4">
            <div class="grid grid-cols-4 items-center gap-4">
                <label for="edit-name" class="text-right">Name</label>
                <Input
                    id="edit-name"
                    class="col-span-3"
                    bind:value={editUnitName}
                    placeholder="Enter unit name"
                    onkeydown={(e) => e.key === 'Enter' && handleEditSubmit()}
                />
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
                <label for="edit-symbol" class="text-right">Symbol</label>
                <Input
                    id="edit-symbol"
                    class="col-span-3"
                    bind:value={editUnitSymbol}
                    placeholder="Enter unit symbol (optional)"
                    onkeydown={(e) => e.key === 'Enter' && handleEditSubmit()}
                />
            </div>
        </div>
        <Dialog.Footer>
            <Button variant="outline" onclick={() => crud.showEditDialog = false} disabled={crud.isSubmitting}>
                Cancel
            </Button>
            <Button onclick={handleEditSubmit} disabled={!editUnitName.trim() || crud.isSubmitting}>
                {crud.isSubmitting ? 'Saving...' : 'Save Changes'}
            </Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>

<!-- Delete Unit Dialog -->
<Dialog.Root bind:open={crud.showDeleteDialog}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>Delete Unit</Dialog.Title>
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