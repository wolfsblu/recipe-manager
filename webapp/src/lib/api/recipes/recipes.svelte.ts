import {client} from "$lib/api/client";

export const getRecipes = async () => {
    const response = await client.GET("/recipes")
    if (response.error) {
        throw response.error
    }
    return response.data
}

export const addRecipe = async (recipe) => {
    const response = await client.POST("/recipes", {
        body: recipe
    })
    if (response.error) {
        throw response.error
    }
    return response.data
}