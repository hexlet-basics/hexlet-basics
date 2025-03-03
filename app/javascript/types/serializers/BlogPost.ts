// Typelizer digest 585f16c3dfa2484f04b8d22ada862136
//
// DO NOT MODIFY: This file was automatically generated by Typelizer.

type BlogPost = {
  id: number;
  name: string | null;
  slug: string | null;
  description: string | null;
  body: string | null;
  created_at: string;
  state: string | null;
  locale: string | null;
  cover_thumb_variant: string | null;
  cover_list_variant: string | null;
  state_events: Array<[string, string]>;
}

export default BlogPost;
