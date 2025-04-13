<script lang="ts">
    import { uploadFile } from '$lib/upload/tus'
    import { Fileupload, SpeedDial, SpeedDialButton } from "flowbite-svelte";
    import {
        ShoppingBagSolid,
        ReceiptSolid,
        RulerCombinedOutline
    } from "flowbite-svelte-icons";
    import type { ChangeEventHandler } from "svelte/elements";

    const onFileChange: ChangeEventHandler<HTMLInputElement> = async (e) => {
        const files = (e.target as HTMLInputElement)?.files
        if (files == null) {
            return
        }
        for (let i = 0; i < files.length; ++i) {
            const file = files[i]
            await uploadFile(file)
        }
    }
    
    const iconClass = "w-6 h-6"
</script>

<Fileupload accept="image/*" multiple onchange={onFileChange} />

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