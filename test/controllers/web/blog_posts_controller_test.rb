require "test_helper"

class Web::BlogPostsControllerTest < ActionDispatch::IntegrationTest
  test "should get index" do
    get web_blog_posts_index_url
    assert_response :success
  end

  test "should get show" do
    get web_blog_posts_show_url
    assert_response :success
  end
end
