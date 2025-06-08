<script lang="ts">
    import {assetClasses} from "$lib/state/AssetClasses.svelte";
    import {onMount} from "svelte";
    import {page} from "$app/state"
    import {Spinner} from "flowbite-svelte";
    import {m} from "$lib/paraglide/messages.js";
    import NotFound from "$lib/components/NotFound.svelte";

    let assetClass = $derived.by(() => {
        return assetClasses.assetClasses.find((c) => c.id === page.params.assetClassId)
    })

    onMount(() => {
        assetClasses.load()
    })
</script>

<svelte:head>
    <title>{m.app_name()} | {assetClass?.name} {m.assets()}</title>
    <meta content={m.app_default_description()} name="description" />
</svelte:head>

{#if !assetClass && assetClasses.loading}
        <div class="w-full h-full flex justify-center items-center">
            <Spinner />
        </div>
{:else}
        {#if !assetClass}
            <NotFound
                    message={m.asset_class_not_found()}
            />
        {:else}
            <span>found</span>
        {/if}
{/if}