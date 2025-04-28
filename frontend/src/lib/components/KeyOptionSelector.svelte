<script lang="ts">
    export let hasExistingKey: boolean;
</script>

<div class="mb-6">
    <p class="mb-4 text-lg">Do you have a private key or want to generate a new one?</p>
    
    <fieldset class="flex gap-6 mb-6">
        <label class="flex items-center p-4 border rounded-lg hover:bg-gray-50 cursor-pointer">
            <input
                type="radio"
                bind:group={hasExistingKey}
                value={true}
                class="mr-3"
            />
            <span>I have a private key</span>
        </label>
        
        <label class="flex items-center p-4 border rounded-lg hover:bg-gray-50 cursor-pointer">
            <input
                type="radio"
                bind:group={hasExistingKey}
                value={false}
                class="mr-3"
            />
            <span>Generate new key</span>
        </label>
    </fieldset>
</div> 