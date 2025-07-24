import * as z from "zod";

export const formSchema = z.object({
    name: z.string().nonempty(),
    servings: z.number().positive(),
    minutes: z.number().positive(),
    description: z.string().nonempty(),
});

export type FormSchema = typeof formSchema;