<script lang="ts">
    import { Button } from '$lib/components/ui/button';
    import { Card, CardContent } from '$lib/components/ui/card';
    import ShoppingCartIcon from '@lucide/svelte/icons/shopping-cart';
    import PlusIcon from '@lucide/svelte/icons/plus';
    import * as m from "$lib/paraglide/messages.js";

    interface Props {
        onCreateList: () => void;
    }

    let { onCreateList }: Props = $props();
</script>

<div class="flex-grow flex items-center justify-center px-5">
    <Card class="max-w-md w-full">
        <CardContent>
            <div class="text-center">
                <div class="p-4 bg-muted rounded-full mx-auto w-fit mb-4">
                    <ShoppingCartIcon class="h-8 w-8 text-muted-foreground" />
                </div>
                <h3 class="text-lg font-semibold mb-2">{m.shopping_empty_title()}</h3>
                <p class="text-muted-foreground mb-4">
                    {m.shopping_empty_description()}
                </p>
                <Button onclick={onCreateList}>
                    <PlusIcon />
                    {m.shopping_empty_button()}
                </Button>
            </div>
        </CardContent>
    </Card>
</div>