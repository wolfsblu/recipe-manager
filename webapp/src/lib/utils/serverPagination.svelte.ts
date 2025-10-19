type FetchPageFn<T> = (cursor: string | null) => Promise<{
    data: T[];
    nextCursor?: string | null;
    hasMore: boolean;
}>;

type ServerPaginationOptions<T> = {
    fetchPage: FetchPageFn<T>;
    pageSize?: number;
};

export function useServerPagination<T>(options: ServerPaginationOptions<T>) {
    const { fetchPage, pageSize = 30 } = options;

    let allItems = $state<T[]>([]);
    let pages = $state<Array<{ data: T[]; cursor: string | null; nextCursor: string | null }>>([]);
    let currentPageIndex = $state(0);
    let loading = $state(false);
    let hasMore = $state(true);

    async function loadPage(cursor: string | null) {
        if (loading) return;

        loading = true;
        try {
            const response = await fetchPage(cursor);
            const newPage = {
                data: response.data,
                cursor,
                nextCursor: response.nextCursor ?? null
            };

            pages = [...pages, newPage];
            allItems = pages.flatMap(p => p.data);
            hasMore = response.hasMore;

            return newPage;
        } finally {
            loading = false;
        }
    }

    async function loadInitialPage() {
        if (pages.length === 0) {
            await loadPage(null);
        }
    }

    async function nextPage() {
        if (!canGoNext() || loading) return;

        const nextPageIndex = currentPageIndex + 1;

        // If the next page is already loaded, just navigate to it
        if (nextPageIndex < pages.length) {
            currentPageIndex = nextPageIndex;
            return;
        }

        // Otherwise, fetch the next page
        const currentPage = pages[currentPageIndex];
        if (currentPage?.nextCursor) {
            const newPage = await loadPage(currentPage.nextCursor);
            if (newPage) {
                currentPageIndex = nextPageIndex;
            }
        }
    }

    function previousPage() {
        if (!canGoPrevious()) return;
        currentPageIndex--;
    }

    function goToFirstPage() {
        currentPageIndex = 0;
    }

    function canGoPrevious() {
        return currentPageIndex > 0;
    }

    function canGoNext() {
        if (loading) return false;

        // Can go next if there's a next page already loaded
        if (currentPageIndex < pages.length - 1) return true;

        // Or if there's more data to fetch
        const currentPage = pages[currentPageIndex];
        return hasMore && currentPage?.nextCursor != null;
    }

    function getCurrentPageData() {
        return pages[currentPageIndex]?.data ?? [];
    }

    function getTotalLoadedItems() {
        return allItems.length;
    }

    function getCurrentPageNumber() {
        return currentPageIndex + 1;
    }

    function getTotalLoadedPages() {
        return pages.length;
    }

    return {
        get items() { return allItems; },
        get currentPageData() { return getCurrentPageData(); },
        get loading() { return loading; },
        get hasMore() { return hasMore; },
        get totalLoadedItems() { return getTotalLoadedItems(); },
        get currentPageNumber() { return getCurrentPageNumber(); },
        get totalLoadedPages() { return getTotalLoadedPages(); },
        get canGoNext() { return canGoNext(); },
        get canGoPrevious() { return canGoPrevious(); },
        loadInitialPage,
        nextPage,
        previousPage,
        goToFirstPage
    };
}
