<script lang="ts">
    import placeholder from '$lib/components/recipes/fruits.jpg'
    import {
        Button,
    } from "flowbite-svelte";
    import {TrashBinOutline} from "flowbite-svelte-icons";

    interface Props {
        title: string
    }

    let {
        title
    }: Props = $props()
</script>

<div class="flex flex-col justify-between w-full h-36 bg-cover rounded-md shadow-md" style="background-image: url({placeholder})">
    <div class="flex justify-end p-1">
        <Button class="h-7 w-7 p-2" color="red" outline pill>
            <TrashBinOutline class="w-4 h-4"/>
        </Button>
    </div>
    <div class="flex flex-col justify-end py-3 h-1/2 text-white bg-gradient-to-b from-transparent to-gray-800 rounded-b-md">
        <h1 class="font-semibold">
            {title}
        </h1>
    </div>
</div>