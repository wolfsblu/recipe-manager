import {client} from "$lib/api/client";

wdexport const getUnits = async (params?: { cursor?: string | null; limit?: number }) => {
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