/*
	Installed from @ieedan/shadcn-svelte-extras
*/

import type { HTMLInputAttributes } from 'svelte/elements';

export type TagsInputPropsWithoutHTML = {
	value?: number[];
	validate?: (val: { id: number; label: string }, tags: number[]) => number | undefined;
	suggestions?: { id: number | string; label: string }[];
};

export type TagsInputProps = TagsInputPropsWithoutHTML & Omit<HTMLInputAttributes, 'value'>;
