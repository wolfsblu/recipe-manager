import {client} from "$lib/api/client";

export const getUnits = async () => {
    const response = await client.GET("/units")
    if (response.error) {
        throw response.error
    }
    return response.data
}

export const addUnit = async (unit: { name: string; symbol?: string }) => {
    const response = await client.POST("/units", {
        body: unit
    })
    if (response.error) {
        throw response.error
    }
    return response.data
}

export const updateUnit = async (id: number, unit: { name: string; symbol?: string }) => {
    const response = await client.PUT("/units/{unitId}", {
        params: {
            path: {
                unitId: id
            }
        },
        body: unit
    })
    if (response.error) {
        throw response.error
    }
    return response.data
}

export const deleteUnit = async (id: number) => {
    const response = await client.DELETE("/units/{unitId}", {
        params: {
            path: {
                unitId: id
            }
        }
    })
    if (response.error) {
        throw response.error
    }
    return response.data
}