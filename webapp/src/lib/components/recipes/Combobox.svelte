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
    import VirtualList from '@sveltejs/svelte-virtual-list';

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
        const contentHeight = displayOptions.length * itemHeight;
        return Math.min(contentHeight, maxHeight);
    });

    // Filter options based on search query
    let filteredOptions = $derived.by(() => {
        if (!searchQuery) {
            return options; // Show all options when no search query
        }

        const query = searchQuery.toLowerCase();
        const filtered = options.filter(option =>
            option.label.toLowerCase().includes(query)
        );

        return filtered.slice(0, maxResults);
    });

    // Show selected option even if not in filtered results
    let displayOptions = $derived.by(() => {
        const selectedOption = options.find(option => option.value === value);
        const filtered = filteredOptions;

        if (selectedOption && !filtered.some(option => option.value === value)) {
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
            document.getElementById(triggerId)?.focus();
        });
    }
    const triggerId = useId();
</script>

<Form.Field {form} {name} class="space-y-0 w-full">
    <Popover.Root bind:open>
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
                    {options.find((f) => f.value === value)?.label ?? placeholder}
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
                {#if displayOptions.length === 0}
                    <Command.Empty>
                        {empty}
                    </Command.Empty>
                {:else}
                    <div class="w-full" style="height: {dynamicHeight}px;">
                        <VirtualList items={displayOptions} height="{dynamicHeight}px" itemHeight={itemHeight} let:item>
                            <Command.Item
                                    value={item.value.toString()}
                                    onSelect={() => {
                                        value = item.value;
                                        closeAndFocusTrigger(triggerId);
                                    }}
                            >
                                {item.label}
                                <CheckIcon class={cn("ml-auto h-4 w-4", item.value !== value && "text-transparent")}/>
                            </Command.Item>
                        </VirtualList>
                    </div>
                {/if}
            </Command.Root>
        </Popover.Content>
    </Popover.Root>
    <Form.FieldErrors />
</Form.Field>