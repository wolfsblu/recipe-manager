<script lang="ts">
    import { Button } from '$lib/components/ui/button';
    import { displaySize } from '$lib/components/ui/file-drop-zone';
    import CloseIcon from '@lucide/svelte/icons/x';

    let {
        url,
        name,
        size,
        onRemove
    } = $props<{
        url: string;
        name: string;
        size: number;
        onRemove: () => void;
    }>();
</script>

<div class="group relative aspect-[4/3] overflow-hidden rounded-lg">
    <img
        src={url}
        alt={name}
        class="h-full w-full object-cover transition-transform group-hover:scale-105"
    />
    <div class="absolute inset-0 bg-black/0 transition-colors group-hover:bg-black/20"></div>
    <Button
        variant="destructive"
        size="icon"
        class="absolute top-2 right-2 h-6 w-6 opacity-0 transition-opacity group-hover:opacity-100"
        onclick={onRemove}
    >
        <CloseIcon class="h-3 w-3" />
    </Button>
    <div class="absolute bottom-0 left-0 right-0 bg-black/50 p-2 text-white text-xs opacity-0 transition-opacity group-hover:opacity-100">
        <div class="truncate">{name}</div>
        <div>{displaySize(size)}</div>
    </div>
</div>