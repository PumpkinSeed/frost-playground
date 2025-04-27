'use client'

import {
  Input,
  Button,
  VStack,
} from '@chakra-ui/react'
import { 
    FormLabel, 
    FormControl,
    FormHelperText
} from '@chakra-ui/form-control'
import { useColorModeValue } from '@chakra-ui/color-mode'

interface Props {
  value: string
  onChange: (value: string) => void
  onBack: () => void
}

export function ImportKey({ value, onChange, onBack }: Props) {
  const borderColor = useColorModeValue('gray.100', 'gray.700')

  return (
    <FormControl>
      <FormLabel>Enter your private key</FormLabel>
      <Input
        type="text"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder="Enter your private key here"
        size="lg"
        bg={useColorModeValue('white', 'gray.800')}
        borderColor={borderColor}
        _focus={{
          borderColor: 'blue.500',
          boxShadow: 'outline',
        }}
      />
      <FormHelperText>
        Please enter your private key carefully
      </FormHelperText>
      <VStack gap={4} mt={6}>
        <Button
          colorScheme="blue"
          size="lg"
          w="full"
          disabled={!value.trim()}
        >
          Continue to Step 2
        </Button>
        <Button
          variant="ghost"
          onClick={onBack}
        >
          Go Back
        </Button>
      </VStack>
    </FormControl>
  )
} 