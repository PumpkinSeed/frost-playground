<script lang="ts">
    export let privateKey: string;
</script>

<div class="mb-6">
    <label for="privateKey" class="block mb-2 text-lg">Enter your private key:</label>
    <input
        type="text"
        id="privateKey"
        bind:value={privateKey}
        class="w-full p-3 border rounded-lg text-lg"
        placeholder="Enter your private key here"
    />
</div> 