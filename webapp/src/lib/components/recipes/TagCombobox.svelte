<script lang="ts">
    import { Badge } from '$lib/components/ui/badge';
    import PlusIcon from '@lucide/svelte/icons/plus';
    import CheckIcon from "@lucide/svelte/icons/check";
    import { tick } from "svelte";
    import * as Command from "$lib/components/ui/command/index.js";
    import * as Popover from "$lib/components/ui/popover/index.js";
    import { useId } from "bits-ui";
    import VirtualList from 'svelte-tiny-virtual-list';
    import AddTagPlaceholder from './AddTagPlaceholder.svelte';

    let {
        empty = 'No tags found.',
        search = 'Search tags...',
        options = [],
        maxResults = 200,
        onSelect,
        open = $bindable(false)
    } = $props();

    let searchQuery = $state('')

    const itemHeight = 34
    const maxHeight = 300
    let dynamicHeight = $derived.by(() => {
        const contentHeight = (filteredOptions?.length || 0) * itemHeight
        return Math.min(contentHeight, maxHeight)
    });

    let filteredOptions = $derived.by(() => {
        if (!Array.isArray(options) || options.length === 0) {
            return []
        }

        if (!searchQuery) {
            return options
        }

        const query = searchQuery.toLowerCase();
        const filtered = options.filter(option =>
            option && option.label && typeof option.label === 'string' &&
            option.label.toLowerCase().includes(query)
        );

        return filtered.slice(0, maxResults);
    });

    const triggerId = useId();

    function closeAndFocusTrigger() {
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

    const handleSelect = (option) => {
        if (option && option.value !== undefined) {
            onSelect?.(option.value);
            closeAndFocusTrigger();
        }
    }
</script>

<Popover.Root bind:open onOpenChangeComplete={resetSearchQuery}>
    <Popover.Trigger id={triggerId} asChild class="flex cursor-pointer">
        <Badge class="rounded-full" {onclick}>
            <PlusIcon /> Tag
        </Badge>
    </Popover.Trigger>
    <Popover.Content class="w-64 p-0">
        <Command.Root shouldFilter={false}>
            <Command.Input
                autofocus
                placeholder={search}
                class="h-9"
                bind:value={searchQuery}
            />
            {#if !filteredOptions || filteredOptions.length === 0}
                <Command.Empty>
                    {empty}
                </Command.Empty>
            {:else}
                <div class="my-1">
                    <VirtualList
                        height={dynamicHeight}
                        itemCount={filteredOptions.length}
                        itemSize={itemHeight}
                    >
                        {#snippet item({ index, style })}
                            {@const option = filteredOptions?.[index]}
                            {#if option && option.value !== undefined && option.label !== undefined}
                                <div {style}>
                                    <Command.Item
                                        class="mx-1"
                                        value={option.value.toString()}
                                        onSelect={() => handleSelect(option)}
                                    >
                                        {option.label ?? ''}
                                        <CheckIcon class="ml-auto h-4 w-4 text-transparent"/>
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