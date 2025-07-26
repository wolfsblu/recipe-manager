<script lang="ts">
    import CheckIcon from "@lucide/svelte/icons/check";
    import ChevronsUpDownIcon from "@lucide/svelte/icons/chevrons-up-down";
    import { tick } from "svelte";
    import * as Command from "$lib/components/ui/command/index.js";
    import * as Popover from "$lib/components/ui/popover/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import { cn } from "$lib/utils.js";

    const frameworks = [
        {
            value: "greensalad",
            label: "Green Salad"
        },
        {
            value: "carrots",
            label: "Carrots"
        },
        {
            value: "chicken",
            label: "Chicken"
        },
        {
            value: "salt",
            label: "Salt"
        },
        {
            value: "pepper",
            label: "Pepper"
        }
    ];

    let open = $state(false);
    let value = $state("");
    let triggerRef = $state<HTMLButtonElement>(null!);

    let {
        class: className
    } = $props()

    const selectedValue = $derived(
        frameworks.find((f) => f.value === value)?.label
    );

    // We want to refocus the trigger button when the user selects
    // an item from the list so users can continue navigating the
    // rest of the form with the keyboard.
    function closeAndFocusTrigger() {
        open = false;
        tick().then(() => {
            triggerRef.focus();
        });
    }
</script>

<Popover.Root bind:open>
    <Popover.Trigger bind:ref={triggerRef}>
        {#snippet child({ props })}
            <Button
                    {...props}
                    variant="outline"
                    class={cn("w-[200px] justify-between", className)}
                    role="combobox"
                    aria-expanded={open}
            >
                {selectedValue || "Select an ingredient..."}
                <ChevronsUpDownIcon class="opacity-50" />
            </Button>
        {/snippet}
    </Popover.Trigger>
    <Popover.Content class="w-[var(--bits-popover-anchor-width)] min-w-[var(--bits-popover-anchor-width)] p-0">
        <Command.Root>
            <Command.Input placeholder="Search framework..." />
            <Command.List>
                <Command.Empty>No framework found.</Command.Empty>
                <Command.Group value="frameworks">
                    {#each frameworks as framework (framework.value)}
                        <Command.Item
                                value={framework.value}
                                onSelect={() => {
        value = framework.value;
        closeAndFocusTrigger();
       }}
                        >
                            <CheckIcon
                                    class={cn(value !== framework.value && "text-transparent")}
                            />
                            {framework.label}
                        </Command.Item>
                    {/each}
                </Command.Group>
            </Command.List>
        </Command.Root>
    </Popover.Content>
</Popover.Root>