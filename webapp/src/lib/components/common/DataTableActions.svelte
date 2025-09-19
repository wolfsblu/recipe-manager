<script lang="ts" generics="T">
    import { Button } from "$lib/components/ui/button/index.js";
    import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
    import MoreHorizontalIcon from '@lucide/svelte/icons/more-horizontal';
    import EditIcon from '@lucide/svelte/icons/edit';
    import TrashIcon from '@lucide/svelte/icons/trash-2';

    interface Props {
        item: T;
        onEdit: (item: T) => void;
        onDelete: (item: T) => void;
    }

    let { item, onEdit, onDelete }: Props = $props();
</script>

<div class="flex justify-end">
    <DropdownMenu.Root>
        <DropdownMenu.Trigger>
            {#snippet child({ props })}
                <Button {...props} variant="ghost" class="h-8 w-8 p-0">
                    <span class="sr-only">Open menu</span>
                    <MoreHorizontalIcon class="h-4 w-4" />
                </Button>
            {/snippet}
        </DropdownMenu.Trigger>
        <DropdownMenu.Content align="end">
            <DropdownMenu.Item onclick={() => onEdit(item)}>
                <EditIcon class="mr-2 h-4 w-4" />
                Edit
            </DropdownMenu.Item>
            <DropdownMenu.Separator />
            <DropdownMenu.Item onclick={() => onDelete(item)} class="text-destructive">
                <TrashIcon class="mr-2 h-4 w-4" />
                Delete
            </DropdownMenu.Item>
        </DropdownMenu.Content>
    </DropdownMenu.Root>
</div>