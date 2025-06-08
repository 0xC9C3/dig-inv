<script lang="ts">
    import '../../app.css';
    import {
        Avatar,
        DarkMode,
        Navbar,
        NavBrand,
        NavHamburger,
        Search,
        Sidebar,
        SidebarButton,
        SidebarDropdownWrapper,
        SidebarGroup,
        SidebarItem,
        ToolbarButton,
        uiHelpers
    } from "flowbite-svelte";
    import {
        AdjustmentsHorizontalSolid,
        ArrowRightToBracketOutline,
        ChartOutline,
        InfoCircleOutline,
        OpenDoorOutline,
        SearchOutline,
        ToolsOutline,
        UsersGroupOutline,
        UserSolid
    } from "flowbite-svelte-icons";
    import {fade} from "svelte/transition";
    import {page} from "$app/state";
    import auth from "$lib/state/Auth.svelte";
    import {m} from '$lib/paraglide/messages.js';
    import {assetClasses} from "$lib/state/AssetClasses.svelte";
    import DynamicIcon from "$lib/components/DynamicIcon.svelte";

    let activeUrl = $state(page.url.pathname);
    const iconClass = "h-5 w-5 text-gray-500 transition duration-75 group-hover:text-gray-900 dark:text-gray-400 dark:group-hover:text-white";
    const sidebar = uiHelpers();
    let isDemoOpen = $state(false);
    const closeDemoSidebar = sidebar.close;
    $effect(() => {
        isDemoOpen = sidebar.isOpen;
        activeUrl = page.url.pathname;
    });


    let { children } = $props();
</script>

<div class="h-screen flex flex-col">
    <Navbar>
        
        {#snippet children({ hidden, toggle })}
            <SidebarButton class="mb-2" onclick={sidebar.toggle} />
            <NavBrand href="/">
                <img src="/favicon.png" alt="Logo" class="h-8 mr-3" />
                <h1 class="self-center text-xl font-semibold whitespace-nowrap dark:text-white">{m.app_name()}</h1>
            </NavBrand>
            <div class="flex md:order-2">
                <ToolbarButton class="block md:hidden" onclick={toggle}>
                    <SearchOutline class="h-5 w-5 text-gray-500 dark:text-gray-400" />
                </ToolbarButton>
                <div class="hidden md:block">
                    <Search size="md" class="ms-auto" placeholder="Search..." />
                </div>

                <div class="ml-3">
                    <DarkMode />
                </div>

                <ToolbarButton class="block" onclick={() => auth.logout()}>
                    <OpenDoorOutline class="h-5 w-5 text-gray-500 dark:text-gray-400" />
                </ToolbarButton>

                <NavHamburger />
            </div>
            {#if !hidden}
                <div class="mt-2 w-full md:hidden" transition:fade>
                    <Search size="md" placeholder="Search..." />
                </div>
            {/if}
        {/snippet}
    </Navbar>

    <div class="relative grow">
        <Sidebar activeClass="p-2" {activeUrl} backdrop={false} class="z-50 h-full" closeSidebar={closeDemoSidebar} isOpen={isDemoOpen} nonActiveClass="p-2" params={{ x: -50, duration: 50 }} position="absolute">
            <div class="pb-2 flex items-center justify-start gap-2">
                <Avatar >
                    <UserSolid />
                </Avatar>
                <span class="ml-2 text-lg font-semibold dark:text-white overflow-hidden text-ellipsis whitespace-nowrap">
                    {auth.getUserInfo()?.email || '-'}
                </span>
            </div>

            <hr class="pb-2" />

            <SidebarGroup>
                <SidebarItem href="/" label={m.overview()}>
                    {#snippet icon()}
                        <ChartOutline class={iconClass} />
                    {/snippet}
                </SidebarItem>

                {#each assetClasses.assetClasses as asset (asset.id)}
                    <SidebarItem href={`/assets/${asset.id}`} label={asset.name} class="overflow-hidden text-ellipsis whitespace-nowrap">
                        {#snippet icon()}
                            <DynamicIcon iconName={asset.icon} />
                        {/snippet}
                    </SidebarItem>
                {/each}

                <SidebarDropdownWrapper btnClass="p-2" label={m.configuration()}>
                    {#snippet icon()}
                        <ToolsOutline class={iconClass} />
                    {/snippet}
                    <SidebarItem href="/config/user-groups" label={m.groups()}>
                        {#snippet icon()}
                            <UsersGroupOutline class={iconClass} />
                        {/snippet}
                    </SidebarItem>

                    <SidebarItem href="/config/asset-classes" label={m.asset_classes()}>
                        {#snippet icon()}
                            <AdjustmentsHorizontalSolid class={iconClass} />
                        {/snippet}
                    </SidebarItem>
                </SidebarDropdownWrapper>

                <hr />

                <SidebarItem class="cursor-pointer" label={m.logout()} onclick={() => auth.logout()}>
                    {#snippet icon()}
                        <ArrowRightToBracketOutline class={iconClass} />
                    {/snippet}
                </SidebarItem>
                <SidebarItem href="/about" label={m.about()}>
                    {#snippet icon()}
                        <InfoCircleOutline class={iconClass} />
                    {/snippet}
                </SidebarItem>
            </SidebarGroup>
        </Sidebar>

        <div class="overflow-auto px-4 md:ml-64 h-full overflow-y-auto">
            {@render children()}
        </div>
    </div>
</div>