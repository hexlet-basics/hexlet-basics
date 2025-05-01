require "test_helper"
require "capybara/cuprite"

class ApplicationSystemTestCase < ActionDispatch::SystemTestCase
  driven_by :cuprite

  # https://guides.rubyonrails.org/active_storage_overview.html#serving-files
  # parallelize_setup do |i|
  #   ActiveStorage::Blob.service.root = "#{ActiveStorage::Blob.service.root}-#{i}"
  # end

  # def after_teardown
  #   super
  #   FileUtils.rm_rf(ActiveStorage::Blob.service.root)
  # end

  def sign_in_as(name)
    user = users(name)

    visit new_session_url
    email_element = find('[data-testid="email"]')
    password_element = find('[data-testid="password"]')
    log_in_element = find('[data-testid="submit"]')

    email_element.fill_in with: user.email
    password_element.fill_in with: "password"
    log_in_element.click
  end
end
