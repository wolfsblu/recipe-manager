<script>
    import FilePond, { registerPlugin, supported } from 'svelte-filepond';
    import * as tus from 'tus-js-client'
    // Import the Image EXIF Orientation and Image Preview plugins
    // Note: These need to be installed separately
    // `npm i filepond-plugin-image-preview filepond-plugin-image-exif-orientation --save`
    import FilePondPluginImageExifOrientation from 'filepond-plugin-image-exif-orientation';
    import FilePondPluginImageResize from 'filepond-plugin-image-resize';
    import FilePondPluginImageCrop from 'filepond-plugin-image-crop';
    import FilePondPluginImagePreview from 'filepond-plugin-image-preview';
    import 'filepond/dist/filepond.css'
    import 'filepond-plugin-image-preview/dist/filepond-plugin-image-preview.css'
    import './filepond.css'


    // Register the plugins
    registerPlugin(FilePondPluginImageExifOrientation, FilePondPluginImagePreview, FilePondPluginImageResize, FilePondPluginImageCrop);

    // a reference to the component, used to call FilePond methods
    let pond;

    // pond.getFiles() will return the active files

    // the name to use for the internal file input
    let name = 'filepond';

    // handle filepond events
    function handleInit() {
        console.log('FilePond has initialised');
    }

    function handleAddFile(err, fileItem) {
        console.log('A file has been added', fileItem);
    }

    const server = {
        process: (fieldName, file, metadata, load, error, progress, abort) => {
            var upload = new tus.Upload(file, {
                endpoint: "http://127.0.0.1:8080/upload/",
                retryDelays: [0, 1000, 3000, 5000],
                metadata: {
                    filename: file.name,
                    filetype: file.type
                },
                onError: (err) => error(err),
                onProgress: (bytesUploaded, bytesTotal) => progress(true, bytesUploaded, bytesTotal),
                onSuccess: () => load(upload.url?.split('/').pop())
            })
            upload.findPreviousUploads().then((previousUploads) => {
                if (previousUploads.length) {
                    upload.resumeFromPreviousUpload(previousUploads[0])
                }
                upload.start()
            })
            return {
                abort() {
                    upload.abort()
                    abort()
                }
            }
        }
    }
</script>

<FilePond
        bind:this={pond}
        {name}
        server={server}
        allowMultiple={true}
        allowReorder={true}
        maxParallelUploads={10}
        oninit={handleInit}
        onaddfile={handleAddFile}
        credits={false}
/>