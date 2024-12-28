import type { PropsWithChildren } from 'react'

import Layout from '../../../components/Layout'
import type { Language, LanguageCategory } from '../../../types/serializers';

type Props = PropsWithChildren & {
	language_categories: LanguageCategory[];
	languages: Language[];
};

export default function Index({
  language_categories,
  languages,
}: Props) {

  return (
    <Layout languageCategories={language_categories} languages={languages} >
      <p>jopa</p>
    </Layout>
  )
}
