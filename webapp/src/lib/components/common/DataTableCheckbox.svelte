<script lang="ts">
    import { Checkbox } from "$lib/components/ui/checkbox/index.js";

    interface Props {
        checked?: boolean;
        indeterminate?: boolean;
        onCheckedChange?: (checked: boolean) => void;
        "aria-label"?: string;
    }

    let {
        checked = false,
        indeterminate = false,
        onCheckedChange,
        "aria-label": ariaLabel
    }: Props = $props();
</script>

<Checkbox
    {checked}
    {indeterminate}
    onCheckedChange={(value) => onCheckedChange?.(!!value)}
    aria-label={ariaLabel}
/>