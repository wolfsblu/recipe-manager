<script lang="ts">
    import CheckIcon from "@lucide/svelte/icons/check";
    import ChevronsUpDownIcon from "@lucide/svelte/icons/chevrons-up-down";
    import { tick } from "svelte";
    import * as Command from "$lib/components/ui/command/index.js";
    import * as Form from "$lib/components/ui/form/index.js";
    import * as Popover from "$lib/components/ui/popover/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import { buttonVariants } from "$lib/components/ui/button/index.js";
    import { cn } from "$lib/utils.js";
    import { useId } from "bits-ui";

    let {
        empty = 'No results found.',
        form,
        name,
        placeholder = 'Select option',
        search = 'Search options...',
        value = $bindable(),
        options,
    } = $props();

    let open = $state(false);

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
                <input hidden value={value} name={props.name} />
            {/snippet}
        </Form.Control>
        <Popover.Content class="w-[var(--bits-popover-anchor-width)] min-w-[var(--bits-popover-anchor-width)] p-0">
            <Command.Root>
                <Command.Input
                        autofocus
                        placeholder={search}
                        class="h-9"
                />
                <Command.Empty>
                    {empty}
                </Command.Empty>
                <Command.Group value="items">
                    {#each options as option (option.value)}
                        <Command.Item
                                value={option.value}
                                onSelect={() => {
                                    value = option.value;
                                    closeAndFocusTrigger(triggerId);
                                }}
                        >
                            {option.label}
                            <CheckIcon class={cn("ml-auto", option.value !== value && "text-transparent")}/>
                        </Command.Item>
                    {/each}
                </Command.Group>
            </Command.Root>
        </Popover.Content>
    </Popover.Root>
    <Form.FieldErrors />
</Form.Field>