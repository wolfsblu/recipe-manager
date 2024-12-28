import {
    confirmAccount as apiConfirmAccount,
    fetchProfile as apiFetchProfile,
    login as apiLogin,
    logout as apiLogout,
    register as apiRegister,
    resetPassword as apiResetPassword,
    updatePassword as apiUpdatePassword,
} from "../api/client";

let profile: User | null = $state(null)

export const createUser = () => {
    const confirmAccount = async (token: string) => {
        const response = await apiConfirmAccount(token)
        if (response.error) {
            throw response.error
        }
    }

    const login = async (credentials: Credentials) => {
        const response = await apiLogin(credentials)
        if (response.error) {
            throw response.error
        } else {
            profile = response.data
        }
    }

    const logout = async () => {
        const response = await apiLogout()
        if (response.error) {
            throw response.error
        } else {
            profile = null
        }
    }

    const fetchProfile = async () => {
        const response = await apiFetchProfile()
        if (response.error) {
            throw response.error
        } else {
            profile = response.data
        }
    }

    const register = async (credentials: Credentials) => {
        const response = await apiRegister(credentials)
        if (response.error) {
            throw response.error
        }
    }

    const resetPassword = async (email: string) => {
        const response = await apiResetPassword(email)
        if (response.error) {
            throw response.error
        }
    }
    const updatePassword = async (password: string, token: string) => {
        const response = await apiUpdatePassword(password, token)
        if (response.error) {
            throw response.error
        }
    }

    return {
        get profile() {
            return profile
        },
        confirmAccount,
        fetchProfile,
        login,
        logout,
        register,
        resetPassword,
        updatePassword,
    }
}