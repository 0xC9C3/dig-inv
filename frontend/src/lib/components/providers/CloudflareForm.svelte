<script lang="ts" module>
    import {Input, Label} from 'flowbite-svelte';
    import {m} from '$lib/paraglide/messages.js';

    export {CloudflareFormSnippet}
</script>

{#snippet CloudflareFormSnippet()}
    <Label class="space-y-2">
        <span>{m.api_token()}</span>
        <Input
               name="api_token"
               placeholder={m.api_token()} required/>
    </Label>
{/snippet}


