<script lang="ts">
    // Add any necessary imports or logic here
</script>

<div class="container mx-auto px-6 py-8">
    <h1 class="text-3xl font-bold mb-6">FROST Explanation</h1>
    
    <div class="prose max-w-none">
        <h2 class="text-2xl font-semibold mt-6 mb-4">What is FROST?</h2>    
        <p class="mb-4">
            This page provides an explanation of the FROST (Flexible Round-Optimized Schnorr Threshold) protocol.
        </p>

        <h2 class="text-2xl font-semibold mt-6 mb-4">Curve algorithm</h2>
        <p class="mb-4">
            The curve algorithm is a specific equation that is used to generate a curve. For example, the <b>secp256k1</b> curve is one of the most common curves used in cryptocurrencies like Bitcoin.
        </p>
        <p>The <b>secp256k1</b> curve equation with p prime field is:</p>
        <div class="bg-gray-100 p-4 rounded-lg my-4 font-mono text-center">
            y² = x³ + ax + b (mod p)
        </div>
        <p>Equation on a=0, b=7:</p>
        <div class="bg-gray-100 p-4 rounded-lg my-4 font-mono text-center">
            y² = x³ + 7
        </div>
        <p>I don't want to get into the the details of what is a <a href="https://en.wikipedia.org/wiki/Field_(mathematics)" target="_blank">Mathematical Field</a> or Prime Field, but this curve is defined over a the following prime field:</p>
        <div class="bg-gray-100 p-4 rounded-lg my-4 font-mono text-center">
            p = 2²⁵⁶ - 2³² - 2⁹ - 2⁸ - 2⁷ - 2⁶ - 2⁴ - 1
        </div>
    </div>
</div> 