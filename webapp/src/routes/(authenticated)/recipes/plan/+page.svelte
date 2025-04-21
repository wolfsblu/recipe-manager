<script lang="ts">
    import { Alert } from "flowbite-svelte";
    import { InfoCircleSolid } from "flowbite-svelte-icons";

    // TODO: Use https://wicky.nillia.ms/cally/components/calendar-range/ for selecting the date range
    // TODO: Display as calendar next to schedule https://tailwindcss.com/plus/ui-blocks/application-ui/data-display/calendars#component-5bb74a1ac1c8cf96e7b3d1ea2bf4570b
</script>

<Alert border color="blue" class="items-start!">
    <span slot="icon">
        <InfoCircleSolid class="w-5 h-5" />
        <span class="sr-only">Info</span>
    </span>
    <p class="font-medium">
        Work in progress! This page will show a calendar for planning upcoming
        meals.
    </p>
</Alert>

