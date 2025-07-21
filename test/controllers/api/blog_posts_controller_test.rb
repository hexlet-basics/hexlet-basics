require "test_helper"

class Api::BlogPostsControllerTest < ActionDispatch::IntegrationTest
  test "should get show" do
    post = blog_posts("from-full")
    get api_blog_post_url(post.id)
    assert_response :success
  end
end
