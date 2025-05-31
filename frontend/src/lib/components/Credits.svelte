<script>
    import Attributions from "./Attributions.svelte";
    import {m} from '$lib/paraglide/messages.js';
</script>

<div class="p-2 flex flex-col items-center justify-center gap-2">
    <Attributions />
    <a
        class="link underline"
        href="github.com/0xC9C3"
        target="_blank"
        >
        {m.app_name()} by 0xC9C3
    </a>
</div>
