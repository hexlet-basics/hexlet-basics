# typed: strict

# DO NOT EDIT MANUALLY
# This file was pulled from a central RBI files repository.
# Please run `bin/tapioca annotations` to update it.

class ActiveRecord::Base
  # @shim: this is included at runtime https://github.com/rails/globalid/blob/v1.0.0/lib/global_id/railtie.rb#L38
  include GlobalID::Identification
end
