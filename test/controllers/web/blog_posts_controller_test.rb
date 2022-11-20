require 'test_helper'

class Web::BlogPostsControllerTest < ActionDispatch::IntegrationTest
  test 'should get index' do
    get blog_posts_url
    assert_response :success
  end

  test 'should get show' do
    post = blog_posts(:js)
    get blog_post_url(id: post.slug)
    assert_response :success
  end
end
