class Web::BlogPosts::LikesController < Web::BlogPosts::ApplicationController
  def create
    session[:blog_post_likes] ||= {}

    key = resource_blog_post.id.to_s
    if !session[:blog_post_likes][key]
      like = resource_blog_post.likes.build user_id: current_user.id
      like.save!
      session[:blog_post_likes][key] = true
      f(:success)
    else
      f(:notice)
    end

    redirect_to blog_post_path(resource_blog_post.slug)
  end
end
