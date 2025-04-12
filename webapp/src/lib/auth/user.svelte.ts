import { client } from '$lib/api/client'

export interface Credentials {
    email: string
    password: string
}

export interface User {
    id: number
    email: string
}

const getDefaultUser = (): User => ({
    id: 0,
    email: '',
})

export const user = $state(getDefaultUser())

export const isAuthenticated = () => user.id !== getDefaultUser().id

export const getProfile = async () => {
    const response = await client.GET("/user/profile")
    if (response.error) {
        throw response.error
    }
    Object.assign(user, response.data)
    return response.data
}

export const login = async (credentials: Credentials) => {
    const response = await client.POST("/login", {
        body: {
            email: credentials.email,
            password: credentials.password,
        }
    })
    if (response.error) {
        throw response.error
    }
    Object.assign(user, response.data)
    return response.data
}

export const logout = async () => {
    await client.POST("/logout")
    Object.assign(user, getDefaultUser())
}

export const register = async (credentials: Credentials) => {
    const response = await client.POST("/register", {
        body: {
            email: credentials.email,
            password: credentials.password,
        }
    })
    if (response.error) {
        throw response.error
    }
    return response.data
}