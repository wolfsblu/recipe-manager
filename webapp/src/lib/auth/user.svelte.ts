import { client } from '$lib/api/client'

export const userState = $state({

})

export interface Credentials {
    email: string
    password: string
}

export const login = async (user: Credentials) => {
    const response = await client.POST("/login", {
        body: {
            email: user.email,
            password: user.password,
        }
    })
    if (response.error) {
        throw response.error
    }
    return response.data
}