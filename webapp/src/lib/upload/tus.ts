import { Upload } from 'tus-js-client'
import { PUBLIC_API_URL } from "$env/static/public";

interface UploadResult {
    file: File
    url: string | null
}

export const uploadFile = async (file: File): Promise<UploadResult> => {
    return new Promise(async (resolve, reject) => {
        const upload = new Upload(file, {
            endpoint: PUBLIC_API_URL + '/uploads/',
            retryDelays: [0, 3000, 5000, 10000, 20000],
            metadata: {
                filename: file.name,
                filetype: file.type,
            },
            onBeforeRequest: (req) => {
                const xhr = req.getUnderlyingObject()
                xhr.withCredentials = true
            },
            onSuccess: () => {
                resolve({
                    file: upload.file,
                    url: upload.url
                })
            },
            onError: (err) => {
                reject(err)
            }
        })

        try {
            const previousUploads = await upload.findPreviousUploads()
            if (previousUploads.length) {
                upload.resumeFromPreviousUpload(previousUploads[0])
            }

            upload.start()
        } catch (e) {
            reject(e)
        }
    })
}