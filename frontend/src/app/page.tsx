'use client'

import {
  Container,
  Box,
  Heading,
  Text,
  Button,
  VStack,
  Input,
} from '@chakra-ui/react'
import {
  FormControl,
  FormLabel,
  FormHelperText,
} from '@chakra-ui/form-control'
import { useState } from 'react'

export default function Home() {
  const [selectedOption, setSelectedOption] = useState<'import' | 'generate' | null>(null)
  const [privateKey, setPrivateKey] = useState('')

  return (
    <Container maxW="container.md" py={10}>
      <Box p={8} bg="white" borderRadius="lg" boxShadow="lg">
        <VStack gap={6}>
          <Heading textAlign="center">Welcome to Frost Playground</Heading>
          <Text fontSize="lg" textAlign="center">
            Do you have a private key or would you like to generate a new one?
          </Text>

          <VStack gap={4} width="100%">
            <Button
              size="lg"
              colorScheme="blue"
              onClick={() => setSelectedOption('import')}
              width="100%"
            >
              I have a private key
            </Button>
            <Button
              size="lg"
              colorScheme="green"
              onClick={() => setSelectedOption('generate')}
              width="100%"
            >
              Generate new key
            </Button>
          </VStack>
        </VStack>
      </Box>

      {selectedOption && (
        <Box mt={6} p={8} bg="white" borderRadius="lg" boxShadow="lg">
          <VStack align="stretch" spacing={6}>
            <Heading size="md">
              {selectedOption === 'import' ? 'Import Your Private Key' : 'Generate New Key'}
            </Heading>
            
            {selectedOption === 'import' ? (
              <FormControl>
                <FormLabel>Enter your private key</FormLabel>
                <Input
                  type="text"
                  value={privateKey}
                  onChange={(e) => setPrivateKey(e.target.value)}
                  placeholder="Enter your private key here"
                  size="lg"
                />
                <FormHelperText>
                  Please enter your private key carefully
                </FormHelperText>
                <Button
                  mt={4}
                  colorScheme="blue"
                  width="100%"
                  isDisabled={!privateKey.trim()}
                >
                  Continue
                </Button>
              </FormControl>
            ) : (
              <Button colorScheme="green" size="lg">
                Generate Key
              </Button>
            )}
          </VStack>
        </Box>
      )}
    </Container>
  )
}
