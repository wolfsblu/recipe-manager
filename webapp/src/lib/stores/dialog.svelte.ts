let addIngredientDialogOpen = $state(false);
let addUnitDialogOpen = $state(false);

export const dialogStore = {
    get addIngredientDialogOpen() { return addIngredientDialogOpen; },
    set addIngredientDialogOpen(value: boolean) { addIngredientDialogOpen = value; },
    
    get addUnitDialogOpen() { return addUnitDialogOpen; },
    set addUnitDialogOpen(value: boolean) { addUnitDialogOpen = value; },
    
    openAddIngredientDialog() {
        addIngredientDialogOpen = true;
    },
    
    openAddUnitDialog() {
        addUnitDialogOpen = true;
    }
};