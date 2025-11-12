require "test_helper"

class Api::BlogPostsControllerTest < ActionDispatch::IntegrationTest
  def test_should_get_show
    post = blog_posts("from-full")
    get api_blog_post_url(post.id)
    assert_response :success
  end
end
