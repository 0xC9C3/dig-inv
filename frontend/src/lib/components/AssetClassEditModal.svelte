<script lang="ts">
    import {Modal} from "flowbite-svelte";
    import {onMount} from "svelte";
    import AssetClassForm from "$lib/components/AssetClassForm.svelte";
    import type {DigInvAssetClass} from "$lib/api";
    import {m} from '$lib/paraglide/messages.js';

    let {
        onclose,
        assetClass,
        actionName = m.save(),
        onSubmit = $bindable(async () => {})
    }: {
        onclose: () => void;
        assetClass?: DigInvAssetClass;
        actionName?: string;
        onSubmit?: (assetClass: DigInvAssetClass, e: SubmitEvent) => Promise<void>;
    } = $props();

    // this seems to be needed, since spawning the modal with open = true does not work
    let open = $state(false);
    onMount(() => {
        open = true;
    })

    const closeModal = () => {
        open = false;
        onclose();
    };
</script>

<Modal bind:open={open} onclose={onclose} size="xl">
    <AssetClassForm
        actionName={actionName}
        assetClass={assetClass}
        onCancel={closeModal}
        onSubmit={async (assetClass: DigInvAssetClass, e: SubmitEvent) => {
            await onSubmit(assetClass, e);
            open = false;
            onclose();
        }}
    />
</Modal>