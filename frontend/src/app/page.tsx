'use client'

import {
  Container,
  Box,
  Heading,
  Text,
  VStack,
  Badge,
} from '@chakra-ui/react'
import { useColorModeValue } from '@chakra-ui/color-mode'
import { useState } from 'react'
import { OptionButtons } from './components/KeySetup/OptionButtons'
import { ImportKey } from './components/KeySetup/ImportKey'
import { GenerateKey } from './components/KeySetup/GenerateKey'
import { KeyOption } from './components/KeySetup/types'

export default function Home() {
  const [selectedOption, setSelectedOption] = useState<KeyOption>(null)
  const [privateKey, setPrivateKey] = useState('')

  const bgColor = useColorModeValue('white', 'gray.800')
  const borderColor = useColorModeValue('gray.100', 'gray.700')
  const stepBadgeBg = useColorModeValue('blue.50', 'blue.900')

  const handleGenerate = (key: string) => {
    setPrivateKey(key)
    // Here you can add logic to move to the next step
  }

  return (
    <Box minH="100vh" bg={useColorModeValue('gray.50', 'gray.900')} py={12}>
      <Container maxW="container.md">
        <Box 
          w="full"
          bg={bgColor}
          borderRadius="2xl"
          boxShadow="xl"
          overflow="hidden"
          border="1px"
          borderColor={borderColor}
        >
          <Box 
            bg={stepBadgeBg} 
            p={4} 
            borderBottom="1px" 
            borderColor={borderColor}
          >
            <Badge 
              colorScheme="blue" 
              fontSize="md" 
              px={3} 
              py={1} 
              borderRadius="full"
            >
              Step 1 - Set Up Your Private Key
            </Badge>
          </Box>

          <Box p={8}>
            <VStack gap={8}>
              {!selectedOption ? (
                <VStack gap={6} w="full">
                  <VStack gap={2}>
                    <Heading size="lg" textAlign="center">
                      Welcome to Frost Playground
                    </Heading>
                    <Text 
                      fontSize="lg" 
                      textAlign="center" 
                      color={useColorModeValue('gray.600', 'gray.400')}
                    >
                      Do you have a private key or would you like to generate a new one?
                    </Text>
                  </VStack>

                  <OptionButtons onSelect={setSelectedOption} />
                </VStack>
              ) : (
                <VStack gap={6} w="full">
                  <Heading size="md">
                    {selectedOption === 'import' ? 'Import Your Private Key' : 'Generate New Key'}
                  </Heading>
                  
                  {selectedOption === 'import' ? (
                    <ImportKey
                      value={privateKey}
                      onChange={setPrivateKey}
                      onBack={() => setSelectedOption(null)}
                    />
                  ) : (
                    <GenerateKey
                      onGenerate={handleGenerate}
                      onBack={() => setSelectedOption(null)}
                    />
                  )}
                </VStack>
              )}
            </VStack>
          </Box>
        </Box>
      </Container>
    </Box>
  )
}
