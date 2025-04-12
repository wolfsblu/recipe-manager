import { redirect } from "@sveltejs/kit";
import { isAuthenticated } from "$lib/auth/user.svelte";

export const load = ({ params }) => {
    if (!isAuthenticated()) {
        redirect(307, '/auth/login')
    }
}