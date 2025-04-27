'use client'

import { VStack, Button, Icon, Text } from '@chakra-ui/react'
import { FaKey, FaPlus } from 'react-icons/fa'
import { KeyOption } from './types'

interface Props {
  onSelect: (option: KeyOption) => void
}

export function OptionButtons({ onSelect }: Props) {
  return (
    <VStack gap={4} w="full">
      <Button
        size="lg"
        colorScheme="blue"
        onClick={() => onSelect('import')}
        w="full"
        h={20}
        display="flex"
        flexDirection="column"
        gap={2}
        variant="outline"
        _hover={{ transform: 'translateY(-2px)', boxShadow: 'lg' }}
        transition="all 0.2s"
      >
        <Icon as={FaKey} boxSize={6} />
        <Text>I have a private key</Text>
      </Button>
      <Button
        size="lg"
        colorScheme="blue"
        onClick={() => onSelect('generate')}
        w="full"
        h={20}
        display="flex"
        flexDirection="column"
        gap={2}
        variant="outline"
        _hover={{ transform: 'translateY(-2px)', boxShadow: 'lg' }}
        transition="all 0.2s"
      >
        <Icon as={FaPlus} boxSize={6} />
        <Text>Generate new key</Text>
      </Button>
    </VStack>
  )
} 