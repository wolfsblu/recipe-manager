import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
    const breadcrumbs = [
        { link: "/", name: "Home" },
        { link: "/recipes", name: "Recipes" },
    ]

    return {
        breadcrumbs,
    };
};