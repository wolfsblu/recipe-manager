import { redirect } from "@sveltejs/kit";
import { isAuthenticated } from "$lib/auth/user.svelte";

export const load = async ({ parent }) => {
    await parent()
    if (!isAuthenticated()) {
        redirect(307, '/auth/login')
    }
}