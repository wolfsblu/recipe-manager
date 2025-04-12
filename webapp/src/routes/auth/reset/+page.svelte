<script>
    import { Card, Button, Label, Input } from 'flowbite-svelte';
</script>

<Card class="self-center my-auto">
    <form class="flex flex-col space-y-6" action="/">
        <h3 class="text-xl font-medium text-gray-900 dark:text-white">
            Reset your password
        </h3>
        <Label class="space-y-2">
            <span>Email</span>
            <Input type="email" name="email" placeholder="name@company.com" required />
        </Label>
        <Button type="submit" class="w-full">Reset your password</Button>
    </form>
</Card>