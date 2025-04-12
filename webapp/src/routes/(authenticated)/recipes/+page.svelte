<script lang="ts">
    import {SpeedDial, SpeedDialButton} from "flowbite-svelte";
    import {
        ShoppingBagSolid,
        ReceiptSolid,
        RulerCombinedOutline
    } from "flowbite-svelte-icons";
    
    const iconClass = "w-6 h-6"
</script>

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