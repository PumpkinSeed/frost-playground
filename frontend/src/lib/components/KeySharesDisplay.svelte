<script lang="ts">
    export let keyShares: {
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
    }[];
    export let verificationKeys: string[];
</script>

<div class="space-y-6">
    <h3 class="text-xl font-semibold mb-4">Key Shares</h3>
    
    <div class="grid gap-4">
        {#each keyShares as share}
            <div class="border rounded-lg p-4 bg-gray-50">
                <h4 class="font-medium mb-2">Share {share.details.id}</h4>
                <div class="grid gap-2 text-sm">
                    <div>
                        <span class="font-medium">Secret:</span>
                        <span class="font-mono break-all">{share.sk}</span>
                    </div>
                    <div>
                        <span class="font-medium">Public Key:</span>
                        <span class="font-mono break-all">{share.public}</span>
                    </div>
                </div>
            </div>
        {/each}
    </div>

    <div class="mt-8">
        <h3 class="text-xl font-semibold mb-4">Verification Keys</h3>
        <div class="border rounded-lg p-4 bg-gray-50">
            <div class="grid gap-2 text-sm">
                {#each verificationKeys as key, i}
                    <div>
                        <span class="font-medium">Key {i + 1}:</span>
                        <span class="font-mono break-all">{key}</span>
                    </div>
                {/each}
            </div>
        </div>
    </div>
</div> 