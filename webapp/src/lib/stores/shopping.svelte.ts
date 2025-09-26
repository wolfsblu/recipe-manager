import type { ReadShoppingList, ReadShoppingListItem } from "$lib/api/shopping/shopping.svelte";

interface ShoppingState {
    lists: ReadShoppingList[];
    currentListId: number | null;
    isLoading: boolean;
    error: string | null;
}

class ShoppingStore {
    private state = $state<ShoppingState>({
        lists: [],
        currentListId: null,
        isLoading: false,
        error: null
    });

    get lists() {
        return this.state.lists;
    }

    get currentListId() {
        return this.state.currentListId;
    }

    get currentList() {
        if (!this.state.currentListId) return null;
        return this.state.lists.find(list => list.id === this.state.currentListId) || null;
    }

    get isLoading() {
        return this.state.isLoading;
    }

    get error() {
        return this.state.error;
    }

    get hasLists() {
        return this.state.lists.length > 0;
    }

    setLists(lists: ReadShoppingList[]) {
        this.state.lists = lists;
        // If no current list is selected and we have lists, select the first one
        if (!this.state.currentListId && lists.length > 0) {
            this.state.currentListId = lists[0].id;
        }
        // If current list was deleted, switch to first available list
        if (this.state.currentListId && !lists.find(l => l.id === this.state.currentListId)) {
            this.state.currentListId = lists.length > 0 ? lists[0].id : null;
        }
    }

    setCurrentListId(listId: number | null) {
        this.state.currentListId = listId;
    }

    setLoading(loading: boolean) {
        this.state.isLoading = loading;
    }

    setError(error: string | null) {
        this.state.error = error;
    }

    addList(list: ReadShoppingList) {
        this.state.lists = [...this.state.lists, list];
        // Set as current list if it's the first one
        if (this.state.lists.length === 1) {
            this.state.currentListId = list.id;
        }
    }

    updateList(updatedList: ReadShoppingList) {
        this.state.lists = this.state.lists.map(list => 
            list.id === updatedList.id ? updatedList : list
        );
    }

    removeList(listId: number) {
        this.state.lists = this.state.lists.filter(list => list.id !== listId);
        // If we deleted the current list, select another one
        if (this.state.currentListId === listId) {
            this.state.currentListId = this.state.lists.length > 0 ? this.state.lists[0].id : null;
        }
    }

    updateListItems(listId: number, items: ReadShoppingListItem[]) {
        this.state.lists = this.state.lists.map(list => 
            list.id === listId ? { ...list, items } : list
        );
    }

    addItemToList(listId: number, item: ReadShoppingListItem) {
        this.state.lists = this.state.lists.map(list => 
            list.id === listId ? { ...list, items: [...list.items, item] } : list
        );
    }

    updateItemInList(listId: number, updatedItem: ReadShoppingListItem) {
        this.state.lists = this.state.lists.map(list => 
            list.id === listId 
                ? { 
                    ...list, 
                    items: list.items.map(item => 
                        item.id === updatedItem.id ? updatedItem : item
                    ) 
                } 
                : list
        );
    }

    removeItemFromList(listId: number, itemId: number) {
        this.state.lists = this.state.lists.map(list => 
            list.id === listId 
                ? { ...list, items: list.items.filter(item => item.id !== itemId) }
                : list
        );
    }

    reset() {
        this.state.lists = [];
        this.state.currentListId = null;
        this.state.isLoading = false;
        this.state.error = null;
    }
}

export const shoppingStore = new ShoppingStore();