<script lang="ts">
    export let number: number;
    export let title: string;
</script>

<header class="flex items-center mb-6">
    <span class="w-8 h-8 rounded-full bg-blue-500 text-white flex items-center justify-center font-bold">
        {number}
    </span>
    <h2 class="text-xl font-semibold ml-3">{title}</h2>
</header> 