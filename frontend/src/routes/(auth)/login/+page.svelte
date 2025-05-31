<script lang="ts">
    import {Button} from "flowbite-svelte";
    import {ArrowUpRightFromSquareOutline} from "flowbite-svelte-icons";
    import auth from "$lib/state/Auth.svelte";
</script>

<Button
    class="mb-4 gap-2 my-4"
    color="primary"
    onclick={() => auth.initialize()}
    size="lg"
    >
    <span>
        Login
    </span>
    <ArrowUpRightFromSquareOutline />
</Button>

