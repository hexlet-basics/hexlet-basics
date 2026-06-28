# typed: strict

class Web::BlogPosts::ApplicationController < Web::ApplicationController
  sig { returns(T.untyped) }
  def resource_blog_post
    @resource_blog_post ||= T.let(BlogPost.find_by!(slug: params[:blog_post_id]), T.untyped)
  end
end
