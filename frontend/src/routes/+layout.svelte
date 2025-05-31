<script lang="ts">
    import auth from '$lib/state/Auth.svelte'
    import {DarkMode, Spinner} from 'flowbite-svelte';
    import ToastOutlet from "$lib/components/ToastOutlet.svelte";

    let { children } = $props();
</script>

{#if auth.loading}
    <div class="h-screen flex flex-col gap-4 items-center justify-center">
        <img src="/favicon.png" alt="Logo" class="w-24 h-24 mb-4" />
        <DarkMode />
        <Spinner />
    </div>
{:else}
    {@render children()}
{/if}

<div class="fixed bottom-0 right-0 p-4">
    <ToastOutlet />
</div>