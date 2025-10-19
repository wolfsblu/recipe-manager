import type { PageLoad } from './$types';
import { getIngredients } from "$lib/api/ingredients/ingredients.svelte";

export const load: PageLoad = async () => {
    const breadcrumbs = [
        { link: "/", name: "Home" },
        { link: "/ingredients", name: "Ingredients" },
    ]

    // Fetch ingredients with pagination (getting a large batch for now)
    const response = await getIngredients({ limit: 100 });
    const ingredients = response?.data ?? [];

    return {
        breadcrumbs,
        ingredients,
    };
};