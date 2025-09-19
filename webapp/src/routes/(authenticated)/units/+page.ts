import type { PageLoad } from './$types';
import { getUnits } from "$lib/api/units/units.svelte";

export const load: PageLoad = async () => {
    const breadcrumbs = [
        { link: "/", name: "Home" },
        { link: "/units", name: "Units" },
    ]

    const units = await getUnits()
    return {
        breadcrumbs,
        units,
    };
};