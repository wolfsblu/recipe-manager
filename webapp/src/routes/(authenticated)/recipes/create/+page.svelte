<script lang="ts">
    import { Alert } from "flowbite-svelte";
    import { InfoCircleSolid } from "flowbite-svelte-icons";
</script>

<Alert border color="blue" class="items-start!">
    {#snippet icon()}
        <InfoCircleSolid class="w-5 h-5" />
        <span class="sr-only">Info</span>
    {/snippet}
    <p class="font-medium">
        Work in progress! This page will allow you to create new recipes.
    </p>
</Alert>
