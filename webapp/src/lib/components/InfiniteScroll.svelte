<script lang="ts">
	import { onMount } from 'svelte';

	interface Props {
		hasMore: boolean;
		loading: boolean;
		onLoadMore: () => void | Promise<void>;
		threshold?: number;
	}

	let { hasMore, loading, onLoadMore, threshold = 0.5 }: Props = $props();

	let sentinel: HTMLDivElement;
	let observer: IntersectionObserver | null = null;

	onMount(() => {
		observer = new IntersectionObserver(
			(entries) => {
				const entry = entries[0];
				if (entry.isIntersecting && hasMore && !loading) {
					onLoadMore();
				}
			},
			{
				threshold,
				rootMargin: '100px'
			}
		);

		if (sentinel) {
			observer.observe(sentinel);
		}

		return () => {
			if (observer) {
				observer.disconnect();
			}
		};
	});
</script>

<div bind:this={sentinel} class="sentinel">
	{#if loading}
		<div class="loading-indicator">
			<div class="spinner"></div>
			<p>Loading more...</p>
		</div>
	{:else if !hasMore}
		<div class="end-message">
			<p>No more items</p>
		</div>
	{/if}
</div>

<style>
	.sentinel {
		min-height: 1px;
		padding: 2rem 0;
		display: flex;
		justify-content: center;
		align-items: center;
	}

	.loading-indicator,
	.end-message {
		text-align: center;
		color: var(--color-text-secondary, #666);
	}

	.loading-indicator {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.5rem;
	}

	.spinner {
		width: 2rem;
		height: 2rem;
		border: 3px solid var(--color-border, #e0e0e0);
		border-top-color: var(--color-primary, #3b82f6);
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.end-message p {
		font-size: 0.875rem;
		opacity: 0.7;
	}
</style>
