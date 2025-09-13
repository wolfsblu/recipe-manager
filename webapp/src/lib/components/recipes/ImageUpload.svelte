<script lang="ts">
    import { Button } from '$lib/components/ui/button';
    import {
        displaySize,
        FileDropZone,
        MEGABYTE,
        type FileDropZoneProps
    } from '$lib/components/ui/file-drop-zone';
    import { Progress } from '$lib/components/ui/progress';
    import CloseIcon from '@lucide/svelte/icons/x';
    import { onDestroy } from 'svelte';
    import { toast } from 'svelte-sonner';
    import { SvelteDate } from 'svelte/reactivity';

    const onUpload: FileDropZoneProps['onUpload'] = async (files) => {
        await Promise.allSettled(files.map((file) => uploadFile(file)));
    };

    const onFileRejected: FileDropZoneProps['onFileRejected'] = async ({ reason, file }) => {
        toast.error(`${file.name} failed to upload!`, { description: reason });
    };

    const uploadFile = async (file: File) => {
        // don't upload duplicate files
        if (files.find((f) => f.name === file.name)) return;

        const urlPromise = new Promise<string>((resolve) => {
            // add some fake loading time
            resolve(URL.createObjectURL(file));
        });

        files.push({
            name: file.name,
            type: file.type,
            size: file.size,
            uploadedAt: Date.now(),
            url: urlPromise
        });

        // we await since we don't want the onUpload to be complete until the files are actually uploaded
        await urlPromise;
    };

    type UploadedFile = {
        name: string;
        type: string;
        size: number;
        uploadedAt: number;
        url: Promise<string>;
    };

    let files = $state<UploadedFile[]>([]);
    let date = new SvelteDate();

    onDestroy(async () => {
        for (const file of files) {
            URL.revokeObjectURL(await file.url);
        }
    });

    $effect(() => {
        const interval = setInterval(() => {
            date.setTime(Date.now());
        }, 10);

        return () => {
            clearInterval(interval);
        };
    });
</script>

<div class="flex w-full flex-col gap-4">
    {#if files.length === 0}
        <FileDropZone
                {onUpload}
                {onFileRejected}
                maxFileSize={2 * MEGABYTE}
                accept="image/*"
                maxFiles={5}
                fileCount={files.length}
        />
    {:else}
        <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4">
            {#each files as file, i (file.name)}
                {#await file.url}
                    <div class="relative aspect-[4/3] overflow-hidden rounded-lg bg-gray-100 dark:bg-gray-800">
                        <div class="absolute inset-0 flex items-center justify-center">
                            <Progress
                                    class="h-2 w-3/4"
                                    value={((date.getTime() - file.uploadedAt) / 1000) * 100}
                                    max={100}
                            />
                        </div>
                    </div>
                {:then src}
                    <div class="group relative aspect-[4/3] overflow-hidden rounded-lg">
                        <img
                                {src}
                                alt={file.name}
                                class="h-full w-full object-cover transition-transform group-hover:scale-105"
                        />
                        <div class="absolute inset-0 bg-black/0 transition-colors group-hover:bg-black/20"></div>
                        <Button
                                variant="destructive"
                                size="icon"
                                class="absolute top-2 right-2 h-6 w-6 opacity-0 transition-opacity group-hover:opacity-100"
                                onclick={() => {
                                    URL.revokeObjectURL(src);
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
                {/await}
            {/each}
            
            {#if files.length < 5}
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