<script lang="ts">
    import { createRawSnippet } from "svelte";
    import { onMount } from "svelte";
    import * as Dialog from "$lib/components/ui/dialog/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import { Input } from "$lib/components/ui/input/index.js";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import { getIngredients, addIngredient, updateIngredient, deleteIngredient } from "$lib/api/ingredients/ingredients.svelte";
    import { useCrud } from "$lib/hooks/useCrud.svelte";
    import { useServerPagination } from "$lib/utils/serverPagination.svelte";
    import { renderSnippet } from "$lib/components/ui/data-table/index.js";
    import { dialogStore, addIngredientDialogOpen } from "$lib/stores/dialog.svelte";
    import type { PageProps } from './$types';
    import * as m from "$lib/paraglide/messages.js";

    type Ingredient = {
        id: number;
        name: string;
    };

    const { data }: PageProps = $props();

    const pagination = useServerPagination<Ingredient>({
        fetchPage: async (cursor) => {
            const response = await getIngredients({ cursor, limit: 30 });
            return {
                data: response?.data ?? [],
                nextCursor: response?.nextCursor,
                hasMore: response?.hasMore ?? false
            };
        }
    });

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
        if ($addIngredientDialogOpen) {
            crud.openAddDialog();
            addIngredientDialogOpen.set(false);
        }
    });

    onMount(() => {
        pagination.loadInitialPage();
    });
</script>

<DataTable
    data={pagination.currentPageData}
    searchColumn="name"
    searchPlaceholder={m.ingredients_searchPlaceholder()}
    entityName={m.ingredients_entityName()}
    {columns}
    onAdd={crud.openAddDialog}
    onEdit={openEditDialog}
    onDelete={crud.openDeleteDialog}
    onBulkDelete={crud.openBulkDeleteDialog}
    serverPagination={{
        canGoNext: pagination.canGoNext,
        canGoPrevious: pagination.canGoPrevious,
        onNext: pagination.nextPage,
        onPrevious: pagination.previousPage,
        onFirst: pagination.goToFirstPage,
        loading: pagination.loading,
        currentPage: pagination.currentPageNumber,
        totalLoaded: pagination.totalLoadedItems
    }}
/>

<!-- Add Ingredient Dialog -->
<Dialog.Root bind:open={crud.showAddDialog}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>{m.ingredients_addDialog_title()}</Dialog.Title>
        </Dialog.Header>
        <div class="grid gap-4 py-4">
            <div class="grid grid-cols-4 items-center gap-4">
                <label for="name" class="text-right">{m.ingredients_addDialog_nameLabel()}</label>
                <Input
                    id="name"
                    class="col-span-3"
                    bind:value={newIngredientName}
                    placeholder={m.ingredients_addDialog_namePlaceholder()}
                    onkeydown={(e) => e.key === 'Enter' && handleAddSubmit()}
                />
            </div>
        </div>
        <Dialog.Footer>
            <Button variant="outline" onclick={() => crud.showAddDialog = false} disabled={crud.isSubmitting}>
                {m.ingredients_addDialog_cancel()}
            </Button>
            <Button onclick={handleAddSubmit} disabled={!newIngredientName.trim() || crud.isSubmitting}>
                {crud.isSubmitting ? m.ingredients_addDialog_adding() : m.ingredients_addDialog_button()}
            </Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>

<!-- Edit Ingredient Dialog -->
<Dialog.Root bind:open={crud.showEditDialog}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>{m.ingredients_editDialog_title()}</Dialog.Title>
        </Dialog.Header>
        <div class="grid gap-4 py-4">
            <div class="grid grid-cols-4 items-center gap-4">
                <label for="edit-name" class="text-right">{m.ingredients_editDialog_nameLabel()}</label>
                <Input
                    id="edit-name"
                    class="col-span-3"
                    bind:value={editIngredientName}
                    placeholder={m.ingredients_editDialog_namePlaceholder()}
                    onkeydown={(e) => e.key === 'Enter' && handleEditSubmit()}
                />
            </div>
        </div>
        <Dialog.Footer>
            <Button variant="outline" onclick={() => crud.showEditDialog = false} disabled={crud.isSubmitting}>
                {m.ingredients_editDialog_cancel()}
            </Button>
            <Button onclick={handleEditSubmit} disabled={!editIngredientName.trim() || crud.isSubmitting}>
                {crud.isSubmitting ? m.ingredients_editDialog_saving() : m.ingredients_editDialog_button()}
            </Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>

<!-- Delete Ingredient Dialog -->
<Dialog.Root bind:open={crud.showDeleteDialog}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>{m.ingredients_deleteDialog_title()}</Dialog.Title>
        </Dialog.Header>
        <div class="py-4">
            <p>{m.ingredients_deleteDialog_description({ name: crud.deletingItem?.name || '' })}</p>
        </div>
        <Dialog.Footer>
            <Button variant="outline" onclick={() => crud.showDeleteDialog = false} disabled={crud.isSubmitting}>
                {m.ingredients_deleteDialog_cancel()}
            </Button>
            <Button variant="destructive" onclick={crud.handleDelete} disabled={crud.isSubmitting}>
                {crud.isSubmitting ? m.ingredients_deleteDialog_deleting() : m.ingredients_deleteDialog_button()}
            </Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>