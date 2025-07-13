import { getProfile } from "$lib/auth/user.svelte";
import type { LayoutLoad } from './$types';

export const prerender = true
export const ssr = false

const breadcrumbs = [
    { link: "/", name: "Home" },
]

export const load: LayoutLoad = async () => {
    try {
        await getProfile()
    } catch (e) {
        return { breadcrumbs }
    }
    return { breadcrumbs }
}