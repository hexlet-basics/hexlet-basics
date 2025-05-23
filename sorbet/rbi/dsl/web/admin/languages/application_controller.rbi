# typed: true

# DO NOT EDIT MANUALLY
# This is an autogenerated file for dynamic methods in `Web::Admin::Languages::ApplicationController`.
# Please instead update this file by running `bin/tapioca dsl Web::Admin::Languages::ApplicationController`.


class Web::Admin::Languages::ApplicationController
  sig { returns(HelperProxy) }
  def helpers; end

  module HelperMethods
    include ::InertiaRails::Helper
    include ::Ransack::Helpers::FormHelper
    include ::ViteRails::TagHelpers
    include ::ActionController::Base::HelperMethods
    include ::LanguageCategoriesHelper
    include ::Web::Admin::LanguageLandingPagesHelper
    include ::Pundit::Helper
    include ::ApplicationController::HelperMethods

    sig { returns(T.untyped) }
    def resource_language; end
  end

  class HelperProxy < ::ActionView::Base
    include HelperMethods
  end
end
