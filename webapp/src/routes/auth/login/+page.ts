import type { PageLoad } from './$types';
import {superValidate} from "sveltekit-superforms";
import {zod4} from "sveltekit-superforms/adapters";
import {formSchema} from "./schema";

const breadcrumbs = [
    { link: "/", name: "Home" },
    { link: "/auth/login", name: "Login" },
]

export const load: PageLoad = async () => {
    return {
        breadcrumbs,
        form: await superValidate(zod4(formSchema)),
    };
};