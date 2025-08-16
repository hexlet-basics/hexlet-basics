class Web::BlogPosts::ApplicationController < Web::ApplicationController
  def resource_blog_post
    @resource_blog_post ||= BlogPost.find_by!(slug: params[:blog_post_id])
  end
end
