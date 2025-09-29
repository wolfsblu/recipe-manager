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
    import * as m from "$lib/paraglide/messages.js";

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
            header: m.units_name(),
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
            header: m.units_symbol(),
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
    searchPlaceholder={m.units_searchPlaceholder()}
    entityName={m.units_entityName()}
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
            <Dialog.Title>{m.units_addDialog_title()}</Dialog.Title>
        </Dialog.Header>
        <div class="grid gap-4 py-4">
            <div class="grid grid-cols-4 items-center gap-4">
                <label for="name" class="text-right">{m.units_addDialog_nameLabel()}</label>
                <Input
                    id="name"
                    class="col-span-3"
                    bind:value={newUnitName}
                    placeholder={m.units_addDialog_namePlaceholder()}
                    onkeydown={(e) => e.key === 'Enter' && handleAddSubmit()}
                />
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
                <label for="symbol" class="text-right">{m.units_addDialog_symbolLabel()}</label>
                <Input
                    id="symbol"
                    class="col-span-3"
                    bind:value={newUnitSymbol}
                    placeholder={m.units_addDialog_symbolPlaceholder()}
                    onkeydown={(e) => e.key === 'Enter' && handleAddSubmit()}
                />
            </div>
        </div>
        <Dialog.Footer>
            <Button variant="outline" onclick={() => crud.showAddDialog = false} disabled={crud.isSubmitting}>
                {m.units_addDialog_cancel()}
            </Button>
            <Button onclick={handleAddSubmit} disabled={!newUnitName.trim() || crud.isSubmitting}>
                {crud.isSubmitting ? m.units_addDialog_adding() : m.units_addDialog_button()}
            </Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>

<!-- Edit Unit Dialog -->
<Dialog.Root bind:open={crud.showEditDialog}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>{m.units_editDialog_title()}</Dialog.Title>
        </Dialog.Header>
        <div class="grid gap-4 py-4">
            <div class="grid grid-cols-4 items-center gap-4">
                <label for="edit-name" class="text-right">{m.units_editDialog_nameLabel()}</label>
                <Input
                    id="edit-name"
                    class="col-span-3"
                    bind:value={editUnitName}
                    placeholder={m.units_editDialog_namePlaceholder()}
                    onkeydown={(e) => e.key === 'Enter' && handleEditSubmit()}
                />
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
                <label for="edit-symbol" class="text-right">{m.units_editDialog_symbolLabel()}</label>
                <Input
                    id="edit-symbol"
                    class="col-span-3"
                    bind:value={editUnitSymbol}
                    placeholder={m.units_editDialog_symbolPlaceholder()}
                    onkeydown={(e) => e.key === 'Enter' && handleEditSubmit()}
                />
            </div>
        </div>
        <Dialog.Footer>
            <Button variant="outline" onclick={() => crud.showEditDialog = false} disabled={crud.isSubmitting}>
                {m.units_editDialog_cancel()}
            </Button>
            <Button onclick={handleEditSubmit} disabled={!editUnitName.trim() || crud.isSubmitting}>
                {crud.isSubmitting ? m.units_editDialog_saving() : m.units_editDialog_button()}
            </Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>

<!-- Delete Unit Dialog -->
<Dialog.Root bind:open={crud.showDeleteDialog}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>{m.units_deleteDialog_title()}</Dialog.Title>
        </Dialog.Header>
        <div class="py-4">
            <p>{m.units_deleteDialog_description({ name: crud.deletingItem?.name || '' })}</p>
        </div>
        <Dialog.Footer>
            <Button variant="outline" onclick={() => crud.showDeleteDialog = false} disabled={crud.isSubmitting}>
                {m.units_deleteDialog_cancel()}
            </Button>
            <Button variant="destructive" onclick={crud.handleDelete} disabled={crud.isSubmitting}>
                {crud.isSubmitting ? m.units_deleteDialog_deleting() : m.units_deleteDialog_button()}
            </Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>