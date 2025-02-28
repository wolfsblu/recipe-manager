<script lang="ts">
    import {onMount, type Snippet} from "svelte";

    interface Props {
        button: Snippet
        children: Snippet
    }

    let {button, children}: Props = $props()
    let isOpen = $state(false)

    let dropdownNode: Node;
    const handleClick = (e: MouseEvent) => {
        const target = e.target as HTMLElement
        isOpen = dropdownNode && dropdownNode.contains(target);
    }
    const handleKeyDown = (e: KeyboardEvent) => {
        if (e.key === "Escape" && isOpen) {
            isOpen = false;
        }
    }
    onMount(() => {
        window.addEventListener("click", handleClick)
        window.addEventListener("keydown", handleKeyDown)
        return () => {
            window.removeEventListener("keydown", handleKeyDown)
            window.removeEventListener("click", handleClick)
        }
    })

    const onOpen = (e: MouseEvent) => {
        e.preventDefault()
        e.stopPropagation()
        isOpen = !isOpen;
    }
</script>

<div class="relative inline-block text-left" bind:this={dropdownNode}>
    <button onclick={onOpen}>
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
             min-w-56
             max-w-96
             z-10
        "
         role="menu" aria-orientation="vertical" aria-labelledby="menu-button" tabindex="-1">
        {@render children()}
    </div>
</div>