import * as z from "zod";

export const formSchema = z.object({
    password: z.string().nonempty()
});

export type FormSchema = typeof formSchema;