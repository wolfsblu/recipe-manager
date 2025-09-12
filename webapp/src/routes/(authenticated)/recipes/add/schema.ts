import * as z from "zod";

export const formSchema = z.object({
    name: z.string().nonempty(),
    servings: z.number().positive().or(z.literal('')).default(''),
    minutes: z.number().positive().or(z.literal('')).default(''),
    tags: z.array(z.string().nonempty()),
    description: z.string().nonempty(),
    steps: z.array(z.object({
        ingredients: z.array(z.object({
            unitId: z.number().positive(),
            ingredientId: z.number().positive(),
            amount: z.number().positive().or(z.literal('')).default(''),
        })),
        instructions: z.string().nonempty()
    })).nonempty(),
});

export type FormSchema = typeof formSchema;