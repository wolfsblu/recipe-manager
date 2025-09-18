import { redirect } from "@sveltejs/kit";
import { isAuthenticated } from "$lib/api/auth/user.svelte";
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async () => {
    if (isAuthenticated()) {
        redirect(307, '/');
    }
}