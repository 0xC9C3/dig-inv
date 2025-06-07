<script lang="ts">
    import {Button, Checkbox, Input, Label, Modal} from "flowbite-svelte";
    import {onMount} from "svelte";

    let {
        onclose,
    }: {
        onclose: () => void;
    } = $props();

    // this seems to be needed, since spawning the modal with open = true does not work
    let open = $state(false);
    onMount(() => {
        open = true;
    })
</script>

<Modal bind:open={open} onclose={onclose} size="xl">
    <form action="#" class="flex flex-col space-y-6" method="dialog">
        <h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">Sign in to our platform</h3>
        <Label class="space-y-2">
            <span>Email</span>
            <Input name="email" placeholder="name@company.com" required type="email" />
        </Label>
        <Label class="space-y-2">
            <span>Your password</span>
            <Input name="password" placeholder="•••••" required type="password" />
        </Label>
        <div class="flex items-start">
            <Checkbox>Remember me</Checkbox>
            <a class="text-primary-700 dark:text-primary-500 ms-auto text-sm hover:underline" href="/">Lost password?</a>
        </div>
        <Button class="w-full1" type="submit">Login to your account</Button>
        <div class="text-sm font-medium text-gray-500 dark:text-gray-300">
            Not registered? <a class="text-primary-700 dark:text-primary-500 hover:underline" href="/">Create account</a>
        </div>
    </form>
</Modal>