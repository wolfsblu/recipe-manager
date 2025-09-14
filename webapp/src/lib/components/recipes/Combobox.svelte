<script lang="ts">
    import CheckIcon from "@lucide/svelte/icons/check";
    import ChevronsUpDownIcon from "@lucide/svelte/icons/chevrons-up-down";
    import { tick } from "svelte";
    import * as Command from "$lib/components/ui/command/index.js";
    import * as Form from "$lib/components/ui/form/index.js";
    import * as Popover from "$lib/components/ui/popover/index.js";
    import { buttonVariants } from "$lib/components/ui/button/index.js";
    import { cn } from "$lib/utils.js";
    import { useId } from "bits-ui";
    import VirtualList from 'svelte-tiny-virtual-list';

    let {
        empty = 'No results found.',
        form,
        name,
        placeholder = 'Select option',
        search = 'Search options...',
        value = $bindable(),
        options,
        maxResults = 200,
    } = $props();

    let open = $state(false);
    let searchQuery = $state('');
    
    // Dynamic height calculation
    const itemHeight = 32;
    const maxHeight = 300;
    let dynamicHeight = $derived.by(() => {
        const contentHeight = (displayOptions?.length || 0) * itemHeight;
        return Math.min(contentHeight, maxHeight);
    });

    // Filter options based on search query
    let filteredOptions = $derived.by(() => {
        if (!Array.isArray(options) || options.length === 0) {
            return [];
        }

        if (!searchQuery) {
            return options; // Show all options when no search query
        }

        const query = searchQuery.toLowerCase();
        const filtered = options.filter(option =>
            option && option.label && typeof option.label === 'string' &&
            option.label.toLowerCase().includes(query)
        );

        return filtered.slice(0, maxResults);
    });

    // Show selected option even if not in filtered results
    let displayOptions = $derived.by(() => {
        if (!Array.isArray(options) || options.length === 0) {
            return [];
        }

        const selectedOption = options.find(option =>
            option && option.value !== undefined && option.value === value
        );
        const filtered = filteredOptions || [];

        if (selectedOption && !filtered.some(option =>
            option && option.value !== undefined && option.value === value
        )) {
            return [selectedOption, ...filtered];
        }

        return filtered;
    });

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

    const resetSearchQuery = () => {
        searchQuery = ''
    }

    const triggerId = useId();
</script>

<Form.Field {form} {name} class="space-y-0 w-full">
    <Popover.Root bind:open onOpenChangeComplete={resetSearchQuery}>
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
                    {Array.isArray(options) ? options.find((f) => f && f.value === value)?.label ?? placeholder : placeholder}
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
                />
                {#if !displayOptions || displayOptions.length === 0}
                    <Command.Empty>
                        {empty}
                    </Command.Empty>
                {:else}
                    <div class="w-full" style="height: {dynamicHeight}px;">
                        <VirtualList
                            width="100%"
                            height={dynamicHeight}
                            itemCount={displayOptions.length}
                            itemSize={itemHeight}
                        >
                            {#snippet item({ index, style })}
                                {@const option = displayOptions?.[index]}
                                {#if option && option.value !== undefined && option.label !== undefined}
                                    <div {style}>
                                        <Command.Item
                                                value={option.value.toString()}
                                                onSelect={() => {
                                                    if (option && option.value !== undefined) {
                                                        value = option.value;
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