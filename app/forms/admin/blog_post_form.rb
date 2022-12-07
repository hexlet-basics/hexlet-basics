# frozen_string_literal: true

class Admin::BlogPostForm < BlogPost
  include ActiveFormModel

  permit :creator_id, :language_id, :body, :slug, :locale, :name, :cover
end
