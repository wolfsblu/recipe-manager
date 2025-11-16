<script lang="ts">
    import CheckIcon from "@lucide/svelte/icons/check";
    import ChevronsUpDownIcon from "@lucide/svelte/icons/chevrons-up-down";
    import { tick, untrack } from "svelte";
    import * as Command from "$lib/components/ui/command/index.js";
    import * as Form from "$lib/components/ui/form/index.js";
    import * as Popover from "$lib/components/ui/popover/index.js";
    import { buttonVariants } from "$lib/components/ui/button/index.js";
    import { cn } from "$lib/utils.js";
    import { useId } from "bits-ui";
    import VirtualList from 'svelte-tiny-virtual-list';
    import { useDebounce } from 'runed';

    type Option = { value: number; label: string };
    type FetchResult = {
        options: Option[];
        nextCursor?: string | null;
        hasMore: boolean;
    };

    let {
        empty = 'No results found.',
        form,
        name,
        placeholder = 'Select option',
        search = 'Search options...',
        value = $bindable(),
        fetchOptions,
        minSearchLength = 2,
        debounceMs = 300,
    }: {
        empty?: string;
        form: any;
        name: string;
        placeholder?: string;
        search?: string;
        value?: number;
        fetchOptions: (searchQuery: string, cursor?: string | null) => Promise<FetchResult>;
        minSearchLength?: number;
        debounceMs?: number;
    } = $props();

    let open = $state(false);
    let searchQuery = $state('');
    let loading = $state(false);
    let options = $state<Option[]>([]);
    let selectedOption = $state<Option | undefined>(undefined);
    let nextCursor = $state<string | null | undefined>(undefined);
    let hasMore = $state(false);

    const itemHeight = 34;
    const maxHeight = 300;
    let dynamicHeight = $derived.by(() => {
        const contentHeight = (options?.length || 0) * itemHeight;
        return Math.min(contentHeight, maxHeight);
    });

    // Search function that only updates options when results arrive
    async function performSearch(query: string) {
        loading = true;
        try {
            const result = await fetchOptions(query);
            // Only update if dropdown is still open
            if (open) {
                options = result.options;
                nextCursor = result.nextCursor;
                hasMore = result.hasMore;
            }
        } catch (error) {
            console.error('Failed to fetch options:', error);
            if (open) {
                options = [];
                nextCursor = undefined;
                hasMore = false;
            }
        } finally {
            loading = false;
        }
    }

    // Load more results for infinite scroll
    async function loadMore() {
        if (!hasMore || !nextCursor) return;

        try {
            const result = await fetchOptions(searchQuery, nextCursor);
            if (open) {
                options = [...options, ...result.options];
                nextCursor = result.nextCursor;
                hasMore = result.hasMore;
            }
        } catch (error) {
            console.error('Failed to load more options:', error);
        }
    }

    // Debounced search for user typing
    const debouncedSearch = useDebounce(performSearch, debounceMs);

    // Handle search input changes
    function handleSearchInput(query: string) {
        if (query.length >= minSearchLength || query.length === 0) {
            loading = true;
            if (query.length === 0) {
                // Immediate search for empty query
                performSearch('');
            } else {
                // Debounced for user typing
                debouncedSearch(query);
            }
        }
    }

    // Handle dropdown open/close
    function handleOpenChange(isOpen: boolean) {
        open = isOpen;

        if (isOpen) {
            // Opening - load fresh initial results
            searchQuery = '';
            // Only show loading if we truly have no options (first time)
            if (options.length === 0) {
                loading = true;
            }
            // Fetch fresh results - they'll replace old ones when ready
            performSearch('');
        } else {
            // Closing - just reset loading state
            loading = false;
        }
    }

    // Load selected option on mount if editing
    $effect(() => {
        if (value && !selectedOption) {
            loadSelectedOption(value);
        }
    });

    async function loadSelectedOption(val: number) {
        // Try to find in current options first
        const found = options.find(opt => opt.value === val);
        if (found) {
            selectedOption = found;
            return;
        }

        // If not found, we need to fetch all to find it
        // Or you could add a getById API endpoint
        try {
            const result = await fetchOptions('');
            const option = result.options.find(opt => opt.value === val);
            if (option) {
                selectedOption = option;
            }
        } catch (error) {
            console.error('Failed to load selected option:', error);
        }
    }

    // We want to refocus the trigger button when the user selects
    // an item from the list so users can continue navigating the
    // rest of the form with the keyboard.
    function closeAndFocusTrigger(triggerId: string) {
        open = false;
        tick().then(() => {
            const trigger = document.getElementById(triggerId);
            if (trigger) {
                trigger.focus();
            }
        });
    }

    const triggerId = useId();
</script>

<Form.Field {form} {name} class="space-y-0 w-full">
    <Popover.Root bind:open onOpenChange={(isOpen) => handleOpenChange(isOpen ?? false)}>
        <Form.Control id={triggerId}>
            {#snippet children({ props })}
                <Popover.Trigger
                    class={cn(
                        buttonVariants({ variant: "outline" }),
                        "w-full justify-between",
                        !value && "text-muted-foreground"
                    )}
                    role="combobox"
                    {...props}
                >
                    {selectedOption?.label ?? placeholder}
                    <ChevronsUpDownIcon class="opacity-50" />
                </Popover.Trigger>
                <input hidden bind:value={value} name={props.name} />
            {/snippet}
        </Form.Control>
        <Popover.Content class="w-[var(--bits-popover-anchor-width)] min-w-[var(--bits-popover-anchor-width)] p-0">
            <Command.Root shouldFilter={false}>
                <Command.Input
                    autofocus
                    placeholder={search}
                    class="h-9"
                    bind:value={searchQuery}
                    oninput={() => handleSearchInput(searchQuery)}
                />
                {#if !options || options.length === 0}
                    {#if loading}
                        <div class="py-6 text-center text-sm text-muted-foreground">
                            Loading...
                        </div>
                    {:else}
                        <Command.Empty>
                            {searchQuery.length > 0 && searchQuery.length < minSearchLength
                                ? `Type at least ${minSearchLength} characters to search`
                                : empty}
                        </Command.Empty>
                    {/if}
                {:else}
                    <div class="my-1">
                        <VirtualList
                            height={dynamicHeight}
                            itemCount={options.length}
                            itemSize={itemHeight}
                            onAfterScroll={(e) => {
                                console.log(e)
                                const target = e.event.target;
                                if (target == null) return
                                const scrollTop = target.scrollTop;
                                const scrollHeight = target.scrollHeight;
                                const clientHeight = target.clientHeight;

                                // Load more when scrolled to within 100px of bottom
                                if (scrollHeight - scrollTop - clientHeight < 100) {
                                    loadMore();
                                }
                            }}
                        >
                            {#snippet item({ index, style })}
                                {@const option = options?.[index]}
                                {#if option && option.value !== undefined && option.label !== undefined}
                                    <div {style}>
                                        <Command.Item
                                            class="mx-1"
                                            value={option.value.toString()}
                                            onSelect={() => {
                                                if (option && option.value !== undefined) {
                                                    value = option.value;
                                                    selectedOption = option;
                                                    closeAndFocusTrigger(triggerId);
                                                }
                                            }}
                                        >
                                            {option.label ?? ''}
                                            <CheckIcon class={cn("ml-auto h-4 w-4", option.value !== value && "text-transparent")}/>
                                        </Command.Item>
                                    </div>
                                {/if}
                            {/snippet}
                        </VirtualList>
                    </div>
                {/if}
            </Command.Root>
        </Popover.Content>
    </Popover.Root>
    <Form.FieldErrors />
</Form.Field>
