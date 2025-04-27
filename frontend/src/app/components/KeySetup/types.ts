export type KeyOption = 'import' | 'generate' | null

export interface KeySetupProps {
  onComplete: (key: string) => void
}

export interface OptionButtonProps {
  icon: React.ComponentType
  text: string
  onClick: () => void
} 