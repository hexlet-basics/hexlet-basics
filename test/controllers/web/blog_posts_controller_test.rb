require "test_helper"

class Web::BlogPostsControllerTest < ActionDispatch::IntegrationTest
  def test_should_get_index
    get blog_posts_url
    assert_response :success
  end

  def test_should_get_show
    post = blog_posts("from-full")
    get blog_post_url(post.slug)
    assert_response :success
  end
end
