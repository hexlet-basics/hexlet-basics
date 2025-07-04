import { Anchor, AnchorProps } from '@mantine/core';
import { Link } from '@inertiajs/react';
import { PropsWithChildren } from 'react';
import useConfirmation from '@/hooks/useConfirmation';
import { noop } from 'es-toolkit';
import { Method } from '@inertiajs/core';

type AppAnchorProps = AnchorProps & PropsWithChildren & {
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
  method = "get",
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
        onClick={() => {
          window.location.href = href;
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

  // default: Inertia-powered
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

