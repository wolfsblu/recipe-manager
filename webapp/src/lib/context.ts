import { Context } from 'runed';
import type { UseBoolean } from './hooks/use-boolean.svelte';

const MenuKey = 'command-menu-context'

export const commandContext = new Context<UseBoolean>(MenuKey);