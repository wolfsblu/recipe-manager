import {client} from "$lib/api/client";

export const getMealPlan = async (from?: string, until?: string) => {
    const response = await client.GET("/mealplan", {
        params: {
            query: {
                ...(from && { from }),
                ...(until && { until })
            }
        }
    })
    if (response.error) {
        throw response.error
    }
    return response.data
}

export const createMealPlan = async (recipeId: number, date: string) => {
    const response = await client.POST("/mealplan", {
        body: {
            recipeId,
            date
        }
    })
    if (response.error) {
        throw response.error
    }
    return response.data
}

export const deleteMealPlan = async (recipeId: number, date: string) => {
    const response = await client.DELETE("/mealplan/{recipeId}", {
        params: {
            path: {
                recipeId
            },
            query: {
                date
            }
        }
    })
    if (response.error) {
        throw response.error
    }
    return response.data
}