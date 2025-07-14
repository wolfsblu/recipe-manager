<script lang="ts">
    import * as Dialog from '$lib/components/ui/dialog';
    import * as Command from '$lib/components/ui/command';
    import { goto } from '$app/navigation';
    import { commandContext } from '$lib/context';

    const commandState = commandContext.get();
</script>

<Dialog.Root bind:open={commandState.current}>
    <Dialog.Content class="p-0" hideClose>
        <Command.Root>
            <Command.Input placeholder="Search for extras..." />
            <Command.List class="min-h-[300px]">
                <Command.Empty>No results found.</Command.Empty>
                    <Command.Group heading="Test">
                            <Command.Item>
                                Lala
                            </Command.Item>
                    </Command.Group>
            </Command.List>
        </Command.Root>
    </Dialog.Content>
</Dialog.Root>