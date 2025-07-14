import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
    const breadcrumbs = [
        { link: "/", name: "Home" },
        { link: "/recipes", name: "Recipes" },
        { link: `/recipes/${params.id}`, name: "My super duper tasty recipe" },
    ]

    return {
        breadcrumbs,
    };
};