'use client'

import { ChakraProvider, type ChakraProviderProps, defaultSystem } from '@chakra-ui/react'

export default function ClientLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <ChakraProvider value={defaultSystem}>
      {children}
    </ChakraProvider>
  );
}