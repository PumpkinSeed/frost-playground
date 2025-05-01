<script lang="ts">
    import type { KeyShare, CommitKeyShareResponse } from '$lib/types';

    export let selectedShares: KeyShare[];
    export let commitResponses: Map<number, CommitKeyShareResponse>;
    export let isCommitting: Map<number, boolean>;
    export let onCommit: (share: KeyShare) => void;
</script>

<div class="space-y-4">
    <h3 class="text-lg font-medium">Selected Shares: {selectedShares.length}</h3>
    
    <div class="grid gap-4">
        {#each selectedShares as share (share.details.id)}
            <div class="border rounded-lg p-4 bg-gray-50">
                <div class="flex justify-between items-start">
                    <div>
                        <h4 class="font-medium">Share {share.details.id}</h4>
                        <div class="text-sm text-gray-600 mt-1">
                            <div class="font-mono break-all">
                                Secret: {share.sk.substring(0, 16)}...
                            </div>
                        </div>
                    </div>
                    
                    <button 
                        on:click={() => onCommit(share)}
                        class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition-colors flex items-center gap-2"
                        disabled={isCommitting.get(share.details.id) || commitResponses.has(share.details.id)}
                    >
                        {#if isCommitting.get(share.details.id)}
                            Committing...
                        {:else if commitResponses.has(share.details.id)}
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                            </svg>
                            Committed
                        {:else}
                            Commit Share
                        {/if}
                    </button>
                </div>

                {#if commitResponses.has(share.details.id)}
                    <div class="mt-4 bg-gray-100 p-3 rounded">
                        <h5 class="font-medium text-sm mb-2">Commit Response:</h5>
                        <div class="text-sm space-y-1">
                            <div>
                                <span class="font-medium">Signer:</span>
                                <span class="font-mono break-all">{commitResponses.get(share.details.id)?.signer}</span>
                            </div>
                            <div>
                                <span class="font-medium">Commitment:</span>
                                <span class="font-mono break-all">{commitResponses.get(share.details.id)?.commitment}</span>
                            </div>
                        </div>
                    </div>
                {/if}
            </div>
        {/each}
    </div>
</div> 