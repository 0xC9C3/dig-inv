import {Provider} from "$lib/providers/AbstractProvider";
import type {Snippet} from "svelte";
import {CloudflareFormSnippet} from "$lib/components/providers/CloudflareForm.svelte";

class Cloudflare extends Provider {
    protected providerKey: string = 'cloudflare';

    public getForm(): Snippet {
        return CloudflareFormSnippet
    }
}

export const cloudflare = new Cloudflare();

