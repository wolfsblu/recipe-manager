import * as z from "zod";

export const formSchema = z.object({
    name: z.string().nonempty(),
    servings: z.number().positive().or(z.literal('')).default(''),
    minutes: z.number().positive().or(z.literal('')).default(''),
    tags: z.array(z.string().nonempty()),
    description: z.string().nonempty(),
    ingredients: z.array(z.string().nonempty()),
});

export type FormSchema = typeof formSchema;