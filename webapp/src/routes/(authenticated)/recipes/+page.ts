import type { PageLoad } from './$types';
import { getTags } from "$lib/api/recipes/recipes.svelte";

export const load: PageLoad = async () => {
    const breadcrumbs = [
        { link: "/", name: "Home" },
        { link: "/recipes", name: "Recipes" },
    ]

    // Fetch all tags (small dataset, not paginated for now)
    // We'll fetch all pages to get complete tag list
    const tagsResponse = await getTags({ limit: 100 });
    const tags = tagsResponse?.data ?? [];

    return {
        breadcrumbs,
        tags,
    };
};