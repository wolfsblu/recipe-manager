import type { PageLoad } from './$types';
import { getMealPlan } from "$lib/api/mealplan/mealplan.svelte";
import { getTags } from "$lib/api/recipes/recipes.svelte";
import { getShoppingLists } from "$lib/api/shopping/shopping.svelte";

export const load: PageLoad = async ({ url }) => {
    const breadcrumbs = [
        { link: "/", name: "Home" },
        { link: "/mealplan", name: "Meal Plan" },
    ];

    // Get date range from query parameters or default to current week
    const today = new Date();
    const nextWeek = new Date(today);
    nextWeek.setDate(today.getDate() + 7);

    const from = url.searchParams.get('from') || today.toISOString().split('T')[0];
    const until = url.searchParams.get('until') || nextWeek.toISOString().split('T')[0];

    const [mealPlanResponse, tagsResponse, shoppingListsResponse] = await Promise.all([
        getMealPlan({ from, until, limit: 100 }),
        getTags({ limit: 100 }),
        getShoppingLists({ limit: 100 }).catch(() => ({ data: [], hasMore: false, nextCursor: null }))
    ]);

    return {
        breadcrumbs,
        mealPlan: mealPlanResponse?.data ?? [],
        tags: tagsResponse?.data ?? [],
        shoppingLists: shoppingListsResponse?.data ?? [],
        from,
        until
    };
};