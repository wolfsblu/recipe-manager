import * as z from "zod";

export const formSchema = z.object({
    name: z.string().nonempty(),
    images: z.array(z.string().url()),
    servings: z.number().positive(),
    minutes: z.number().positive(),
    tags: z.array(z.string().nonempty()),
    description: z.string().nonempty(),
    steps: z.array(z.object({
        ingredients: z.array(z.object({
            unitId: z.number({ message: "Invalid unit" }).positive({ message: "Invalid unit" }),
            ingredientId: z.number({ message: "Invalid ingredient" }).positive({ message: "Invalid ingredient" }),
            amount: z.number({ message: "Invalid amount" }).positive({ message: "Invalid amount" }),
        })),
        instructions: z.string().nonempty()
    })).nonempty(),
});

export type FormSchema = typeof formSchema;