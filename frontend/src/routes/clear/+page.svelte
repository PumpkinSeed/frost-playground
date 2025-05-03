<script lang="ts">
    import StepIndicator from '$lib/components/StepIndicator.svelte';
    import { onMount } from 'svelte';

    interface StorageData {
        key: string;
        value: string;
        displayValue: any;
        parsedValue?: any;
    }

    let storageData: StorageData[] = [];
    let error = '';

    function parseValue(value: string, key: string): any {
        try {
            if (!value) return null;
            
            // Try to parse as JSON first
            try {
                return JSON.parse(value);
            } catch {
                // If not JSON, return as is
                return value;
            }
        } catch (err) {
            console.error(`Error parsing ${key}:`, err);
            return value;
        }
    }

    function formatDate(dateStr: string): string {
        try {
            return new Date(dateStr).toLocaleString();
        } catch {
            return dateStr;
        }
    }

    function loadStorageData() {
        try {
            const keys = [
                'main_private_key',
                'main_key_shares',
                'verification_keys',
                'verification_key',
                'commit_data',
                'message',
                'aggregated_signature'
            ];

            storageData = keys.map(key => {
                const value = localStorage.getItem(key) || '';
                const parsedValue = parseValue(value, key);
                return {
                    key,
                    value,
                    displayValue: value || 'Not set',
                    parsedValue
                };
            });
        } catch (err) {
            error = 'Failed to load storage data';
            console.error('Error loading storage data:', err);
        }
    }

    function clearAllData() {
        try {
            localStorage.removeItem('main_private_key');
            localStorage.removeItem('main_key_shares');
            localStorage.removeItem('verification_keys');
            localStorage.removeItem('verification_key');
            localStorage.removeItem('commit_data');
            localStorage.removeItem('message');
            localStorage.removeItem('aggregated_signature');
            loadStorageData();
        } catch (err) {
            error = 'Failed to clear storage data';
            console.error('Error clearing storage:', err);
        }
    }

    onMount(loadStorageData);
</script>

<main class="container mx-auto p-6">
    <section class="bg-white shadow-md rounded-lg p-6 max-w-4xl mx-auto">
        <StepIndicator number={1} title="Stored Data Overview" />

        {#if error}
            <div class="text-red-500 mb-4">
                {error}
            </div>
        {/if}

        <div class="space-y-6">
            <div class="flex justify-between items-center">
                <h3 class="text-lg font-medium">Local Storage Data</h3>
                <button 
                    on:click={clearAllData}
                    class="text-red-500 hover:text-red-600 px-4 py-2 rounded 
                           border border-red-500 hover:border-red-600 
                           transition-colors flex items-center gap-2"
                >
                    <svg 
                        xmlns="http://www.w3.org/2000/svg" 
                        class="h-5 w-5" 
                        viewBox="0 0 20 20" 
                        fill="currentColor"
                    >
                        <path 
                            fill-rule="evenodd" 
                            d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" 
                            clip-rule="evenodd" 
                        />
                    </svg>
                    Clear All Data
                </button>
            </div>

            <div class="grid gap-4 max-h-[70vh] overflow-y-auto pr-2">
                {#each storageData as data}
                    <div class="border rounded-lg p-4 bg-gray-50">
                        <h4 class="font-medium text-lg mb-2">{data.key}</h4>
                        <div class="bg-white p-3 rounded border" style="word-wrap: break-word; max-width: 50vh;">
                            {#if data.key === 'commit_data' && Array.isArray(data.parsedValue)}
                                <div class="space-y-2">
                                    {#each data.parsedValue as commit}
                                        <div class="p-2 border rounded" style="word-wrap: break-word; max-width: 50vh; padding: 10px; margin-bottom: 10px;">
                                            <div style="margin-bottom: 10px;"><span class="font-semibold">Signer:</span> {commit.signer}</div>
                                            <div style="margin-bottom: 10px;"><span class="font-semibold">Secret:</span> {commit.secret}</div>
                                            <div><span class="font-semibold">Saved:</span> {formatDate(commit.savedAt)}</div>
                                        </div>
                                    {/each}
                                </div>
                            {:else if data.key === 'main_key_shares' && Array.isArray(data.parsedValue)}
                                <div class="space-y-2">
                                    {#each data.parsedValue as share}
                                        <div class="p-2 border rounded" style="word-wrap: break-word; max-width: 50vh; padding: 10px; margin-bottom: 10px;">
                                            <div style="margin-bottom: 10px;"><span class="font-semibold">Secret:</span> {share.secret}</div>
                                            <div><span class="font-semibold">Public Key:</span> {share.public}</div>
                                        </div>
                                    {/each}
                                </div>
                            {:else if Array.isArray(data.parsedValue)}
                                <div class="overflow-x-auto">
                                    <pre class="text-sm font-mono whitespace-pre-wrap break-words" style="max-width: 50vh;">
                                        {JSON.stringify(data.parsedValue, null, 2)}
                                    </pre>
                                </div>
                            {:else}
                                <div class="font-mono text-sm break-words" style="word-wrap: break-word; max-width: 50vh;">
                                    {data.displayValue}
                                </div>
                            {/if}
                        </div>
                    </div>
                {/each}
            </div>
        </div>
    </section>
</main>