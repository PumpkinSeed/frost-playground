<script lang="ts">
    import type { KeyShare, KeyData } from '$lib/types';

    export let savedKeyData: KeyData;
    export let onShareSelect: (share: KeyShare, isSelected: boolean) => void;
</script>

<div class="grid gap-4">
    {#each savedKeyData.keyShares as share}
        <div class="border rounded-lg p-4 bg-gray-50">
            <label class="flex items-start gap-3">
                <input 
                    type="checkbox"
                    class="mt-1"
                    on:change={(e) => onShareSelect(share, e.currentTarget.checked)}
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