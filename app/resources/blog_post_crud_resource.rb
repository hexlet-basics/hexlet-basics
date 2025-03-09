class BlogPostCrudResource
  include Rails.application.routes.url_helpers
  include Alba::Resource
  include Typelizer::DSL

  typelize_from BlogPost
  root_key :blog_post

  attributes :id, :name, :slug, :description, :body

  # typelize :string, nullable: true
  # attribute :cover_thumb_variant do |post|
  #   rails_representation_url(post.cover.variant(:thumb)) if post.cover.attached?
  # end

  # typelize :string, nullable: true
  # attribute :cover_list_variant do |post|
  #   rails_representation_url(post.cover.variant(:list)) if post.cover.attached?
  # end

  typelize_meta meta: "{ cover_thumb_url: string, state_events: Array<[string, string]>}"
  meta do
    {
      state_events: object.aasm.events_for_select,
      cover_thumb_variant_url: object.cover.attached? ?
        rails_representation_url(object.cover.variant(:thumb)) : nil
    }
  end
end
