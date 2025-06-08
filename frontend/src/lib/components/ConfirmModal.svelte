<script lang="ts">
    import {Button, Modal} from "flowbite-svelte";
    import {onMount} from "svelte";
    import {m} from '$lib/paraglide/messages.js';

    let {
        onclose,
        onConfirm = $bindable(async () => {}),
        text = "?"
    }: {
        onclose: () => void;
        onConfirm?: (e: SubmitEvent) => Promise<void>;
        text?: string;
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

<Modal bind:open={open} onclose={onclose} size="xs">
    <form class="flex flex-col space-y-6" method="dialog"
          onsubmit={async (e: SubmitEvent) => {
              e.preventDefault();
              await onConfirm(e);
              open = false;
              onclose();
          }}>
        <p class="text-gray-500 dark:text-gray-400">{text}</p>
        <div class="flex justify-end space-x-2">
            <Button color="gray" onclick={closeModal} type="button">{m.cancel()}</Button>
            <Button color="red" type="submit">{m.confirm()}</Button>
        </div>
    </form>
</Modal>