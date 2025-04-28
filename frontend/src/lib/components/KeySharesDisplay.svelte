<script lang="ts">
    interface KeyShare {
        group: string;
        sk: string;
        public: string;
    }

    export let keyShares: KeyShare[] = [];
</script>

<div class="space-y-4">
    <h3 class="text-lg font-medium">Generated Key Shares:</h3>
    
    <div class="space-y-4">
        {#each keyShares as share, i}
            <div class="border rounded-lg p-4 bg-gray-50">
                <h4 class="font-medium mb-2">Share {i + 1}</h4>
                <div class="space-y-2">
                    <div>
                        <span class="font-medium">Group:</span>
                        <code class="block bg-white p-2 rounded border mt-1 break-all">
                            {share.group}
                        </code>
                    </div>
                    <div>
                        <span class="font-medium">SK:</span>
                        <code class="block bg-white p-2 rounded border mt-1 break-all">
                            {share.sk}
                        </code>
                    </div>
                    <div>
                        <span class="font-medium">Public:</span>
                        <code class="block bg-white p-2 rounded border mt-1 break-all">
                            {share.public}
                        </code>
                    </div>
                </div>
            </div>
        {/each}
    </div>
</div> 