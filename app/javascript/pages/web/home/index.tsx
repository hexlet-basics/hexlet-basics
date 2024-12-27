import { useState } from 'react'

import Layout from '../../../components/Layout'

export default function InertiaExample({ one }: { one: string }) {
  const [count, setCount] = useState(0)

  return (
    <Layout>
      <p>{one}</p>
    </Layout>
  )
}
