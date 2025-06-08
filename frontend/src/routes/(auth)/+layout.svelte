<script lang="ts">
    import '../../app.css';
    import {DarkMode} from "flowbite-svelte";
    import {m} from '$lib/paraglide/messages.js';
    import Logo from "$lib/components/Logo.svelte";

    let { children } = $props();
</script>

<svelte:head>
    <title>{m.app_name()}| {m.app_default_title()}</title>
    <meta content={m.app_default_description()} name="description" />
</svelte:head>

<div class="h-screen flex flex-col justify-center items-center">
    <Logo />
    <DarkMode />
    {@render children()}
</div>