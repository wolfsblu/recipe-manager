<script lang="ts">
    import { Dropdown, DropdownItem, Avatar } from "flowbite-svelte";
    import { EyeSolid } from "flowbite-svelte-icons";
</script>

<Dropdown
    placement="right-end"
    class="w-full max-w-sm rounded-sm divide-y divide-gray-100 shadow-sm dark:bg-gray-800 dark:divide-gray-700"
>
    <div slot="header" class="text-center py-2 font-bold">Notifications</div>
    <DropdownItem class="flex space-x-4 rtl:space-x-reverse">
        <Avatar
            src="/images/profile-picture-1.webp"
            dot={{ color: "bg-gray-300" }}
            rounded
        />
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
    </DropdownItem>
    <DropdownItem class="flex space-x-4 rtl:space-x-reverse">
        <Avatar
            src="/images/profile-picture-2.webp"
            dot={{ color: "bg-red-400" }}
            rounded
        />
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
    </DropdownItem>
    <DropdownItem class="flex space-x-4 rtl:space-x-reverse">
        <Avatar
            src="/images/profile-picture-3.webp"
            dot={{ color: "bg-green-400" }}
            rounded
        />
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
    </DropdownItem>
    <a
        slot="footer"
        href="/"
        class="block py-2 -my-1 text-sm font-medium text-center text-gray-900 bg-gray-50 hover:bg-gray-100 dark:bg-gray-800 dark:hover:bg-gray-700 dark:text-white"
    >
        <div class="inline-flex items-center">
            <EyeSolid class="me-2 w-4 h-4 text-gray-500 dark:text-gray-400" />
            View all
        </div>
    </a>
</Dropdown>
