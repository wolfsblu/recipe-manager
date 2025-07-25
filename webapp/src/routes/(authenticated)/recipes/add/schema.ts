import * as z from "zod";

export const formSchema = z.object({
    name: z.string().nonempty(),
    servings: z.number().positive(),
    minutes: z.number().positive(),
    tags: z.array(z.string().nonempty()),
    description: z.string().nonempty(),
});

export type FormSchema = typeof formSchema;