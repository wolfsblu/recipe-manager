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
        form,
        name,
        value = $bindable(),
    } = $props();

    const ingredients = [
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

<Form.Field {form} {name} class="w-full">
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
                    {ingredients.find((f) => f.value === value)?.label ??
                    "Select ingredient"}
                    <ChevronsUpDownIcon class="opacity-50" />
                </Popover.Trigger>
                <input hidden value={value} name={props.name} />
            {/snippet}
        </Form.Control>
        <Popover.Content class="w-[var(--bits-popover-anchor-width)] min-w-[var(--bits-popover-anchor-width)] p-0">
            <Command.Root>
                <Command.Input
                        autofocus
                        placeholder="Search ingredient..."
                        class="h-9"
                />
                <Command.Empty>No ingredient found.</Command.Empty>
                <Command.Group value="ingredients">
                    {#each ingredients as ingredient (ingredient.value)}
                        <Command.Item
                                value={ingredient.value}
                                onSelect={() => {
                                    value = ingredient.value;
                                    closeAndFocusTrigger(triggerId);
                                }}
                        >
                            {ingredient.label}
                            <CheckIcon class={cn("ml-auto", ingredient.value !== value && "text-transparent")}/>
                        </Command.Item>
                    {/each}
                </Command.Group>
            </Command.Root>
        </Popover.Content>
    </Popover.Root>
    <Form.FieldErrors />
</Form.Field>