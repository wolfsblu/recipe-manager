import {client} from "$lib/api/client";

export const getIngredients = async (params?: { cursor?: string | null; limit?: number }) => {
    const response = await client.GET("/ingredients", {
        params: {
            query: {
                cursor: params?.cursor ?? undefined,
                limit: params?.limit ?? undefined
            }
        }
    })
    if (response.error) {
        throw response.error
    }
    return response.data
}

export const addIngredient = async (ingredient: { name: string }) => {
    const response = await client.POST("/ingredients", {
        body: ingredient
    })
    if (response.error) {
        throw response.error
    }
    return response.data
}

export const updateIngredient = async (id: number, ingredient: { name: string }) => {
    const response = await client.PUT("/ingredients/{ingredientId}", {
        params: {
            path: {
                ingredientId: id
            }
        },
        body: ingredient
    })
    if (response.error) {
        throw response.error
    }
    return response.data
}

export const deleteIngredient = async (id: number) => {
    const response = await client.DELETE("/ingredients/{ingredientId}", {
        params: {
            path: {
                ingredientId: id
            }
        }
    })
    if (response.error) {
        throw response.error
    }
    return response.data
}