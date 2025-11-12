require "test_helper"

class Web::BlogPosts::LikesControllerTest < ActionDispatch::IntegrationTest
  def test_create
    post = blog_posts("full-python-ru")
    post blog_post_likes_url(post.slug)
    assert_response :redirect
  end

  def test_create_signed_in
    sign_in_as("full")

    post = blog_posts("full-python-ru")
    post blog_post_likes_url(post.slug)
    assert_response :redirect
  end
end
