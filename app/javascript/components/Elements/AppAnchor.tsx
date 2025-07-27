import type { Method } from '@inertiajs/core';
import { Link, router } from '@inertiajs/react';
import { Anchor, type AnchorProps } from '@mantine/core';
import { noop } from 'es-toolkit';
import type { PropsWithChildren } from 'react';
import useConfirmation from '@/hooks/useConfirmation';

type AppAnchorProps = AnchorProps &
  React.AriaAttributes &
  PropsWithChildren & {
    /**
     * If true, disables href and renders a pseudo-link (SEO-unfriendly).
     */
    pseudo?: boolean;

    /**
     * If true, opens in a new tab and disables referrer
     */
    external?: boolean;
    withConfirmation?: boolean;

    /**
     * HTTP method for Inertia-powered links
     */
    method?: Method;

    /**
     * Always required
     */
    href: string;
  };

export default function AppAnchor({
  pseudo,
  external,
  method = 'get',
  withConfirmation,
  href,
  ...props
}: AppAnchorProps) {
  const confirmDeleting = useConfirmation();

  if (pseudo) {
    return (
      <Anchor
        role="button"
        rel="nofollow"
        onClick={(e) => {
          // Если middle click или cmd/ctrl + click
          if (e.metaKey || e.ctrlKey || e.button === 1) {
            window.open(href, '_blank', 'noopener,noreferrer');
            return;
          }

          // e.preventDefault();
          router.visit(href);
        }}
        onAuxClick={(e) => {
          if (e.button === 1) {
            window.open(href, '_blank', 'noopener,noreferrer');
          }
        }}
        {...props}
      />
    );
  }

  if (external) {
    return (
      <Anchor
        href={href}
        target="_blank"
        rel="noreferrer nofollow"
        {...props}
      />
    );
  }

  return (
    <Anchor
      component={Link}
      href={href}
      method={method}
      onClick={withConfirmation ? confirmDeleting : noop}
      {...props}
    />
  );
}
