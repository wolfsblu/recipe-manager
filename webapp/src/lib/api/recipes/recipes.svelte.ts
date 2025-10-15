import {client} from "$lib/api/client";

export const getRecipes = async () => {
    const response = await client.GET("/recipes")
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

export const getIngredients = async () => {
    const response = await client.GET("/ingredients")
    if (response.error) {
        throw response.error
    }
    return response.data
}

export const getUnits = async () => {
    const response = await client.GET("/units")
    if (response.error) {
        throw response.error
    }
    return response.data
}

export const getTags = async () => {
    const response = await client.GET("/tags")
    if (response.error) {
        throw response.error
    }
    return response.data
}