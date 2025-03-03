# typed: true

# DO NOT EDIT MANUALLY
# This is an autogenerated file for types exported from the `i18n-js` gem.
# Please instead update this file by running `bin/tapioca gem i18n-js`.


# source://i18n-js//lib/i18n-js/schema.rb#3
module I18nJS
  class << self
    # source://i18n-js//lib/i18n-js/plugin.rb#6
    def available_plugins; end

    # source://i18n-js//lib/i18n-js.rb#22
    def call(config_file: T.unsafe(nil), config: T.unsafe(nil)); end

    # source://i18n-js//lib/i18n-js/clean_hash.rb#4
    def clean_hash(hash); end

    # source://i18n-js//lib/i18n-js.rb#45
    def export_group(group); end

    # source://i18n-js//lib/i18n-js/plugin.rb#28
    def initialize_plugins!(config:); end

    # source://i18n-js//lib/i18n-js.rb#103
    def load_config_file(config_file); end

    # source://i18n-js//lib/i18n-js/plugin.rb#22
    def load_plugins!; end

    # source://i18n-js//lib/i18n-js/plugin.rb#18
    def plugin_files; end

    # source://i18n-js//lib/i18n-js/plugin.rb#10
    def plugins; end

    # source://i18n-js//lib/i18n-js/plugin.rb#14
    def register_plugin(plugin); end

    # source://i18n-js//lib/i18n-js/sort_hash.rb#4
    def sort_hash(hash); end

    # source://i18n-js//lib/i18n-js.rb#94
    def translations; end

    # source://i18n-js//lib/i18n-js.rb#73
    def write_file(file_path, translations); end
  end
end

# source://i18n-js//lib/i18n-js.rb#20
class I18nJS::MissingConfigError < ::StandardError; end

# source://i18n-js//lib/i18n-js/plugin.rb#34
class I18nJS::Plugin
  # @return [Plugin] a new instance of Plugin
  #
  # source://i18n-js//lib/i18n-js/plugin.rb#42
  def initialize(config:); end

  # This method is called whenever `I18nJS.call(**kwargs)` finishes exporting
  # JSON files based on your configuration.
  #
  # You can use it to further process exported files, or generate new files
  # based on the translations that have been exported.
  #
  # source://i18n-js//lib/i18n-js/plugin.rb#100
  def after_export(files:); end

  # Return the plugin configuration
  #
  # source://i18n-js//lib/i18n-js/plugin.rb#60
  def config; end

  # Infer the config key name out of the class.
  # If you plugin is called `MySamplePlugin`, the key will be `my_sample`.
  #
  # source://i18n-js//lib/i18n-js/plugin.rb#49
  def config_key; end

  # Check whether plugin is enabled or not.
  # A plugin is enabled when the plugin configuration has `enabled: true`.
  #
  # @return [Boolean]
  #
  # source://i18n-js//lib/i18n-js/plugin.rb#66
  def enabled?; end

  # The configuration that's being used to export translations.
  #
  # source://i18n-js//lib/i18n-js/plugin.rb#36
  def main_config; end

  # The `I18nJS::Schema` instance that can be used to validate your plugin's
  # configuration.
  #
  # source://i18n-js//lib/i18n-js/plugin.rb#40
  def schema; end

  # This method must set up the basic plugin configuration, like adding the
  # config's root key in case your plugin accepts configuration (defined via
  # the config file).
  #
  # If you don't add this key, the linter will prevent non-default keys from
  # being added to the configuration file.
  #
  # source://i18n-js//lib/i18n-js/plugin.rb#92
  def setup; end

  # This method is responsible for transforming the translations. The
  # translations you'll receive may be already be filtered by other plugins
  # and by the default filtering itself. If you need to access the original
  # translations, use `I18nJS.translations`.
  #
  # source://i18n-js//lib/i18n-js/plugin.rb#74
  def transform(translations:); end

  # In case your plugin accepts configuration, this is where you must validate
  # the configuration, making sure only valid keys and type is provided.
  # If the configuration contains invalid data, then you must raise an
  # exception using something like
  # `raise I18nJS::Schema::InvalidError, error_message`.
  #
  # source://i18n-js//lib/i18n-js/plugin.rb#83
  def validate_schema; end
end

# source://i18n-js//lib/i18n-js/schema.rb#4
class I18nJS::Schema
  # @return [Schema] a new instance of Schema
  #
  # source://i18n-js//lib/i18n-js/schema.rb#33
  def initialize(target); end

  # source://i18n-js//lib/i18n-js/schema.rb#158
  def expect_array_with_items(path:); end

  # source://i18n-js//lib/i18n-js/schema.rb#170
  def expect_required_keys(keys:, path:); end

  # source://i18n-js//lib/i18n-js/schema.rb#135
  def expect_type(path:, types:); end

  # source://i18n-js//lib/i18n-js/schema.rb#207
  def prepare_path(path:); end

  # @raise [InvalidError]
  #
  # source://i18n-js//lib/i18n-js/schema.rb#130
  def reject(error_message, node = T.unsafe(nil)); end

  # source://i18n-js//lib/i18n-js/schema.rb#188
  def reject_extraneous_keys(keys:, path:); end

  # Returns the value of attribute target.
  #
  # source://i18n-js//lib/i18n-js/schema.rb#31
  def target; end

  # source://i18n-js//lib/i18n-js/schema.rb#37
  def validate!; end

  # source://i18n-js//lib/i18n-js/schema.rb#93
  def validate_lint_scripts; end

  # source://i18n-js//lib/i18n-js/schema.rb#78
  def validate_lint_translations; end

  # source://i18n-js//lib/i18n-js/schema.rb#56
  def validate_plugins; end

  # source://i18n-js//lib/i18n-js/schema.rb#69
  def validate_root; end

  # source://i18n-js//lib/i18n-js/schema.rb#115
  def validate_translation(_translation, index); end

  # source://i18n-js//lib/i18n-js/schema.rb#107
  def validate_translations; end

  # source://i18n-js//lib/i18n-js/schema.rb#212
  def value_for(path:); end

  class << self
    # source://i18n-js//lib/i18n-js/schema.rb#21
    def required_root_keys; end

    # source://i18n-js//lib/i18n-js/schema.rb#12
    def root_keys; end

    # source://i18n-js//lib/i18n-js/schema.rb#25
    def validate!(target); end
  end
end

# source://i18n-js//lib/i18n-js/schema.rb#5
class I18nJS::Schema::InvalidError < ::StandardError; end

# source://i18n-js//lib/i18n-js/schema.rb#8
I18nJS::Schema::REQUIRED_LINT_SCRIPTS_KEYS = T.let(T.unsafe(nil), Array)

# source://i18n-js//lib/i18n-js/schema.rb#7
I18nJS::Schema::REQUIRED_LINT_TRANSLATIONS_KEYS = T.let(T.unsafe(nil), Array)

# source://i18n-js//lib/i18n-js/schema.rb#9
I18nJS::Schema::REQUIRED_TRANSLATION_KEYS = T.let(T.unsafe(nil), Array)

# source://i18n-js//lib/i18n-js/schema.rb#10
I18nJS::Schema::TRANSLATION_KEYS = T.let(T.unsafe(nil), Array)

# source://i18n-js//lib/i18n-js/version.rb#4
I18nJS::VERSION = T.let(T.unsafe(nil), String)
