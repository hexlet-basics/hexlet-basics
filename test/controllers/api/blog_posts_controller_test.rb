require "test_helper"

class Api::BlogPostsControllerTest < ActionDispatch::IntegrationTest
  test "should get show" do
    get api_blog_posts_show_url
    assert_response :success
  end
end
