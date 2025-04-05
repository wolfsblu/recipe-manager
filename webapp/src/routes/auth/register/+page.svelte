<script>
    import { Card, Button, Label, Input, Checkbox } from 'flowbite-svelte';
</script>

<Card class="self-center my-auto">
    <form class="flex flex-col space-y-6" action="/">
        <h3 class="text-xl font-medium text-gray-900 dark:text-white">Register on our platform</h3>
        <Label class="space-y-2">
            <span>Email</span>
            <Input type="email" name="email" placeholder="name@company.com" required />
        </Label>
        <Label class="space-y-2">
            <span>Your password</span>
            <Input type="password" name="password" placeholder="•••••" required />
        </Label>
        <Button type="submit" class="w-full">Create a new account</Button>
    </form>
</Card>