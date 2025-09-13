<script lang="ts">
    import {Button} from '$lib/components/ui/button';
    import {displaySize, FileDropZone, type FileDropZoneProps, MEGABYTE} from '$lib/components/ui/file-drop-zone';
    import {Progress} from '$lib/components/ui/progress';
    import {onDestroy} from 'svelte';
    import {toast} from 'svelte-sonner';
    import {type DragDropState, draggable, droppable} from "@thisux/sveltednd";
    import {flip} from 'svelte/animate';
    import {fade} from 'svelte/transition';
    import {cn} from "$lib/utils/utils";
    import * as tus from 'tus-js-client';
    import {PUBLIC_UPLOAD_URL} from '$env/static/public';
    import CloseIcon from '@lucide/svelte/icons/x'

    let {
        value = $bindable(),
        class: className,
        maxFiles = 4,
        maxFileSize = 2 * MEGABYTE,
    } = $props()


    type UploadedFile = {
        id: number;
        name: string;
        type: string;
        size: number;
        uploadedAt: number;
        progress: number;
        status: 'uploading' | 'completed' | 'error';
        url: string | null;
        upload?: tus.Upload;
    };

    let files = $state<UploadedFile[]>([]);

    const onUpload: FileDropZoneProps['onUpload'] = async (files) => {
        await Promise.allSettled(files.map((file) => uploadFile(file)));
    };

    const onFileRejected: FileDropZoneProps['onFileRejected'] = async ({ reason, file }) => {
        toast.error(`${file.name} failed to upload!`, { description: reason });
    };

    const uploadFile = async (file: File) => {
        if (files.find((f) => f.name === file.name)) return;

        const fileData: UploadedFile = {
            id: files.length,
            name: file.name,
            type: file.type,
            size: file.size,
            uploadedAt: Date.now(),
            progress: 0,
            status: 'uploading',
            url: null
        };

        files.push(fileData);

        const upload = new tus.Upload(file, {
            endpoint: PUBLIC_UPLOAD_URL,
            retryDelays: [0, 3000, 5000, 10000, 20000],
            metadata: {
                filename: file.name,
                filetype: file.type,
            },
            onError: () => {
                fileData.status = 'error';
                files[0] = fileData
            },
            onProgress: (bytesUploaded, bytesTotal) => {
                fileData.progress = Math.round((bytesUploaded / bytesTotal) * 100);
                files[0] = fileData
            },
            onSuccess: () => {
                fileData.status = 'completed';
                fileData.url = upload.url || '';
                files[0] = fileData
            }
        });

        fileData.upload = upload;
        upload.start();
    };

    $effect(() => {
        value = files
            .filter(file => file.status === 'completed' && file.url)
            .map(file => file.url!)
    });

    onDestroy(() => {
        for (const file of files) {
            if (file.upload && file.status === 'uploading') {
                file.upload.abort();
            }
        }
    });

    const handleDrop = (state: DragDropState<UploadedFile>) => {
        const { draggedItem, targetContainer } = state;
        const dragIndex = files.findIndex((img) => img.id === draggedItem.id);
        const dropIndex = parseInt(targetContainer ?? '0');

        if (dragIndex !== -1 && !isNaN(dropIndex)) {
            const [item] = files.splice(dragIndex, 1);
            files.splice(dropIndex, 0, item);
        }
    }
</script>

<div class={cn('flex w-full flex-col gap-4', className)}>
    {#if files.length === 0}
        <FileDropZone
                accept="image/*"
                fileCount={files.length}
                {onUpload}
                {onFileRejected}
                {maxFiles}
                {maxFileSize}
        />
    {:else}
        <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4">
            {#each files as file, i (file.id)}
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
                        <div class="relative aspect-[4/3] overflow-hidden rounded-lg bg-gray-100 dark:bg-gray-800">
                            <div class="absolute inset-0 flex flex-col items-center justify-center gap-2">
                                <Progress
                                        class="h-2 w-3/4"
                                        value={file.progress}
                                        max={100}
                                />
                                <div class="text-sm text-gray-600 dark:text-gray-400">
                                    {file.progress}% uploaded
                                </div>
                            </div>
                        </div>
                    {:else if file.status === 'error'}
                        <div class="relative aspect-[4/3] overflow-hidden rounded-lg bg-red-100 dark:bg-red-900/20">
                            <div class="absolute inset-0 flex flex-col items-center justify-center gap-2">
                                <div class="text-red-600 dark:text-red-400">Upload failed</div>
                                <Button
                                    size="sm"
                                    onclick={() => {
                                        // Retry upload
                                        if (file.upload) {
                                            file.status = 'uploading';
                                            file.progress = 0;
                                            file.upload.start();
                                            files = files.slice();
                                        }
                                    }}
                                >
                                    Retry
                                </Button>
                            </div>
                        </div>
                    {:else if file.status === 'completed' && file.url}
                        <div class="group relative aspect-[4/3] overflow-hidden rounded-lg">
                            <img
                                    src={file.url}
                                    alt={file.name}
                                    class="h-full w-full object-cover transition-transform group-hover:scale-105"
                            />
                            <div class="absolute inset-0 bg-black/0 transition-colors group-hover:bg-black/20"></div>
                            <Button
                                    variant="destructive"
                                    size="icon"
                                    class="absolute top-2 right-2 h-6 w-6 opacity-0 transition-opacity group-hover:opacity-100"
                                    onclick={() => {
                                        // Cancel upload if still in progress
                                        if (file.upload && file.status === 'uploading') {
                                            file.upload.abort();
                                        }
                                        files = [...files.slice(0, i), ...files.slice(i + 1)];
                                    }}
                            >
                                <CloseIcon class="h-3 w-3" />
                            </Button>
                            <div class="absolute bottom-0 left-0 right-0 bg-black/50 p-2 text-white text-xs opacity-0 transition-opacity group-hover:opacity-100">
                                <div class="truncate">{file.name}</div>
                                <div>{displaySize(file.size)}</div>
                            </div>
                        </div>
                    {/if}
                </div>
            {/each}
            
            {#if files.length < maxFiles}
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