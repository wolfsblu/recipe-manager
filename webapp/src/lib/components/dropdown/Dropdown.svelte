<script lang="ts">
    import type {Snippet} from "svelte";
    import type {MenuItem} from "./types";

    interface Props {
        button: Snippet
        menu: MenuItem[]
    }

    let {button, menu}: Props = $props()
    let isOpen = $state(false)

    let dropdownNode: Node;
    window.addEventListener("click", (e: MouseEvent) => {
        const target = e.target as HTMLElement
        isOpen = dropdownNode && dropdownNode.contains(target);
    })
    const onClick = (e: MouseEvent) => {
        e.preventDefault()
        e.stopPropagation()
        isOpen = !isOpen;
    }

    const iconClass = "relative h-5 w-5"
    const menuItemClass = "flex gap-x-3 items-center justify-self-stretch px-4 py-2 text-gray-700 text-sm"
</script>

<div class="relative inline-block text-left" bind:this={dropdownNode}>
    <button onclick={onClick}>
        {@render button()}
    </button>

    <div class:hidden={!isOpen}
         class="
             absolute
             bg-white
             divide-gray-100
             divide-y
             focus:outline-none
             mt-2
             origin-top-right
             right-0
             ring-1
             ring-black/5
             rounded-md
             shadow-lg
             w-56
             z-10
        "
         role="menu" aria-orientation="vertical" aria-labelledby="menu-button" tabindex="-1">
        <div class="py-1" role="none">
            {#each menu as item}
                {#if item.href}
                    <a href={item.href}
                       tabindex="-1"
                       class={`
                            cursor-pointer
                            hover:bg-orange-50
                            ${menuItemClass}
                            ${item.class}
                    `}>
                        <span class={iconClass}>
                            <item.icon />
                        </span>
                        {item.label}
                    </a>
                {:else if item.onClick}
                    <button onclick={item.onClick}
                            class={`
                            hover:bg-orange-50
                            text-left
                            w-full
                            ${menuItemClass}
                            ${item.class}
                    `}>
                        <span class={iconClass}>
                            <item.icon />
                        </span>
                        {item.label}
                    </button>
                {:else}
                    <span class={`
                            ${menuItemClass}
                            ${item.class}
                    `}>
                        {item.label}
                    </span>
                {/if}
            {/each}
        </div>
    </div>
</div>