import type { PageLoad } from './$types';
import { getShoppingLists } from "$lib/api/shopping/shopping.svelte";

export const load: PageLoad = async () => {
    const breadcrumbs = [
        { link: "/", name: "Home" },
        { link: "/shopping", name: "Shopping Lists" },
    ];

    try {
        const shoppingLists = await getShoppingLists();
        return {
            breadcrumbs,
            shoppingLists,
        };
    } catch (error) {
        return {
            breadcrumbs,
            shoppingLists: [],
            error: 'Failed to load shopping lists'
        };
    }
};