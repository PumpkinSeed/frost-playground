<script lang="ts">
    export let onSign: (message: string, shares: string[]) => void;
    
    let message = '';
    let shares: string[] = [''];

    function addShare() {
        shares = [...shares, ''];
    }

    function removeShare(index: number) {
        shares = shares.filter((_, i) => i !== index);
    }

    function handleSubmit() {
        const validShares = shares.filter(share => share.trim() !== '');
        if (message.trim() && validShares.length > 0) {
            onSign(message, validShares);
        }
    }
</script>

<div class="space-y-6">
    <div>
        <label for="message" class="block mb-2 text-lg">Message:</label>
        <input
            type="text"
            id="message"
            bind:value={message}
            class="w-full p-3 border rounded-lg text-lg"
            placeholder="Enter message to sign"
        />
    </div>

    <div>
        <label class="block mb-2 text-lg">Key Shares:</label>
        <div class="space-y-3">
            {#each shares as _, i}
                <div class="flex gap-2">
                    <input
                        type="text"
                        bind:value={shares[i]}
                        placeholder="Enter key share"
                        class="flex-1 p-3 border rounded-lg text-lg"
                    />
                    {#if shares.length > 1}
                        <button
                            type="button"
                            on:click={() => removeShare(i)}
                            class="px-4 py-2 text-red-500 hover:text-red-700"
                        >
                            Remove
                        </button>
                    {/if}
                </div>
            {/each}
        </div>

        <button
            type="button"
            on:click={addShare}
            class="mt-3 text-blue-500 hover:text-blue-700"
        >
            + Add another share
        </button>
    </div>

    <button 
        on:click={handleSubmit}
        class="w-full bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600 transition-colors text-lg font-medium"
    >
        Sign Message
    </button>
</div> 