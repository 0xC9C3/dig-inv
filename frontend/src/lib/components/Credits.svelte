<script>
    import Attributions from "./Attributions.svelte";
    import {m} from '$lib/paraglide/messages.js';
    import Logo from "$lib/components/Logo.svelte";
</script>

<div class="p-2 flex flex-col items-center justify-center gap-2 grow h-full">
    <div class="w-24 h-24 mb-14 logo-animate">
        <Logo />
    </div>
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

<style>
    .logo-animate {
        animation: logo-rotate-bounce 3s infinite ease-in-out;
    }
    @keyframes logo-rotate-bounce {
        0% {
            transform: translateY(0);
        }
        50% {
            transform: translateY(15px);
        }
        100% {
            transform: translateY(0px);
        }
    }
</style>