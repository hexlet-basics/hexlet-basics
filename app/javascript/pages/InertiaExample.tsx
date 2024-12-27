import { useState } from 'react'

import Layout from '../components/Layout'

export default function InertiaExample({ name }: { name: string }) {
  const [count, setCount] = useState(0)

  return (
    <Layout>
      <p>JOPA</p>
    </Layout>
  )
}
