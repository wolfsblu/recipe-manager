<script lang="ts">
    import {Progressbar, Toast} from "flowbite-svelte";
    import {ExclamationCircleOutline} from "flowbite-svelte-icons";
    import {blur} from 'svelte/transition'
    import {linear} from "svelte/easing";
    import {type Toast as ToastType, toasts} from "$lib/components/notifications/toasts";

    type ToastColor = "blue" | "red" | "green" | "yellow"

    const getColor = (toast: ToastType): ToastColor => {
        switch (toast.type) {
            case "error":
                return "red";
            case "info":
                return "blue";
            case "success":
                return "green";
            case "warning":
                return "yellow";
            default:
                return "blue";
        }
    }
</script>

<section class="p-3 absolute top-0 start-0 md:start-auto md:end-0 w-full md:max-w-1/2 lg:max-w-1/3">
    {#each $toasts as toast (toast.id)}
        <Toast bind:toastStatus={toast.visible}
               color={getColor(toast)}
               divClass="md:w-full p-4 text-gray-500 bg-white shadow-sm dark:text-gray-400 dark:bg-gray-800 gap-3"
               transition={blur}
               params={{ amount: 10 }}
        >
            <svelte:fragment slot="icon">
                <ExclamationCircleOutline class="w-5 h-5"/>
                <span class="sr-only">Error icon</span>
            </svelte:fragment>

            {toast.message}
        </Toast>


        {#if toast.visible}
            <div transition:blur={{ amount: 10 }}>
                <Progressbar progress={100}
                             animate
                             color={getColor(toast)}
                             class="mb-3"
                             size="h-0.5"
                             easing={linear}
                             tweenDuration={toast.timeout}
                />
            </div>
        {/if}
    {/each}
</section>