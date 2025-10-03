<script lang="ts">
    import ClockIcon from '@lucide/svelte/icons/clock'
    import ImageIcon from '@lucide/svelte/icons/image'
    import { Badge } from "$lib/components/ui/badge/index.js";
    import { ScrollArea } from "$lib/components/ui/scroll-area/index.js";
    import { Separator } from "$lib/components/ui/separator/index.js";
    import { patchRecipe } from "$lib/api/recipes/recipes.svelte.js";
    import { toast } from "svelte-sonner";
    import TagCombobox from './TagCombobox.svelte';
    import * as m from "$lib/paraglide/messages.js";

    let { recipe, availableTags, overlay } = $props()

    // Reactive state for tags
    let tags = $state(recipe.tags || [])

    // Calculate total number of ingredients
    const totalIngredients = $derived(
        recipe.steps?.reduce((total, step) =>
            total + (step.ingredients?.length || 0), 0
        ) || 0
    )

    const formatMinutesAsHours = (minutes: number) => {
        if (!minutes || minutes < 0) {
            return m.common_time_notAvailable();
        }

        if (minutes < 60) {
            return m.common_time_minutes({ count: minutes });
        }

        const hours = minutes / 60;

        // Round to a maximum of 2 decimal places.
        // We multiply by 100, round the result, then divide by 100.
        // This avoids trailing zeros for numbers like 1.5 or 2.0.
        const roundedHours = Math.round(hours * 100) / 100;

        return m.common_time_hours({ count: roundedHours });
    }
    
    const handleTagAdded = async (tagId: number) => {
        const newTag = availableTags.find(tag => tag.id === tagId);
        if (!newTag || tags.some(tag => tag.id === tagId)) {
            return;
        }

        const previousTags = [...tags];
        tags = [...tags, newTag];

        try {
            const updatedTagIds = tags.map(tag => tag.id);
            await patchRecipe(recipe.id, { tags: updatedTagIds });
            toast.success(m.recipes_tags_addSuccess({ tagName: newTag.name }));
        } catch (error) {
            tags = previousTags;
            toast.error(m.recipes_tags_addError());
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
                <div class="w-full h-52 bg-muted flex items-center justify-center">
                    <div class="text-center text-muted-foreground">
                        <ImageIcon class="w-12 h-12 mx-auto mb-2" />
                        <p class="text-sm">{m.recipes_card_noImage()}</p>
                    </div>
                </div>
            {/if}
        </a>

        {#if overlay}
            <div class="absolute top-2 right-2 z-10 flex flex-col items-center gap-1 bg-secondary rounded-full px-1 py-2 opacity-90">
                {@render overlay?.()}
            </div>
        {/if}
    </div>
    <div class="pt-3 px-4 pb-2 space-y-2">
        <a href="/recipes/{recipe.id}" class="inline-block font-semibold text-base">
            {recipe.name}
        </a>
        <div class="text-sm text-muted-foreground grid space-y-1">
            <div class="flex justify-between">
                <span>{m.recipes_card_servings()}</span>
                <Badge>{recipe.servings}</Badge>
            </div>
            <Separator />
            <div class="flex justify-between">
                <span>{m.recipes_card_ingredients()}</span>
                <Badge>{totalIngredients}</Badge>
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
                search={m.recipes_card_searchTags()}
                empty={m.recipes_card_noTagsFound()}
                onSelect={handleTagAdded}
            />
        </div>
    </ScrollArea>
</div>

