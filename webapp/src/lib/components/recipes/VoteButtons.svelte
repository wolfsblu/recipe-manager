<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
    import { Separator } from "$lib/components/ui/separator/index.js";
    import ChevronUpIcon from '@lucide/svelte/icons/chevron-up';
    import ChevronDownIcon from '@lucide/svelte/icons/chevron-down';
    import { addVote, removeVote } from "$lib/api/recipes/recipes.svelte";
    import { toast } from "svelte-sonner";
    import type { components } from "$lib/api/schema";
    import * as m from "$lib/paraglide/messages.js";

    type RecipeVotes = components["schemas"]["RecipeVotes"];

    interface Props {
        recipeId: number;
        votes?: RecipeVotes;
        size?: "sm" | "md" | "lg";
        orientation?: "vertical" | "horizontal";
    }

    let { recipeId, votes = { total: 0, user: 0 }, size = "md", orientation = "vertical" }: Props = $props();

    let isLoading = $state(false);
    let currentVotes = $state(votes);

    // Convert backend vote values to frontend state
    const getUserVoteState = (userVote: number): 'up' | 'down' | null => {
        if (userVote === 1) return 'up';
        if (userVote === -1) return 'down';
        return null;
    };

    // Convert frontend state to backend vote values
    const getVoteValue = (voteType: 'up' | 'down'): 1 | -1 => {
        return voteType === 'up' ? 1 : -1;
    };

    let userVote = $derived(getUserVoteState(currentVotes.user));

    const handleVote = async (voteType: 'up' | 'down') => {
        if (isLoading) return;

        isLoading = true;
        const previousVotes = { ...currentVotes };

        try {
            if (userVote === voteType) {
                // Remove vote if clicking the same button
                const response = await removeVote(recipeId);
                currentVotes = response;
                toast.success(m.recipes_voting_voteRemoved());
            } else {
                // Add or change vote
                const voteValue = getVoteValue(voteType);
                const response = await addVote(recipeId, voteValue);
                currentVotes = response;
                toast.success(voteType === 'up' ? m.recipes_voting_upvoted() : m.recipes_voting_downvoted());
            }
        } catch (error) {
            console.error('Voting error:', error);
            currentVotes = previousVotes;
            const errorMessage = userVote === voteType ? m.recipes_voting_removeVoteError() : m.recipes_voting_voteError();
            toast.error(errorMessage);
        } finally {
            isLoading = false;
        }
    };

    const iconSize = size === "sm" ? "w-4 h-4" : size === "lg" ? "w-6 h-6" : "w-5 h-5";
    const buttonSize = size === "sm" ? "sm" : size === "lg" ? "default" : "sm";
    const textSize = size === "sm" ? "text-sm" : size === "lg" ? "text-lg" : "text-base";
</script>

{#if orientation === "vertical"}
    <div class="flex flex-col items-center gap-2">
        <Button
            variant="ghost"
            size={buttonSize}
            onclick={() => handleVote('up')}
            disabled={isLoading}
            class={userVote === 'up' ? 'text-green-600 hover:text-green-700' : ''}
            title={m.recipes_voting_upvote()}
        >
            <ChevronUpIcon class={iconSize} />
        </Button>

        <span class="{textSize} font-bold text-foreground">
            {currentVotes.total}
        </span>

        <Button
            variant="ghost"
            size={buttonSize}
            onclick={() => handleVote('down')}
            disabled={isLoading}
            class={userVote === 'down' ? 'text-red-600 hover:text-red-700' : ''}
            title={m.recipes_voting_downvote()}
        >
            <ChevronDownIcon class={iconSize} />
        </Button>
    </div>
{:else}
    <div class="flex items-center gap-2">
        <Button
            variant="ghost"
            size={buttonSize}
            onclick={() => handleVote('up')}
            disabled={isLoading}
            class={userVote === 'up' ? 'text-green-600 hover:text-green-700' : ''}
            title={m.recipes_voting_upvote()}
        >
            <ChevronUpIcon class={iconSize} />
        </Button>

        <span class="{textSize} font-bold text-foreground min-w-8 text-center">
            {currentVotes.total}
        </span>

        <Button
            variant="ghost"
            size={buttonSize}
            onclick={() => handleVote('down')}
            disabled={isLoading}
            class={userVote === 'down' ? 'text-red-600 hover:text-red-700' : ''}
            title={m.recipes_voting_downvote()}
        >
            <ChevronDownIcon class={iconSize} />
        </Button>
    </div>
{/if}