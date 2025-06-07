<script lang="ts">
    import {Button} from "flowbite-svelte";
    import toasts from "$lib/state/Toast.svelte";

    import {m} from '$lib/paraglide/messages.js';
</script>

<svelte:head>
    <title>{m.app_name()}| {m.app_default_title()}</title>
    <meta content={m.app_default_description()} name="description" />
</svelte:head>

<h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the documentation</p>

<Button
        onclick={() => toasts.addToast("This is a toast message!", "error")}
>
    toast
</Button>