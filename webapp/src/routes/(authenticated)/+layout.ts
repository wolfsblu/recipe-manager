import { redirect } from "@sveltejs/kit";
import { isAuthenticated } from "$lib/api/auth/user.svelte";

export const load = async ({ data, parent }) => {
    await parent()
    if (!isAuthenticated()) {
        redirect(307, '/auth/login')
    }
    return data
}