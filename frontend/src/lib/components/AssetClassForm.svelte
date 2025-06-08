<script lang="ts">
    import IconPicker from "$lib/components/IconPicker.svelte";
    import ColorPicker from "svelte-awesome-color-picker";
    import {Button, Input, Label, Select, Spinner, Textarea} from "flowbite-svelte";
    import {m} from '$lib/paraglide/messages.js';
    import type {DigInvAssetClass} from "$lib/api";
    import {providerKVMap, providers} from "$lib/providers";
    import type {Snippet} from "svelte";
    import {EmptySnippet} from "./EmptyComponent.svelte";

    let {
        onSubmit = $bindable(async () => {
        }),
        onCancel = $bindable(() => {
        }),
        assetClass,
        actionName = m.save()
    }: {
        onSubmit?: (assetClass: DigInvAssetClass, e: SubmitEvent) => Promise<void>;
        onCancel?: () => void;
        assetClass?: DigInvAssetClass;
        actionName?: string;
    } = $props();
    let loading = $state(false);
    let formAssetClass: DigInvAssetClass = $state(assetClass || {
        name: "",
        description: "",
        order: 0,
        icon: "PlusOutline",
        color: "#ffffff"
    });

    const submitForm = (e: SubmitEvent) => {
        e.preventDefault();
        loading = true;
        onSubmit(formAssetClass, e)
            .finally(() => {
                loading = false;
            });
    };

    const renderProviderForm = (providerName: string): Snippet => {
        const provider = providers.find(p => p.key === providerName);
        if (provider) {
            return provider.getForm()
        }

        console.warn(`Provider component for ${providerName} not found.`);
        return EmptySnippet;
    };
</script>

<form action="#" class="flex flex-col space-y-6" method="dialog"
      onsubmit={submitForm}
>
    <h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">{actionName}</h3>
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <Label class="space-y-2">
            <span>{m.name()}</span>
            <Input bind:value={formAssetClass.name}
                    name="name"
                   placeholder={m.name()} required/>
        </Label>

        <div class="flex flex-col gap-4">
            <Label class="space-y-2">
                <span>{m.provider()}</span>
                <Select bind:value={formAssetClass.provider}
                        items={providerKVMap}
                        placeholder={m.choose_a_provider()} />
            </Label>

            {#if formAssetClass.provider}
                {@render renderProviderForm(formAssetClass.provider)()}
            {/if}
        </div>

        <Label class="space-y-2">
            <span>{m.description()}</span>
            <Textarea bind:value={formAssetClass.description}
                    name="description"
                   placeholder={m.description()}/>
        </Label>


        <Label class="space-y-2">
            <span>{m.order()}</span>
            <Input bind:value={formAssetClass.order}
                   min="0"
                   name="order"
                   placeholder={m.order()} required type="number"/>
        </Label>

        <Label class="space-y-2">
            <span>{m.icon()}</span>

            <div>
                <IconPicker
                        bind:selectedIcon={formAssetClass.icon}
                />
            </div>
        </Label>
        <Label class="space-y-2">
            <span>{m.color()}</span>

            <div class="colorPicker flex flex-col justify-center items-center">
                <ColorPicker
                        bind:hex={formAssetClass.color}
                        isDialog={false}
                        position="responsive"
                />
            </div>
        </Label>
    </div>

    <Button class="w-full1" disabled={loading} type="submit">
        {#if loading}
            <span class="flex items-center justify-center">
                <Spinner size="5" />
            </span>
        {/if}
        {actionName}
    </Button>
    <Button class="w-full1" color="gray" onclick={() => onCancel()} type="button">{m.cancel()}</Button>
</form>