import type { PageLoad } from './$types';
import {superValidate} from "sveltekit-superforms";
import {zod} from "sveltekit-superforms/adapters";
import {formSchema} from "./schema";
import {getIngredients, getUnits, getTags} from "$lib/api/recipes/recipes.svelte";

export const load: PageLoad = async () => {
    const breadcrumbs = [
        { link: "/", name: "Home" },
        { link: "/recipes", name: "Recipes" },
        { link: "/recipes/add", name: "New Recipe" },
    ]

    const ingredients = await getIngredients()
    const units = await getUnits()
    const tags = await getTags()

    return {
        breadcrumbs,
        ingredients,
        units,
        tags,
        form: await superValidate(zod(formSchema)),
    };
};