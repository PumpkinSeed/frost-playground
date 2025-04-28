<script lang="ts">
    import StepIndicator from '$lib/components/StepIndicator.svelte';
    import SignMessageForm from '$lib/components/SignMessageForm.svelte';

    let isLoading = false;
    let error = '';
    let signature = '';

    async function handleSign(message: string, shares: string[]) {
        isLoading = true;
        error = '';
        signature = '';
        
        try {
            const response = await fetch('http://localhost:3000/sign', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                },
                mode: 'cors',
                credentials: 'include',
                body: JSON.stringify({
                    message,
                    shares
                })
            });

            if (!response.ok) {
                throw new Error(`Failed to sign message: ${response.statusText}`);
            }

            const data = await response.json();
            signature = data.signature;
        } catch (err) {
            error = err instanceof Error ? err.message : 'Failed to sign message';
            console.error('Error signing message:', err);
        } finally {
            isLoading = false;
        }
    }
</script>

<main class="container mx-auto p-6">
    <section class="bg-white shadow-md rounded-lg p-6 max-w-4xl mx-auto mb-8">
        <StepIndicator number={1} title="Sign Message" />
        
        {#if error}
            <div class="text-red-500 mb-4">
                {error}
            </div>
        {/if}

        <SignMessageForm 
            onSign={handleSign}
        />

        {#if isLoading}
            <div class="text-center text-gray-600 mt-4">
                Signing message...
            </div>
        {/if}

        {#if signature}
            <div class="mt-6">
                <h3 class="text-lg font-medium mb-2">Signature:</h3>
                <div class="bg-gray-50 p-4 rounded-lg border break-all font-mono">
                    {signature}
                </div>
            </div>
        {/if}
    </section>
</main> 