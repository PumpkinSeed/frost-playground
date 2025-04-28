<script lang="ts">
    export let threshold = 2;
    export let total = 4;
    export let onSplit: (t: number, n: number) => void;
    
    function handleSubmit() {
        onSplit(threshold, total);
    }
</script>

<div class="space-y-4">
    <div class="grid grid-cols-2 gap-4">
        <div>
            <label for="threshold" class="block mb-2 text-lg">Threshold:</label>
            <input
                type="number"
                id="threshold"
                bind:value={threshold}
                min="1"
                class="w-full p-3 border rounded-lg text-lg"
            />
        </div>
        
        <div>
            <label for="total" class="block mb-2 text-lg">Total Shares:</label>
            <input
                type="number"
                id="total"
                bind:value={total}
                min="1"
                class="w-full p-3 border rounded-lg text-lg"
            />
        </div>
    </div>

    <button 
        on:click={handleSubmit}
        class="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600 transition-colors text-lg font-medium w-full"
    >
        Split Key
    </button>
</div> 