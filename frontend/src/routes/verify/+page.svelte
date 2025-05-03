<script lang="ts">
    import { onMount } from 'svelte';
    import StepIndicator from '$lib/components/StepIndicator.svelte';
    import { API_URL } from '$lib/config';

    interface SignatureVerificationResponse {
        status: string;
        message: string;
    }

    let savedMessage = '';
    let savedSignature = '';
    let savedKeyData: { verificationKeys: string[] } | null = null;
    let verificationResult: SignatureVerificationResponse | null = null;
    let isVerifying = false;
    let error = '';

    onMount(() => {
        try {
            savedMessage = localStorage.getItem('message') || '';
            savedSignature = localStorage.getItem('aggregated_signature') || '';
            const verificationKey = localStorage.getItem('verification_key');
            
            if (verificationKey) {
                savedKeyData = {
                    verificationKeys: [verificationKey]
                };
            }
        } catch (err) {
            error = 'Failed to load saved data';
            console.error('Error loading from localStorage:', err);
        }
    });

    async function verifySignature() {
        if (!savedMessage || !savedSignature || !savedKeyData?.verificationKeys[0]) {
            error = 'Missing required data for verification';
            return;
        }

        isVerifying = true;
        error = '';
        verificationResult = null;

        try {
            const response = await fetch(`${API_URL}/signature-verification`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                },
                mode: 'cors',
                credentials: 'include',
                body: JSON.stringify({
                    message: savedMessage,
                    signature: savedSignature,
                    verification_key: savedKeyData.verificationKeys[0]
                })
            });

            if (!response.ok) {
                throw new Error(`Failed to verify signature: ${response.statusText}`);
            }

            verificationResult = await response.json();
        } catch (err) {
            error = err instanceof Error ? err.message : 'Failed to verify signature';
            console.error('Error verifying signature:', err);
        } finally {
            isVerifying = false;
        }
    }
</script>

<main class="container mx-auto p-6">
    <section class="bg-white shadow-md rounded-lg p-6 max-w-4xl mx-auto">
        <StepIndicator number={1} title="Verify Signature" />

        {#if error}
            <div class="text-red-500 mb-4">
                {error}
            </div>
        {/if}

        {#if !savedMessage || !savedSignature || !savedKeyData}
            <div class="text-gray-600 p-4 bg-gray-100 rounded-lg">
                No saved signature data found. Please sign a message first.
            </div>
        {:else}
            <div class="space-y-6">
                <div>
                    <h3 class="text-lg font-medium mb-2">Message</h3>
                    <div class="bg-gray-50 p-4 rounded-lg break-all font-mono">
                        {savedMessage}
                    </div>
                </div>

                <div>
                    <h3 class="text-lg font-medium mb-2">Signature</h3>
                    <div class="bg-gray-50 p-4 rounded-lg break-all font-mono text-sm">
                        {savedSignature}
                    </div>
                </div>

                <div>
                    <h3 class="text-lg font-medium mb-2">Verification Key</h3>
                    <div class="bg-gray-50 p-4 rounded-lg break-all font-mono text-sm">
                        {savedKeyData.verificationKeys[0]}
                    </div>
                </div>

                <div class="flex items-center gap-4">
                    <button 
                        on:click={verifySignature}
                        class="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600 transition-colors text-lg font-medium flex items-center gap-2"
                        disabled={isVerifying}
                    >
                        {#if isVerifying}
                            <svg class="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                            </svg>
                            Verifying...
                        {:else}
                            Verify Signature
                        {/if}
                    </button>
                </div>

                {#if verificationResult}
                    <div class="mt-6">
                        {#if verificationResult.status === 'valid'}
                            <div class="bg-green-50 border border-green-200 p-4 rounded-lg">
                                <div class="flex items-center gap-2 text-green-700">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
                                    </svg>
                                    <span class="font-medium">Valid Signature</span>
                                </div>
                                <p class="mt-2 text-green-600">{verificationResult.message}</p>
                            </div>
                        {:else}
                            <div class="bg-red-50 border border-red-200 p-4 rounded-lg">
                                <div class="flex items-center gap-2 text-red-700">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
                                    </svg>
                                    <span class="font-medium">Invalid Signature</span>
                                </div>
                                <p class="mt-2 text-red-600">{verificationResult.message}</p>
                            </div>
                        {/if}
                    </div>
                {/if}
            </div>
        {/if}
    </section>
</main> 