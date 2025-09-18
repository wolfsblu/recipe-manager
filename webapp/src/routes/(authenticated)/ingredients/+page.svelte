<script lang="ts">
    import ChevronDownIcon from "@lucide/svelte/icons/chevron-down";
    import ChevronLeftIcon from "@lucide/svelte/icons/chevron-left";
    import ChevronRightIcon from "@lucide/svelte/icons/chevron-right";
    import ChevronsLeftIcon from "@lucide/svelte/icons/chevrons-left";
    import ChevronsRightIcon from "@lucide/svelte/icons/chevrons-right";
    import PlusIcon from "@lucide/svelte/icons/plus";
    import TrashIcon from "@lucide/svelte/icons/trash-2";
    import {
        type ColumnDef,
        type ColumnFiltersState,
        type PaginationState,
        type RowSelectionState,
        type VisibilityState,
        getCoreRowModel,
        getFilteredRowModel,
        getPaginationRowModel
    } from "@tanstack/table-core";
    import { createRawSnippet } from "svelte";
    import DataTableCheckbox from "$lib/components/ingredients/data-table-checkbox.svelte";
    import DataTableActions from "$lib/components/ingredients/data-table-actions.svelte";
    import * as Table from "$lib/components/ui/table/index.js";
    import * as Dialog from "$lib/components/ui/dialog/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
    import { Input } from "$lib/components/ui/input/index.js";
    import {
        FlexRender,
        createSvelteTable,
        renderComponent,
        renderSnippet
    } from "$lib/components/ui/data-table/index.js";
    import { addIngredient, updateIngredient, deleteIngredient } from "$lib/api/ingredients/ingredients.svelte";
    import type { PageProps } from './$types';
    import { invalidateAll } from '$app/navigation';

    type Ingredient = {
        id: number;
        name: string;
    };

    const { data }: PageProps = $props();

    let showAddDialog = $state(false);
    let showEditDialog = $state(false);
    let showDeleteDialog = $state(false);
    let editingIngredient = $state<Ingredient | null>(null);
    let deletingIngredient = $state<Ingredient | null>(null);
    let newIngredientName = $state('');
    let editIngredientName = $state('');
    let isSubmitting = $state(false);
    let showBulkDeleteDialog = $state(false);

    const handleAddIngredient = async () => {
        if (!newIngredientName.trim() || isSubmitting) return;
        
        isSubmitting = true;
        try {
            await addIngredient({ name: newIngredientName.trim() });
            newIngredientName = '';
            showAddDialog = false;
            await invalidateAll();
        } catch (error) {
            console.error('Failed to add ingredient:', error);
        } finally {
            isSubmitting = false;
        }
    };

    const handleEditIngredient = async () => {
        if (!editingIngredient || !editIngredientName.trim() || isSubmitting) return;
        
        isSubmitting = true;
        try {
            await updateIngredient(editingIngredient.id, { name: editIngredientName.trim() });
            showEditDialog = false;
            editingIngredient = null;
            editIngredientName = '';
            await invalidateAll();
        } catch (error) {
            console.error('Failed to update ingredient:', error);
        } finally {
            isSubmitting = false;
        }
    };

    const handleDeleteIngredient = async () => {
        if (!deletingIngredient || isSubmitting) return;
        
        isSubmitting = true;
        try {
            await deleteIngredient(deletingIngredient.id);
            showDeleteDialog = false;
            deletingIngredient = null;
            await invalidateAll();
        } catch (error) {
            console.error('Failed to delete ingredient:', error);
        } finally {
            isSubmitting = false;
        }
    };

    const openEditDialog = (ingredient: Ingredient) => {
        editingIngredient = ingredient;
        editIngredientName = ingredient.name;
        showEditDialog = true;
    };

    const openDeleteDialog = (ingredient: Ingredient) => {
        deletingIngredient = ingredient;
        showDeleteDialog = true;
    };

    const handleBulkDelete = async () => {
        const selectedRows = table.getFilteredSelectedRowModel().rows;
        if (selectedRows.length === 0 || isSubmitting) return;
        
        isSubmitting = true;
        try {
            for (const row of selectedRows) {
                await deleteIngredient(row.original.id);
            }
            showBulkDeleteDialog = false;
            table.resetRowSelection();
            await invalidateAll();
        } catch (error) {
            console.error('Failed to delete ingredients:', error);
        } finally {
            isSubmitting = false;
        }
    };

    const getSelectedCount = () => {
        return table.getFilteredSelectedRowModel().rows.length;
    };

    const columns: ColumnDef<Ingredient>[] = [
        {
            id: "select",
            header: ({ table }) =>
                renderComponent(DataTableCheckbox, {
                    checked: table.getIsAllRowsSelected(),
                    indeterminate:
                        table.getIsSomeRowsSelected() &&
                        !table.getIsAllRowsSelected(),
                    onCheckedChange: (value) => table.toggleAllRowsSelected(!!value),
                    "aria-label": "Select all"
                }),
            cell: ({ row }) =>
                renderComponent(DataTableCheckbox, {
                    checked: row.getIsSelected(),
                    onCheckedChange: (value) => row.toggleSelected(!!value),
                    "aria-label": "Select row"
                }),
            enableSorting: false,
            enableHiding: false
        },
        {
            accessorKey: "name",
            header: "Name",
            cell: ({ row }) => {
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
            id: "actions",
            header: () => {
                const actionsHeaderSnippet = createRawSnippet(() => {
                    return {
                        render: () => `<div class="text-right">Actions</div>`
                    };
                });
                return renderSnippet(actionsHeaderSnippet, "");
            },
            enableHiding: false,
            cell: ({ row }) =>
                renderComponent(DataTableActions, { 
                    ingredient: row.original,
                    onEdit: openEditDialog,
                    onDelete: openDeleteDialog
                })
        }
    ];

    let pagination = $state<PaginationState>({ pageIndex: 0, pageSize: 10 });
    let columnFilters = $state<ColumnFiltersState>([]);
    let rowSelection = $state<RowSelectionState>({});
    let columnVisibility = $state<VisibilityState>({});

    const table = createSvelteTable({
        get data() {
            return data.ingredients;
        },
        columns,
        state: {
            get pagination() {
                return pagination;
            },
            get columnVisibility() {
                return columnVisibility;
            },
            get rowSelection() {
                return rowSelection;
            },
            get columnFilters() {
                return columnFilters;
            }
        },
        getCoreRowModel: getCoreRowModel(),
        getPaginationRowModel: getPaginationRowModel(),
        getFilteredRowModel: getFilteredRowModel(),
        onPaginationChange: (updater) => {
            if (typeof updater === "function") {
                pagination = updater(pagination);
            } else {
                pagination = updater;
            }
        },
        onColumnFiltersChange: (updater) => {
            if (typeof updater === "function") {
                columnFilters = updater(columnFilters);
            } else {
                columnFilters = updater;
            }
        },
        onColumnVisibilityChange: (updater) => {
            if (typeof updater === "function") {
                columnVisibility = updater(columnVisibility);
            } else {
                columnVisibility = updater;
            }
        },
        onRowSelectionChange: (updater) => {
            if (typeof updater === "function") {
                rowSelection = updater(rowSelection);
            } else {
                rowSelection = updater;
            }
        }
    });
</script>

<div class="flex flex-col gap-5 h-full p-5">
    <div class="w-full">
        <div class="flex items-center justify-between gap-2 pb-4">
            <Input
                placeholder="Filter ingredients..."
                value={(table.getColumn("name")?.getFilterValue() as string) ?? ""}
                oninput={(e) =>
                    table.getColumn("name")?.setFilterValue(e.currentTarget.value)}
                onchange={(e) => {
                    table.getColumn("name")?.setFilterValue(e.currentTarget.value);
                }}
                class="max-w-sm"
            />
            <div class="flex items-center gap-2">
                {#if getSelectedCount() > 0}
                    <Button variant="destructive" onclick={() => showBulkDeleteDialog = true}>
                        <TrashIcon />
                        <span class="hidden sm:inline">Delete {getSelectedCount()}</span>
                    </Button>
                {/if}
                <Button onclick={() => showAddDialog = true}>
                    <PlusIcon />
                    <span class="hidden sm:inline">Add Ingredient</span>
                </Button>
            </div>
        </div>
        <div class="flex items-center justify-between text-sm text-muted-foreground pb-4">
            <div>
                {table.getFilteredSelectedRowModel().rows.length} of
                {table.getFilteredRowModel().rows.length} row(s) selected.
            </div>
            <div>
                Page {table.getState().pagination.pageIndex + 1} of {table.getPageCount()}
            </div>
        </div>
        <div class="rounded-md border">
            <Table.Root>
                <Table.Header>
                    {#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
                        <Table.Row>
                            {#each headerGroup.headers as header (header.id)}
                                <Table.Head class="[&:has([role=checkbox])]:pl-3">
                                    {#if !header.isPlaceholder}
                                        <FlexRender
                                            content={header.column.columnDef.header}
                                            context={header.getContext()}
                                        />
                                    {/if}
                                </Table.Head>
                            {/each}
                        </Table.Row>
                    {/each}
                </Table.Header>
                <Table.Body>
                    {#each table.getRowModel().rows as row (row.id)}
                        <Table.Row data-state={row.getIsSelected() && "selected"}>
                            {#each row.getVisibleCells() as cell (cell.id)}
                                <Table.Cell class="[&:has([role=checkbox])]:pl-3">
                                    <FlexRender
                                        content={cell.column.columnDef.cell}
                                        context={cell.getContext()}
                                    />
                                </Table.Cell>
                            {/each}
                        </Table.Row>
                    {:else}
                        <Table.Row>
                            <Table.Cell colspan={columns.length} class="h-24 text-center">
                                No ingredients found.
                            </Table.Cell>
                        </Table.Row>
                    {/each}
                </Table.Body>
            </Table.Root>
        </div>
        <div class="flex items-center justify-center pt-4">
            <div class="flex items-center space-x-2">
                <Button
                    variant="outline"
                    size="sm"
                    onclick={() => table.setPageIndex(0)}
                    disabled={!table.getCanPreviousPage()}
                >
                    <ChevronsLeftIcon class="h-4 w-4" />
                </Button>
                <Button
                    variant="outline"
                    size="sm"
                    onclick={() => table.previousPage()}
                    disabled={!table.getCanPreviousPage()}
                >
                    <ChevronLeftIcon class="h-4 w-4" />
                </Button>
                <Button
                    variant="outline"
                    size="sm"
                    onclick={() => table.nextPage()}
                    disabled={!table.getCanNextPage()}
                >
                    <ChevronRightIcon class="h-4 w-4" />
                </Button>
                <Button
                    variant="outline"
                    size="sm"
                    onclick={() => table.setPageIndex(table.getPageCount() - 1)}
                    disabled={!table.getCanNextPage()}
                >
                    <ChevronsRightIcon class="h-4 w-4" />
                </Button>
            </div>
        </div>
    </div>
</div>

<!-- Add Ingredient Dialog -->
<Dialog.Root bind:open={showAddDialog}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>Add Ingredient</Dialog.Title>
        </Dialog.Header>
        <div class="grid gap-4 py-4">
            <div class="grid grid-cols-4 items-center gap-4">
                <label for="name" class="text-right">Name</label>
                <Input
                    id="name"
                    class="col-span-3"
                    bind:value={newIngredientName}
                    placeholder="Enter ingredient name"
                    onkeydown={(e) => e.key === 'Enter' && handleAddIngredient()}
                />
            </div>
        </div>
        <Dialog.Footer>
            <Button variant="outline" onclick={() => showAddDialog = false} disabled={isSubmitting}>
                Cancel
            </Button>
            <Button onclick={handleAddIngredient} disabled={!newIngredientName.trim() || isSubmitting}>
                {isSubmitting ? 'Adding...' : 'Add Ingredient'}
            </Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>

<!-- Edit Ingredient Dialog -->
<Dialog.Root bind:open={showEditDialog}>
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
                    onkeydown={(e) => e.key === 'Enter' && handleEditIngredient()}
                />
            </div>
        </div>
        <Dialog.Footer>
            <Button variant="outline" onclick={() => showEditDialog = false} disabled={isSubmitting}>
                Cancel
            </Button>
            <Button onclick={handleEditIngredient} disabled={!editIngredientName.trim() || isSubmitting}>
                {isSubmitting ? 'Saving...' : 'Save Changes'}
            </Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>

<!-- Delete Ingredient Dialog -->
<Dialog.Root bind:open={showDeleteDialog}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>Delete Ingredient</Dialog.Title>
        </Dialog.Header>
        <div class="py-4">
            <p>Are you sure you want to delete "{deletingIngredient?.name}"? This action cannot be undone.</p>
        </div>
        <Dialog.Footer>
            <Button variant="outline" onclick={() => showDeleteDialog = false} disabled={isSubmitting}>
                Cancel
            </Button>
            <Button variant="destructive" onclick={handleDeleteIngredient} disabled={isSubmitting}>
                {isSubmitting ? 'Deleting...' : 'Delete'}
            </Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>

<!-- Bulk Delete Confirmation Dialog -->
<Dialog.Root bind:open={showBulkDeleteDialog}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>Delete Ingredients</Dialog.Title>
        </Dialog.Header>
        <div class="py-4">
            <p>Are you sure you want to delete {getSelectedCount()} selected ingredient{getSelectedCount() > 1 ? 's' : ''}? This action cannot be undone.</p>
        </div>
        <Dialog.Footer>
            <Button variant="outline" onclick={() => showBulkDeleteDialog = false} disabled={isSubmitting}>
                Cancel
            </Button>
            <Button variant="destructive" onclick={handleBulkDelete} disabled={isSubmitting}>
                {isSubmitting ? 'Deleting...' : `Delete ${getSelectedCount()} Ingredient${getSelectedCount() > 1 ? 's' : ''}`}
            </Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>