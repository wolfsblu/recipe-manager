<script lang="ts">
    import { Button } from '$lib/components/ui/button';

    let { onRetry } = $props<{
        onRetry: () => void;
    }>();
</script>

<div class="relative aspect-[4/3] overflow-hidden rounded-lg bg-red-100 dark:bg-red-900/20">
    <div class="absolute inset-0 flex flex-col items-center justify-center gap-2">
        <div class="text-red-600 dark:text-red-400">Upload failed</div>
        <Button size="sm" onclick={onRetry}>
            Retry
        </Button>
    </div>
</div>