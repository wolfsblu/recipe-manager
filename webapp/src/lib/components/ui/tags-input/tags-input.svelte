<!--
	Installed from @ieedan/shadcn-svelte-extras
-->

<script lang="ts">
	import { cn } from '$lib/utils/utils';
	import type { TagsInputProps } from './types';
	import TagsInputTag from './tags-input-tag.svelte';
	import { untrack } from 'svelte';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import * as Command from '$lib/components/ui/command/index.js';

	const defaultValidate: TagsInputProps['validate'] = (val, tags) => {
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
		suggestions = [],
		...rest
	}: TagsInputProps = $props();

	let inputValue = $state('');
	let tagIndex = $state<number>();
	let invalid = $state(false);
	let isComposing = $state(false);
	let showDropdown = $state(false);
	let selectedSuggestionIndex = $state<number>(-1);
	let inputElement: HTMLInputElement;
	let containerElement: HTMLDivElement;

	$effect(() => {
		// whenever input value changes reset invalid and update dropdown
		// eslint-disable-next-line @typescript-eslint/no-unused-expressions
		inputValue;

		untrack(() => {
			invalid = false;
			showDropdown = inputValue.length >= 2 && filteredSuggestions.length > 0;
			selectedSuggestionIndex = -1;
		});
	});

	const filteredSuggestions = $derived.by(() => {
		if (!inputValue || inputValue.length < 2 || !suggestions) return [];

		const query = inputValue.toLowerCase();
		return suggestions.filter(suggestion =>
			suggestion.label.toLowerCase().includes(query) &&
			!value.includes(Number(suggestion.id))
		).slice(0, 10);
	});

	const selectedTags = $derived.by(() => {
		if (!suggestions || !value) return [];
		return value.map(tagId => {
			const suggestion = suggestions.find(s => Number(s.id) === tagId);
			return {
				id: tagId,
				label: suggestion?.label || `Tag ${tagId}`
			};
		});
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

	const selectSuggestion = (suggestion: { id: number | string; label: string }) => {
		const tagObj = { id: Number(suggestion.id), label: suggestion.label };
		const validated = validate(tagObj, value);
		if (validated !== undefined) {
			value = [...value, validated];
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
			if (showDropdown && selectedSuggestionIndex >= 0 && selectedSuggestionIndex < filteredSuggestions.length) {
				selectSuggestion(filteredSuggestions[selectedSuggestionIndex]);
				return;
			}

			enter();
			return;
		}

		// Handle arrow keys for suggestion navigation
		if (showDropdown && filteredSuggestions.length > 0) {
			if (e.key === 'ArrowDown') {
				e.preventDefault();
				selectedSuggestionIndex = Math.min(selectedSuggestionIndex + 1, filteredSuggestions.length - 1);
				return;
			}
			if (e.key === 'ArrowUp') {
				e.preventDefault();
				selectedSuggestionIndex = Math.max(selectedSuggestionIndex - 1, -1);
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

{#if showDropdown && filteredSuggestions.length > 0}
	<div class="absolute top-full left-0 right-0 z-50 mt-1">
		<div class="bg-popover text-popover-foreground overflow-hidden rounded-md border shadow-md">
			<div class="p-1 max-h-[200px] overflow-y-auto">
				{#each filteredSuggestions as suggestion, index}
					<div
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
		</div>
	</div>
{/if}

</div>
