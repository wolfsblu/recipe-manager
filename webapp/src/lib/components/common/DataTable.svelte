<script lang="ts" generics="T extends Record<string, any>">
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
    import DataTableCheckbox from "./DataTableCheckbox.svelte";
    import DataTableActions from "./DataTableActions.svelte";
    import * as Table from "$lib/components/ui/table/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import { Input } from "$lib/components/ui/input/index.js";
    import {
        FlexRender,
        createSvelteTable,
        renderComponent,
        renderSnippet
    } from "$lib/components/ui/data-table/index.js";

    interface Props {
        data: T[];
        searchColumn: string;
        searchPlaceholder: string;
        entityName: string;
        onAdd: () => void;
        onEdit: (item: T) => void;
        onDelete: (item: T) => void;
        onBulkDelete: (items: T[]) => void;
        columns: ColumnDef<T>[];
        // Server-side pagination props
        serverPagination?: {
            canGoNext: boolean;
            canGoPrevious: boolean;
            onNext: () => void;
            onPrevious: () => void;
            onFirst: () => void;
            loading: boolean;
            currentPage: number;
            totalLoaded: number;
        };
    }

    let {
        data,
        searchColumn,
        searchPlaceholder,
        entityName,
        onAdd,
        onEdit,
        onDelete,
        onBulkDelete,
        columns: userColumns,
        serverPagination
    }: Props = $props();

    let paginationState = $state<PaginationState>({ pageIndex: 0, pageSize: 10 });
    let columnFilters = $state<ColumnFiltersState>([]);
    let rowSelection = $state<RowSelectionState>({});
    let columnVisibility = $state<VisibilityState>({});

    // Use large page size for server pagination to show all items on current page
    $effect(() => {
        if (serverPagination) {
            paginationState = { pageIndex: 0, pageSize: 1000 };
        }
    });

    const columns: ColumnDef<T>[] = [
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
            enableHiding: false,
            size: 40,
            minSize: 40,
        },
        ...userColumns,
        {
            id: "actions",
            enableHiding: false,
            cell: ({ row }) =>
                renderComponent(DataTableActions, { 
                    item: row.original,
                    onEdit,
                    onDelete
                }),
        }
    ];

    const table = createSvelteTable({
        get data() {
            return data;
        },
        columns,
        state: {
            get pagination() {
                return paginationState;
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
        columnResizeMode: 'onChange',
        onPaginationChange: (updater) => {
            if (typeof updater === "function") {
                paginationState = updater(paginationState);
            } else {
                paginationState = updater;
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

    const getSelectedCount = () => {
        return table.getFilteredSelectedRowModel().rows.length;
    };

    const handleBulkDelete = () => {
        const selectedItems = table.getFilteredSelectedRowModel().rows.map(row => row.original);
        onBulkDelete(selectedItems);
        table.resetRowSelection();
    };
</script>

<div class="flex flex-col gap-5 h-full p-5">
    <div class="w-full">
        <div class="flex items-center justify-between gap-2 pb-4">
            <Input
                placeholder={searchPlaceholder}
                value={(table.getColumn(searchColumn)?.getFilterValue() as string) ?? ""}
                oninput={(e) =>
                    table.getColumn(searchColumn)?.setFilterValue(e.currentTarget.value)}
                onchange={(e) => {
                    table.getColumn(searchColumn)?.setFilterValue(e.currentTarget.value);
                }}
                class="max-w-sm"
            />
            <div class="flex items-center gap-2">
                {#if getSelectedCount() > 0}
                    <Button variant="destructive" onclick={handleBulkDelete}>
                        <TrashIcon />
                        <span class="hidden sm:inline">Delete {getSelectedCount()}</span>
                    </Button>
                {/if}
                <Button onclick={onAdd}>
                    <PlusIcon />
                    <span class="hidden sm:inline">Add {entityName}</span>
                </Button>
            </div>
        </div>
        <div class="flex items-center justify-between text-sm text-muted-foreground pb-4">
            <div>
                {table.getFilteredSelectedRowModel().rows.length} of
                {table.getFilteredRowModel().rows.length} row(s) selected
            </div>
            <div>
                {#if serverPagination}
                    Page {serverPagination.currentPage} · {serverPagination.totalLoaded} loaded
                    {#if serverPagination.loading}
                        <span class="ml-1">(loading...)</span>
                    {/if}
                {:else}
                    Page {table.getState().pagination.pageIndex + 1} of {table.getPageCount()}
                {/if}
            </div>
        </div>
        <div class="rounded-md border">
            <Table.Root style="table-layout: fixed; width: 100%;">
                <Table.Header>
                    {#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
                        <Table.Row>
                            {#each headerGroup.headers as header (header.id)}
                                <Table.Head 
                                    class="[&:has([role=checkbox])]:pl-3"
                                    style="width: {header.getSize()}px;"
                                >
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
                                No results found.
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
                    onclick={() => serverPagination ? serverPagination.onFirst() : table.setPageIndex(0)}
                    disabled={serverPagination ? !serverPagination.canGoPrevious || serverPagination.loading : !table.getCanPreviousPage()}
                >
                    <ChevronsLeftIcon class="h-4 w-4" />
                </Button>
                <Button
                    variant="outline"
                    size="sm"
                    onclick={() => serverPagination ? serverPagination.onPrevious() : table.previousPage()}
                    disabled={serverPagination ? !serverPagination.canGoPrevious || serverPagination.loading : !table.getCanPreviousPage()}
                >
                    <ChevronLeftIcon class="h-4 w-4" />
                </Button>
                <Button
                    variant="outline"
                    size="sm"
                    onclick={() => serverPagination ? serverPagination.onNext() : table.nextPage()}
                    disabled={serverPagination ? !serverPagination.canGoNext || serverPagination.loading : !table.getCanNextPage()}
                >
                    <ChevronRightIcon class="h-4 w-4" />
                </Button>
                {#if !serverPagination}
                    <Button
                        variant="outline"
                        size="sm"
                        onclick={() => table.setPageIndex(table.getPageCount() - 1)}
                        disabled={!table.getCanNextPage()}
                    >
                        <ChevronsRightIcon class="h-4 w-4" />
                    </Button>
                {/if}
            </div>
        </div>
    </div>
</div>