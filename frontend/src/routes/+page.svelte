<script lang="ts">
    let hasExistingKey = false;
    let privateKey = '';
</script>

<div class="container mx-auto p-6">
    <h1 class="text-2xl font-bold mb-6">Frost Step-by-Step Guide</h1>
    
    <div class="bg-white shadow-md rounded-lg p-6 max-w-4xl mx-auto">
        <div class="flex items-center mb-6">
            <div class="w-8 h-8 rounded-full bg-blue-500 text-white flex items-center justify-center font-bold">
                1
            </div>
            <h2 class="text-xl font-semibold ml-3">Private Key Setup</h2>
        </div>
        
        <div class="mb-6">
            <p class="mb-4 text-lg">Do you have a private key or want to generate a new one?</p>
            
            <div class="flex gap-6 mb-6">
                <label class="flex items-center p-4 border rounded-lg hover:bg-gray-50 cursor-pointer">
                    <input
                        type="radio"
                        bind:group={hasExistingKey}
                        value={true}
                        class="mr-3"
                    />
                    I have a private key
                </label>
                
                <label class="flex items-center p-4 border rounded-lg hover:bg-gray-50 cursor-pointer">
                    <input
                        type="radio"
                        bind:group={hasExistingKey}
                        value={false}
                        class="mr-3"
                    />
                    Generate new key
                </label>
            </div>
        </div>

        {#if hasExistingKey}
            <div class="mb-6">
                <label for="privateKey" class="block mb-2 text-lg">Enter your private key:</label>
                <input
                    type="text"
                    id="privateKey"
                    bind:value={privateKey}
                    class="w-full p-3 border rounded-lg text-lg"
                    placeholder="Enter your private key here"
                />
            </div>
        {/if}

        <button 
            class="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600 transition-colors text-lg font-medium"
        >
            {hasExistingKey ? 'Continue' : 'Generate New Key'}
        </button>
    </div>
</div>
