import type { PageLoad } from './$types';
import {superValidate} from "sveltekit-superforms";
import {zod} from "sveltekit-superforms/adapters";
import {formSchema} from "./schema";

const breadcrumbs = [
    { link: "/", name: "Home" },
    { link: "/auth/register", name: "Register" },
]

export const load: PageLoad = async () => {
    return {
        breadcrumbs,
        form: await superValidate(zod(formSchema)),
    };
};