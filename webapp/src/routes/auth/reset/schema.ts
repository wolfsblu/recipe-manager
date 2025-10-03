import * as z from "zod";

export const formSchema = z.object({
    email: z.email(),
});

export type FormSchema = typeof formSchema;