import * as z from "zod";

export const formSchema = z.object({
    email: z.email(),
    password: z.string().nonempty(),
    locale: z.string().nonempty(),
});

export type FormSchema = typeof formSchema;