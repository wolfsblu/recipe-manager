import type { PageLoad } from './$types';
import { getUnits } from "$lib/api/units/units.svelte";

export const load: PageLoad = async () => {
    const breadcrumbs = [
        { link: "/", name: "Home" },
        { link: "/units", name: "Units" },
    ]

    const response = await getUnits({ limit: 100 });
    const units = response?.data ?? [];

    return {
        breadcrumbs,
        units,
    };
};