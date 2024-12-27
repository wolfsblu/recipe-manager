<script lang="ts">
    import type {Snippet} from "svelte";

    interface Props {
        button: Snippet
        children: Snippet
    }

    let {button, children}: Props = $props()
    let isOpen = $state(false)

    let dropdownNode: Node;
    window.addEventListener("click", (e: MouseEvent) => {
        if (!dropdownNode) {
            return;
        }
        const target = e.target as HTMLElement
        isOpen = dropdownNode.contains(target);
    })
</script>

<div class="relative inline-block text-left" bind:this={dropdownNode}>
    <button>
        {@render button()}
    </button>

    <!--
      Dropdown menu, show/hide based on menu state.

      Entering: "transition ease-out duration-100"
        From: "transform opacity-0 scale-95"
        To: "transform opacity-100 scale-100"
      Leaving: "transition ease-in duration-75"
        From: "transform opacity-100 scale-100"
        To: "transform opacity-0 scale-95"
    -->
    <div class:hidden={!isOpen} class="absolute right-0 z-10 mt-2 w-56 origin-top-right divide-y divide-gray-100 rounded-md bg-white shadow-lg ring-1 ring-black/5 focus:outline-none"
         role="menu" aria-orientation="vertical" aria-labelledby="menu-button" tabindex="-1">
        <div class="py-1" role="none">
            <a href="#" class="block px-4 py-2 text-sm text-gray-700" role="menuitem" tabindex="-1" id="menu-item-0">Profile</a>
        </div>
        <div class="py-1" role="none">
            <a href="#" class="block px-4 py-2 text-sm text-gray-700" role="menuitem" tabindex="-1" id="menu-item-2">Sign Out</a>
        </div>
    </div>
</div>