<script lang="ts">
    import * as m from "$lib/paraglide/messages"
    import { PUBLIC_API_URL } from '$env/static/public'
</script>

<h1>{m.example_message({username: "Frank"})}</h1>

<p>
    Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the documentation
</p>
<p>
    Go to <a href="about">about page</a>
</p>

<p>
    {PUBLIC_API_URL}
</p>