<script lang="ts">
    import toasts from "$lib/state/Toast.svelte.js";
    import {CheckCircleSolid, CloseCircleSolid, ExclamationCircleSolid, FireOutline} from "flowbite-svelte-icons";
    import {Toast} from "flowbite-svelte";

    const getToastColor = (type: string) => {
        switch (type) {
            case 'success':
                return 'green';
            case 'error':
                return 'red'
            case 'warning':
                return 'yellow'
            default:
                return 'blue';
        }
    };
</script>

{#each toasts.toasts as toast (toast.id)}
    <Toast
            color={getToastColor(toast.type)}
            class="mb-4">
        {#snippet icon()}
            {#if toast.type === 'success'}
                <CheckCircleSolid class="h-6 w-6" />
            {:else if toast.type === 'error'}
                <CloseCircleSolid class="h-6 w-6" />
            {:else if toast.type === 'warning'}
                <ExclamationCircleSolid class="h-6 w-6" />
            {:else}
                <FireOutline class="h-6 w-6" />
            {/if}
        {/snippet}

        {toast.message}
    </Toast>
{/each}