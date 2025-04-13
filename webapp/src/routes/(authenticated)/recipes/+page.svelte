<script lang="ts">
    import { uploadFile } from '$lib/upload/tus'
    import {type ImgType, Fileupload, Gallery, SpeedDial, SpeedDialButton} from "flowbite-svelte";
    import {
        ShoppingBagSolid,
        ReceiptSolid,
        RulerCombinedOutline
    } from "flowbite-svelte-icons";
    import type { ChangeEventHandler } from "svelte/elements";

    const images = $state<ImgType[]>([]);

    const onFileChange: ChangeEventHandler<HTMLInputElement> = async (e) => {
        const files = (e.target as HTMLInputElement)?.files
        if (files == null) {
            return
        }
        const uploadPromises = []
        for (let i = 0; i < files.length; ++i) {
            const file = files[i]
            uploadPromises.push(uploadFile(file))
        }

        const uploads = await Promise.allSettled(uploadPromises)
        for (let i = 0; i < uploads.length; ++i) {
            const upload = uploads[i]
            if (upload.status === "fulfilled" && upload.value.url) {
                images.push({
                    alt: "Image",
                    src: upload.value.url,
                })
            }
        }
    }
    
    const iconClass = "w-6 h-6"
</script>

<Fileupload accept="image/*" multiple onchange={onFileChange} />

<Gallery items={images} class="gap-4 grid-cols-2 md:grid-cols-4">

</Gallery>

<SpeedDial>
    <SpeedDialButton name="Measurement">
        <RulerCombinedOutline class={iconClass} />
    </SpeedDialButton>
    <SpeedDialButton name="Ingredient">
        <ShoppingBagSolid class={iconClass} />
    </SpeedDialButton>
    <SpeedDialButton name="Recipe">
        <ReceiptSolid class={iconClass} />
    </SpeedDialButton>
</SpeedDial>