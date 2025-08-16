require "test_helper"

class Web::BlogPosts::LikesControllerTest < ActionDispatch::IntegrationTest
  test "create" do
    post = blog_posts("full-python-ru")
    post blog_post_likes_url(post.slug)
    assert_response :redirect
  end

  test "create (signed in)" do
    sign_in_as("full")

    post = blog_posts("full-python-ru")
    post blog_post_likes_url(post.slug)
    assert_response :redirect
  end
end
