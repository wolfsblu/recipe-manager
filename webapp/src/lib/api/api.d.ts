interface Credentials {
    email: string
    password: string
}

interface APIError {
    message: string
}

interface User {
    id: number
    email: string
}

interface Recipe {
    name: string
    servings: number | null
    minutes: number | null
}