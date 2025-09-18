import type { PageLoad } from './$types';
import { getIngredients } from "$lib/api/ingredients/ingredients.svelte";

export const load: PageLoad = async () => {
    const breadcrumbs = [
        { link: "/", name: "Home" },
        { link: "/ingredients", name: "Ingredients" },
    ]

    const ingredients = await getIngredients()
    return {
        breadcrumbs,
        ingredients,
    };
};