<script lang="ts">
    import StepIndicator from '$lib/components/StepIndicator.svelte';
    import { onMount } from 'svelte';

    interface KeyShare {
        group: string;
        sk: string;
        public: string;
        details: {
            secret: string;
            verificationKey: string;
            publicKey: string;
            vssCommitment: string[];
            id: number;
            group: number;
        };
    }

    interface KeyData {
        privateKey: string;
        keyShares: KeyShare[];
        verificationKeys: string[];
        savedAt: string;
    }

    interface CommitKeyShareResponse {
        signer: string;
        commitment: string;
    }

    let currentStep = 1;
    let selectedShares: KeyShare[] = [];
    let savedKeyData: KeyData | null = null;
    let error = '';
    let commitResponses: Map<number, CommitKeyShareResponse> = new Map();
    let isCommitting: Map<number, boolean> = new Map();

    function clearStorage() {
        try {
            localStorage.removeItem('keyData');
            savedKeyData = null;
        } catch (err) {
            error = 'Failed to clear saved key data';
            console.error('Error clearing localStorage:', err);
        }
    }

    onMount(() => {
        try {
            const storedData = localStorage.getItem('keyData');
            if (storedData) {
                savedKeyData = JSON.parse(storedData);
            }
        } catch (err) {
            error = 'Failed to load saved key data';
            console.error('Error loading from localStorage:', err);
        }
    });

    function handleShareSelection(share: KeyShare, isSelected: boolean) {
        if (isSelected) {
            selectedShares = [...selectedShares, share];
        } else {
            selectedShares = selectedShares.filter(s => s.details.id !== share.details.id);
        }
    }

    async function handleCommitShare(share: KeyShare) {
        isCommitting.set(share.details.id, true);
        error = '';

        try {
            const response = await fetch('http://localhost:3000/commit-key-share', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                },
                mode: 'cors',
                credentials: 'include',
                body: JSON.stringify({
                    verification_key: savedKeyData?.verificationKeys[0],
                    secret_key_share: share.sk,
                    public_key_shares: savedKeyData?.keyShares.map(share => share.public)
                })
            });

            if (!response.ok) {
                throw new Error(`Failed to commit share: ${response.statusText}`);
            }

            const data: CommitKeyShareResponse = await response.json();
            commitResponses.set(share.details.id, data);
            commitResponses = commitResponses; // Trigger Svelte reactivity
        } catch (err) {
            error = err instanceof Error ? err.message : 'Failed to commit share';
            console.error('Error committing share:', err);
        } finally {
            isCommitting.set(share.details.id, false);
            isCommitting = isCommitting; // Trigger Svelte reactivity
        }
    }
</script>

<main class="container mx-auto p-6">
    <!-- Step 1 -->
    <section class="bg-white shadow-md rounded-lg p-6 max-w-4xl mx-auto mb-8">
        <StepIndicator number={1} title="Select Key Shares" />

        {#if error}
            <div class="text-red-500 mb-4">
                {error}
            </div>
        {/if}

        {#if savedKeyData}
            <div class="mb-4">
                <div class="flex justify-between items-center mb-4">
                    <p class="text-gray-600">
                        Found saved key data from {new Date(savedKeyData.savedAt).toLocaleString()}
                    </p>
                    <button 
                        on:click={clearStorage}
                        class="text-red-500 hover:text-red-600 px-4 py-2 rounded border border-red-500 hover:border-red-600 transition-colors flex items-center gap-2"
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                            <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
                        </svg>
                        Clear Saved Data
                    </button>
                </div>
                <p class="text-sm text-gray-500 mb-4">
                    Select the key shares you want to use for signing:
                </p>

                <div class="grid gap-4">
                    {#each savedKeyData.keyShares as share}
                        <div class="border rounded-lg p-4 bg-gray-50">
                            <label class="flex items-start gap-3">
                                <input 
                                    type="checkbox"
                                    class="mt-1"
                                    on:change={(e) => handleShareSelection(share, e.currentTarget.checked)}
                                />
                                <div>
                                    <h4 class="font-medium">Share {share.details.id}</h4>
                                    <div class="text-sm text-gray-600 mt-1">
                                        <div class="font-mono break-all">
                                            Secret: {share.sk.substring(0, 16)}...
                                        </div>
                                    </div>
                                </div>
                            </label>
                        </div>
                    {/each}
                </div>
            </div>

            <div class="mt-6">
                <button 
                    class="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600 transition-colors text-lg font-medium"
                    disabled={selectedShares.length === 0}
                    on:click={() => currentStep = 2}
                >
                    Continue with Selected Shares
                </button>
            </div>
        {:else}
            <div class="text-gray-600 p-4 bg-gray-100 rounded-lg">
                No saved key data found. Please generate and save key shares first.
            </div>
        {/if}
    </section>

    <!-- Step 2 -->
    <section class="bg-white shadow-md rounded-lg p-6 max-w-4xl mx-auto mb-8 {currentStep < 2 ? 'opacity-50' : ''}">
        <StepIndicator number={2} title="Commit Shares" />
        
        {#if currentStep >= 2}
            {#if error}
                <div class="text-red-500 mb-4">
                    {error}
                </div>
            {/if}

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
                                    on:click={() => handleCommitShare(share)}
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
        {/if}
    </section>
</main> 