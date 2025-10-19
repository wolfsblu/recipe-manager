import {client} from "$lib/api/client";

export const getRecipes = async (params?: { cursor?: string | null; limit?: number }) => {
    const response = await client.GET("/recipes", {
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

export const getRecipe = async (id: number) => {
    const response = await client.GET("/recipes/{recipeId}", {
        params: {
            path: {
                recipeId: id
            }
        }
    })
    if (response.error) {
        throw response.error
    }
    return response.data
}

export const addRecipe = async (recipe: any) => {
    const response = await client.POST("/recipes", {
        body: recipe
    })
    if (response.error) {
        throw response.error
    }
    return response.data
}

export const updateRecipe = async (id: number, recipe: any) => {
    const response = await client.POST("/recipes/{recipeId}", {
        params: {
            path: {
                recipeId: id
            }
        },
        body: recipe
    })
    if (response.error) {
        throw response.error
    }
    return response.data
}

export const patchRecipe = async (id: number, patch: any) => {
    const response = await client.PATCH("/recipes/{recipeId}", {
        params: {
            path: {
                recipeId: id
            }
        },
        body: patch
    })
    if (response.error) {
        throw response.error
    }
    return response.data
}

export const deleteRecipe = async (id: number) => {
    const response = await client.DELETE("/recipes/{recipeId}", {
        params: {
            path: {
                recipeId: id
            }
        }
    })
    if (response.error) {
        throw response.error
    }
    return response.data
}

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

export const getUnits = async (params?: { cursor?: string | null; limit?: number }) => {
    const response = await client.GET("/units", {
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

export const getTags = async (params?: { cursor?: string | null; limit?: number }) => {
    const response = await client.GET("/tags", {
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