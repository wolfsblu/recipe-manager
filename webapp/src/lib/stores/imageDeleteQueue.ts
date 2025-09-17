import { writable } from 'svelte/store';

export const imageDeleteQueue = writable<string[]>([]);

export const queueImageDeletion = (url: string) => {
    imageDeleteQueue.update(urls => {
        if (!urls.includes(url)) {
            return [...urls, url];
        }
        return urls;
    });
};

export const commitImageDeletions = async () => {
    let urlsToDelete: string[] = [];
    
    imageDeleteQueue.update(urls => {
        urlsToDelete = [...urls];
        return [];
    });

    const deletePromises = urlsToDelete.map(async (url) => {
        try {
            await fetch(url, {
                method: 'DELETE',
                headers: {
                    'Tus-Resumable': '1.0.0',
                }
            });
        } catch (error) {
            console.warn(`Error deleting file from server: ${url}`, error);
        }
    });

    await Promise.allSettled(deletePromises);
};

export const clearImageDeletions = () => {
    imageDeleteQueue.set([]);
};