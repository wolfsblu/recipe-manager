import createClient from "openapi-fetch";
import type { paths } from './schema'
import { PUBLIC_API_URL } from "$env/static/public";

export const client = createClient<paths>({
    baseUrl: PUBLIC_API_URL ?? '/api',
    credentials: "include",
})