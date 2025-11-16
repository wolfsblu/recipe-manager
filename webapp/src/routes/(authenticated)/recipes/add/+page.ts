import type { PageLoad } from './$types';
import {superValidate} from "sveltekit-superforms";
import {zod4} from "sveltekit-superforms/adapters";
import {formSchema} from "./schema";
import {getIngredients, getUnits, getTags} from "$lib/api/recipes/recipes.svelte";

export const load: PageLoad = async () => {
    const breadcrumbs = [
        { link: "/", name: "Home" },
        { link: "/recipes", name: "Recipes" },
        { link: "/recipes/add", name: "New Recipe" },
    ]

    const [ingredientsResponse, unitsResponse, tagsResponse] = await Promise.all([
        getIngredients({ limit: 100 }),
        getUnits({ limit: 100 }),
        getTags({ limit: 100 })
    ])

    return {
        breadcrumbs,
        ingredients: ingredientsResponse.data,
        units: unitsResponse.data,
        tags: tagsResponse.data,
        form: await superValidate(zod4(formSchema)),
    };
};