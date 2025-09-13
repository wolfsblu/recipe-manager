<script lang="ts">
    import { FileDropZone, type FileDropZoneProps, MEGABYTE } from '$lib/components/ui/file-drop-zone';
    import {onDestroy} from 'svelte';
    import { toast } from 'svelte-sonner';
    import { type DragDropState, draggable, droppable } from "@thisux/sveltednd";
    import { flip } from 'svelte/animate';
    import { fade } from 'svelte/transition';
    import { cn } from "$lib/utils/utils";
    import { UploadService } from '$lib/upload/uploadService.svelte';
    import UploadProgress from '$lib/components/ui/upload/UploadProgress.svelte';
    import UploadError from '$lib/components/ui/upload/UploadError.svelte';
    import UploadSuccess from '$lib/components/ui/upload/UploadSuccess.svelte';

    let {
        value = $bindable(),
        class: className,
        maxFiles = 4,
        maxFileSize = 2 * MEGABYTE,
    } = $props<{
        value?: string[];
        class?: string;
        maxFiles?: number;
        maxFileSize?: number;
    }>();

    const uploadService = new UploadService();

    const onUpload: FileDropZoneProps['onUpload'] = async (files) => {
        const promises = files.map(async (file) => {
            try {
                await uploadService.uploadFile(file);
            } catch (error) {
                toast.error(`${file.name} failed to upload!`);
            }
        });
        await Promise.allSettled(promises);
    };

    const onFileRejected: FileDropZoneProps['onFileRejected'] = async ({ reason, file }) => {
        toast.error(`${file.name} failed to upload!`, { description: reason });
    };

    const handleDrop = (state: DragDropState<any>) => {
        const { draggedItem, targetContainer } = state;
        const dragIndex = uploadService.uploadedFiles.findIndex((img) => img.id === draggedItem.id);
        const dropIndex = parseInt(targetContainer ?? '0');

        if (dragIndex !== -1 && !isNaN(dropIndex)) {
            uploadService.reorderFiles(dragIndex, dropIndex);
        }
    };

    // Sync completed URLs with bindable value
    $effect(() => {
        value = uploadService.completedUrls;
    });
</script>

<div class={cn('flex w-full flex-col gap-4', className)}>
    {#if uploadService.uploadedFiles.length === 0}
        <FileDropZone
            accept="image/*"
            fileCount={uploadService.uploadedFiles.length}
            {onUpload}
            {onFileRejected}
            {maxFiles}
            {maxFileSize}
        />
    {:else}
        <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4">
            {#each uploadService.uploadedFiles as file, i (file.id)}
                <div class="svelte-dnd-touch-feedback"
                    use:draggable={{
                        container: i.toString(),
                        dragData: file,
                        interactive: ['button'],
                    }}
                    use:droppable={{
                        container: i.toString(),
                        callbacks: { onDrop: handleDrop }
                    }}
                    in:fade={{ duration: 150 }}
                    out:fade={{ duration: 150 }}
                    animate:flip={{duration: 300 }}
                >
                    {#if file.status === 'uploading'}
                        <UploadProgress progress={file.progress} />
                    {:else if file.status === 'error'}
                        <UploadError onRetry={() => uploadService.retryUpload(file.id)} />
                    {:else if file.status === 'completed' && file.url}
                        <UploadSuccess
                            url={file.url}
                            name={file.name}
                            size={file.size}
                            onRemove={() => uploadService.removeFile(file.id)}
                        />
                    {/if}
                </div>
            {/each}

            {#if uploadService.uploadedFiles.length < maxFiles}
                <button
                    class="flex aspect-[4/3] items-center justify-center rounded-lg border-2 border-dashed border-gray-300 bg-gray-50 transition-colors hover:border-gray-400 hover:bg-gray-100 dark:border-gray-600 dark:bg-gray-800 dark:hover:border-gray-500 dark:hover:bg-gray-700"
                    onclick={() => {
                        const input = document.createElement('input');
                        input.type = 'file';
                        input.accept = 'image/*';
                        input.multiple = true;
                        input.onchange = async (e) => {
                            const target = e.target as HTMLInputElement;
                            if (target.files) {
                                await onUpload(Array.from(target.files));
                            }
                        };
                        input.click();
                    }}
                >
                    <div class="flex flex-col items-center gap-2 text-gray-500 dark:text-gray-400">
                        <svg class="h-8 w-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                        </svg>
                        <span class="text-sm">Add Image</span>
                    </div>
                </button>
            {/if}
        </div>
    {/if}
</div>

<style lang="postcss">
    @reference "tailwindcss";

    :global(.dragging) {
        @apply opacity-50 shadow-lg ring-2 ring-blue-400 rounded-lg;
    }

    :global(.drag-over) {
        @apply bg-blue-50;
    }
</style>