<script lang="ts">
    import {Button, Listgroup, ListgroupItem, Modal} from "flowbite-svelte";
    import {EyeSolid} from "flowbite-svelte-icons";

    interface Props {
        open: boolean
    }

    let {
        open = $bindable(),
    }: Props = $props()
</script>

<Modal bind:open={open} autoclose outsideclose classBody="p-0 md:p-0">
    <h1 class="font-bold" slot="header">
        Notifications
    </h1>
    <Listgroup class="rounded-none border-0">
        <ListgroupItem class="flex space-x-4 rtl:space-x-reverse">
            <div class="ps-3 w-full">
                <div class="text-gray-500 text-sm mb-1.5 dark:text-gray-400">
                    New message from <span
                        class="font-semibold text-gray-900 dark:text-white"
                >Jese Leos</span
                >
                    : "Hey, what's up? All set for the presentation?"
                </div>
                <div class="text-xs text-primary-600 dark:text-primary-500">
                    a few moments ago
                </div>
            </div>
        </ListgroupItem>
        <ListgroupItem class="flex space-x-4 rtl:space-x-reverse">
            <div class="ps-3 w-full">
                <div class="text-gray-500 text-sm mb-1.5 dark:text-gray-400">
                <span class="font-semibold text-gray-900 dark:text-white"
                >Joseph Mcfall</span
                >
                    and
                    <span class="font-medium text-gray-900 dark:text-white"
                    >5 others</span
                    >
                    started following you.
                </div>
                <div class="text-xs text-primary-600 dark:text-primary-500">
                    10 minutes ago
                </div>
            </div>
        </ListgroupItem>
        <ListgroupItem class="flex space-x-4 rtl:space-x-reverse">
            <div class="ps-3 w-full">
                <div class="text-gray-500 text-sm mb-1.5 dark:text-gray-400">
                <span class="font-semibold text-gray-900 dark:text-white"
                >Bonnie Green</span
                >
                    and
                    <span class="font-medium text-gray-900 dark:text-white"
                    >141 others</span
                    >
                    love your story. See it and view more stories.
                </div>
                <div class="text-xs text-primary-600 dark:text-primary-500">
                    44 minutes ago
                </div>
            </div>
        </ListgroupItem>
    </Listgroup>
    <svelte:fragment slot="footer">
        <Button color="alternative" href="/user/notifications" onclick={() => open = false}>
            <EyeSolid class="me-2 w-4 h-4"/>
            View all
        </Button>
    </svelte:fragment>
</Modal>
