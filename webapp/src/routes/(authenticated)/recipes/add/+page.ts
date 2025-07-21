import type { PageLoad } from './$types';
import {superValidate} from "sveltekit-superforms";
import {zod} from "sveltekit-superforms/adapters";
import {formSchema} from "./schema";

export const load: PageLoad = async () => {
    const breadcrumbs = [
        { link: "/", name: "Home" },
        { link: "/recipes", name: "Recipes" },
        { link: "/recipes/add", name: "New Recipe" },
    ]

    return {
        breadcrumbs,
        form: await superValidate(zod(formSchema)),
    };
};