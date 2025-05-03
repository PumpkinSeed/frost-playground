<script lang="ts">
    import StepIndicator from '$lib/components/StepIndicator.svelte';
    import KeyOptionSelector from '$lib/components/KeyOptionSelector.svelte';
    import PrivateKeyInput from '$lib/components/PrivateKeyInput.svelte';
    import ActiveKeyDisplay from '$lib/components/ActiveKeyDisplay.svelte';
    import SplitKeyForm from '$lib/components/SplitKeyForm.svelte';
    import KeySharesDisplay from '$lib/components/KeySharesDisplay.svelte';
    import { onMount } from 'svelte';

    let currentStep = 1;
    let hasExistingKey = false;
    let privateKey = '';
    let activeKey = '';
    let isLoading = false;
    let error = '';
    let keyShares: {
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
    }[] = [];
    let verificationKeys: string[] = [];
    let isSaved = false;

    async function generateNewKey() {
        isLoading = true;
        error = '';
        
        try {
            const response = await fetch('http://localhost:3000/private-key', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                },
                mode: 'cors',
                credentials: 'include',
            });

            if (!response.ok) {
                throw new Error(`Failed to generate key: ${response.statusText}`);
            }

            const data = await response.json();
            activeKey = data.private_key_hex;
        } catch (err) {
            error = err instanceof Error ? err.message : 'Failed to generate key';
            console.error('Error generating key:', err);
        } finally {
            isLoading = false;
        }
    }

    async function handleContinue() {
        if (hasExistingKey) {
            activeKey = privateKey;
        } else {
            await generateNewKey();
        }
        
        if (!error) {
            currentStep = 2;
        }
    }

    async function handleSplit(threshold: number, total: number) {
        isLoading = true;
        error = '';
        
        try {
            const response = await fetch('http://localhost:3000/split-private-key', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                },
                mode: 'cors',
                credentials: 'include',
                body: JSON.stringify({
                    private_key_hex: activeKey,
                    threshold,
                    total
                })
            });

            if (!response.ok) {
                throw new Error(`Failed to split key: ${response.statusText}`);
            }

            const data = await response.json();
            keyShares = data.key_shars;
            verificationKeys = data.verification_keys;
            currentStep = 3;
        } catch (err) {
            error = err instanceof Error ? err.message : 'Failed to split key';
            console.error('Error splitting key:', err);
        } finally {
            isLoading = false;
        }
    }

    function saveToLocalStorage() {
        try {
            // Save private key
            localStorage.setItem('main_private_key', activeKey);

            // Save key shares in simplified format
            const simplifiedShares = keyShares.map(share => ({
                secret: share.sk,
                public: share.public
            }));
            localStorage.setItem('main_key_shares', JSON.stringify(simplifiedShares));

            // Save all verification keys
            localStorage.setItem('verification_keys', JSON.stringify(verificationKeys));
            
            // Save first verification key separately
            if (verificationKeys.length > 0) {
                localStorage.setItem('verification_key', verificationKeys[0]);
            }

            isSaved = true;
        } catch (err) {
            error = 'Failed to save key data to local storage';
            console.error('Error saving to localStorage:', err);
        }
    }

    // Update loading from storage if needed
    onMount(() => {
        try {
            const storedPrivateKey = localStorage.getItem('main_private_key');
            const storedShares = localStorage.getItem('main_key_shares');
            const storedVerificationKeys = localStorage.getItem('verification_keys');

            if (storedPrivateKey && storedShares && storedVerificationKeys) {
                activeKey = storedPrivateKey;
                keyShares = JSON.parse(storedShares).map(share => ({
                    ...share,
                    group: '', // Add any required default values
                    details: {
                        id: 0,
                        group: 0,
                        secret: share.secret,
                        verificationKey: '',
                        publicKey: share.public,
                        vssCommitment: []
                    }
                }));
                verificationKeys = JSON.parse(storedVerificationKeys);
            }
        } catch (err) {
            error = 'Failed to load saved data';
            console.error('Error loading from localStorage:', err);
        }
    });
</script>

<main class="container mx-auto p-6">
    <!-- Step 1 -->
    <section class="bg-white shadow-md rounded-lg p-6 max-w-4xl mx-auto mb-8">
        <StepIndicator number={1} title="Private Key Setup" />
        
        <KeyOptionSelector bind:hasExistingKey />

        {#if hasExistingKey}
            <PrivateKeyInput bind:privateKey />
        {/if}

        {#if error}
            <div class="text-red-500 mb-4">
                {error}
            </div>
        {/if}

        <button 
            on:click={handleContinue}
            class="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600 transition-colors text-lg font-medium"
            class:opacity-50={currentStep > 1 || isLoading}
            disabled={currentStep > 1 || isLoading}
        >
            {#if isLoading}
                Generating...
            {:else}
                {hasExistingKey ? 'Continue' : 'Generate New Key'}
            {/if}
        </button>
    </section>

    <!-- Step 2 -->
    <section class="bg-white shadow-md rounded-lg p-6 max-w-4xl mx-auto mb-8 {currentStep < 2 ? 'opacity-50' : ''}">
        <StepIndicator number={2} title="Key Split Configuration" />
        
        <ActiveKeyDisplay {activeKey} />

        {#if currentStep >= 2}
            <div class="mt-6">
                <SplitKeyForm 
                    onSplit={handleSplit}
                />
            </div>
        {/if}
    </section>

    <!-- Step 3 -->
    <section class="bg-white shadow-md rounded-lg p-6 max-w-4xl mx-auto {currentStep < 3 ? 'opacity-50' : ''}">
        <StepIndicator number={3} title="Key Shares" />
        
        {#if currentStep >= 3}
            <KeySharesDisplay {keyShares} {verificationKeys} />
            
            <div class="mt-6 flex items-center gap-4">
                <button 
                    on:click={saveToLocalStorage}
                    class="bg-green-500 text-white px-6 py-3 rounded-lg hover:bg-green-600 transition-colors text-lg font-medium flex items-center gap-2"
                    disabled={isSaved}
                >
                    {#if isSaved}
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                            <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                        </svg>
                        Saved
                    {:else}
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                            <path d="M7.707 10.293a1 1 0 10-1.414 1.414l3 3a1 1 0 001.414 0l3-3a1 1 0 00-1.414-1.414L11 11.586V6h-2v5.586l-1.293-1.293z"/>
                            <path d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z"/>
                        </svg>
                        Save Key Data
                    {/if}
                </button>
                
                {#if isSaved}
                    <span class="text-green-600">Successfully saved to local storage!</span>
                {/if}
            </div>
        {/if}
    </section>
</main>

<style>
    /* Add any component-specific styles here */
</style>
