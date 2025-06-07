<script lang="ts">
    import * as Icons from "flowbite-svelte-icons";
    import {Button, Input, Label} from "flowbite-svelte";
    import DynamicIcon from "$lib/components/DynamicIcon.svelte";

    let {
        selectedIcon = $bindable("PlusOutline"),
    }: {
        selectedIcon?: string;
    } = $props();

    let searchString = $state("");
</script>

<div class="flex flex-col gap-4">
    <Label class="block">
        <Input
            bind:value={searchString}
            class="mt-2"
            placeholder="Type to search icons..."
        />
    </Label>
    <div class="grid auto-cols-auto grid-cols-2 sm:grid-cols-1 md:grid-cols-2 gap-4 h-96 max-h-96 overflow-y-auto">
        {#each Object.keys(Icons).filter(
            iconName =>
            iconName.toLowerCase().includes(searchString.toLowerCase())
        ) as iconName (iconName)}
            <Button
                class="p-2 border rounded justify-center items-center flex flex-col"
                color={
                    selectedIcon === iconName ? "dark" : "light"
                }
                onclick={() => selectedIcon = iconName}
            >
                <DynamicIcon
                    iconName={iconName}
                    />

                <span class="block text-sm text-center mt-1 overflow-hidden text-ellipsis whitespace-nowrap max-w-full">
                    {iconName}
                </span>
            </Button>
        {/each}
    </div>
</div>


