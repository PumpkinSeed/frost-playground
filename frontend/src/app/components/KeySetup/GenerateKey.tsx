'use client'

import { useState } from 'react'
import {
  Button,
  VStack,
  Alert,
} from '@chakra-ui/react'
import { AlertIcon } from '@chakra-ui/alert'

interface Props {
  onGenerate: (key: string) => void
  onBack: () => void
}

export function GenerateKey({ onGenerate, onBack }: Props) {
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleGenerate = async () => {
    try {
      setIsLoading(true)
      setError(null)
      const response = await fetch('http://localhost:3000/private-key')
      if (!response.ok) throw new Error('Failed to generate key')
      const key = await response.text()
      onGenerate(key)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to generate key')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <VStack gap={4} w="full">
      {error && (
        <Alert status="error" borderRadius="md">
          <AlertIcon />
          {error}
        </Alert>
      )}
      <Button
        colorScheme="blue"
        size="lg"
        w="full"
        onClick={handleGenerate}
        isLoading={isLoading}
      >
        Generate New Key
      </Button>
      <Button
        variant="ghost"
        onClick={onBack}
        isDisabled={isLoading}
      >
        Go Back
      </Button>
    </VStack>
  )
} 