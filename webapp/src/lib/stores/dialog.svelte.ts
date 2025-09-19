import { writable } from 'svelte/store';

export const addIngredientDialogOpen = writable(false);
export const addUnitDialogOpen = writable(false);

export const dialogStore = {
    openAddIngredientDialog() {
        addIngredientDialogOpen.set(true);
    },
    
    openAddUnitDialog() {
        addUnitDialogOpen.set(true);
    }
};