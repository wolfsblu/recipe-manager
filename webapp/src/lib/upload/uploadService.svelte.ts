import * as tus from 'tus-js-client';
import { PUBLIC_UPLOAD_URL } from '$env/static/public';

export type UploadFile = {
    id: number;
    name: string;
    type: string;
    size: number;
    progress: number;
    status: 'uploading' | 'completed' | 'error';
    url: string | null;
    upload?: tus.Upload;
};

export class UploadService {
    private files = $state<UploadFile[]>([]);

    get uploadedFiles() {
        return this.files;
    }

    get completedUrls() {
        return this.files
            .filter(file => file.status === 'completed' && file.url)
            .map(file => file.url!);
    }

    async uploadFile(file: File): Promise<UploadFile> {
        if (this.files.find((f) => f.name === file.name)) {
            throw new Error('File already exists');
        }

        const fileData: UploadFile = {
            id: this.files.length,
            name: file.name,
            type: file.type,
            size: file.size,
            progress: 0,
            status: 'uploading',
            url: null
        };

        const length = this.files.push(fileData);

        const upload = new tus.Upload(file, {
            endpoint: PUBLIC_UPLOAD_URL,
            retryDelays: [0, 3000, 5000, 10000, 20000],
            metadata: {
                filename: file.name,
                filetype: file.type,
            },
            onError: () => {
                fileData.status = 'error';
                this.files[length - 1] = fileData
            },
            onProgress: (bytesUploaded, bytesTotal) => {
                fileData.progress = Math.round((bytesUploaded / bytesTotal) * 100);
                this.files[length - 1] = fileData
            },
            onSuccess: () => {
                fileData.status = 'completed';
                fileData.url = upload.url || '';
                this.files[length - 1] = fileData
            }
        });

        fileData.upload = upload;
        upload.start();

        return fileData;
    }

    async removeFile(id: number): Promise<void> {
        const fileIndex = this.files.findIndex(f => f.id === id);
        if (fileIndex === -1) return;

        const file = this.files[fileIndex];

        if (file.upload && file.status === 'uploading') {
            await file.upload.abort();
        }

        if (file.status === 'completed' && file.url) {
            await this.deleteFromServer(file.url);
        }

        this.files.splice(fileIndex, 1);
    }

    async retryUpload(id: number): Promise<void> {
        const file = this.files.find(f => f.id === id);
        if (!file || !file.upload) return;

        file.status = 'uploading';
        file.progress = 0;
        file.upload.start();
    }

    reorderFiles(dragIndex: number, dropIndex: number): void {
        if (dragIndex === dropIndex) return;

        const [item] = this.files.splice(dragIndex, 1);
        this.files.splice(dropIndex, 0, item);
    }

    async cleanup(): Promise<void> {
        this.files.forEach(file => {
            if (file.upload && file.status === 'uploading') {
                file.upload.abort();
            }
        });

        await Promise.allSettled(
            this.completedUrls.map(url => this.deleteFromServer(url))
        );

        this.files = [];
    }

    private async deleteFromServer(url: string): Promise<void> {
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
    }
}