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
        for (let i = 0; i < files.length; ++i) {
            const file = files[i]
            const upload = await uploadFile(file)
            images.push({
                alt: "Image",
                src: upload.url ?? "",
            })
        }
    }
    
    const iconClass = "w-6 h-6"
</script>

<Fileupload accept="image/*" multiple onchange={onFileChange} />

<Gallery items={images} class="gap-4 grid-cols-1 sm:grid-cols-2 md:grid-cols-4 lg:grid-cols-8">

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