import type { PageLoad } from './$types';
import {superValidate} from "sveltekit-superforms";
import {zod4} from "sveltekit-superforms/adapters";
import {formSchema} from "./schema";
import {getLocale} from "$lib/paraglide/runtime";

const breadcrumbs = [
    { link: "/", name: "Home" },
    { link: "/auth/register", name: "Register" },
]

export const load: PageLoad = async () => {
    return {
        breadcrumbs,
        form: await superValidate({ locale: getLocale() }, zod4(formSchema)),
    };
};