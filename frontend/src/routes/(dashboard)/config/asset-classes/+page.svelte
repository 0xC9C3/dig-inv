<script lang="ts">
    import {
        Button,
        Input,
        Label,
        Modal,
        Spinner,
        Table,
        TableBody,
        TableBodyCell,
        TableBodyRow,
        TableHead,
        TableHeadCell
    } from "flowbite-svelte";
    import {PlusOutline} from "flowbite-svelte-icons";
    import {m} from '$lib/paraglide/messages.js';
    import {mount, onMount} from "svelte";
    import AssetClassEditModal from "$lib/components/AssetClassEditModal.svelte";
    import {assetClasses, createAssetClasses} from "$lib/state/AssetClasses.svelte";
    import toasts from "$lib/state/Toast.svelte";
    import ColorPicker from 'svelte-awesome-color-picker';
    import IconPicker from "$lib/components/IconPicker.svelte";
    import DynamicIcon from "$lib/components/DynamicIcon.svelte";

    let createAssetClassModal = $state(false);
    let modalBase: HTMLDivElement;
    let editModalInstance: ReturnType<typeof mount> | undefined = $state(undefined);
    let hex: string = $state("#ffffff");
    let icon: string = $state("PlusOutline");

    const loadAssetClasses = () => {
        assetClasses.load()
            .catch(error => {
                console.error("Failed to load asset classes:", error);
                toasts.addToast(m.asset_classes_load_error(), "error");
            });
    };

    onMount(() => {
        loadAssetClasses();
    })

    const openEditModal = () => {
        editModalInstance = mount(
            AssetClassEditModal, {
                target: modalBase,
                props: {
                    onclose: () => {
                        if (editModalInstance) {
                            // lifecycle_double_unmount?
                            // unmount(editModalInstance);
                            editModalInstance = undefined;
                        }
                    },
                },
        })
    };

    const saveNewAssetClass = async (e: Event) => {
        e.preventDefault();
        const formData = new FormData(e.target as HTMLFormElement);
        const name = formData.get("name") as string;
        const description = formData.get("description") as string;

        try {
            await createAssetClasses.create({
                    name,
                    description,
                    icon: icon,
                    color: hex
            });
            createAssetClassModal = false;
            loadAssetClasses();
        } catch (error) {
            console.error("Failed to create asset class:", error);
            toasts.addToast(m.asset_classes_create_error(), "error");
        }
    };

</script>


<svelte:head>
    <title>{m.app_name()} | {m.asset_classes()}</title>
    <meta content={m.app_default_description()} name="description" />
</svelte:head>

<div bind:this={modalBase}></div>

<div class="gap-4">
    <h1 class="text-2xl font-bold">{m.asset_classes()}</h1>
    <p class="text-gray-600 dark:text-gray-400">{m.asset_classes_subtitle()}</p>
    <div class="flex py-4 gap-4 justify-between">
        <div class="flex flex-col justify-center">
            <Button
                color="primary"
                disabled={createAssetClasses.loading}
                onclick={() => createAssetClassModal = true}
                >
                {#if createAssetClasses.loading}
                    <Spinner size="5" />
                {:else}
                    <PlusOutline class="mr-2" />
                {/if}
                {m.asset_class_create()}
            </Button>
        </div>

        <div class:hidden={!assetClasses.loading}>
            <Spinner size="8" />
        </div>
    </div>
</div>

<Modal bind:open={createAssetClassModal} size="xs">
    <form action="#" class="flex flex-col space-y-6" method="dialog"
        onsubmit={async (e) => saveNewAssetClass(e)}
    >
        <h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">Create Asset Class</h3>
        <Label class="space-y-2">
            <span>Name</span>
            <Input name="name" placeholder="Asset Class Name" required />
        </Label>
        <Label class="space-y-2">
            <span>Description</span>
            <Input name="description" placeholder="Description of the asset class" />
        </Label>
        <Label class="space-y-2">
            <span>Icon</span>

            <div>
                <IconPicker
                    bind:selectedIcon={icon}
                    />
            </div>
        </Label>
        <Label class="space-y-2">
            <span>Color</span>

            <div class="colorPicker flex flex-col justify-center items-center">
                <ColorPicker
                        bind:hex
                        isDialog={false}
                        position="responsive"
                />
            </div>
        </Label>

        <Button class="w-full1" type="submit">
            Create Asset Class
        </Button>
        <Button class="w-full1" color="gray" onclick={() => createAssetClassModal = false} type="button">Cancel</Button>
    </form>
</Modal>

<Table hoverable={true}>
    <TableHead>
    <TableHeadCell>{m.icon()}</TableHeadCell>
    <TableHeadCell>{m.color()}</TableHeadCell>
        <TableHeadCell>{m.name()}</TableHeadCell>
        <TableHeadCell>{m.asset_count()}</TableHeadCell>
    </TableHead>
    <TableBody>
        {#each assetClasses.assetClasses as assetClass (assetClass.id)}
            <TableBodyRow onclick={openEditModal} class="cursor-pointer">
                <TableBodyCell>
                    <DynamicIcon iconName={assetClass.icon} />
                </TableBodyCell>
                <TableBodyCell>
                    <div class="w-6 h-6 rounded-full" style="background-color: {assetClass.color};"></div>
                </TableBodyCell>
                <TableBodyCell>{assetClass.name}</TableBodyCell>
                <TableBodyCell>Sliver</TableBodyCell>
            </TableBodyRow>
        {/each}
    </TableBody>
</Table>
