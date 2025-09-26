<script lang="ts">
    import fruits from './fruits.jpg'
    import ClockIcon from '@lucide/svelte/icons/clock'
    import ChevronUp from '@lucide/svelte/icons/chevron-up'
    import ChevronDown from '@lucide/svelte/icons/chevron-down'
    import { Badge } from "$lib/components/ui/badge/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import { ScrollArea } from "$lib/components/ui/scroll-area/index.js";
    import { Separator } from "$lib/components/ui/separator/index.js";
    import { addVote, removeVote } from "$lib/api/recipes/recipes.svelte.js";
    import { toast } from "svelte-sonner";
    import { scale } from "svelte/transition";
    import { quintOut } from "svelte/easing";
    import TagCombobox from './TagCombobox.svelte';

    let { recipe, availableTags } = $props()

    // Reactive state for voting
    let votes = $state(recipe.votes)
    let isVoting = $state(false)
    
    // Reactive state for tags
    let tags = $state(recipe.tags || [])

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

    const formatMinutesAsHours = (minutes: number) => {
        if (!minutes || minutes < 0) {
            return "N/A";
        }

        if (minutes < 60) {
            return `${minutes}m`;
        }

        const hours = minutes / 60;

        // Round to a maximum of 2 decimal places.
        // We multiply by 100, round the result, then divide by 100.
        // This avoids trailing zeros for numbers like 1.5 or 2.0.
        const roundedHours = Math.round(hours * 100) / 100;

        return `${roundedHours}h`;
    }
    
    const handleTagAdded = (tagId: number) => {
        const newTag = availableTags.find(tag => tag.id === tagId);
        if (newTag && !tags.some(tag => tag.id === tagId)) {
            tags = [...tags, newTag];
        }
    }
</script>

<div class="bg-card rounded-lg shadow-sm overflow-hidden transition-all duration-300">
    <div class="relative">
        <a href="/recipes/{recipe.id}" class="block">
            {#if recipe.minutes}
                <Badge variant="secondary" class="absolute top-2 left-2 z-10 opacity-90 rounded-full">
                    <ClockIcon /> {formatMinutesAsHours(recipe.minutes)}
                </Badge>
            {/if}
            {#if recipe.images && recipe.images.length > 0}
                <img loading="lazy" src={recipe.images[0]} alt="Recipe" class="w-full h-52 object-cover" />
            {:else}
                <img loading="lazy" src={fruits} alt="Fruits" class="w-full h-52 object-cover" />
            {/if}
        </a>

        <!-- Voting overlay -->
        <div class="absolute top-2 right-2 z-10 flex flex-col items-center gap-1 bg-secondary rounded-full px-1 py-2 opacity-90">
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
                <ChevronUp class={`h-3 w-3 ${votes.user === 1 ? 'text-green-600' : 'text-muted-foreground'}`} />
            </Button>
            {#if votes.total !== 0}
                <span
                    class="text-sm font-bold text-muted-foreground min-w-4 text-center"
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
                <ChevronDown class={`h-3 w-3 ${votes.user === -1 ? 'text-red-600' : 'text-muted-foreground'}`} />
            </Button>
        </div>
    </div>
    <div class="pt-3 px-4 pb-2 space-y-2">
        <a href="/recipes/{recipe.id}" class="inline-block font-semibold text-base">
            {recipe.name}
        </a>
        <div class="text-sm text-muted-foreground grid space-y-1">
            <div class="flex justify-between">
                <span>Servings</span>
                <Badge>{recipe.servings}</Badge>
            </div>
            <Separator />
            <div class="flex justify-between">
                <span>Ingredients</span>
                <Badge>12</Badge>
            </div>
                <Separator />
        </div>
    </div>
    <ScrollArea class="whitespace-nowrap pb-4 px-4" orientation="horizontal">
        <div class="flex w-max gap-x-1">
            {#each tags as tag}
                <Badge variant="secondary" class="text-muted-foreground rounded-full">
                    {tag.name}
                </Badge>
            {/each}
            <TagCombobox
                options={availableTags
                    .filter(tag => !tags.some(t => t.id === tag.id))
                    .map(tag => ({ value: tag.id, label: tag.name }))}
                search="Search tags..."
                empty="No tags found"
                onSelect={handleTagAdded}
            />
        </div>
    </ScrollArea>
</div>

