<script lang="ts">
    import '../../app.css';
    import {DarkMode} from "flowbite-svelte";
    import {m} from '$lib/paraglide/messages.js';

    let { children } = $props();
</script>

<svelte:head>
    <title>{m.app_name()}| {m.app_default_title()}</title>
    <meta content={m.app_default_description()} name="description" />
</svelte:head>

<div class="h-screen flex flex-col justify-center items-center">
    <img alt="Logo" class="mb-4 w-24 h-24" src="/favicon.png" />
    <DarkMode />
    {@render children()}
</div>