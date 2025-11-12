# frozen_string_literal: true

require "test_helper"

class Web::Admin::BlogPostsControllerTest < ActionDispatch::IntegrationTest
  setup do
    @user = sign_in_as(:admin)
  end

  def test_index
    get admin_blog_posts_url
    assert_response :success
  end

  def test_new
    get new_admin_blog_post_url
    assert_response :success
  end

  def test_create
    lang = languages(:php)

    attrs = attributes_for(:blog_post, language_id: lang.id)
    post admin_blog_posts_url, params: { blog_post: attrs }
    assert_response :redirect

    assert { lang.blog_posts.find_by slug: attrs[:slug] }
  end

  def test_edit
    blog_post = blog_posts("from-full")

    get edit_admin_blog_post_url(blog_post)
    assert_response :success
  end

  def test_update
    blog_post = blog_posts("from-full")

    patch admin_blog_post_url(blog_post), params: { blog_post: { slug: "mumu" } }
    assert_response :redirect

    blog_post.reload

    assert { blog_post.slug == "mumu" }
  end
end
