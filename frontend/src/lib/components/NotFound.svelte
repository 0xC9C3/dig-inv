<script lang="ts">
    import {m} from '$lib/paraglide/messages.js';

    let {message = $bindable(m.not_found_default_message())}: {message?: string} = $props();
</script>

<div class="flex flex-col items-center justify-center w-full h-full">
    <h1 class="text-2xl font-bold mb-4">😔 {message}</h1>
</div>