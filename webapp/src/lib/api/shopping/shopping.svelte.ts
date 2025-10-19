import { client } from "$lib/api/client";
import type { components } from "$lib/api/schema";

export type ReadShoppingList = components["schemas"]["ReadShoppingList"];
export type WriteShoppingList = components["schemas"]["WriteShoppingList"];
export type ReadShoppingListItem = components["schemas"]["ReadShoppingListItem"];
export type WriteShoppingListItem = components["schemas"]["WriteShoppingListItem"];

export const getShoppingLists = async (params?: { cursor?: string | null; limit?: number }) => {
    const response = await client.GET("/shopping-lists", {
        params: {
            query: {
                cursor: params?.cursor ?? undefined,
                limit: params?.limit ?? undefined
            }
        }
    });
    if (response.error) {
        throw response.error;
    }
    return response.data;
};

export const createShoppingList = async (shoppingList: WriteShoppingList) => {
    const response = await client.POST("/shopping-lists", {
        body: shoppingList
    });
    if (response.error) {
        throw response.error;
    }
    return response.data;
};

export const getShoppingList = async (id: number) => {
    const response = await client.GET("/shopping-lists/{shoppingListId}", {
        params: {
            path: {
                shoppingListId: id
            }
        }
    });
    if (response.error) {
        throw response.error;
    }
    return response.data;
};

export const updateShoppingList = async (id: number, shoppingList: WriteShoppingList) => {
    const response = await client.PUT("/shopping-lists/{shoppingListId}", {
        params: {
            path: {
                shoppingListId: id
            }
        },
        body: shoppingList
    });
    if (response.error) {
        throw response.error;
    }
    return response.data;
};

export const deleteShoppingList = async (id: number) => {
    const response = await client.DELETE("/shopping-lists/{shoppingListId}", {
        params: {
            path: {
                shoppingListId: id
            }
        }
    });
    if (response.error) {
        throw response.error;
    }
    return response.data;
};

export const addShoppingListItem = async (shoppingListId: number, item: WriteShoppingListItem) => {
    const response = await client.POST("/shopping-lists/{shoppingListId}/items", {
        params: {
            path: {
                shoppingListId: shoppingListId
            }
        },
        body: item
    });
    if (response.error) {
        throw response.error;
    }
    return response.data;
};

export const updateShoppingListItem = async (shoppingListId: number, itemId: number, item: WriteShoppingListItem) => {
    const response = await client.PUT("/shopping-lists/{shoppingListId}/items/{itemId}", {
        params: {
            path: {
                shoppingListId: shoppingListId,
                itemId: itemId
            }
        },
        body: item
    });
    if (response.error) {
        throw response.error;
    }
    return response.data;
};

export const deleteShoppingListItem = async (shoppingListId: number, itemId: number) => {
    const response = await client.DELETE("/shopping-lists/{shoppingListId}/items/{itemId}", {
        params: {
            path: {
                shoppingListId: shoppingListId,
                itemId: itemId
            }
        }
    });
    if (response.error) {
        throw response.error;
    }
    return response.data;
};