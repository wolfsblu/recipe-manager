<script lang="ts">
	import { cn } from '$lib/utils/utils';
	import TagsInputTag from './tags-input-tag.svelte';
	import { untrack } from 'svelte';
	import { useDebounce } from 'runed';

	type Suggestion = { id: number; label: string };

	// Define defaultValidate before using it
	const defaultValidate = (val: Suggestion, tags: number[]) => {
		// disallow duplicates
		if (tags.find((t) => t === val.id)) return undefined;
		return val.id;
	};

	let {
		value = $bindable([]),
		placeholder,
		class: className,
		disabled = false,
		validate = defaultValidate,
		fetchSuggestions,
		minSearchLength = 2,
		debounceMs = 300,
		id,
		...rest
	}: {
		value?: number[];
		placeholder?: string;
		class?: string;
		disabled?: boolean;
		validate?: (val: Suggestion, tags: number[]) => number | undefined;
		fetchSuggestions: (searchQuery: string) => Promise<Suggestion[]>;
		minSearchLength?: number;
		debounceMs?: number;
		id?: string;
	} = $props();

	let inputValue = $state('');
	let tagIndex = $state<number>();
	let invalid = $state(false);
	let isComposing = $state(false);
	let showDropdown = $state(false);
	let selectedSuggestionIndex = $state<number>(-1);
	let inputElement: HTMLInputElement;
	let containerElement: HTMLDivElement;
	let loading = $state(false);
	let suggestions = $state<Suggestion[]>([]);

	// Debounced search function
	const debouncedSearch = useDebounce(async (query: string) => {
		if (query.length < minSearchLength) {
			suggestions = [];
			return;
		}

		loading = true;
		try {
			const results = await fetchSuggestions(query);
			// Filter out already selected tags
			suggestions = results.filter(s => !value.includes(s.id));
		} catch (error) {
			console.error('Failed to fetch suggestions:', error);
			suggestions = [];
		} finally {
			loading = false;
		}
	}, debounceMs);

	$effect(() => {
		// whenever input value changes reset invalid and update dropdown
		// eslint-disable-next-line @typescript-eslint/no-unused-expressions
		inputValue;

		untrack(() => {
			invalid = false;
			if (inputValue.length >= minSearchLength) {
				debouncedSearch(inputValue);
				showDropdown = true;
			} else {
				showDropdown = false;
				suggestions = [];
			}
			selectedSuggestionIndex = -1;
		});
	});

	$effect(() => {
		// Scroll highlighted item into view
		if (selectedSuggestionIndex >= 0) {
			const highlightedElement = document.getElementById(`suggestion-${selectedSuggestionIndex}`);
			if (highlightedElement) {
				highlightedElement.scrollIntoView({
					behavior: 'smooth',
					block: 'nearest'
				});
			}
		}
	});

	// Get selected tag labels by fetching all tags (or maintain a cache)
	let selectedTagsLabels = $state<Map<number, string>>(new Map());

	// Load selected tag labels when value changes
	$effect(() => {
		// Track value reactively
		const tagIds = value;

		untrack(async () => {
			// Find any tags we don't have labels for
			const missingIds = tagIds.filter(id => !selectedTagsLabels.has(id));
			if (missingIds.length > 0) {
				try {
					// Fetch all to find missing labels
					const allTags = await fetchSuggestions('');
					allTags.forEach(tag => {
						selectedTagsLabels.set(tag.id, tag.label);
					});
					// Trigger reactivity
					selectedTagsLabels = new Map(selectedTagsLabels);
				} catch (error) {
					console.error('Failed to load tag labels:', error);
				}
			}
		});
	});

	const selectedTags = $derived.by(() => {
		return value.map(tagId => ({
			id: tagId,
			label: selectedTagsLabels.get(tagId) || `Tag ${tagId}`
		}));
	});

	const enter = () => {
		if (isComposing) return;

		// Try to find a matching suggestion
		const matchingSuggestion = suggestions?.find(s =>
			s.label.toLowerCase() === inputValue.toLowerCase().trim()
		);

		if (matchingSuggestion) {
			selectSuggestion(matchingSuggestion);
		} else {
			invalid = true;
		}
	};

	const compositionStart = () => {
		isComposing = true;
	};

	const compositionEnd = () => {
		isComposing = false;
	};

	const selectSuggestion = (suggestion: Suggestion) => {
		const validated = validate(suggestion, value);
		if (validated !== undefined) {
			value = [...value, validated];
			selectedTagsLabels.set(suggestion.id, suggestion.label);
			selectedTagsLabels = new Map(selectedTagsLabels);
			inputValue = '';
			showDropdown = false;
			selectedSuggestionIndex = -1;
		}
	};

	const keydown = (e: KeyboardEvent) => {
		const target = e.target as HTMLInputElement;

		if (e.key === 'Enter') {
			// prevent form submit
			e.preventDefault();

			if (isComposing) return;

			// If dropdown is open and suggestion is selected
			if (showDropdown && selectedSuggestionIndex >= 0 && selectedSuggestionIndex < suggestions.length) {
				selectSuggestion(suggestions[selectedSuggestionIndex]);
				return;
			}

			enter();
			return;
		}

		// Handle arrow keys for suggestion navigation
		if (showDropdown && suggestions.length > 0) {
			if (e.key === 'ArrowDown') {
				e.preventDefault();
				selectedSuggestionIndex = selectedSuggestionIndex >= suggestions.length - 1
					? -1
					: selectedSuggestionIndex + 1;
				return;
			}
			if (e.key === 'ArrowUp') {
				e.preventDefault();
				selectedSuggestionIndex = selectedSuggestionIndex <= -1
					? suggestions.length - 1
					: selectedSuggestionIndex - 1;
				return;
			}
			if (e.key === 'Escape') {
				e.preventDefault();
				showDropdown = false;
				selectedSuggestionIndex = -1;
				return;
			}
		}

		const isAtBeginning = target.selectionStart === 0 && target.selectionEnd === 0;

		let shouldResetIndex = true;

		if (e.key === 'Backspace') {
			if (isAtBeginning) {
				e.preventDefault();

				if (tagIndex !== undefined) {
					deleteIndex(tagIndex);

					// focus previous
					const prev = tagIndex - 1;

					if (prev < 0) {
						tagIndex = undefined;
					} else {
						tagIndex = prev;
					}
				} else {
					tagIndex = selectedTags.length - 1;
				}

				shouldResetIndex = false;
			}
		}

		if (e.key === 'Delete') {
			if (isAtBeginning) {
				if (inputValue.length === 0) {
					if (tagIndex !== undefined) {
						e.preventDefault();

						deleteIndex(tagIndex);

						// stay focused on the same index unless selectedTags.length === 0
						if (selectedTags.length === 0) tagIndex = undefined;

						shouldResetIndex = false;
					}
				}
			}
		}

		// controls for tag selection
		if (isAtBeginning) {
			// left
			if (e.key === 'ArrowLeft') {
				if (tagIndex !== undefined) {
					const prev = tagIndex - 1;

					if (prev < 0) {
						tagIndex = 0;
					} else {
						tagIndex = prev;
					}
				} else {
					// set initial index
					tagIndex = selectedTags.length - 1;
				}

				shouldResetIndex = false;
			}

			// right
			// we can only move right if the value is empty
			if (inputValue.length === 0) {
				if (e.key === 'ArrowRight') {
					if (tagIndex !== undefined) {
						const next = tagIndex + 1;

						if (next > selectedTags.length - 1) {
							tagIndex = undefined;
						} else {
							tagIndex = next;
						}

						shouldResetIndex = false;
					}
				}
			}
		}

		// reset the tag index to undefined
		if (shouldResetIndex) {
			tagIndex = undefined;
		}
	};

	const deleteValue = (val: number) => {
		const index = value.findIndex((v) => val === v);

		if (index === -1) return;

		deleteIndex(index);
	};

	const deleteIndex = (index: number) => {
		value = [...value.slice(0, index), ...value.slice(index + 1)];
	};

	const blur = (e: FocusEvent) => {
		// Check if focus is moving to dropdown
		const relatedTarget = e.relatedTarget as HTMLElement;
		if (relatedTarget && containerElement?.contains(relatedTarget)) {
			return;
		}
		tagIndex = undefined;
		showDropdown = false;
		selectedSuggestionIndex = -1;
	};

</script>

<div bind:this={containerElement} class="relative w-full">
<div
	class={cn(
		'border-input bg-background selection:bg-primary dark:bg-input/30 flex min-h-[36px] w-full flex-wrap place-items-center gap-1 rounded-md border py-0.5 pr-1 pl-1 disabled:opacity-50 aria-disabled:cursor-not-allowed',
		className
	)}
	aria-disabled={disabled}
>
	{#each selectedTags as tag, i (tag.id)}
		<TagsInputTag value={tag.id} label={tag.label} {disabled} onDelete={deleteValue} active={i === tagIndex} />
	{/each}
	<input
		{...rest}
		bind:this={inputElement}
		bind:value={inputValue}
		onblur={blur}
		oncompositionstart={compositionStart}
		oncompositionend={compositionEnd}
		{disabled}
		{placeholder}
		data-invalid={invalid}
		onkeydown={keydown}
		class="placeholder:text-muted-foreground min-w-16 shrink grow basis-0 border-none bg-transparent px-2 outline-hidden focus:outline-hidden disabled:cursor-not-allowed data-[invalid=true]:text-red-500 md:text-sm"
	/>
</div>

{#if showDropdown}
	<div class="absolute top-full left-0 right-0 z-50 mt-1">
		<div class="bg-popover text-popover-foreground overflow-hidden rounded-md border shadow-md">
			{#if loading}
				<div class="p-2 text-center text-sm text-muted-foreground">
					Loading...
				</div>
			{:else if inputValue.length > 0 && inputValue.length < minSearchLength}
				<div class="p-2 text-center text-sm text-muted-foreground">
					Type at least {minSearchLength} characters
				</div>
			{:else if suggestions.length === 0}
				<div class="p-2 text-center text-sm text-muted-foreground">
					No results found
				</div>
			{:else}
				<div class="p-1 max-h-[200px] overflow-y-auto">
					{#each suggestions as suggestion, index}
						<div
							id="suggestion-{index}"
							role="option"
							aria-selected={index === selectedSuggestionIndex}
							tabindex="-1"
							class={cn(
								'relative flex cursor-pointer select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground',
								index === selectedSuggestionIndex && 'bg-accent text-accent-foreground'
							)}
							onmousedown={(e) => {
								e.preventDefault();
								selectSuggestion(suggestion);
							}}
							onmouseenter={() => selectedSuggestionIndex = index}
						>
							{suggestion.label}
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>
{/if}

</div>
