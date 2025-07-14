import type { PageLoad } from './$types';

export const prerender = false

export const load: PageLoad = async () => {
    const breadcrumbs = [
        { link: "/", name: "Home" },
        { link: "/recipes", name: "Recipes" },
        { link: "", name: "My super duper tasty recipe" },
    ]

    return {
        breadcrumbs,
    };
};