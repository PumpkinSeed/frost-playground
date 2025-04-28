<script lang="ts">
    import { page } from '$app/stores';

    interface MenuItem {
        text: string;
        href: string;
    }

    const menuItems: MenuItem[] = [
        { text: 'Setup', href: '/' },
        { text: 'Sign', href: '/sign' },
        { text: 'Verify Signature', href: '/verify' }
    ];

    $: currentPath = $page.url.pathname;
    $: isActive = (href: string) => {
        if (href === '/') {
            return currentPath === '/';
        }
        return currentPath.startsWith(href);
    };
</script>

<nav class="bg-white shadow-md">
    <div class="container mx-auto px-6">
        <div class="flex items-center justify-between h-16">
            <div class="flex-shrink-0">
                <span class="text-xl font-bold">Frost Step-by-Step Guide</span>
            </div>
            
            <div class="flex items-center space-x-4">
                <div class="hidden md:block">
                    <div class="flex space-x-4">
                        {#each menuItems as item}
                            <a
                                href={item.href}
                                class="px-3 py-2 rounded-md text-sm font-medium {isActive(item.href)
                                    ? 'bg-blue-500 text-white' 
                                    : 'text-gray-700 hover:bg-gray-100'}"
                            >
                                {item.text}
                            </a>
                        {/each}
                    </div>
                </div>

                <a 
                    href="https://github.com/PumpkinSeed/frost-playground"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="text-gray-700 hover:text-gray-900 transition-colors"
                    aria-label="GitHub Repository"
                >
                    <svg 
                        viewBox="0 0 24 24" 
                        class="w-6 h-6 fill-current"
                    >
                        <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/>
                    </svg>
                </a>
            </div>
        </div>
    </div>
</nav> 