# typed: strict

class Web::BlogPosts::ApplicationController < Web::ApplicationController
  sig { returns(BlogPost) }
  def resource_blog_post
    @resource_blog_post ||= T.let(BlogPost.find_by!(slug: params.expect(:blog_post_id)), T.nilable(BlogPost))
  end
end
