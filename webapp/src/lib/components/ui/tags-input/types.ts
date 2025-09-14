/*
	Installed from @ieedan/shadcn-svelte-extras
*/

import type { HTMLInputAttributes } from 'svelte/elements';

export type TagsInputPropsWithoutHTML = {
	value?: string[];
	validate?: (val: string, tags: string[]) => string | undefined;
	suggestions?: { id: number | string; label: string }[];
};

export type TagsInputProps = TagsInputPropsWithoutHTML & Omit<HTMLInputAttributes, 'value'>;
