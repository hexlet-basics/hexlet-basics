# frozen_string_literal: true

# https://github.com/trailblazer/reform-rails/issues/86
# Temporary fix to show errors in form
# rubocop:disable all

module Reform
  class Contract < Disposable::Twin
    class Result

      private

      # this doesn't do nested errors (e.g. )
      def filter_for(method, *args)
        @results.collect { |r| r.public_send(method, *args).to_h }
                .inject({}) { |hah, err| hah.merge(err) { |key, old_v, new_v| (new_v.is_a?(Array) ? (old_v |= new_v) : old_v.merge(new_v)) } }
                .find_all { |k, v| # filter :nested=>{:something=>["too nested!"]} #DISCUSS: do we want that here?
                  if v.is_a?(Hash)
                    nested_errors = v.select { |attr_key, val| attr_key.is_a?(Integer) && val.is_a?(Array) && val.any? }
                    v = nested_errors.to_a if nested_errors.any?
                  end
                  v.is_a?(ActiveModel::DeprecationHandlingMessageArray)
                }.to_h
      end
    end
  end
end
# rubocop:enable all
