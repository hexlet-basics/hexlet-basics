import { Anchor, AnchorProps } from '@mantine/core';
import { Link } from '@inertiajs/react';
import { ComponentPropsWithoutRef } from 'react';

type AppAnchorProps = AnchorProps & ComponentPropsWithoutRef<typeof Link>;

/**
 * A unified Anchor that works with Inertia routing under the hood.
 */
export default function AppAnchor(props: AppAnchorProps) {
  return <Anchor component={Link} {...props} />;
}
