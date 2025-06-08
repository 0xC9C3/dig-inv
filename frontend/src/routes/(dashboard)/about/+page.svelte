<script>
    import Credits from "$lib/components/Credits.svelte";
    import {m} from "$lib/paraglide/messages.js";
</script>

<svelte:head>
    <title>{m.app_name()} | {m.about()}</title>
    <meta content={m.app_default_description()} name="description" />
</svelte:head>


<Credits />