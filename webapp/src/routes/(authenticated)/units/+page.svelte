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
        code?: string;
    };

    const { data }: PageProps = $props();

    let newUnitName = $state('');
    let newUnitCode = $state('');
    let editUnitName = $state('');
    let editUnitCode = $state('');

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
            accessorKey: "code",
            header: "Code",
            cell: ({ row }: any) => {
                const codeSnippet = createRawSnippet<[string]>((getCode) => {
                    const code = getCode();
                    return {
                        render: () => `<div class="text-muted-foreground">${code || '-'}</div>`
                    };
                });
                return renderSnippet(codeSnippet, row.getValue("code"));
            },
            enableSorting: false
        }
    ];

    const handleAddSubmit = async () => {
        if (!newUnitName.trim()) return;
        
        const unitData: Omit<Unit, 'id'> = { name: newUnitName.trim() };
        if (newUnitCode.trim()) {
            unitData.code = newUnitCode.trim();
        }
        
        const success = await crud.handleAdd(unitData);
        if (success) {
            newUnitName = '';
            newUnitCode = '';
        }
    };

    const handleEditSubmit = async () => {
        if (!editUnitName.trim()) return;
        
        const unitData: Omit<Unit, 'id'> = { name: editUnitName.trim() };
        if (editUnitCode.trim()) {
            unitData.code = editUnitCode.trim();
        }
        
        const success = await crud.handleEdit(unitData);
        if (success) {
            editUnitName = '';
            editUnitCode = '';
        }
    };

    const openEditDialog = (unit: Unit) => {
        editUnitName = unit.name;
        editUnitCode = unit.code || '';
        crud.openEditDialog(unit);
    };

    $effect(() => {
        if (crud.editingItem) {
            editUnitName = crud.editingItem?.name || '';
            editUnitCode = crud.editingItem?.code || '';
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
                <label for="code" class="text-right">Code</label>
                <Input
                    id="code"
                    class="col-span-3"
                    bind:value={newUnitCode}
                    placeholder="Enter unit code (optional)"
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
                <label for="edit-code" class="text-right">Code</label>
                <Input
                    id="edit-code"
                    class="col-span-3"
                    bind:value={editUnitCode}
                    placeholder="Enter unit code (optional)"
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