<script lang="ts">
    import { Button } from '$lib/components/ui/button';
    import * as Dialog from '$lib/components/ui/dialog';
    import { Input } from '$lib/components/ui/input';
    import { Label } from '$lib/components/ui/label';
    import { createShoppingList } from '$lib/api/shopping/shopping.svelte';
    import { shoppingStore } from '$lib/stores/shopping.svelte';
    import { toast } from 'svelte-sonner';

    interface Props {
        open: boolean;
        onOpenChange: (open: boolean) => void;
    }

    let { open, onOpenChange }: Props = $props();

    let listName = $state('');
    let isCreating = $state(false);

    const handleSubmit = async (e: Event) => {
        e.preventDefault();
        
        if (!listName.trim()) {
            toast.error('Please enter a shopping list name');
            return;
        }

        isCreating = true;
        try {
            const newList = await createShoppingList({ name: listName.trim() });
            shoppingStore.addList(newList);
            shoppingStore.setCurrentListId(newList.id);
            
            toast.success('Shopping list created successfully');
            listName = '';
            onOpenChange(false);
        } catch (error) {
            toast.error('Failed to create shopping list');
        } finally {
            isCreating = false;
        }
    };

    const handleOpenChange = (newOpen: boolean) => {
        if (!newOpen) {
            listName = '';
        }
        onOpenChange(newOpen);
    };
</script>

<Dialog.Root {open} onOpenChange={handleOpenChange}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>New Shopping List</Dialog.Title>
            <Dialog.Description>
                Give your shopping list a name to help you organize your items.
            </Dialog.Description>
        </Dialog.Header>
        
        <form onsubmit={handleSubmit} class="space-y-4">
            <div class="space-y-2">
                <Input
                    id="list-name"
                    bind:value={listName}
                    placeholder="e.g., Weekly Groceries, Party Shopping"
                    disabled={isCreating}
                    required
                />
            </div>
            
            <Dialog.Footer>
                <Button
                    type="button"
                    variant="outline"
                    onclick={() => onOpenChange(false)}
                    disabled={isCreating}
                >
                    Cancel
                </Button>
                <Button type="submit" disabled={isCreating || !listName.trim()}>
                    {isCreating ? 'Creating...' : 'Create List'}
                </Button>
            </Dialog.Footer>
        </form>
    </Dialog.Content>
</Dialog.Root>