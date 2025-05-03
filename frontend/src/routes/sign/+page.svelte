<script lang="ts">
    import StepIndicator from '$lib/components/StepIndicator.svelte';
    import ShareSelector from '$lib/components/ShareSelector.svelte';
    import CommitSharesDisplay from '$lib/components/CommitSharesDisplay.svelte';
    import { onMount } from 'svelte';
    import type { KeyShare, KeyData, CommitKeyShareResponse, AggregateSignaturesResponse } from '$lib/types';
    import { getAllCommitData } from '$lib/utils/commitStorage';
    import { API_URL } from '$lib/config';

    let currentStep = 1;
    let selectedShares: KeyShare[] = [];
    let savedKeyData: KeyData | null = null;
    let error = '';
    let commitResponses: Map<number, CommitKeyShareResponse> = new Map();
    let isCommitting: Map<number, boolean> = new Map();
    let isCommitDataSaved = false;
    let message = '';
    let signatures: Map<string, string> = new Map();
    let isSigningMap: Map<string, boolean> = new Map();
    let signError = '';
    let aggregatedSignature = '';
    let isAggregating = false;
    let aggregationError = '';
    let aggregationSuccess = false;

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
            // Load key shares and verification key
            const storedShares = localStorage.getItem('main_key_shares');
            const verificationKey = localStorage.getItem('verification_key');
            
            if (storedShares && verificationKey) {
                const shares = JSON.parse(storedShares);
                savedKeyData = {
                    privateKey: localStorage.getItem('main_private_key') || '',
                    keyShares: shares.map((share, index) => ({
                        group: '',
                        sk: share.secret,
                        public: share.public,
                        details: {
                            secret: share.secret,
                            verificationKey: verificationKey,
                            publicKey: share.public,
                            vssCommitment: [],
                            id: index + 1,
                            group: 0
                        }
                    })),
                    verificationKeys: [verificationKey],
                    savedAt: new Date().toISOString()
                };
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
            const response = await fetch(`${API_URL}/commit-key-share`, {
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
                    public_key_shares: savedKeyData?.keyShares.map(share => share.public),
                    threshold: selectedShares.length,
                    total: savedKeyData?.keyShares.length
                })
            });

            if (!response.ok) {
                throw new Error(`Failed to commit share: ${response.statusText}`);
            }

            const data: CommitKeyShareResponse = await response.json();
            commitResponses.set(share.details.id, data);
            commitResponses = commitResponses;
        } catch (err) {
            error = err instanceof Error ? err.message : 'Failed to commit share';
            console.error('Error committing share:', err);
        } finally {
            isCommitting.set(share.details.id, false);
            isCommitting = isCommitting;
        }
    }

    function saveCommitResponses() {
        try {
            const commitData = Array.from(commitResponses.entries()).map(([shareId, response]) => {
                const share = selectedShares.find(s => s.details.id === shareId);
                return {
                    secret: share?.sk,
                    signer: response.signer,
                    commitment: response.commitment,
                    savedAt: new Date().toISOString()
                };
            });
            
            // Always overwrite previous data
            localStorage.setItem('commit_data', JSON.stringify(commitData));
            isCommitDataSaved = true;
            error = '';
            return true;
        } catch (err) {
            error = 'Failed to save commit responses';
            console.error('Error saving commit responses:', err);
            return false;
        }
    }

    async function handleSign(signer: string, commitments: string[]) {
        isSigningMap.set(signer, true);
        signError = '';

        try {
            const response = await fetch(`${API_URL}/sign-key-share`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                },
                mode: 'cors',
                credentials: 'include',
                body: JSON.stringify({
                    message,
                    signer,
                    commitments
                })
            });

            if (!response.ok) {
                throw new Error(`Failed to sign message: ${response.statusText}`);
            }

            const data = await response.json();
            signatures.set(signer, data.signature);
            signatures = signatures; // Trigger Svelte reactivity

            // Save signature shares to localStorage
            const signatureShares = Array.from(signatures.entries()).map(([signer, signature]) => ({
                signer,
                signature,
                savedAt: new Date().toISOString()
            }));
            localStorage.setItem('shared_signatures', JSON.stringify(signatureShares));

        } catch (err) {
            signError = err instanceof Error ? err.message : 'Failed to sign message';
            console.error('Error signing message:', err);
        } finally {
            isSigningMap.set(signer, false);
            isSigningMap = isSigningMap; // Trigger Svelte reactivity
        }
    }

    async function handleAggregateSignatures() {
        isAggregating = true;
        aggregationError = '';

        try {
            const signatureArray = Array.from(signatures.values());
            const response = await fetch(`${API_URL}/aggregate-signatures`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                },
                mode: 'cors',
                credentials: 'include',
                body: JSON.stringify({
                    signatures: signatureArray,
                    message,
                    commitments: commitData.map(d => d.commitment),
                    verification_key: savedKeyData?.verificationKeys[0],
                    public_key_shares: savedKeyData?.keyShares.map(share => share.public),
                    threshold: commitData.length,
                    total: savedKeyData?.keyShares.length
                })
            });

            if (!response.ok) {
                throw new Error(`Failed to aggregate signatures: ${response.statusText}`);
            }

            const data: AggregateSignaturesResponse = await response.json();
            aggregatedSignature = data.signature;

            // Save message and aggregated signature to localStorage
            localStorage.setItem('message', message);
            localStorage.setItem('aggregated_signature', data.signature);

            // Show success message
            aggregationSuccess = true;
        } catch (err) {
            aggregationError = err instanceof Error ? err.message : 'Failed to aggregate signatures';
            console.error('Error aggregating signatures:', err);
        } finally {
            isAggregating = false;
        }
    }

    $: commitData = isCommitDataSaved ? getAllCommitData() : [];
    $: currentStep = isCommitDataSaved ? (signatures.size === commitData.length && commitData.length > 0 ? 4 : 3) : 2;
    $: canAggregate = signatures.size === commitData.length && commitData.length > 0;
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

                <ShareSelector {savedKeyData} onShareSelect={handleShareSelection} />
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

            <CommitSharesDisplay 
                {selectedShares}
                {commitResponses}
                {isCommitting}
                onCommit={handleCommitShare}
            />

            {#if commitResponses.size > 0}
                <div class="mt-6 flex items-center gap-4">
                    <button 
                        on:click={() => {
                            if (saveCommitResponses()) {
                                error = '';
                            }
                        }}
                        class="bg-green-500 text-white px-6 py-3 rounded-lg hover:bg-green-600 transition-colors text-lg font-medium flex items-center gap-2"
                        disabled={isCommitDataSaved}
                    >
                        {#if isCommitDataSaved}
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                            </svg>
                            Saved
                        {:else}
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                <path d="M7.707 10.293a1 1 0 10-1.414 1.414l3 3a1 1 0 001.414 0l3-3a1 1 0 00-1.414-1.414L11 11.586V6h-2v5.586l-1.293-1.293z"/>
                                <path d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z"/>
                            </svg>
                            Save Commit Data
                        {/if}
                    </button>
                    
                    {#if isCommitDataSaved}
                        <span class="text-green-600 flex items-center gap-2">
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
                            </svg>
                            Successfully saved commit data!
                        </span>
                    {/if}
                </div>
            {/if}
        {/if}
    </section>

    <!-- Step 3 -->
    <section class="bg-white shadow-md rounded-lg p-6 max-w-4xl mx-auto mb-8 {currentStep < 3 ? 'opacity-50' : ''}">
        <StepIndicator number={3} title="Sign Message" />
        
        {#if currentStep >= 3}
            {#if signError}
                <div class="text-red-500 mb-4">
                    {signError}
                </div>
            {/if}

            <div class="space-y-6">
                <div>
                    <h3 class="text-lg font-medium mb-4">Enter Message to Sign</h3>
                    <div class="mb-6">
                        <textarea
                            id="message"
                            bind:value={message}
                            class="w-full p-3 border rounded-lg"
                            rows="3"
                            placeholder="Enter your message here..."
                        ></textarea>
                    </div>

                    <h3 class="text-lg font-medium mb-4">Committed Shares</h3>
                    <div class="grid gap-4">
                        {#each commitData as data}
                            <div class="border rounded-lg p-4 bg-gray-50">
                                <div class="flex justify-between items-start gap-4">
                                    <div class="grid gap-2 flex-grow">
                                        <div>
                                            <span class="font-medium">Secret:</span>
                                            <span class="font-mono break-all">{data.secret}</span>
                                        </div>
                                        <div>
                                            <span class="font-medium">Signer:</span>
                                            <span class="font-mono break-all">{data.signer}</span>
                                        </div>
                                        <div>
                                            <span class="font-medium">Commitment:</span>
                                            <span class="font-mono break-all">{data.commitment}</span>
                                        </div>
                                        {#if signatures.has(data.signer)}
                                            <div>
                                                <span class="font-medium">Signature:</span>
                                                <span class="font-mono break-all">{signatures.get(data.signer)}</span>
                                            </div>
                                        {/if}
                                    </div>

                                    <button 
                                        on:click={() => handleSign(data.signer, commitData.map(d => d.commitment))}
                                        class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition-colors flex items-center gap-2 h-fit min-w-[120px] justify-center"
                                        disabled={!message || isSigningMap.get(data.signer) || signatures.has(data.signer)}
                                    >
                                        {#if isSigningMap.get(data.signer)}
                                            Signing...
                                        {:else if signatures.has(data.signer)}
                                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                                            </svg>
                                            Signed
                                        {:else}
                                            Sign
                                        {/if}
                                    </button>
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>
            </div>
        {/if}
    </section>

    <!-- Step 4 -->
    <section class="bg-white shadow-md rounded-lg p-6 max-w-4xl mx-auto mb-8 {currentStep < 4 ? 'opacity-50' : ''}">
        <StepIndicator number={4} title="Aggregate Signatures" />
        
        {#if currentStep >= 4}
            {#if aggregationError}
                <div class="text-red-500 mb-4">
                    {aggregationError}
                </div>
            {/if}

            <div class="space-y-6">
                <div>
                    <h3 class="text-lg font-medium mb-4">Message</h3>
                    <div class="bg-gray-50 p-4 rounded-lg break-all font-mono">
                        {message}
                    </div>
                </div>

                <div>
                    <h3 class="text-lg font-medium mb-4">Individual Signatures</h3>
                    <div class="grid gap-4">
                        {#each Array.from(signatures.entries()) as [signer, signature]}
                            <div class="bg-gray-50 p-4 rounded-lg">
                                <div class="font-medium mb-2">Signer: {signer.substring(0, 16)}...</div>
                                <div class="font-mono break-all text-sm">{signature}</div>
                            </div>
                        {/each}
                    </div>
                </div>

                <div class="flex items-center gap-4">
                    <button 
                        on:click={handleAggregateSignatures}
                        class="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600 transition-colors text-lg font-medium flex items-center gap-2"
                        disabled={!canAggregate || isAggregating || aggregatedSignature}
                    >
                        {#if isAggregating}
                            <svg class="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                            </svg>
                            Aggregating...
                        {:else if aggregatedSignature}
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                            </svg>
                            Aggregated
                        {:else}
                            Aggregate Signatures
                        {/if}
                    </button>

                    {#if aggregationSuccess}
                        <span class="text-green-600 flex items-center gap-2">
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
                            </svg>
                            Successfully saved message and signature!
                        </span>
                    {/if}
                </div>

                {#if aggregatedSignature}
                    <div>
                        <h3 class="text-lg font-medium mb-4">Aggregated Signature</h3>
                        <div class="bg-green-50 border border-green-200 p-4 rounded-lg">
                            <div class="font-mono break-all">{aggregatedSignature}</div>
                        </div>
                    </div>
                {/if}
            </div>
        {/if}
    </section>
</main> 