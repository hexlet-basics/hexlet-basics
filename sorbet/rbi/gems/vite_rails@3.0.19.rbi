# typed: true

# DO NOT EDIT MANUALLY
# This is an autogenerated file for types exported from the `vite_rails` gem.
# Please instead update this file by running `bin/tapioca gem vite_rails`.


# source://vite_rails//lib/vite_rails/version.rb#3
module ViteRails; end

# source://vite_rails//lib/vite_rails/cli.rb#5
module ViteRails::CLI; end

# source://vite_rails//lib/vite_rails/cli.rb#8
module ViteRails::CLI::Build
  # source://vite_rails//lib/vite_rails/cli.rb#9
  def call(**options); end

  # Internal: Attempts to initialize the Rails application.
  #
  # source://vite_rails//lib/vite_rails/cli.rb#15
  def ensure_rails_init; end
end

# Internal: Extends the base installation script from Vite Ruby to work for a
# typical Rails app.
#
# source://vite_rails//lib/vite_rails/cli.rb#24
module ViteRails::CLI::Install
  # Override: Create a sample JS file and attempt to inject it in an HTML template.
  #
  # source://vite_rails//lib/vite_rails/cli.rb#60
  def install_sample_files; end

  # Override: Setup a typical apps/web Hanami app to use Vite.
  #
  # source://vite_rails//lib/vite_rails/cli.rb#28
  def setup_app_files; end

  # Internal: Configure CSP rules that allow to load @vite/client correctly.
  #
  # source://vite_rails//lib/vite_rails/cli.rb#38
  def setup_content_security_policy(csp_file); end
end

# source://vite_rails//lib/vite_rails/cli.rb#25
ViteRails::CLI::Install::RAILS_TEMPLATES = T.let(T.unsafe(nil), Pathname)

# source://vite_rails//lib/vite_rails/config.rb#3
module ViteRails::Config
  private

  # Override: Default values for a Rails application.
  #
  # source://vite_rails//lib/vite_rails/config.rb#7
  def config_defaults; end
end

# source://vite_rails//lib/vite_rails/engine.rb#5
class ViteRails::Engine < ::Rails::Engine
  class << self
    private

    # source://activesupport/8.0.2/lib/active_support/class_attribute.rb#15
    def __class_attr___callbacks; end

    # source://activesupport/8.0.2/lib/active_support/class_attribute.rb#17
    def __class_attr___callbacks=(new_value); end
  end
end

# Public: Allows to render HTML tags for scripts and styles processed by Vite.
#
# source://vite_rails//lib/vite_rails/tag_helpers.rb#4
module ViteRails::TagHelpers
  # Public: Resolves the path for the specified Vite asset.
  #
  # Example:
  #   <%= vite_asset_path 'calendar.css' %> # => "/vite/assets/calendar-1016838bab065ae1e122.css"
  #
  # source://vite_rails//lib/vite_rails/tag_helpers.rb#25
  def vite_asset_path(name, **options); end

  # Public: Resolves the url for the specified Vite asset.
  #
  # Example:
  #   <%= vite_asset_url 'calendar.css' %> # => "https://example.com/vite/assets/calendar-1016838bab065ae1e122.css"
  #
  # source://vite_rails//lib/vite_rails/tag_helpers.rb#33
  def vite_asset_url(name, **options); end

  # Public: Renders a script tag for vite/client to enable HMR in development.
  #
  # source://vite_rails//lib/vite_rails/tag_helpers.rb#6
  def vite_client_tag(crossorigin: T.unsafe(nil), **options); end

  # Public: Renders an <img> tag for the specified Vite asset.
  #
  # source://vite_rails//lib/vite_rails/tag_helpers.rb#72
  def vite_image_tag(name, **options); end

  # Public: Renders a <script> tag for the specified Vite entrypoints.
  #
  # source://vite_rails//lib/vite_rails/tag_helpers.rb#38
  def vite_javascript_tag(*names, type: T.unsafe(nil), asset_type: T.unsafe(nil), skip_preload_tags: T.unsafe(nil), skip_style_tags: T.unsafe(nil), crossorigin: T.unsafe(nil), media: T.unsafe(nil), **options); end

  # Public: Renders a <picture> tag with one or more Vite asset sources.
  #
  # source://vite_rails//lib/vite_rails/tag_helpers.rb#83
  def vite_picture_tag(*sources, &block); end

  # Public: Renders a script tag to enable HMR with React Refresh.
  #
  # source://vite_rails//lib/vite_rails/tag_helpers.rb#13
  def vite_react_refresh_tag(**options); end

  # Public: Renders a <link> tag for the specified Vite entrypoints.
  #
  # source://vite_rails//lib/vite_rails/tag_helpers.rb#63
  def vite_stylesheet_tag(*names, **options); end

  # Public: Renders a <script> tag for the specified Vite entrypoints.
  #
  # source://vite_rails//lib/vite_rails/tag_helpers.rb#58
  def vite_typescript_tag(*names, **options); end

  private

  # Internal: Returns the current manifest loaded by Vite Ruby.
  #
  # source://vite_rails//lib/vite_rails/tag_helpers.rb#98
  def vite_manifest; end

  # Internal: Renders a modulepreload link tag.
  #
  # source://vite_rails//lib/vite_rails/tag_helpers.rb#103
  def vite_preload_tag(*sources, crossorigin:, **options); end
end

# source://vite_rails//lib/vite_rails/version.rb#4
ViteRails::VERSION = T.let(T.unsafe(nil), String)

class ViteRuby::CLI::Build < ::ViteRuby::CLI::Vite
  include ::ViteRails::CLI::Build
end

class ViteRuby::CLI::Install < ::Dry::CLI::Command
  include ::ViteRails::CLI::Install
end

class ViteRuby::Config
  extend ::ViteRails::Config
end
