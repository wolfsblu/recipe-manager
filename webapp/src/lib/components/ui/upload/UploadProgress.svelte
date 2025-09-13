<script lang="ts">
    import { Progress } from '$lib/components/ui/progress';

    let { progress = 0 } = $props();
</script>

<div class="relative aspect-[4/3] overflow-hidden rounded-lg bg-gray-100 dark:bg-gray-800">
    <div class="absolute inset-0 flex flex-col items-center justify-center gap-2">
        <Progress class="h-2 w-3/4" value={progress} max={100} />
        <div class="text-sm text-gray-600 dark:text-gray-400">
            {progress}% uploaded
        </div>
    </div>
</div>