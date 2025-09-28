<script lang="ts">
    import ChevronUp from '@lucide/svelte/icons/chevron-up'
    import ChevronDown from '@lucide/svelte/icons/chevron-down'
    import { Button } from "$lib/components/ui/button/index.js";
    import { addVote, removeVote } from "$lib/api/recipes/recipes.svelte.js";
    import { toast } from "svelte-sonner";
    import { scale } from "svelte/transition";
    import { quintOut } from "svelte/easing";

    let { recipe } = $props()

    let votes = $state(recipe.votes)
    let isVoting = $state(false)

    const handleVote = async (voteValue: 1 | -1) => {
        if (isVoting) return

        isVoting = true
        const previousVotes = { ...votes }

        try {
            const data = await addVote(recipe.id, voteValue)
            votes = data
            toast.success(`${voteValue === 1 ? 'Upvoted' : 'Downvoted'} recipe`)
        } catch (error) {
            votes = previousVotes
            toast.error("Failed to vote. Please try again.")
        } finally {
            isVoting = false
        }
    }

    const handleRemoveVote = async () => {
        if (isVoting) return

        isVoting = true
        const previousVotes = { ...votes }

        try {
            const data = await removeVote(recipe.id)
            votes = data
            toast.success("Vote removed")
        } catch (error) {
            votes = previousVotes
            toast.error("Failed to remove vote. Please try again.")
        } finally {
            isVoting = false
        }
    }
</script>

<Button
    variant="ghost"
    size="sm"
    class="h-6 w-6 p-0 hover:bg-transparent"
    disabled={isVoting}
    onclick={(e) => {
        e.preventDefault();
        e.stopPropagation();
        votes.user === 1 ? handleRemoveVote() : handleVote(1);
    }}
>
    <ChevronUp class={`${votes.user === 1 ? 'text-green-600' : ''}`} />
</Button>
{#if votes.total !== 0}
    <span
        class="text-sm font-bold min-w-4 text-center"
        transition:scale={{ duration: 300, start: 0.5, easing: quintOut }}
    >
        {votes.total}
    </span>
{/if}
<Button
    variant="ghost"
    size="sm"
    class="h-6 w-6 p-0 hover:bg-transparent"
    disabled={isVoting}
    onclick={(e) => {
        e.preventDefault();
        e.stopPropagation();
        votes.user === -1 ? handleRemoveVote() : handleVote(-1);
    }}
>
    <ChevronDown class={`${votes.user === -1 ? 'text-red-600' : ''}`} />
</Button>
