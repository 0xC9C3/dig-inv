<script>
    import Attributions from "./Attributions.svelte";
    import {m} from '$lib/paraglide/messages.js';
    import Logo from "$lib/components/Logo.svelte";
</script>

<div class="p-2 flex flex-col items-center justify-center gap-2 grow h-full">
    <Logo />
    <Attributions />
    <a
        class="link underline"
        href="https://github.com/0xC9C3"
        rel="noopener noreferrer external"
        target="_blank"
        >
        {m.app_name()} by 0xC9C3
    </a>
</div>
