<script lang="ts">
    import {
        Button,
        Spinner,
        Table,
        TableBody,
        TableBodyCell,
        TableBodyRow,
        TableHead,
        TableHeadCell
    } from "flowbite-svelte";
    import {EditSolid, PlusOutline, TrashBinSolid} from "flowbite-svelte-icons";
    import {m} from '$lib/paraglide/messages.js';
    import {type Component, mount, onMount} from "svelte";
    import AssetClassEditModal from "$lib/components/AssetClassEditModal.svelte";
    import {assetClasses, createAssetClass, updateAssetClass, deleteAssetClass} from "$lib/state/AssetClasses.svelte";
    import toasts from "$lib/state/Toast.svelte";
    import DynamicIcon from "$lib/components/DynamicIcon.svelte";
    import type {DigInvAssetClass} from "$lib/api";
    import ConfirmModal from "$lib/components/ConfirmModal.svelte";
    import {fromKey} from "$lib/providers";

    let modalBase: HTMLDivElement;
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    let modalInstance: ReturnType<typeof mount> | undefined = $state(undefined);

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

    const openModal = <T extends Record<string, unknown>>(
        component: Component<T>,
        props: Omit<T, 'onclose'> & { onclose?: () => void; } = {} as never
    ) => {
        modalInstance = mount(
            component, {
                target: modalBase,
                props: {
                    ...props,
                    onclose: () => {
                        // lifecycle_double_unmount?
                        // unmount(modalInstance);
                        modalInstance = undefined;
                    },
                },
            } as never)
    }

    const openEditModal = (editClass: DigInvAssetClass) => {
        openModal(AssetClassEditModal, {
            actionName: m.asset_class_edit(),
            assetClass: JSON.parse(JSON.stringify(editClass)),
            onSubmit: async (updatedClass: DigInvAssetClass) => {
                await updateAssetClassAction(updatedClass);
            },
        });
    };

    const openNewAssetClassModal = () => {
        openModal(AssetClassEditModal, {
            actionName: m.asset_class_create(),
            onSubmit: async (newClass: DigInvAssetClass) => {
                await saveNewAssetClassAction(newClass);
            },
        });
    };

    const openDeleteConfirmationModal = (assetClass: DigInvAssetClass) => {
        openModal(ConfirmModal, {
            text: m.confirm_delete({name: assetClass.name || m.asset_class()}),
            onConfirm: async () => {
                await deleteAssetClassAction(assetClass);
            },
        });
    };

    const assetClassAction = (
        method: (targetClass: DigInvAssetClass) => Promise<void>,
        errorMessage: string
    ) =>  async (targetClass: DigInvAssetClass) => {
        try {
            await method(targetClass);
            loadAssetClasses();
        } catch (error) {
            console.error(errorMessage, error);
            toasts.addToast(errorMessage, "error");
        }
    };

    const saveNewAssetClassAction = assetClassAction(
        async (assetClass) => createAssetClass.create(assetClass),
        m.asset_classes_create_error()
    );

    const deleteAssetClassAction = assetClassAction(
        async (assetClass) => assetClass?.id ? deleteAssetClass.delete(assetClass.id) : Promise.reject(new Error("Asset class ID is required")),
        m.asset_classes_delete_error()
    );

    const updateAssetClassAction = assetClassAction(
        async (assetClass) => updateAssetClass.update(assetClass),
        m.asset_classes_update_error()
    );
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
                disabled={createAssetClass.loading}
                onclick={() => openNewAssetClassModal()}
                >
                {#if createAssetClass.loading}
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

<Table hoverable={true}>
    <TableHead>
        <TableHeadCell>{m.icon()}</TableHeadCell>
        <TableHeadCell>{m.order()}</TableHeadCell>
        <TableHeadCell>{m.provider()}</TableHeadCell>
        <TableHeadCell>{m.color()}</TableHeadCell>
        <TableHeadCell>{m.name()}</TableHeadCell>
        <TableHeadCell></TableHeadCell>
    </TableHead>
    <TableBody>
        {#each assetClasses.assetClasses as assetClass (assetClass.id)}
            <TableBodyRow>
                <TableBodyCell>
                    <DynamicIcon iconName={assetClass.icon} />
                </TableBodyCell>
                <TableBodyCell>{assetClass.order}</TableBodyCell>
                <TableBodyCell>
                    {fromKey(assetClass.provider)?.name || assetClass.provider}
                </TableBodyCell>
                <TableBodyCell>
                    <div class="w-6 h-6 rounded-full" style="background-color: {assetClass.color};"></div>
                </TableBodyCell>
                <TableBodyCell class="w-full">{assetClass.name}</TableBodyCell>
                <TableBodyCell class="flex gap-4">
                    <Button
                            color="primary"
                            class="gap-2"
                            onclick={() => openEditModal(assetClass)}
                    >
                        <EditSolid />
                        {m.edit()}
                    </Button>
                    <Button
                            color="red"
                            class="gap-2"
                            onclick={() => openDeleteConfirmationModal(assetClass)}
                    >
                        <TrashBinSolid />
                        {m.delete()}
                    </Button>
                </TableBodyCell>
            </TableBodyRow>
        {/each}
    </TableBody>
</Table>
