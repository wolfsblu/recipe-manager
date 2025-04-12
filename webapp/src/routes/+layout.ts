import { getProfile } from "$lib/auth/user.svelte";

export const prerender = true
export const ssr = false

export const load = async () => {
    try {
        await getProfile()
    } catch (e) {
        return null
    }
}