<script lang="ts">
    import { Button } from '$lib/components/ui/button';
    import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
    import ShoppingCartIcon from '@lucide/svelte/icons/shopping-cart';
    import PlusIcon from '@lucide/svelte/icons/plus';

    interface Props {
        onCreateList: () => void;
    }

    let { onCreateList }: Props = $props();
</script>

<div class="min-h-[80vh] flex items-center justify-center p-5">
    <Card class="max-w-md w-full">
        <CardHeader class="text-center">
            <div class="flex justify-center mb-4">
                <div class="p-4 bg-muted rounded-full">
                    <ShoppingCartIcon class="h-8 w-8 text-muted-foreground" />
                </div>
            </div>
            <CardTitle class="text-xl">No Shopping Lists</CardTitle>
        </CardHeader>
        <CardContent class="text-center space-y-4">
            <p class="text-muted-foreground">
                You haven't created any shopping lists yet. Create your first shopping list to start organizing your shopping items.
            </p>
            <Button onclick={onCreateList} class="w-full">
                <PlusIcon />
                Create Your First Shopping List
            </Button>
        </CardContent>
    </Card>
</div>