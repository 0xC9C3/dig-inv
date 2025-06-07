<script lang="ts">
    import * as Icons from "flowbite-svelte-icons";

    let {iconName = $bindable("PlusOutline")}: {iconName?: string} = $props();

    const getIconWithProps = (iconName: string) => {
        // @ts-expect-error Icons are dynamically accessed
        const IconComponent = Icons[iconName];
        return IconComponent ? IconComponent : Icons["PlusOutline"];
    };
</script>


{@render getIconWithProps(iconName)({})}