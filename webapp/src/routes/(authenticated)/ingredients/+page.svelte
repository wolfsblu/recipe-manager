<script lang="ts">
    import { onMount } from "svelte";
    import * as Dialog from "$lib/components/ui/dialog/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import { Input } from "$lib/components/ui/input/index.js";
    import * as Table from "$lib/components/ui/table/index.js";
    import { Checkbox } from "$lib/components/ui/checkbox/index.js";
    import { getIngredients, addIngredient, updateIngredient, deleteIngredient } from "$lib/api/ingredients/ingredients.svelte";
    import { useCrud } from "$lib/hooks/useCrud.svelte";
    import { useServerPagination } from "$lib/utils/serverPagination.svelte";
    import { dialogStore, addIngredientDialogOpen } from "$lib/stores/dialog.svelte";
    import {
        createSvelteTable,
        getCoreRowModel,
        type ColumnDef,
        type SortingState,
        type RowSelectionState
    } from "tanstack-table-8-svelte-5";
    import type { PageProps } from './$types';
    import * as m from "$lib/paraglide/messages.js";
    import { ArrowUpDown, Plus, Trash2, Pencil } from "@lucide/svelte";

    type Ingredient = {
        id: number;
        name: string;
    };

    const { data }: PageProps = $props();

    // Sorting and filtering state
    let sortBy = $state<"name" | "id">("name");
    let sortOrder = $state<"asc" | "desc">("asc");
    let searchQuery = $state("");
    let debouncedSearch = $state("");

    // Debounce search input with cleanup
    $effect(() => {
        const timeout = setTimeout(() => {
            debouncedSearch = searchQuery;
        }, 300);

        return () => clearTimeout(timeout);
    });

    // Reset pagination when sort/filter changes (track previous values)
    let initialized = $state(false);
    $effect(() => {
        // Only reset after initial load
        if (initialized) {
            pagination.reset();
        }
        // Track dependencies
        debouncedSearch;
        sortBy;
        sortOrder;
    });

    $effect(() => {
        // Mark as initialized after first load
        if (pagination.currentPageData.length > 0) {
            initialized = true;
        }
    });

    const pagination = useServerPagination<Ingredient>({
        fetchPage: async (cursor) => {
            const response = await getIngredients({
                cursor,
                limit: 30,
                sortBy,
                sortOrder,
                search: debouncedSearch || undefined
            });
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

    // Table state
    let rowSelection = $state<RowSelectionState>({});
    let sorting = $state<SortingState>([{ id: "name", desc: false }]);

    // Column definitions - simplified since we render manually in template
    const columns: ColumnDef<Ingredient>[] = [
        {
            id: "select",
            enableSorting: false
        },
        {
            accessorKey: "name",
            header: "Name",
            enableSorting: true
        },
        {
            id: "actions",
            enableSorting: false
        }
    ];

    const table = createSvelteTable({
        get data() { return pagination.currentPageData; },
        columns,
        getCoreRowModel: getCoreRowModel(),
        manualPagination: true,
        manualSorting: true,
        manualFiltering: true,
        state: {
            get sorting() { return sorting; },
            get rowSelection() { return rowSelection; }
        },
        onSortingChange: (updater) => {
            sorting = typeof updater === "function" ? updater(sorting) : updater;

            // Update server-side sorting based on table state
            if (sorting.length > 0) {
                const sort = sorting[0];
                sortBy = sort.id === "name" ? "name" : "name";
                sortOrder = sort.desc ? "desc" : "asc";
            }
        },
        onRowSelectionChange: (updater) => {
            rowSelection = typeof updater === "function" ? updater(rowSelection) : updater;
        },
        enableRowSelection: true
    });

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

    const handleSort = () => {
        sortOrder = sortOrder === "asc" ? "desc" : "asc";
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

<div class="space-y-4">
    <!-- Toolbar -->
    <div class="flex items-center justify-between gap-4">
        <Input
            bind:value={searchQuery}
            placeholder={m.ingredients_searchPlaceholder()}
            class="max-w-sm"
        />
        <Button onclick={crud.openAddDialog}>
            <Plus class="mr-2 h-4 w-4" />
            Add {m.ingredients_entityName()}
        </Button>
    </div>

    <!-- Table -->
    <div class="rounded-md border">
        <Table.Root>
            <Table.Header>
                {#each table.getHeaderGroups() as headerGroup}
                    <Table.Row>
                        {#each headerGroup.headers as header}
                            <Table.Head>
                                {#if header.id === "name"}
                                    <Button variant="ghost" onclick={handleSort}>
                                        {header.column.columnDef.header}
                                        <ArrowUpDown class="ml-2 h-4 w-4" />
                                    </Button>
                                {:else if header.id === "select"}
                                    <Checkbox
                                        checked={table.getIsAllPageRowsSelected()}
                                        onCheckedChange={(checked) => table.toggleAllPageRowsSelected(!!checked)}
                                    />
                                {:else}
                                    {header.column.columnDef.header}
                                {/if}
                            </Table.Head>
                        {/each}
                    </Table.Row>
                {/each}
            </Table.Header>
            <Table.Body>
                {#each table.getRowModel().rows as row}
                    <Table.Row>
                        {#each row.getVisibleCells() as cell}
                            <Table.Cell>
                                {#if cell.column.id === "select"}
                                    <Checkbox
                                        checked={row.getIsSelected()}
                                        onCheckedChange={(checked) => row.toggleSelected(!!checked)}
                                    />
                                {:else if cell.column.id === "name"}
                                    {row.original.name}
                                {:else if cell.column.id === "actions"}
                                    <div class="flex items-center gap-2">
                                        <Button
                                            variant="ghost"
                                            size="sm"
                                            onclick={() => openEditDialog(row.original)}
                                        >
                                            <Pencil class="h-4 w-4" />
                                        </Button>
                                        <Button
                                            variant="ghost"
                                            size="sm"
                                            onclick={() => crud.openDeleteDialog(row.original)}
                                        >
                                            <Trash2 class="h-4 w-4" />
                                        </Button>
                                    </div>
                                {/if}
                            </Table.Cell>
                        {/each}
                    </Table.Row>
                {/each}
            </Table.Body>
        </Table.Root>
    </div>

    <!-- Pagination -->
    <div class="flex items-center justify-between">
        <div class="text-sm text-muted-foreground">
            Page {pagination.currentPageNumber} · {pagination.totalLoadedItems} loaded
        </div>
        <div class="flex items-center gap-2">
            <Button
                variant="outline"
                size="sm"
                onclick={pagination.goToFirstPage}
                disabled={!pagination.canGoPrevious || pagination.loading}
            >
                First
            </Button>
            <Button
                variant="outline"
                size="sm"
                onclick={pagination.previousPage}
                disabled={!pagination.canGoPrevious || pagination.loading}
            >
                Previous
            </Button>
            <Button
                variant="outline"
                size="sm"
                onclick={pagination.nextPage}
                disabled={!pagination.canGoNext || pagination.loading}
            >
                {pagination.loading ? "Loading..." : "Next"}
            </Button>
        </div>
    </div>

    <!-- Selected rows info -->
    {#if Object.keys(rowSelection).length > 0}
        <div class="flex items-center gap-2">
            <span class="text-sm text-muted-foreground">
                {Object.keys(rowSelection).length} selected
            </span>
            <Button
                variant="destructive"
                size="sm"
                onclick={() => {
                    const selectedIngredients = table.getSelectedRowModel().rows.map(r => r.original);
                    crud.openBulkDeleteDialog(selectedIngredients);
                }}
            >
                <Trash2 class="mr-2 h-4 w-4" />
                Delete selected
            </Button>
        </div>
    {/if}
</div>

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
