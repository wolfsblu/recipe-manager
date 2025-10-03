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

    const [mealPlan, tags, shoppingLists] = await Promise.all([
        getMealPlan(from, until),
        getTags(),
        getShoppingLists().catch(() => [])
    ]);

    return {
        breadcrumbs,
        mealPlan,
        tags,
        shoppingLists,
        from,
        until
    };
};