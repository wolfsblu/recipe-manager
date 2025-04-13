import { Upload } from 'tus-js-client'
import { PUBLIC_API_URL } from "$env/static/public";

export const uploadFile = async (file: File) => {
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
    })

    const previousUploads = await upload.findPreviousUploads()
    if (previousUploads.length) {
        upload.resumeFromPreviousUpload(previousUploads[0])
    }

    upload.start()
}