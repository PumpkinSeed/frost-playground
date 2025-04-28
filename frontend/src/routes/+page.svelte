<script lang="ts">
    import StepIndicator from '$lib/components/StepIndicator.svelte';
    import KeyOptionSelector from '$lib/components/KeyOptionSelector.svelte';
    import PrivateKeyInput from '$lib/components/PrivateKeyInput.svelte';
    import ActiveKeyDisplay from '$lib/components/ActiveKeyDisplay.svelte';
    import SplitKeyForm from '$lib/components/SplitKeyForm.svelte';
    import KeySharesDisplay from '$lib/components/KeySharesDisplay.svelte';

    let currentStep = 1;
    let hasExistingKey = false;
    let privateKey = '';
    let activeKey = '';
    let isLoading = false;
    let error = '';
    let keyShares = [];

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
            currentStep = 3;
        } catch (err) {
            error = err instanceof Error ? err.message : 'Failed to split key';
            console.error('Error splitting key:', err);
        } finally {
            isLoading = false;
        }
    }
</script>

<main class="container mx-auto p-6">
    <h1 class="text-2xl font-bold mb-6">Frost Step-by-Step Guide</h1>
    
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
            <KeySharesDisplay {keyShares} />
        {/if}
    </section>
</main>

<style>
    /* Add any component-specific styles here */
</style>
