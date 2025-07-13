import type { PageLoad } from './$types';
import {superValidate} from "sveltekit-superforms";
import {zod} from "sveltekit-superforms/adapters";
import {formSchema} from "./schema";

export const load: PageLoad = async () => {
    return {
        form: await superValidate(zod(formSchema)),
    };
};