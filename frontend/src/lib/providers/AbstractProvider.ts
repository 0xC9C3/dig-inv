import {m} from '$lib/paraglide/messages.js';
import type {Snippet} from "svelte";

export abstract class Provider {
    protected abstract providerKey: string;

    public get key(): string {
        return this.providerKey;
    }

    public get name(): string {
        // @ts-expect-error needs to be keyof m
        return m[`provider_${this.providerKey}`]() || this.providerKey;
    }

    public abstract getForm(): Snippet;
}
