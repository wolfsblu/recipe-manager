/**
 * Infinite scroll utility for managing paginated data with Svelte 5 runes
 *
 * Usage:
 * ```ts
 * const scroll = useInfiniteScroll({
 *   fetchPage: async (cursor) => {
 *     const response = await getRecipes({ cursor });
 *     return {
 *       data: response.data,
 *       nextCursor: response.nextCursor,
 *       hasMore: response.hasMore
 *     };
 *   }
 * });
 *
 * // In template: bind to scroll.loadMoreTrigger
 * ```
 */

export interface PaginatedResponse<T> {
	data: T[];
	nextCursor: string | null | undefined;
	hasMore: boolean;
}

export interface UseInfiniteScrollOptions<T> {
	fetchPage: (cursor: string | null) => Promise<PaginatedResponse<T>>;
	initialData?: T[];
	onError?: (error: Error) => void;
}

export function useInfiniteScroll<T>(options: UseInfiniteScrollOptions<T>) {
	const { fetchPage, initialData = [], onError } = options;

	let items = $state<T[]>(initialData);
	let cursor = $state<string | null>(null);
	let hasMore = $state(true);
	let loading = $state(false);
	let error = $state<Error | null>(null);

	async function loadMore() {
		if (loading || !hasMore) return;

		loading = true;
		error = null;

		try {
			const response = await fetchPage(cursor);

			items = [...items, ...response.data];
			cursor = response.nextCursor ?? null;
			hasMore = response.hasMore;
		} catch (err) {
			error = err instanceof Error ? err : new Error('Failed to load more items');
			if (onError) {
				onError(error);
			}
		} finally {
			loading = false;
		}
	}

	async function refresh() {
		items = [];
		cursor = null;
		hasMore = true;
		await loadMore();
	}

	return {
		get items() {
			return items;
		},
		get loading() {
			return loading;
		},
		get hasMore() {
			return hasMore;
		},
		get error() {
			return error;
		},
		loadMore,
		refresh
	};
}
