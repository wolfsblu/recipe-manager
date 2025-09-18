<script lang="ts">
    import CheckIcon from "@lucide/svelte/icons/check";
    import MinusIcon from "@lucide/svelte/icons/minus";
    import { cn } from "$lib/utils.js";

    interface Props {
        checked?: boolean;
        indeterminate?: boolean;
        onCheckedChange?: (checked: boolean) => void;
        disabled?: boolean;
        class?: string;
        "aria-label"?: string;
    }

    let {
        checked = false,
        indeterminate = false,
        onCheckedChange,
        disabled = false,
        class: className,
        "aria-label": ariaLabel
    }: Props = $props();

    const handleChange = () => {
        if (disabled) return;
        const newChecked = !checked;
        onCheckedChange?.(newChecked);
    };
</script>

<button
    type="button"
    role="checkbox"
    aria-checked={indeterminate ? "mixed" : checked}
    aria-label={ariaLabel}
    {disabled}
    onclick={handleChange}
    class={cn(
        "peer h-4 w-4 shrink-0 rounded-sm border border-primary ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground align-middle",
        checked && "bg-primary text-primary-foreground",
        className
    )}
    data-state={indeterminate ? "indeterminate" : checked ? "checked" : "unchecked"}
>
    {#if indeterminate}
        <MinusIcon class="h-4 w-4" />
    {:else if checked}
        <CheckIcon class="h-4 w-4" />
    {/if}
</button>