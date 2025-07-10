import { usePage } from '@inertiajs/react';
import { Alert } from '@mantine/core';
import { useState } from 'react';
import { hasObjectKey } from '@/lib/utils';
import type { FlashKey, FlashVariants, SharedProps } from '@/types';
import XssContent from './XssContent';

export default function XFlash() {
  const { flash } = usePage<SharedProps>().props;

  const [dismissed, setDismissed] = useState(false);
  const handleClose = () => setDismissed(true);

  const variants: FlashVariants = {
    error: 'red',
    notice: 'blue',
    success: 'green',
  };

  const flashEntry = (Object.entries(variants) as [FlashKey, string][]).find(
    ([key]) => hasObjectKey(flash, key) && !dismissed,
  );

  if (!flashEntry) return null;

  const [key, color] = flashEntry;
  const value = flash[key];

  return (
    <Alert mb="md" color={color} withCloseButton onClose={handleClose}>
      <XssContent>{value}</XssContent>
    </Alert>
  );
}
