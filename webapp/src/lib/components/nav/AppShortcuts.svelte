<script lang="ts">
    import { shortcut } from '$lib/actions/shortcut.svelte';
    import { commandContext } from "$lib/context";
    const commandState = commandContext.get();
</script>

<svelte:window
        use:shortcut={{
		key: 'k',
		ctrl: true,
		callback: () => commandState.setTrue()
	}}
/>