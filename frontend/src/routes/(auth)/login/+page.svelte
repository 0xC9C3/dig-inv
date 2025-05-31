<script lang="ts">
    import {A, Button, Checkbox, Input, Label} from "flowbite-svelte";
</script>

<form>
    <div class="mb-6">
        <Label class="mb-2" for="email">Email address</Label>
        <Input id="email" placeholder="john.doe@company.com" required type="email" />
    </div>
    <div class="mb-6">
        <Label class="mb-2" for="password">Password</Label>
        <Input id="password" placeholder="•••••••••" required type="password" />
    </div>
    <div class="mb-6">
        <Label class="mb-2" for="confirm_password">Confirm password</Label>
        <Input id="confirm_password" placeholder="•••••••••" required type="password" />
    </div>
    <Checkbox divClass="mb-6 gap-1 rtl:space-x-reverse" required>
        I agree with the <A class="text-primary-700 dark:text-primary-600 hover:underline" href="/">terms and conditions</A>.
    </Checkbox>
    <Button type="submit">Submit</Button>
</form>

