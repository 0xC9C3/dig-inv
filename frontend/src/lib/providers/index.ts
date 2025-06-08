import {cloudflare} from "$lib/providers/cloudflare";

import type {Provider} from "$lib/providers/AbstractProvider";


export const providers: Provider[] = [
    cloudflare
];

export const providerKVMap = providers.reduce((map, provider) => {
    map.push({
        value: provider.key,
        name: provider.name
    })
    return map;
}, [] as {value: string, name: string}[]);


export const fromKey = (key: string | undefined): Provider | undefined => {
    return providers.find(provider => provider.key === key);
}