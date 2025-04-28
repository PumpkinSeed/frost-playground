<script lang="ts">
    export let activeKey: string;
</script>

<div class="bg-gray-50 p-4 rounded-lg border">
    <p class="text-lg mb-2 font-medium">Active Private Key:</p>
    <code class="block bg-white p-3 rounded border font-mono break-all">
        {activeKey || 'No key selected'}
    </code>
</div> 