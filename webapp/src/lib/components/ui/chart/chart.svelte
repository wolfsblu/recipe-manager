<script lang="ts">
	import { cn } from '$lib/utils.js';
	
	let {
		class: className,
		...props
	}: {
		class?: string;
		config: Record<string, any>;
		[key: string]: any;
	} = $props();
</script>

<div
	class={cn('', className)}
	{...props}
>
	{@render children?.()}
</div>