<script lang="ts">
    import StepIndicator from '$lib/components/StepIndicator.svelte';
    import KeyOptionSelector from '$lib/components/KeyOptionSelector.svelte';
    import PrivateKeyInput from '$lib/components/PrivateKeyInput.svelte';
    import ActiveKeyDisplay from '$lib/components/ActiveKeyDisplay.svelte';

    let currentStep = 1;
    let hasExistingKey = false;
    let privateKey = '';
    let activeKey = '';
    let isLoading = false;
    let error = '';

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
    <section class="bg-white shadow-md rounded-lg p-6 max-w-4xl mx-auto {currentStep < 2 ? 'opacity-50' : ''}">
        <StepIndicator number={2} title="Key Confirmation" />
        
        <ActiveKeyDisplay {activeKey} />

        <button 
            class="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600 transition-colors text-lg font-medium mt-6"
            disabled={currentStep < 2}
        >
            Continue to Next Step
        </button>
    </section>
</main>

<style>
    /* Add any component-specific styles here */
</style>
