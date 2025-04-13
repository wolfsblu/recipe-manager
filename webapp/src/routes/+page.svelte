<script lang="ts">
    import * as m from "$lib/paraglide/messages"
    import {Fileupload} from "flowbite-svelte";
    import {Upload} from 'tus-js-client'
    import type {ChangeEventHandler} from "svelte/elements";

    const onFileChange: ChangeEventHandler<HTMLInputElement> = async (e) => {
        const files = (e.target as HTMLInputElement)?.files
        if (files == null || files.length <= 0) {
            return
        }
        const file = files[0]
        const upload = new Upload(file, {
            endpoint: 'http://127.0.0.1:8080/api/uploads/',
            retryDelays: [0, 3000, 5000, 10000, 20000],
            metadata: {
                filename: file.name,
                filetype: file.type,
            },
            onError: (error) => {
                console.error(error)
            },
            onProgress: (bytesUploaded, bytesTotal) => {
                const percentage = ((bytesUploaded / bytesTotal) * 100).toFixed(2)
                console.log(bytesUploaded, bytesTotal, percentage + "%")
            },
            onSuccess: () => {
                console.log("Download %s from %s", upload.file.name, upload.url)
                console.log(upload)
            }
        })

        const previousUploads = await upload.findPreviousUploads()
        if (previousUploads.length) {
            upload.resumeFromPreviousUpload(previousUploads[0])
        }
        upload.start()
    }
</script>

<h1>{m.example_message({username: "Frank"})}</h1>

<p>
    Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the documentation
</p>
<p>
    Go to <a href="about">about page</a>
</p>

<Fileupload accept="image/*" onchange={onFileChange} />