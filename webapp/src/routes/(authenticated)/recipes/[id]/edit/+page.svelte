<script lang="ts">
    import type {PageProps} from './$types';

    let {data}: PageProps = $props();
</script>

<p>Recipe edit page for {data.recipe.name}</p>