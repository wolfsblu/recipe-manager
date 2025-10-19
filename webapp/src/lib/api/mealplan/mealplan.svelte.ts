import {client} from "$lib/api/client";

export const getMealPlan = async (params?: { from?: string; until?: string; cursor?: string | null; limit?: number }) => {
    const response = await client.GET("/mealplan", {
        params: {
            query: {
                ...(params?.from && { from: params.from }),
                ...(params?.until && { until: params.until }),
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