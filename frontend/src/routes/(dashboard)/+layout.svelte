<script lang="ts">
    import '../../app.css';
    import {
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
        ArrowRightToBracketOutline,
        ChartOutline,
        EditSolid,
        GridSolid,
        MailBoxSolid,
        SearchOutline,
        ShoppingBagSolid,
        UserSolid
    } from "flowbite-svelte-icons";
    import {fade} from "svelte/transition";
    import {page} from "$app/state";

    let activeUrl = $state(page.url.pathname);
    const spanClass = "flex-1 ms-3 whitespace-nowrap";
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
        {#snippet children({ hidden, toggle, NavContainer })}
            <SidebarButton class="mb-2" onclick={sidebar.toggle} />
            <NavBrand href="/">
                <img src="/favicon.png" alt="Logo" class="h-8 mr-3" />
                <h1 class="self-center text-xl font-semibold whitespace-nowrap dark:text-white">dig-inv</h1>
            </NavBrand>
            <div class="flex md:order-2">
                <div class="mr-3">
                    <DarkMode />
                </div>
                <ToolbarButton class="block md:hidden" onclick={toggle}>
                    <SearchOutline class="h-5 w-5 text-gray-500 dark:text-gray-400" />
                </ToolbarButton>
                <div class="hidden md:block">
                    <Search size="md" class="ms-auto" placeholder="Search..." />
                </div>
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
            <SidebarGroup>
                <SidebarItem label="Dashboard">
                    {#snippet icon()}
                        <ChartOutline class="h-5 w-5 text-gray-500 transition duration-75 group-hover:text-gray-900 dark:text-gray-400 dark:group-hover:text-white" />
                    {/snippet}
                </SidebarItem>
                <SidebarDropdownWrapper btnClass="p-2" label="E-commerce">
                    {#snippet icon()}
                        <ShoppingBagSolid class="h-5 w-5 text-gray-500 transition duration-75 group-hover:text-gray-900 dark:text-gray-400 dark:group-hover:text-white" />
                    {/snippet}
                    <SidebarItem href="/docs/components/sidebar" label="Sidebar" />
                    <SidebarItem label="Billing" />
                    <SidebarItem label="Invoice" />
                </SidebarDropdownWrapper>
                <SidebarItem href="/" label="Kanban" {spanClass}>
                    {#snippet icon()}
                        <GridSolid class="h-5 w-5 text-gray-500 transition duration-75 group-hover:text-gray-900 dark:text-gray-400 dark:group-hover:text-white" />
                    {/snippet}
                    {#snippet subtext()}
                        <span class="ms-3 inline-flex items-center justify-center rounded-full bg-gray-200 px-2 text-sm font-medium text-gray-800 dark:bg-gray-700 dark:text-gray-300">Pro</span>
                    {/snippet}
                </SidebarItem>
                <SidebarItem href="/" label="Inbox" {spanClass}>
                    {#snippet icon()}
                        <MailBoxSolid class="h-5 w-5 text-gray-500 transition duration-75 group-hover:text-gray-900 dark:text-gray-400 dark:group-hover:text-white" />
                    {/snippet}
                    {#snippet subtext()}
                        <span class="bg-primary-200 text-primary-600 dark:bg-primary-900 dark:text-primary-200 ms-3 inline-flex h-3 w-3 items-center justify-center rounded-full p-3 text-sm font-medium">3</span>
                    {/snippet}
                </SidebarItem>
                <SidebarItem label="Users">
                    {#snippet icon()}
                        <UserSolid class="h-5 w-5 text-gray-500 transition duration-75 group-hover:text-gray-900 dark:text-gray-400 dark:group-hover:text-white" />
                    {/snippet}
                </SidebarItem>
                <SidebarItem label="Sign In">
                    {#snippet icon()}
                        <ArrowRightToBracketOutline class="h-5 w-5 text-gray-500 transition duration-75 group-hover:text-gray-900 dark:text-gray-400 dark:group-hover:text-white" />
                    {/snippet}
                </SidebarItem>
                <SidebarItem label="Sign Up">
                    {#snippet icon()}
                        <EditSolid class="h-5 w-5 text-gray-500 transition duration-75 group-hover:text-gray-900 dark:text-gray-400 dark:group-hover:text-white" />
                    {/snippet}
                </SidebarItem>
            </SidebarGroup>
        </Sidebar>

        <div class="h-96 overflow-auto px-4 md:ml-64">
            {@render children()}
        </div>
    </div>
</div>