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

<section class="absolute top-5 end-5">
    {#each $toasts as toast (toast.id)}
        <Toast bind:toastStatus={toast.visible}
               class="min-w-80"
               color={getColor(toast)}
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
                <Progressbar bind:progress={toast.progress}
                             color={getColor(toast)}
                             animate class="mb-3"
                             size="h-0.5"
                             easing={linear}
                             tweenDuration={toast.timeout}
                />
            </div>
        {/if}
    {/each}
</section>