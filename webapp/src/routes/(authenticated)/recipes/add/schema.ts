import * as z from "zod";

export const formSchema = z.object({
    name: z.string().nonempty(),
    servings: z.number().positive().or(z.literal('')).default(''),
    minutes: z.number().positive().or(z.literal('')).default(''),
    tags: z.array(z.string().nonempty()),
    description: z.string().nonempty(),
    steps: z.array(z.object({
        ingredients: z.array(z.object({
            unit: z.string().nonempty(),
            name: z.string().nonempty(),
            amount: z.number().positive().or(z.literal('')).default(''),
        })).nonempty(),
    })).nonempty(),
});

export type FormSchema = typeof formSchema;