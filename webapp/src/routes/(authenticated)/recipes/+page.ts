import type { PageLoad } from './$types';
import { getRecipes, getTags } from "$lib/api/recipes/recipes.svelte";

export const load: PageLoad = async () => {
    const breadcrumbs = [
        { link: "/", name: "Home" },
        { link: "/recipes", name: "Recipes" },
    ]

    const [recipes, tags] = await Promise.all([
        getRecipes(),
        getTags()
    ]);

    return {
        breadcrumbs,
        recipes,
        tags,
    };
};