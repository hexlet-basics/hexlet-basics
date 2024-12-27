import { PropsWithChildren } from 'react'
import { Link } from '@inertiajs/react'

export default function Layout({ children }: PropsWithChildren) {
  return (
    <main>
      <header className="container mb-4 flex-shrink-0 py-2 border-bottom">
      </header>
      <div>{children}</div>
    </main>
  )
}
