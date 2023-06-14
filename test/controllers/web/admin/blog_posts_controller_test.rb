# frozen_string_literal: true

require 'test_helper'

class Web::Admin::BlogPostsControllerTest < ActionDispatch::IntegrationTest
  setup do
    @user = sign_in_as(:admin)
  end

  test 'index' do
    get admin_blog_posts_url
    assert_response :success
  end

  test 'new' do
    get new_admin_blog_post_url
    assert_response :success
  end

  test 'create' do
    lang = languages(:php)

    attrs = attributes_for(:blog_post, language_id: lang.id)
    post admin_blog_posts_url, params: { admin_blog_post_form: attrs }
    assert_response :redirect

    assert { lang.blog_posts.find_by slug: attrs[:slug] }
  end

  test 'edit' do
    blog_post = blog_posts('from-full')

    get edit_admin_blog_post_url(blog_post)
    assert_response :success
  end

  test 'update' do
    blog_post = blog_posts('from-full')

    patch admin_blog_post_url(blog_post), params: { admin_blog_post_form: { slug: 'mumu' } }
    assert_response :redirect

    blog_post.reload
    assert { blog_post.slug == 'mumu' }
  end
end
